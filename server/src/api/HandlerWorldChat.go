package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
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
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))

	//获得聊天信息
	var msg string
	_ = json.Unmarshal(request.GetData(), &msg)
	if msg == "-1" {
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
	} else {
		//player轮询世界角色管理器发起广播
		for _, anotherPlayer := range game.WorldMgrObj.GetAllPlayers() {
			if anotherPlayer != player {
				anotherPlayer.SendStringMsg(19, msg)
			}
		}
		player.SendStringMsg(9, game.WorldChatStr)
	}
}
