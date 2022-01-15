package ziface

import "net"

//定义链接模块的抽象类

type IConnection interface {
	// Start 启动链接
	Start()

	// Stop 停止链接 结束当前链接的工作
	Stop()

	// GetTCPConnection 获取绑定的socket
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取链接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态
	RemoteAddr() net.Addr

	// SendMsg Send 发送数据
	SendMsg(msgId uint32, data []byte) error

	// SetProperty 设置链接属性
	SetProperty(key string, value interface{})

	// GetProperty 获取链接属性
	GetProperty(key string) (interface{}, error)

	// RemovePorperty 移除链接属性
	RemovePorperty(key string)
}

// HandleFunc 处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
