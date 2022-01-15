package znet

import (
	"fmt"
	"go_linux/src/Genshin_Home_System/zinx/utils"
	ziface2 "go_linux/src/Genshin_Home_System/zinx/ziface"
	"strconv"
)

/*
	消息处理模块的实现
*/

type MsgHandler struct {
	//存放每个msgID所对应的处理方法
	APIs map[uint32]ziface2.IRouter

	//负责worker取任务的消息队列
	TaskQueue []chan ziface2.IRequest

	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化msghandler模块
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs: make(map[uint32]ziface2.IRouter),
		//从全局配置中获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//根据当前池子的大小初始化消息队列
		TaskQueue: make([]chan ziface2.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度。执行对应的router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface2.IRequest) {
	router, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API msgID =", request.GetMsgID(), "is not found! Need register!")
	}
	router.PreHandle(request)
	router.Handler(request)
	router.PostHandler(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface2.IRouter) {
	//1.判断当前的msg绑定的api是否已经存在
	if _, ok := mh.APIs[msgID]; ok {
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID))) //直接宕机
	}
	//2.添加id和api的绑定关系
	mh.APIs[msgID] = router
	fmt.Println("Add api MsgID =", msgID, "succ!")

}

// StartWorkerPool 启动一个worker工作池子,开启工作池子的动作智能发生一次，一个zinx框架只能有一个worker池子，对外暴露的方法
func (mh *MsgHandler) StartWorkerPool() {
	//根据workPoolSize分别开启worker
	var i uint32
	for i = 0; i < mh.WorkerPoolSize; i++ {
		//一个worker被启动
		//1.开辟当前worker对应的消息队列，开辟空间,消息队列的最大值在配置文件中配置
		mh.TaskQueue[i] = make(chan ziface2.IRequest, utils.GlobalObject.MaxTaskQueueLen)
		//2.启动当前的worker，阻塞等待消息从channel中传过来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}

}

//启动一个worker工作流程,不对外暴露的方法
func (mh *MsgHandler) startOneWorker(workID uint32, taskQueue chan ziface2.IRequest) {
	fmt.Println("worker ID =", workID, "is started ...")
	//阻塞的同时不断地轮询消息队列地信息
	for {
		select {
		//如果由消息过来,用for-select而不是for range是为了扩展性留有空间
		case request := <-taskQueue:
			mh.DoMsgHandler(request)

		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue来由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface2.IRequest) {
	//负载均衡算法，可以改进，目前使用单体的轮询
	//根据客户端建立的conn的ID进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), "request MsgID =", request.GetMsgID(), "to WorkerID =", workerID)

	//将消息发送给对应的worker的taskQueue
	mh.TaskQueue[workerID] <- request
}
