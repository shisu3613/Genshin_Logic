package ziface

//将请求的消息封装到message里面，定义一个抽象层模块

type IMessage interface {
	// GetMsgId 获取消息的ID
	GetMsgId() uint32

	// GetMsgLen 获取消息的长度
	GetMsgLen() uint32

	// GetData 获取消息的内容
	GetData() []byte

	// SetMsgId 设置消息的ID
	SetMsgId(uint32)

	// SetData 设置消息的内容
	SetData([]byte)

	// SetDataLen 设置消息的长度
	SetDataLen(uint32)
}
