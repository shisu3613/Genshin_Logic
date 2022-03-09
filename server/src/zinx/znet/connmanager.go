package znet

import (
	"errors"
	"fmt"
	"server/zinx/ziface"
	"sync"
)

//链接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理来链接集合
	connLock    sync.RWMutex                  //读写锁
}

// NewConnManager 初始化函数
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 创造链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入ConnManager中]
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection ID =", conn.GetConnID(), " add to connManager successfully :conn num =", connMgr.Len())
}

//// Remove 删除链接
//func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
//	//保护共享资源
//	//fmt.Println("any problem here?")
//	connMgr.connLock.Lock()
//	defer connMgr.connLock.Unlock()
//
//	//删除链接信息
//	//fmt.Println("any problem here?")
//	delete(connMgr.connections, conn.GetConnID())
//	fmt.Println("connection ID =", conn.GetConnID(), " remove from connManager successfully :conn num =", connMgr.Len())
//
//}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {

	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

// Get 根据链接ID查找链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {

	//保护共享资源,读锁就可以
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection not FOUND!")
	}
}

// Len 总链接个数获取
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)

}

// ClearConn Clear 链接清理（GC机制）
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Connections succ! conn num =", connMgr.Len())
}
