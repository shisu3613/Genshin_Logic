package ziface

/*
	消息管理抽象层
*/

type IMsgHandler interface {
	// DoMsgHandler 调度。执行对应的router消息处理方法
	DoMsgHandler(request IRequest)

	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// StartWorkerPool StartWorkPool  启动一个worker工作池子,开启工作池子的动作智能发生一次，一个zinx框架只能有一个worker池子，对外暴露的方法
	StartWorkerPool()

	// SendMsgToTaskQueue 将消息交给taskqueue来由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
