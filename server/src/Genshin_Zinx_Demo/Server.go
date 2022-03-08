package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"server/api"
	"server/csvs"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
	"time"
)

type PlayerRouter struct {
	znet.BaseRouter
}

// Handler Handler For Genshin Player
func (pr *PlayerRouter) Handler(request ziface.IRequest) {
	//这里是主界面的选择程序
	//找到对应的player
	UserID, err := request.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//根据pID得到player对象
	player := game.WorldMgrObj.GetPlayerByPID(UserID.(int))
	//player.SendStringMsg(2, "欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
	//player.SendStringMsg(2,"欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")

	var msgChoice int
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	switch msgChoice {
	case 1:
		//player.HandleBaseRemote(request)
		player.SendStringMsg(3, "当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")

	case 2:
		player.HandleBag()
	case 3:
		player.SendStringMsg(33, "请输出抽卡次数:")
	case 4:
		player.HandleWishUp()
	case 5:
		player.HandleMap()
	}
}

// DoConnectionBegin 客户端和服务器链接创立成功时候，设置链接的一些属性，和发送开始消息给客户端
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=======>DoConnectionBegin is Called ...")
	//当客户端链接成功后，创建一个玩家
	player := game.NewClientPlayer(conn)
	//告知客户端pID,同步已经生成的玩家ID给客户端
	player.SyncPid()
	conn.SetProperty("PID", player.ModPlayer.UserId)

	//将玩家加入世界管理器中
	game.WorldMgrObj.AddPlayer(player)

	player.SendStringMsg(2, player.ModPlayer.Name+",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
}

func DoConnectionLost(conn ziface.IConnection) {
	pID, _ := conn.GetProperty("PID")

	//根据pID获取对应的玩家对象
	player := game.WorldMgrObj.GetPlayerByPID(pID.(int))

	//触发玩家下线业务
	if player != nil {
		game.WorldMgrObj.RemovePlayerByPID(pID.(int))
	}

	fmt.Println("====> Player ", pID, " left =====")
}

func main() {
	//1.创建server的句柄
	s := znet.NewServer("原神测试工具服务器【V0.1】")
	rand.Seed(time.Now().UnixNano())
	csvs.CheckLoadCsv()

	//启动违禁词库协程
	go game.GetManageBanWord().Run()
	//2.给当前zinx框架添加多个自定义的router
	s.AddRouter(202, &PlayerRouter{})
	s.AddRouter(203, &api.HandlerBase{})
	s.AddRouter(204, &api.HandlerBaseName{})
	s.AddRouter(233, &api.HandlerWishesTest{})
	//注册Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.启动server协程
	s.Serve()
}
