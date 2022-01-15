package znet

import (
	ziface2 "go_linux/src/Genshin_Home_System/zinx/ziface"
)

type Request struct {
	//已经和客户端建立好的链接
	conn ziface2.IConnection

	//客户端请求的数据
	//data []byte
	msg ziface2.IMessage
}

func (r *Request) GetConnection() ziface2.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
