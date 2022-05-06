package game

import (
	"encoding/json"
	"fmt"
	"log"
	"server/DB/RedisTool"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	globalDB = 1
	personDB = 2
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
	PrivateChat map[int]MsgSlice
	player      *Player
	statusFlag  int32
	talkLock    sync.RWMutex //读写锁
}

func (mt *ModTalk) SaveData() {

}

func (mt *ModTalk) LoadData() {
	mt.GlobalChat = mt.GetGlobalHistory()
	mt.loadP2POffline()
}

// init
// @Description: 初始化切片和map
// @receiver mt
func (mt *ModTalk) init(player *Player) {
	mt.player = player
	mt.PrivateChat = make(map[int]MsgSlice)
	atomic.StoreInt32(&mt.statusFlag, 0)
	mt.talkLock = sync.RWMutex{}
}

func (mt *ModTalk) SetFlag(num int32) {
	atomic.StoreInt32(&mt.statusFlag, num)
}

func (mt *ModTalk) ResetFlag() {
	atomic.StoreInt32(&mt.statusFlag, 0)
}

func (mt *ModTalk) CheckFlag(num int32) bool {
	if mt.statusFlag == num {
		return true
	}
	return false
}

// GetGlobalHistory
// @Description: 获取redis里面保存的全局聊天内容
// @receiver mt
// @return []*ChatMsg
func (mt *ModTalk) GetGlobalHistory() MsgSlice {
	var Slice MsgSlice
	b := RedisTool.GetAllKeys(globalDB)
	for _, key := range b {
		m, _ := mt.GetMessage(key, globalDB)
		Slice = append(Slice, &m)
	}
	sort.Stable(Slice)
	return Slice
}

// loadP2POffline
// @Description: 初始化时记录是否有离线消息加载曾经通话对象
// @receiver mt
// @param pattern( example: 1:2:*)
func (mt *ModTalk) loadP2POffline() {
	curUID := mt.player.GetUserID()
	lastLogout := mt.player.GetMod(ModPlay).(*ModPlayer).UpdatedAt.Format("2006-01-02 15:04:05")
	for _, x := range WorldMgrObj.GetAllPlayersUID() {
		var key string
		if x < curUID {
			key = strconv.Itoa(x) + ":" + strconv.Itoa(curUID)
		} else if x > curUID {
			key = strconv.Itoa(x) + ":" + strconv.Itoa(curUID)
		}
		if key != "" && RedisTool.CheckKeyExists(personDB, key) {
			var lastMsg ChatMsg
			msg := RedisTool.GetLastListVal(personDB, key)
			err := json.Unmarshal([]byte(msg), &lastMsg)
			if err != nil {
				panic(err)
			}
			log.Println(lastMsg)
			if lastMsg.IdTime > lastLogout {
				mt.PrivateChat[x] = make(MsgSlice, 1)
			} else {
				mt.PrivateChat[x] = nil
			}
		}
	}
}

// HandlerOfflineMsg
// @Description: 在登录时提示离线消息
// @receiver mt
// @return string
func (mt *ModTalk) HandlerOfflineMsg() string {
	output := ""
	for k, v := range mt.PrivateChat {
		if len(v) > 0 {
			output += fmt.Sprintf("UId:%d在您离线时发送过消息给您\n", k)
		}
	}
	return output
}

func (mt *ModTalk) checkOfflineMessage(lastLogout string) bool {
	return true
}

// PrintGlobalHistory
// @Description: 获得历史聊天记录并用string形式记录
// @receiver mt
// @return string
func (mt *ModTalk) PrintGlobalHistory() string {
	res := "==================历史聊天====================\n"
	mt.talkLock.RLock()
	for _, msg := range mt.GlobalChat {
		res += "时间：" + msg.IdTime + "," + msg.Uid + ":" + msg.Cnt + "\n"
	}
	mt.talkLock.RUnlock()
	res += "==============================================\n"
	return res
}

func (mt *ModTalk) PrintP2PHistory(uid int) string {
	res := "==================历史聊天====================\n"
	curUID := mt.player.GetUserID()
	min, max := uid, curUID
	if curUID < uid {
		min, max = curUID, uid
	}
	MsgList, _ := RedisTool.GetListByKey(personDB, fmt.Sprintf("%d:%d", min, max))
	for _, msg := range MsgList {
		var Msg ChatMsg
		err := json.Unmarshal([]byte(msg), &Msg)
		if err != nil {
			panic(err)
		}
		res += "时间：" + Msg.IdTime + "," + Msg.Uid + ":" + Msg.Cnt + "\n"
	}

	res += "==============================================\n"
	return res
}

func (mt *ModTalk) GetMessage(key string, db int) (ChatMsg, bool) {
	val, err := RedisTool.GetValueByKey(db, key)
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
		return RedisTool.SetRecord(1, msg.Uid+":"+msg.IdTime, data)
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

// SetPrivateMessage
// @Description: 将数据保存到redis数据库里面
// @receiver mt
// @param msg
// @return bool
func (mt *ModTalk) SetPrivateMessage(msg ChatMsg) bool {
	data, err := json.Marshal(msg)
	min, max := msg.Uid, msg.SendTo
	if min > max {
		min, max = max, min
	}
	if err == nil {
		return RedisTool.SetRecordList(personDB, min+":"+max, string(data))
	} else {
		return false
	}
}
