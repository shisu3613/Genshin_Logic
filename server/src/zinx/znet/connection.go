package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"server/zinx/utils"
	"server/zinx/ziface"
	"sync"
)

// Connection 当前链接的模块
type Connection struct {
	//V0,9当前conn是属于哪个Server
	TcpServer ziface.Iserver

	//当前链接的socket
	Conn *net.TCPConn

	//当前的额ID
	ConnID uint32

	//当前状态
	IsClosed bool

	//当前的业务处理api
	HandleAPI ziface.HandleFunc

	//告知退出的channel:由reader告知writer退出
	ExitChan chan bool

	//WriteExitChan chan bool

	//无缓冲的管道，用于读写groutine之间的消息通讯
	msgChan chan []byte

	//该链接处理的方法
	//Router ziface.IRouter

	//消息的管理msgid对应的api对应关系
	MsgHandler ziface.IMsgHandler

	//链接属性集合
	property map[string]interface{}

	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接模块的方法
func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {

	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		//WriteExitChan: make(chan bool, 1),
		msgChan:  make(chan []byte), //无缓冲
		property: make(map[string]interface{}),
	}

	//将Conn加入到connManager中,
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// StartWrite 链接的写业务方法，用户将消息发送给客户端消息的模块
func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	//不断地阻塞的等待channel的消息，进行写给客户端
	//todo 可以添加写业务
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出，此时writer也要退出
			return
		}
	}
}

// StartRead 链接的读业务方法
func (c *Connection) StartRead() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader is exit!],c.connID:", c.ConnID, "remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	//创建一个拆包解包的对象
	dp := NewDataPack()

	//读取客户端的数据到buff中
	for {
		//select {
		//case <-c.ExitChan:
		//	//代表reader已经退出
		//	//c.WriteExitChan <- true
		//	return
		//default:

		//读取msgHead 8bytes
		headData := make([]byte, dp.GetHeadLen())
		//fmt.Println("readFull Test")
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read massage Head error:", err)
			break
		}

		//将msg进行拆包得到msgID和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("msg struct err", err)
			return
		}

		//根据dataLen再次读取data
		var data []byte
		if msg.GetMsgLen() > 0 {
			//说明由数据的，需要进行第二次读取
			//第二次读把data读出来
			data = make([]byte, msg.GetMsgLen())
			//跟据dataLen的长度再次从IO流中读取
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("跟据dataLen的长度再次从IO流中读取失败:", err)
				return
			}
		}
		msg.SetData(data)

		//读取客户端的数据到buff中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receiver buf err", err)
		//	continue
		//}

		//调用当前链接所绑定的handleApi
		//if err:=c.HandleAPI(c.Conn,buf,cnt); err!=nil{
		//	fmt.Println("c.connID:",c.ConnID,"Handler is error",err,"remote addr is", c.RemoteAddr().String())
		//	break
		//}

		//得到当前链接数据的request对象
		req := &Request{
			conn: c,
			msg:  msg,
		}

		//执行注册路由的方法
		//go c.MsgHandler.DoMsgHandler(req)

		//先判断是否开启工作池子
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启工作池子
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			//没有开启工作池的情况
			//执行注册路由的方法
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
	//}
}

// Start 启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID:", c.ConnID)

	//启动从当前链接的读业务
	go c.StartRead()
	//启动当前链接的写业务
	go c.StartWrite()

	//按照开发者传递进来的处理业务
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	fmt.Println("Conn Stop() ... ConnID -", c.ConnID)

	//调用开发者注册的要在销毁链接之前需要执行的业务部分
	c.TcpServer.CallOnConnStop(c)

	c.IsClosed = true
	//告知writer关闭
	c.ExitChan <- true

	//time.Sleep(10 * time.Second)
	//关闭链接，关闭资源
	//fmt.Println("start to close")

	//将链接从连接管理器中删除
	//c.TcpServer.GetConnMgr().Remove(c) //删除conn从ConnManager中
	//fmt.Println("start to close2")

	err := c.Conn.Close()
	if err != nil {
		fmt.Println("close connection failed", err)
		return
	}

	//fmt.Println("all the finish2")
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
	//fmt.Println("all the finish")
}

// GetTCPConnection 获取绑定的socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg Send 发送数据
func (c *Connection) SendMsg(MsgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("connection closed when send msg")
	}
	//将data进行封包
	dp := NewDataPack()
	sendMsg, err := dp.Pack(&Message{
		Id:      MsgId,
		DataLen: uint32(len(data)),
		Data:    data,
	})
	if err != nil {
		fmt.Println("error msgID :", MsgId)
		return errors.New("client pack err")
	}
	//将数据直接发送给客户端，V0.7的时候读写分离，将数据发送给channel
	//_, err = c.Conn.Write(sendMsg)
	//if err != nil {
	//	fmt.Println("Write error msgID :", MsgId)
	//	return errors.New("conn Write err")
	//}
	//fmt.Println("Send Msg Successful")

	c.msgChan <- sendMsg
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个属性
	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if val, ok := c.property[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("No the property named" + key)
	}
}

// RemovePorperty 移除链接属性
func (c *Connection) RemovePorperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
