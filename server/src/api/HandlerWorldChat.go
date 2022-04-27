package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
	"strconv"
	"time"
)

/**
    @author: WangYuding
    @since: 2022/4/25
    @desc: //世界聊天模块，哈哈哈，都做了框架了怎么能不试试世界聊天呢
**/

type HandlerWorldChat struct {
	znet.BaseRouter
}

func (hw *HandlerWorldChat) Handler(request ziface.IRequest) {
	// 第一步根据uid 获得player对象
	PID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	//log.Println("UserID is:" + strconv.Itoa(UserID.(int)))
	player := game.WorldMgrObj.GetPlayerByPID(PID.(int))

	//获得聊天信息
	var msg string
	_ = json.Unmarshal(request.GetData(), &msg)
	if msg == "exit;" {
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
	} else {
		//player轮询世界角色管理器发起广播
		uid := player.GetUserID()
		newMsg := game.ChatMsg{
			Uid:    strconv.Itoa(uid),
			IdTime: time.Now().Format("2006-01-02 15:04:05"),
			Cnt:    msg,
			SendTo: "Global",
		}
		//保存对话信息到数据库
		player.GetMod(game.TalkMod).(*game.ModTalk).SetGlobalMessage(newMsg)
		for _, anotherPlayer := range game.WorldMgrObj.GetAllPlayers() {
			//保存到在线玩家modTalk缓存里
			anotherPlayer.GetMod(game.TalkMod).(*game.ModTalk).AddGlobalMessage(newMsg)
			//处理msg信息
			//如果目标在对话功能里面：
			//直接发送对话信息
			//否则发送：您有一条新信息
			if anotherPlayer.GetMod(game.TalkMod).(*game.ModTalk).CheckFlag() {
				anotherPlayer.SendStringMsg(0, "时间："+newMsg.IdTime+","+newMsg.Uid+":"+newMsg.Cnt)
			} else {
				anotherPlayer.SendStringMsg(0, "您收到一条新的世界聊天")
			}
			//anotherPlayer.SendStringMsg(0,)
		}
		//player.SendStringMsg(9, game.WorldChatStr)
	}
}
