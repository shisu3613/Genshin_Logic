package ziface

//链接管理模块的抽象层
type IConnManager interface {
	// Add 创造链接
	Add(conn IConnection)
	// Remove 删除链接
	Remove(conn IConnection)
	// Get 根据链接ID查找链接
	Get(connID uint32) (IConnection, error)
	// Len 总链接个数获取
	Len() int
	// ClearConn Clear 链接清理（GC机制）
	ClearConn()
}
