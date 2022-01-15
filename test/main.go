package main

import (
	"fmt"
	"go_linux/src/Genshin_Home_System/csvs"
	"go_linux/src/Genshin_Home_System/game"
	"time"
)

func main() {
	fmt.Println("数据测试---start")

	csvs.CheckLoadCsv()
	//加载配置功能
	go game.GetManageBanWord().Run()
	player := game.NewTestPlayer()
	player.RecvSetIcon(1)

	tickerIn := time.NewTicker(time.Second * 3)
	tickerOut := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-tickerIn.C:
			player.RecvName("专业代练")

		case <-tickerOut.C:
			player.RecvName("正常名字")

		}
	}
	//
	//player.RecvSetIcon(0)
	//player.RecvSetIcon(1)
	//player.RecvSetIcon(2)
	//
	//player.RecvSetCard(1)
	//player.RecvSetCard(2)
	//
	//player.RecvName("好人")
	//player.RecvName("出售外挂")
	//player.RecvSign("感觉不如原神")
	//player.RecvName("华人")
	//player.RecvSign("感觉不如原生")

	//公共管理类都会调用一个线程
	//玩家x1

	//200+10
	//内核级的线程每个线程10M起步
	//所以要使用消息队列,所以C++的逻辑会使用一个线程专门处理公共管理类
	//叫号机 200/7 负载均衡

	//golang的协程，即八核八个管理者，一个协程5k，模块的设计理念强制解耦
	//逻辑上没有一个强制的先后顺序（线程锁竞争），没有逻辑的冲突
}
