package main

import (
	"fmt"
	"server/zinx/ziface"
	"server/zinx/znet"
)

//基于Zinx框架开发的服务器应用程序

// PingRouter addRouter 允许自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// PreHandle Tset PreHandle
//func (pr *PingRouter) PreHandle(request ziface.IRequest) {
//
//	fmt.Println("Call Router PreHandler")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before Ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping error", err)
//	}
//
//}

// Handler Test Handler
func (pr *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handler")
	//先读取客户端的数据，然后再回写ping.ping.ping
	fmt.Println("Recv from client: msgID=", request.GetMsgID(),
		",data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

// PreHandle Tset PreHandle
//func (pr *PingRouter) PreHandle(request ziface.IRequest) {
//
//	fmt.Println("Call Router PreHandler")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before Ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping error", err)
//	}
//
//}

// Handler Test Handler
func (pr *HelloRouter) Handler(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handler")
	//先读取客户端的数据，然后再回写ping.ping.ping
	fmt.Println("Recv from client: msgID=", request.GetMsgID(),
		",data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx!"))
	if err != nil {
		fmt.Println(err)
	}
}

//PostHandler Test PostHandler
//func (pr *PingRouter) PostHandler(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandler")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after Ping...\n"))
//	if err != nil {
//		fmt.Println("call back after ping error", err)
//	}
//}

// DoConnectionBegin 创建链接之后使用的HOOK函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=======>DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin")); err != nil {
		fmt.Println(err)
	}
	//在连接创建之前，给链接设置一些属性
	fmt.Println("Set conn Name, Home....")
	conn.SetProperty("Name", "YudingWang")
	conn.SetProperty("Home", "top.yudingwang")

}

// DoConnectionLost 断开链接之前需要执行的Hook
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("=====>DoConnectionLost is Called....")
	fmt.Println("Conn ID =", conn.GetConnID(), " is Lost")

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name:", name)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home:", home)
	}
}

func main() {
	//1.创建server的句柄
	s := znet.NewServer("[zinx V0.1]")

	//2.给当前zinx框架添加多个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//注册Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.启动server
	s.Serve()
}
