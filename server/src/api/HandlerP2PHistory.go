package api

import (
	"encoding/json"
	"fmt"
	DB "server/DB/GORM"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
	"strconv"
)

/**
    @author: WangYuding
    @since: 2022/5/5
    @desc: //打印输出个人聊天记录
**/

type HandlerP2PHistory struct {
	znet.BaseRouter
}

func (p2p *HandlerP2PHistory) Handler(request ziface.IRequest) {
	// 第一步根据uid 获得player对象
	PID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	//log.Println("UserID is:" + strconv.Itoa(UserID.(int)))
	player := game.WorldMgrObj.GetPlayerByPID(PID.(int))

	//获得聊天信息
	var uid int
	_ = json.Unmarshal(request.GetData(), &uid)
	//判断数据库存不存在当前ID
	var result struct {
		Found bool
	}
	DB.GormDB.Raw("SELECT EXISTS(SELECT 1 FROM BasicProfiles WHERE user_id=? AND `deleted_at` IS NULL) As found", uid).Scan(&result)
	if result.Found == false {
		player.SendStringMsg(8, "当前UID的账号不存在；请重新输入")
		return
	}
	//打印聊天历史
	outputStr := player.GetMod(game.TalkMod).(*game.ModTalk).PrintP2PHistory(uid)
	outputStr += fmt.Sprintf("您已经进入聊天室可以同 %d 通话:\n", uid)
	player.SendStringMsg(0, outputStr)
	player.SendStringMsg(9, strconv.Itoa(uid))
	player.GetMod(game.TalkMod).(*game.ModTalk).SetFlag(int32(uid))

}
