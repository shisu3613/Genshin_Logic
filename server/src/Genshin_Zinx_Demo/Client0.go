package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"server/msgJson"
	"server/zinx/znet"
	"sync"
)

type TcpClient struct {
	conn            net.Conn
	PID             int
	OnlineMsg       map[int]string
	closeClientChan chan struct{}
	closedWg        sync.WaitGroup
}

// Message msg结构
//type Message struct {
//	Id      uint32
//	DataLen uint32
//	Data    []byte
//}

// NewTcpClient 创立新连接
func NewTcpClient(ip string, port int) *TcpClient {
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		panic(err)
	}

	client := &TcpClient{
		conn: conn,
		PID:  0,
		// @Modified By WangYuding 2022/4/24 17:13:00
		// @Modified description 添加OnlineMsg做一下go的聊天室，缓存聊天信息
		OnlineMsg:       make(map[int]string),
		closeClientChan: make(chan struct{}),
	}
	return client
}

func (client *TcpClient) start() {

	//  客户端关闭的部分逻辑
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	client.closedWg.Add(1)
	//保持接收信息
	go func() {
		for {
			conn := client.conn
			dp := znet.NewDataPack()
			binaryHead := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(conn, binaryHead); err != nil {
				fmt.Println("Read Head Error:", err)
				break
			}
			//先读取流中的head部分得到ID鹤datalen,再根据datalen
			msgHead, err := dp.Unpack(binaryHead)
			//fmt.Println(msgHead.GetMsgLen())
			if err != nil {
				fmt.Println("Read Head err:", err)
				break
			}
			if msgHead.GetMsgLen() > 0 {
				msg := msgHead.(*znet.Message)
				msg.Data = make([]byte, msgHead.GetMsgLen())
				if _, err := io.ReadFull(conn, msg.Data); err != nil {
					fmt.Println("Read Head Error:", err)
					break
				}
				//fmt.Println("————————>Recv Server Msg : ID =", msg.Id, ",Len = ", msg.DataLen, ",data = ", string(msg.Data))
				// @Modified By WangYuding 2022/4/25 17:13:00
				// @Modified description 合格的客户端可以同时接收处理多条消息，特别是聊天模块
				go client.DoMsg(msg)
			}
		}
	}()

	//go func() {defer client.closedWg.Done();}()
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		client.stop()
	case <-client.closeClientChan:
		client.exitHandler()
		client.closedWg.Done()
	}
	client.closedWg.Wait()
}

// stop
// @Description: 客户端关闭
// @receiver client
func (client *TcpClient) stop() {
	defer client.closedWg.Done()
	close(client.closeClientChan)
	client.exitHandler()
}

func (client *TcpClient) exitHandler() {
	fmt.Println("Closing connection.....")
	err := client.conn.Close()
	if err != nil {
		panic("关闭链接失败")
	}
}

// DoMsg
// @Description: 分析处理收到的msg
// @receiver client
// @param msg
func (client *TcpClient) DoMsg(msg *znet.Message) {
	switch msg.Id {
	case 0:
		//case0 打印返回信息
		client.PrintMsg(msg)
	case 1:
		//case1 同步服务器
		sycnID := new(msgJson.SyncPID)
		_ = json.Unmarshal(msg.Data, sycnID)
		client.PID = sycnID.PID
		fmt.Println("用户连接成功OK------当前UID为：", client.PID)
		fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	case 4:
		//case ：姓名等需要输入string的情况
		client.PrintMsg(msg)
		var modChoose string
		_, err := fmt.Scan(&modChoose)
		if err != nil {
			fmt.Println("Scan error!")
			return
		}
		msgJson.MsgMgrObj.SendMsg(msg.Id+200, modChoose, client.conn)
	case 51: //增加物品的模块
		//case51:特殊模块：addItem需要输入两个值
		type pair struct {
			ItemId  int
			ItemNum int
		}
		scanRes := pair{}
		fmt.Println("物品ID")
		fmt.Scan(&scanRes.ItemId)
		fmt.Println("物品数量")
		fmt.Scan(&scanRes.ItemNum)
		msgJson.MsgMgrObj.SendMsg(msg.Id+200, scanRes, client.conn)
	//case 4294967295:
	//	_ = client.conn.Close()
	//聊天部分的信息处理：
	case 9:
		//进入聊天室界面
		client.PrintMsg(msg)
		//原子操作应该放到后端逻辑里面
		for {
			var msgStr string
			_, err := fmt.Scan(&msgStr)
			if err != nil {
				fmt.Println("Scan error!")
				return
			}
			if msgStr == "exit;" {
				msgJson.MsgMgrObj.SendMsg(202, 0, client.conn)
				break
			} else {
				msgJson.MsgMgrObj.SendMsg(msg.Id+200, msgStr, client.conn)
			}
		}
	case 19:

	default: //输入数字的情况
		client.PrintMsg(msg)
		var modChoose int
		_, err := fmt.Scan(&modChoose)
		if err != nil {
			fmt.Println("Scan error!")
			return
		}
		// @Modified By WangYuding 2022/4/24 22:51:00
		// @Modified description 增加本地关于退出客户端的判断
		if msg.Id == 2 && modChoose == 6 {
			close(client.closeClientChan)
			return
		}
		msgJson.MsgMgrObj.SendMsg(msg.Id+200, modChoose, client.conn)
	}
}

func (client *TcpClient) PrintMsg(msg *znet.Message) {
	var response string
	_ = json.Unmarshal(msg.Data, &response)
	fmt.Println(response)
}

func main() {
	//Client := NewTcpClient("116.62.193.144", 8999)
	Client := NewTcpClient("127.0.0.1", 8999)
	Client.start()
}
