package znet

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type f func()

type sig struct{}

const DefaultExpiryTime = time.Second

type Pool struct {
	capacity int32 //大小
	running  int32 //活着的goroutine 数量

	//expiryDuration 是 worker 的过期时长，在空闲队列中的 worker 的最新一次运行时间与当前时间之差如果大于这个值则表示已过期，定时清理任务会清理掉这个 worker
	expiryDuration time.Duration

	//workers 一个Worker就是一个Goroutine，消息队列切片
	workers []*Worker

	//通知池子关闭，关闭所有worker退出运行以防止goroutine 泄露
	release chan sig

	//lock 同步锁
	lock sync.Mutex

	//once 用于确保全局池子不会重复创建和关闭
	once sync.Once
}

func NewPool(size int) (*Pool, error) {
	return NewTimingPool(size, int(DefaultExpiryTime))
}

func NewTimingPool(size, expiry int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("invalid PoolSize")
	}
	if expiry <= 0 {
		return nil, errors.New("invalid expiry set")
	}
	p := &Pool{
		capacity:       int32(size),
		expiryDuration: time.Duration(expiry) * time.Second,
		release:        make(chan sig, 1),
	}
	return p, nil
}

func (p *Pool) Submit(task f) error {
	if len(p.release) > 0 {
		return errors.New("pool has already closed")
	}
	w := p.getWorker()
	w.task <- task

	return nil
}

func (p *Pool) getWorker() *Worker {
	var w *Worker
	waitingFlag := false
	//调度算法，加锁,后进先出
	p.lock.Lock()
	w = p.getLastWorker()
	if w == nil {
		//判断当前运行worker数量是否已经超出限制
		waitingFlag = p.Running() >= p.Cap()
	}
	p.lock.Unlock()

	//等待资源的逻辑
	//如果当前pool已经满了，新请求等待
	if waitingFlag {
		for {
			p.lock.Lock()
			w = p.getLastWorker()
			if w == nil {
				continue
			}
			p.lock.Unlock()
			break
		}
		//对应源码中当pool还没有满但没有空闲worker的情况
	} else if w == nil {
		w = &Worker{
			pool: p,
			task: make(chan f, 1),
		}
		w.run()
		p.addRunning(1)
	}
	return w
}

func (p *Pool) getLastWorker() (w *Worker) {
	curWorkers := p.workers
	//判断有无可行的worker,如没有将flag置为nil
	l := len(curWorkers)
	if l == 0 {
		//判断当前运行worker数量是否已经超出限制
		return nil
	} else {
		w = curWorkers[l-1]
		//优秀的代码习惯，手动清空资源
		curWorkers[l-1] = nil
		p.workers = curWorkers[:l-1]
	}
	return
}

func (p *Pool) putWorkerBack(worker *Worker) {
	worker.recycleTime = time.Now()
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.lock.Unlock()

}

//====================================Worker==============

type Worker struct {
	pool        *Pool
	task        chan f
	recycleTime time.Time
}

// run
// Description:结合前面的 p.Submit(task f) 和 p.getWorker() ，提交任务到 Pool 之后，获取一个可用 worker，每新建一个 worker 实例之时都需要调用 w.run() 启动一个 goroutine 监听 worker 的任务列表 task ，一有任务提交进来就执行
// receiver w
func (w *Worker) run() {
	go func() {
		for f := range w.task {
			if f == nil {
				w.pool.addRunning(-1)
				return
			}
			f()
			//重点：这里应该加上一个逻辑当当前worker执行完之后，解绑放回POOL
			w.pool.putWorkerBack(w)
		}
	}()

}

//-------------------------------------------- 一些原子操作

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

func (p *Pool) addRunning(delta int) {
	atomic.AddInt32(&p.running, int32(delta))
}
