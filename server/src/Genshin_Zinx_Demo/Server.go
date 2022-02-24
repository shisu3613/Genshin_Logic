package main

import (
	"fmt"
	"math/rand"
	"server/csvs"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
	"time"
)

type PlayerRouter struct {
	znet.BaseRouter
	playerList map[int]*game.Player
}

// Handler Handler For Genshin Player
func (pr *PlayerRouter) Handler(request ziface.IRequest) {
	request.GetConnection()
}

// DoConnectionBegin 客户端和服务器链接创立成功时候，设置链接的一些属性，和发送开始消息给客户端
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=======>DoConnectionBegin is Called ...")
	//当客户端链接成功后，创建一个玩家
	player := game.NewClientPlayer(conn)
	//告知客户端pID,同步已经生成的玩家ID给客户端
	player.SyncPid()
	conn.SetProperty("PID", player.ModPlayer.UserId)
}

func main() {
	//1.创建server的句柄
	s := znet.NewServer("原神测试工具服务器【V0.1】")
	rand.Seed(time.Now().UnixNano())
	csvs.CheckLoadCsv()
	go game.GetManageBanWord().Run()
	//2.给当前zinx框架添加多个自定义的router
	s.AddRouter(0, &PlayerRouter{})

	//注册Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	//s.SetOnConnStop(DoConnectionLost)

	//3.启动server
	s.Serve()
}
