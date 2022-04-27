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
	PID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(PID.(int))

	var msgChoice int
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	switch msgChoice {
	case 0:
		// @Modified By WangYuding 2022/4/27 14:22:00
		// @Modified description 添加逻辑，复位各种状态
		player.GetMod(game.TalkMod).(*game.ModTalk).ResetFlag()
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
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
		// @Modified By WangYuding 2022/4/27 14:31:00
		// @Modified description 说明进入世界聊天状态
		//首先是得到历史聊天的结果
		//player.GetMod(game.TalkMod).(*game.ModTalk).GetGlobalHistory()
		player.SendStringMsg(9, player.GetMod(game.TalkMod).(*game.ModTalk).GetGlobalHistory()+game.WorldChatStr)
		player.GetMod(game.TalkMod).(*game.ModTalk).SetFlag()
	default:
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
	}
}
