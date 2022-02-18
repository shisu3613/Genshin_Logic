package ziface

/*
	IRequest接口：
	实际上是将客户端请求的链接信息和请求的数据包装到一个request中
*/

type IRequest interface {
	// GetConnection 得到当前链接
	GetConnection() IConnection

	// GetData 得到当前数据
	GetData() []byte

	// GetMsgID 添加msg方法
	GetMsgID() uint32
}
