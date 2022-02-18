package znet

import (
	"server/zinx/ziface"
)

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

//GetMsgId 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

//SetMsgId 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//SetDataLen 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func NewMsgPackage(id uint32, data []byte) ziface.IMessage {
	msg := &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	return msg
}
