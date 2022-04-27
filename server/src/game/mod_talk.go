package game

import (
	"encoding/json"
	"server/DB/RedisTool"
	"sort"
	"sync"
	"sync/atomic"
)

/**
    @author: WangYuding
    @since: 2022/4/27
    @desc: //聊天模块加载，用redis加载保存一个月内的聊天记录
**/

type ChatMsg struct {
	Uid    string `json:"uid"` //消息发送者的UID
	IdTime string `json:"idTime"`
	Cnt    string `json:"cnt"`
	SendTo string `json:"sendTo"` //消息到达者的UID,global情况储存到db1，其他储存到db2
}

type MsgSlice []*ChatMsg

func (s MsgSlice) Len() int           { return len(s) }
func (s MsgSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MsgSlice) Less(i, j int) bool { return s[i].IdTime < s[j].IdTime }

type ModTalk struct {
	GlobalChat  MsgSlice
	PrivateChat map[string]MsgSlice
	player      *Player
	atomicFlag  int32
	talkLock    sync.RWMutex //读写锁
}

func (mt *ModTalk) SaveData() {

}

func (mt *ModTalk) LoadData() {
	mt.GlobalChat = mt.GetGlobalMessage()
}

// init
// @Description: 初始化切片和map
// @receiver mt
func (mt *ModTalk) init(player *Player) {
	mt.player = player
	mt.PrivateChat = make(map[string]MsgSlice)
	atomic.StoreInt32(&mt.atomicFlag, 0)
	mt.talkLock = sync.RWMutex{}
}

func (mt *ModTalk) SetFlag() {
	atomic.StoreInt32(&mt.atomicFlag, 1)
}

func (mt *ModTalk) ResetFlag() {
	atomic.StoreInt32(&mt.atomicFlag, 0)
}

func (mt *ModTalk) CheckFlag() bool {
	if mt.atomicFlag == 1 {
		return true
	}
	return false
}

// GetGlobalMessage
// @Description: 获取redis里面保存的全局聊天内容
// @receiver mt
// @return []*ChatMsg
func (mt *ModTalk) GetGlobalMessage() MsgSlice {
	var Slice MsgSlice
	b := RedisTool.GetAllKeys(1)
	for _, idtime := range b {
		m, _ := mt.GetMessage(idtime)
		Slice = append(Slice, &m)
	}
	sort.Stable(Slice)
	return Slice
}

// GetGlobalHistory
// @Description: 获得历史聊天记录并用string形式记录
// @receiver mt
// @return string
func (mt *ModTalk) GetGlobalHistory() string {
	res := "==================历史聊天====================\n"
	mt.talkLock.RLock()
	for _, msg := range mt.GlobalChat {
		res += "时间：" + msg.IdTime + "," + msg.Uid + ":" + msg.Cnt + "\n"
	}
	mt.talkLock.RUnlock()
	res += "==============================================\n"
	return res
}

func (mt *ModTalk) GetMessage(idtime string) (ChatMsg, bool) {
	val, err := RedisTool.GetValueByKey(1, idtime)
	if nil == err {
		var msg ChatMsg
		err := json.Unmarshal([]byte(val), &msg) //反序列化
		RedisTool.CheckError(err)
		return msg, true
		//在无err情况下返回Message，并设置状态为true
		//true表示获取成功
	} else {
		return ChatMsg{"", "", "", ""}, false
		//在err情况下返回空的Message，并设置状态为false
		//false表示获取失败
	}
}

// SetGlobalMessage
// @Description: 将数据保存到redis数据库里面
// @receiver mt
// @param msg
// @return bool
func (mt *ModTalk) SetGlobalMessage(msg ChatMsg) bool {
	data, err := json.Marshal(msg)
	if err == nil {
		return RedisTool.SetRecord(1, msg.IdTime, data)
	} else {
		return false
	}
}

// AddGlobalMessage
// @Description: 添加对话到缓存里面
// @receiver mt
// @param msg
func (mt *ModTalk) AddGlobalMessage(msg ChatMsg) {
	mt.talkLock.Lock()
	defer mt.talkLock.Unlock()
	mt.GlobalChat = append(mt.GlobalChat, &msg)
}
