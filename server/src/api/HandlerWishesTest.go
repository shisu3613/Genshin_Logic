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
	//player.SendStringMsg(0, player.ModWish.DoPoolTest(times))
	player.SendStringMsg(2, player.ModWish.DoPoolTest(times)+"\n"+player.GetUserName()+game.MainLogicStr)
}
