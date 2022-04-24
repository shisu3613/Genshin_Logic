package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type HandlerBase struct {
	znet.BaseRouter
}

func (hb *HandlerBase) Handler(request ziface.IRequest) {
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	var msgChoice int
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	switch msgChoice {
	case 0:
		player.SendStringMsg(2, player.GetUserName()+game.MainLogicStr)
		goto END
	case 1:
		//player.HandleBaseRemote(request)
		//HandleBaseGetInfo
		player.SendStringMsg(0, player.HandleBaseGetInfoServer())
		player.SendStringMsg(3, game.BasicLogicStr)
	case 2:
		//设置名字
		player.SendStringMsg(4, "请输入名字:")
	case 3:
		player.HandleWishTest()
	case 4:
		player.HandleWishUp()
	case 5:
		player.HandleMap()
	}
	//player.SendStringMsg(3, game.BasicLogicStr)
END:
}
