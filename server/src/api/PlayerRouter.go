package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type PlayerRouter struct {
	znet.BaseRouter
}

// Handler Handler For Genshin Player
func (pr *PlayerRouter) Handler(request ziface.IRequest) {
	//这里是主界面的选择程序
	//找到对应的player
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	//player.SendStringMsg(2, "欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
	//player.SendStringMsg(2,"欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")

	var msgChoice int
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	switch msgChoice {
	case 1:
		//player.HandleBaseRemote(request)
		player.SendStringMsg(3, game.BasicLogicStr)
	case 2:
		player.SendStringMsg(5, game.BagLogicStr)
	case 3:
		player.SendStringMsg(33, "请输出抽卡次数:")
	case 4:
		player.SendStringMsg(6, game.WishLogicStr)
		//player.HandleWishUp()
	case 5:
		player.HandleMap()
	case 7:
		player.SendStringMsg(9, game.WorldChatStr)
	default:
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
	}
}
