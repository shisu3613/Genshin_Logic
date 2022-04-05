package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type HandlerWishesTest struct {
	znet.BaseRouter
}

func (hb *HandlerWishesTest) Handler(request ziface.IRequest) {
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	var times int
	_ = json.Unmarshal(request.GetData(), &times)
	player.SendStringMsg(0, player.ModWish.DoPoolTest(times))
	player.SendStringMsg(2, player.GetMod(game.ModPlay).(*game.ModPlayer).Name+",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
}
