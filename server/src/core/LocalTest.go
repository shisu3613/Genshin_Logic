package main

import (
	"math/rand"
	"server/csvs"
	"server/game"
	"sync"
	"time"
)

// main
// @Description 本地测试模块
// @Author WangYuding 2022-04-11 14:37:12
func main() {

	//**********************************************************
	// 加载配置
	//db := DB.NewDBConnection()
	//defer func(db *gorm.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		fmt.Println("Close Database failure:", err)
	//	}
	//}(db)

	//DB.DBtest(db)
	//db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&game.DBPlayer{})
	//
	//return
	rand.Seed(time.Now().UnixNano())
	//time.Sleep(time.Millisecond * 100)
	csvs.CheckLoadCsv()
	wg := sync.WaitGroup{}
	//启动违禁词库协程
	wg.Add(1)
	go game.GetManageBanWord().Run(&wg)

	//fmt.Printf("数据测试----start\n")
	playerTest := game.NewTestPlayer()
	go playerTest.Run()

	select {}
	game.GetManageBanWord().Close()
	wg.Wait()

	//ticker := time.NewTicker(time.Second * 10)
	//for {
	//	select {
	//	case <-ticker.C:
	//		playerTest := game.NewTestPlayer()
	//		go playerTest.Run()
	//	}
	//}

	return
}

func playerLoadConfig(player *game.Player) {
	for i := 0; i < 1000000; i++ {
		config := csvs.ConfigUniqueTaskMap[10001]
		if config != nil {
			println(config.TaskId)
		}
	}
}
