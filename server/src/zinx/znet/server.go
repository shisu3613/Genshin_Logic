package znet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"server/zinx/utils"
	"server/zinx/ziface"
	"syscall"
	"time"
)

const (
	MsgLen   = 4096
	inActive = time.Minute * 20
)

// Server IServer的接口实现，定义一个server的模块
type Server struct {
	//服务器的名称
	Name string

	//服务器的端口
	Port int

	//服务器的IP
	IP string

	//服务器的IP版本
	IPVersion string

	//当前的server添加一个router，server注册的链接对应的处理业务(单router)
	//Router ziface.IRouter

	//当前server的消息管理模块，用来绑定msgid和对应的业务处理模块
	MsgHandler ziface.IMsgHandler

	//链接模块的管理模块
	ConnMgr ziface.IConnManager

	//该Server创建链接之后自动调用Hook函数 -- OnConnStart
	OnConnStart func(conn ziface.IConnection)

	//该Server销毁链接之前自动调用Hook函数 -- OnConnStop
	OnConnStop func(conn ziface.IConnection)

	//增加监听信号功能，syscall
	SignalChan chan os.Signal
}

// CallBackToClient 定义当前客户端所绑定的handlerAPI，目前是写死的，一抹后话应该由用户自定义
//func CallBackToClient(conn *net.TCPConn, data []byte,cnt int) error{
//	fmt.Println("[Conn Handler] CallbackToClient")
//	if _,err := conn.Write(data[:cnt]); err!=nil{
//		fmt.Println("Write back buf err",err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

// Start 启动服务器
func (s *Server) Start() {
	//开启消息队列和worker工作池
	s.MsgHandler.StartWorkerPool()

	fmt.Printf("Server Name : %s,listener at IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("Version : %s,MaxConn : %d, MaxPackageSize:%d is starting\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	//fmt.Printf("[Start] Server Listener at IP : %s,Port %d, is starting\n", s.IP, s.Port)
	go func() {
		//1.获取tcp的addr
		Addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Address get error", err)
		}
		//2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, Addr)
		if err != nil {
			fmt.Println("net.Listen err:", err)
			return
		}
		//close listen socket
		defer listener.Close()

		fmt.Println("start Zinx server ", s.Name, "success, Listening...")
		var cid uint32 = 0
		for {
			//accept
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Listener accept err:", err)
				continue
			}
			//设置最大链接个数的判断，如果超出最大数量，则关闭链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//todo 给用户端一个sendMsg说明超出最大数量
				fmt.Println("Too Many Connection MaxConn =", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新链接的业务方法和conn进行绑定，得到我面的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动业务处理
			dealConn.Start()

		}
	}()

}

// Stop 停止服务器
func (s *Server) Stop() {
	//尝试将服务器的一些资源，状态等停止
	fmt.Println("[STOP] Server Name :", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动sever的功能
	s.Start()

	//阻塞
	signal.Notify(s.SignalChan, syscall.SIGINT)
	select {
	case <-s.SignalChan:
		fmt.Println("[!!!!]Close the server after 3 seconds")
		time.Sleep(time.Second * 3)
		s.Stop()
	}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("add Router successfully")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) GetSignal() chan os.Signal {
	return s.SignalChan
}

// NewServer 初始化Serve模块的方法
func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		Port:      utils.GlobalObject.TcpPort,
		IP:        utils.GlobalObject.Host,
		IPVersion: "tcp4",
		//Router:    nil,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
		SignalChan: make(chan os.Signal),
	}

	return s
}

// SetOnConnStart 注册OnConnStart方法
func (s *Server) SetOnConnStart(hook func(conn ziface.IConnection)) {
	s.OnConnStart = hook
}

// SetOnConnStop 注册OnConnStop方法
func (s *Server) SetOnConnStop(hook func(conn ziface.IConnection)) {
	s.OnConnStop = hook
}

// CallOnConnStart 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用OnConnstop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	s.ConnMgr.Remove(conn)
	if s.OnConnStop != nil {
		fmt.Println("----->Call OnConnStop()...")
		s.OnConnStop(conn)
	}

}
