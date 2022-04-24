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
		player.SendStringMsg(0, "如果祈愿之缘数量不足，请通过背包功能增加祈愿之缘，物品id为1000005")
		if err = player.GetMod(game.BagMod).(*game.ModBag).RemoveItem(1000005, 1); err != nil {
			player.SendStringMsg(0, fmt.Sprint(err))
		} else {
			player.SendStringMsg(0, utils.CaptureOutput(func() {
				player.ModWish.DoPool(1, player)
			}))
		}
		player.SendStringMsg(6, "您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")
		goto END
	case 2:
		player.SendStringMsg(0, "如果祈愿之缘数量不足，请通过背包功能增加祈愿之缘，物品id为1000005")
		if err = player.GetMod(game.BagMod).(*game.ModBag).RemoveItem(1000005, 10); err != nil {
			player.SendStringMsg(0, fmt.Sprint(err))
		} else {
			player.SendStringMsg(0, utils.CaptureOutput(func() {
				player.ModWish.DoPool(10, player)
			}))
		}
		player.SendStringMsg(6, "您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")
		goto END
	case 3:
		player.SendStringMsg(0, player.WishHelper())
		player.SendStringMsg(6, "您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")
		goto END
	default:
		player.SendStringMsg(6, "您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")

	}
END:
}
