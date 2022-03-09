package ziface

import "os"

//定义一个服务器接口

type Iserver interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能给当前辅助注册一个路由方法
	AddRouter(msgID uint32, router IRouter)

	// GetConnMgr 用于conn获得serverManager
	GetConnMgr() IConnManager

	// SetOnConnStart 注册OnConnStart方法
	SetOnConnStart(func(conn IConnection))

	// SetOnConnStop 注册OnConnStop方法
	SetOnConnStop(func(conn IConnection))

	// CallOnConnStart 调用OnConnStart钩子函数的方法
	CallOnConnStart(conn IConnection)

	// CallOnConnStop 调用OnConnstop钩子函数的方法
	CallOnConnStop(conn IConnection)

	// GetSignal 监听signal信号
	GetSignal() chan os.Signal

	//// SignalHandler
	//SignalHandler()
}
