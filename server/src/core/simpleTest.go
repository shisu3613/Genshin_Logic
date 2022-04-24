package main

import (
	"fmt"
	"os"
	"os/signal"
	"server/utils"
	"sync"
	"time"
)

type SyncPID struct {
	PID     int `json:"PID"`
	TestVal int
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			handle()
		}
	}
}

func (t *Task) Stop() {
	close(t.closed)
}

func handle() {
	for i := 0; i < 5; i++ {
		fmt.Print("#")
		time.Sleep(time.Millisecond * 200)
	}
	fmt.Println()
}

func main() {
	//Map := make(map[int]interface{})
	//s := &SyncPID{PID: 99999, TestVal: 1111}
	//Map[1] = s
	//
	//newS := Map[1]
	//

	// @Modified By WangYuding 2022/4/24 16:54:00
	// @Modified description 同步信号的test
	//tests := SyncPID{PID: 88888}
	//msg1, _ := json.Marshal(s)
	//msg2, _ := json.Marshal(tests)
	////json.Unmarshal(msg, newS)
	//fmt.Println(string(msg1))
	//fmt.Println(string(msg2))
	//json.Unmarshal(msg2, newS)
	//fmt.Println(newS)

	//capture屏幕信号的test
	test := func(a int) { fmt.Println(a) }
	fmt.Println("test" + utils.CaptureOutput(func() {
		test(2)
	}))

	//signal捕获crtl+c的测试
	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Second * 2),
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	task.wg.Add(1)
	go func() { defer task.wg.Done(); task.Run() }()
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
		//case <-task.closed:
		//	task.wg.Done()
		//	//case <-task.ticker.C:
		//	//	handle()
	}
	println("wait for stop")
	task.wg.Wait()
}
