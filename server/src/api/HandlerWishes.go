package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/utils"
	"server/zinx/ziface"
	"server/zinx/znet"
)

/**
    @author: WangYuding
    @since: 2022/4/17
    @desc: //需要和背包模块交互的抽卡模块
**/

type HandlerWishes struct {
	znet.BaseRouter
}

func (hb *HandlerWishes) Handler(request ziface.IRequest) {
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
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
		goto END
	case 1:
		// @Modified By WangYuding 2022/4/25 19:11:00
		// @Modified description 重构这部分的逻辑以防止由于网络问题导致的信息错位
		outputStr := ""

		player.SendStringMsg(0, game.WishHint)
		if err = player.GetMod(game.BagMod).(*game.ModBag).RemoveItem(1000005, 1); err != nil {
			//player.SendStringMsg(0, fmt.Sprint(err))
			outputStr += fmt.Sprintln(err)
		} else {
			//player.SendStringMsg(0, utils.CaptureOutput(func() {
			//	player.ModWish.DoPool(1, player)
			//}))
			outputStr += utils.CaptureOutput(func() { player.ModWish.DoPool(1, player) }) + "\n"
		}
		player.SendStringMsg(6, outputStr+game.WishLogicStr)
		goto END
	case 2:
		outputStr := ""

		player.SendStringMsg(0, game.WishHint)
		if err = player.GetMod(game.BagMod).(*game.ModBag).RemoveItem(1000005, 10); err != nil {
			//player.SendStringMsg(0, fmt.Sprint(err))
			outputStr += fmt.Sprintln(err)
		} else {
			//player.SendStringMsg(0, utils.CaptureOutput(func() {
			//	player.ModWish.DoPool(1, player)
			//}))
			outputStr += utils.CaptureOutput(func() { player.ModWish.DoPool(10, player) }) + "\n"
		}
		player.SendStringMsg(6, outputStr+game.WishLogicStr)
		goto END
	case 3:
		//player.SendStringMsg(0, player.WishHelper())
		player.SendStringMsg(6, player.WishHelper()+"\n"+game.WishLogicStr)
		goto END
	default:
		player.SendStringMsg(6, game.WishLogicStr)

	}
END:
}
