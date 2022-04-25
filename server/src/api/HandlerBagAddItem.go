package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/utils"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type HandlerBagAddItem struct {
	znet.BaseRouter
}

func (hb *HandlerBagAddItem) Handler(request ziface.IRequest) {
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	type pair struct {
		ItemId  int
		ItemNum int
	}
	scanRes := pair{}
	_ = json.Unmarshal(request.GetData(), &scanRes)
	//fmt.Println(scanRes)
	//player.SendStringMsg(0, utils.CaptureOutput(func() {
	//	player.AddBagItem(scanRes.ItemId, int64(scanRes.ItemNum))
	//}))
	//utils.CaptureOutput(func() {
	//	player.ModBag.AddItem(scanRes.ItemId, int64(scanRes.ItemNum), player)
	//})
	//player.ModBag.AddItem(scanRes.ItemId, int64(scanRes.ItemNum), player)
	outputString := utils.CaptureOutput(func() {
		player.AddBagItem(scanRes.ItemId, int64(scanRes.ItemNum))
	})
	player.SendStringMsg(5, outputString+"\n"+game.BagLogicStr)

}
