package ziface

// IDataPack 封包拆包的模块，直接面向tcp链接中的数据流，用于处理tcp粘包问题
type IDataPack interface {
	// GetHeadLen 获取报的头部长度方法
	GetHeadLen() uint32

	// Pack 封包方法
	Pack(msg IMessage) ([]byte, error)

	// Unpack 拆包方法
	Unpack([]byte) (IMessage, error)
}
