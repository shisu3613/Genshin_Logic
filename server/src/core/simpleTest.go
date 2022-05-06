package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type SyncPID struct {
	PID     int `json:"PID"`
	TestVal int
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			handle()
		}
	}
}

func (t *Task) Stop() {
	close(t.closed)
}

func handle() {
	for i := 0; i < 5; i++ {
		fmt.Print("#")
		time.Sleep(time.Millisecond * 200)
	}
	fmt.Println()
}

func main() {
	//Map := make(map[int]interface{})
	//s := &SyncPID{PID: 99999, TestVal: 1111}
	//Map[1] = s
	//
	//newS := Map[1]
	//

	// @Modified By WangYuding 2022/4/24 16:54:00
	// @Modified description 同步信号的test
	//tests := SyncPID{PID: 88888}
	//msg1, _ := json.Marshal(s)
	//msg2, _ := json.Marshal(tests)
	////json.Unmarshal(msg, newS)
	//fmt.Println(string(msg1))
	//fmt.Println(string(msg2))
	//json.Unmarshal(msg2, newS)
	//fmt.Println(newS)
	//===================================================================================
	//capture屏幕信号的test
	//test := func(a int) { fmt.Println(a) }
	//fmt.Println("test" + utils.CaptureOutput(func() {
	//	test(2)
	//}))
	//===================================================================================
	////signal捕获crtl+c的测试
	//task := &Task{
	//	closed: make(chan struct{}),
	//	ticker: time.NewTicker(time.Second * 2),
	//}
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt)
	//task.wg.Add(1)
	//go func() { defer task.wg.Done(); task.Run() }()
	//select {
	//case sig := <-c:
	//	fmt.Printf("Got %s signal. Aborting...\n", sig)
	//	task.Stop()
	//	//case <-task.closed:
	//	//	task.wg.Done()
	//	//	//case <-task.ticker.C:
	//	//	//	handle()
	//}
	//println("wait for stop")
	//task.wg.Wait()

	//===================================================================================
	//redis工具测试
	//var ctx = context.Background()
	//client := RedisTool.NewRedis(3)
	////err := client.RPush(ctx, "key", 1, 2, 3, 4, 5).Err()
	////if err != nil {
	////	panic(err)
	////}
	////val, err := client.LIndex(ctx, "key", -1).Result()
	////if err != nil {
	////	println(err)
	////}
	////
	////fmt.Println(val)
	//val, err := client.Exists(ctx, "key1").Result()
	//fmt.Println(val, err)

	//===================================================================================
	//type ChatMsg struct {
	//	Uid    string `json:"uid"` //消息发送者的UID
	//	IdTime string `json:"idTime"`
	//	Cnt    string `json:"cnt"`
	//	SendTo string `json:"sendTo"` //消息到达者的UID,global情况储存到db1，其他储存到db2
	//}
	//
	//type MsgSlice []*ChatMsg
	//PrivateChat := make(map[string]MsgSlice)
	//PrivateChat["1"] = nil
	////PrivateChat["1"] = append(PrivateChat["1"], &ChatMsg{})
	//res, ok := PrivateChat["1"]
	//println(res, ok)
	//fmt.Println(MsgSlice{})
	//for i := range PrivateChat {
	//	fmt.Println(i)
	//}
	//===================================================================================
	//scan的测试
	//var msg string
	//SpaceScan(&msg)
	//fmt.Println(msg)
	var msgStr string
	for msgStr == "" {
		msgStr, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		msgStr = strings.TrimSpace(msgStr)
	}
	fmt.Println(msgStr)

	//var msg string
	//fmt.Scan(&msg)
	//fmt.Println(msg)

	//var strInput string
	//fmt.Println("Enter a string ")
	//scanner := bufio.NewScanner(os.Stdin)
	//
	//if scanner.Scan() {
	//	strInput = scanner.Text()
	//}
	//
	//fmt.Println(strInput)

}

func SpaceScan(msg *string) {
	reader := bufio.NewReader(os.Stdin) // 标准输入输出
	*msg, _ = reader.ReadString('\n')   // 回车结束
	*msg = strings.TrimSpace(*msg)      // 去除最后一个空格
}
