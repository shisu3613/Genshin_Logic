package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type HandlerBag struct {
	znet.BaseRouter
}

func (hb *HandlerBag) Handler(request ziface.IRequest) {
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	var modChoose int
	_ = json.Unmarshal(request.GetData(), &modChoose)
	switch modChoose {
	case 0: //返回操作
		player.SendStringMsg(2, player.GetMod(game.ModPlay).(*game.ModPlayer).Name+game.MainLogicStr)
		goto END
	case 1:
		//增加物品模块
		player.SendStringMsg(51, "")
		//player.SendStringMsg(5, "当前处于背包界面,请选择操作：0返回1增加物品2扣除物品3使用物品")
	default:
		player.SendStringMsg(5, game.BagLogicStr)
	}
END:
}
