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
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	var msgChoice string
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	//player.RecvSetName(msgChoice)
	player.SendStringMsg(0, utils.CaptureOutput(func() {
		player.RecvSetName(msgChoice)
	}))
	player.SendStringMsg(3, "当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")

}
