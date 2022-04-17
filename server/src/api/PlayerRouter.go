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
		player.SendStringMsg(3, "当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")
	case 2:
		player.SendStringMsg(5, "当前处于背包界面,请选择操作：0返回1增加物品2扣除物品3使用物品")
	case 3:
		player.SendStringMsg(33, "请输出抽卡次数:")
	case 4:
		player.SendStringMsg(6, "您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")
		//player.HandleWishUp()
	case 5:
		player.HandleMap()
	default:
		player.SendStringMsg(2, "欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
	}
}
