package api

import (
	"encoding/json"
	"fmt"
	"server/game"
	"server/utils"
	"server/zinx/ziface"
	"server/zinx/znet"
)

type HandlerBaseName struct {
	znet.BaseRouter
}

func (hb *HandlerBaseName) Handler(request ziface.IRequest) {
	PID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(PID.(int))
	var msgChoice string
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	//player.RecvSetName(msgChoice)
	//player.SendStringMsg(0, utils.CaptureOutput(func() {
	//	player.RecvSetName(msgChoice)
	//}))
	outputString := utils.CaptureOutput(func() {
		player.RecvSetName(msgChoice)
	})
	player.SendStringMsg(3, outputString+"\n"+game.BasicLogicStr)

}
