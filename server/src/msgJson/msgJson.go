package msgJson

import (
	"encoding/json"
	"fmt"
	"net"
	"server/zinx/znet"
	"sync"
)

// GlobalMsgManager 用于管理消息数据结构的全局模块
type GlobalMsgManager struct {
	pLock sync.RWMutex //保护玩家的互斥机制
	//MsgStruct map[uint32]interface{}
}

type SyncUID struct {
	UID int `json:"UID"`
}

type SyncContent struct {
	Content string `json:"Content"`
}

func (gm *GlobalMsgManager) SendMsg(msgId uint32, data interface{}, conn net.Conn) {
	binaryData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("json.Marshal msg %d failed,err:%s", msgId, err)
	}
	dp := znet.NewDataPack()
	binaryMsg, err := dp.Pack(znet.NewMsgPackage(msgId, binaryData))
	if err == nil {
		_, _ = conn.Write(binaryMsg)
	} else {
		fmt.Printf("sendMsg %d failed,err:%s", msgId, err)
	}
}

// MsgMgrObj  提供了一个对外的世界管理模块的句柄
var MsgMgrObj *GlobalMsgManager

//提供了worldManager初始化的方法
func init() {
	MsgMgrObj = &GlobalMsgManager{
		pLock: sync.RWMutex{},
		//MsgStruct: make(map[uint32]interface{}),
	}
	//MsgMgrObj.MsgStruct[1] = &SyncPID{}

}
