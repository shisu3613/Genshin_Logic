package main

import (
	"encoding/json"
	"fmt"
	//"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"server/api"
	"server/csvs"
	"server/game"
	"server/zinx/ziface"
	"server/zinx/znet"
	"sync"
	"time"
)

type LoadOrCreatRouter struct {
	znet.BaseRouter
}

// Handler Handler For Genshin Player
func (lc *LoadOrCreatRouter) Handler(request ziface.IRequest) {
	conn := request.GetConnection()
	var msgChoice int
	_ = json.Unmarshal(request.GetData(), &msgChoice)
	player := game.InitClientPlayer(conn)
	switch msgChoice {
	case -1:
		//player := game.InitClientPlayer(conn)
		player.CreateRoleInDB()

		////在数据库中生成对应的记录，根据记录生成对应的user_id
		//DB.GormDB.Create(&player.ModPlayer.DBPlayer)
		////fmt.Println(player.ModPlayer.DBPlayer.ID)
		////告知客户端pID,同步已经生成的玩家ID给客户端
		//DB.GormDB.Model(&player.ModPlayer.DBPlayer).Update("user_id", player.ModPlayer.DBPlayer.ID+100000000)
		//player.SyncPid()
		//conn.SetProperty("PID", player.ModPlayer.UserId)
		//
		////将玩家加入世界管理器中
		//game.WorldMgrObj.AddPlayer(player)
		//
		//player.SendStringMsg(2, player.ModPlayer.Name+",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
	default:
		//player := game.InitClientPlayer(conn)
		conn.SetProperty("PID", msgChoice-100000000)
		//if DB.GormDB.First(&player.ModPlayer.DBPlayer, msgChoice-100000000).RecordNotFound() {
		//	player.SendStringMsg(800, "当前UID不存在；请重新输入")
		//} else {
		//	conn.SetProperty("PID", player.ModPlayer.UserId)
		//	player.SyncPid()
		//	//将玩家加入世界管理器中
		//	game.WorldMgrObj.AddPlayer(player)
		//	player.SendStringMsg(2, player.ModPlayer.Name+",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
		//}
		player.GetMod(game.ModPlay).LoadData()
	}
	//player.LoadElse()
}

// DoConnectionBegin 客户端和服务器链接创立成功时候，设置链接的一些属性，和发送开始消息给客户端
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=======>DoConnectionBegin is Called ...")

	//finish:增加读取存档功能
	//当前客户端连接成功后发送信息给客户端
	data, err := json.Marshal("欢迎来到提瓦特大陆，读取存档请输入UID,新建存档请输入'-1'")
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.SendMsg(800, data); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	pID, err := conn.GetProperty("PID")
	if err != nil {
		return
	}
	//根据pID获取对应的玩家对象
	player := game.WorldMgrObj.GetPlayerByPID(pID.(int))

	//触发玩家下线业务
	if player != nil {
		fmt.Println("Saving the information into the Database.....")
		//DB.GormDB.Save(&player.ModPlayer.DBPlayer)
		for _, v := range player.GetModManager() {
			v.SaveData()
		}
		game.WorldMgrObj.RemovePlayerByPID(pID.(int))
	}

	fmt.Println("====> Player ", pID, " left =====")
}

func main() {
	//1.创建server的句柄
	s := znet.NewServer("原神测试工具服务器【V0.1】")
	rand.Seed(time.Now().UnixNano())
	csvs.CheckLoadCsv()

	//开启数据库
	// 加载配置
	//db := DB.NewDBConnection()
	//defer func(db *gorm.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		fmt.Println("Close Database failure:", err)
	//	}
	//}(DB.GormDB)

	//DB.DBtest(db)
	//DB.GormDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&game.DBPlayer{})
	//

	wg := sync.WaitGroup{}
	//启动违禁词库协程
	wg.Add(1)
	go game.GetManageBanWord().Run(&wg)

	//给当前zinx框架添加多个自定义的router
	s.AddRouter(1000, &LoadOrCreatRouter{})
	s.AddRouter(202, &api.PlayerRouter{})
	s.AddRouter(203, &api.HandlerBase{})
	s.AddRouter(204, &api.HandlerBaseName{})
	s.AddRouter(205, &api.HandlerBag{})
	s.AddRouter(251, &api.HandlerBagAddItem{})
	s.AddRouter(233, &api.HandlerWishesTest{})
	s.AddRouter(206, &api.HandlerWishes{})
	//注册Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.启动server协程
	s.Serve()

	//增加监听信号结束功能
	//集成在serve里面

	//关闭违反禁止词库
	game.GetManageBanWord().Close()
	//time.Sleep(time.Second * 3)
	wg.Wait()
}
