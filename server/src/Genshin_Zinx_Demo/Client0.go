package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"server/msgJson"
	"server/zinx/znet"
)

type TcpClient struct {
	conn            net.Conn
	PID             int
	isOnline        chan bool
	BackToMainLogic chan struct{}
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
		conn:            conn,
		PID:             0,
		isOnline:        make(chan bool),
		BackToMainLogic: make(chan struct{}),
	}
	return client
}

func (client *TcpClient) start() {
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
				client.DoMsg(msg)
			}
		}
	}()

	select {}
}

func (client *TcpClient) DoMsg(msg *znet.Message) {
	switch msg.Id {
	case 0:
		//打印返回信息
		client.PrintMsg(msg)
	case 1:
		sycnID := new(msgJson.SyncPID)
		_ = json.Unmarshal(msg.Data, sycnID)
		client.PID = sycnID.PID
		fmt.Println("用户连接成功OK------当前UID为：", client.PID)
		fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	case 4: //需要输入string的情况
		client.PrintMsg(msg)
		var modChoose string
		_, err := fmt.Scan(&modChoose)
		if err != nil {
			fmt.Println("Scan error!")
			return
		}
		msgJson.MsgMgrObj.SendMsg(msg.Id+200, modChoose, client.conn)
	case 51: //增加物品的模块
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

	default: //输入数字的情况
		client.PrintMsg(msg)
		var modChoose int
		_, err := fmt.Scan(&modChoose)
		if err != nil {
			fmt.Println("Scan error!")
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

//func (client *TcpClient) StartLogic() {
//	for {
//		fmt.Println("欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
//		var modChoose int
//		_, err := fmt.Scan(&modChoose)
//		if err != nil {
//			fmt.Println("Scan error!")
//			return
//		}
//		msgJson.MsgMgrObj.SendMsg(201, modChoose, client.conn)
//		<-client.BackToMainLogic
//		//switch modChoose {
//		//case 1:
//		//	msgJson.MsgMgrObj.SendMsg(201, 1, client.conn)
//		//case 2:
//		//	msgJson.MsgMgrObj.SendMsg(201, 2, client.conn)
//		//case 3:
//		//	msgJson.MsgMgrObj.SendMsg(201, 3, client.conn)
//		//case 4:
//		//	msgJson.MsgMgrObj.SendMsg(201, 4, client.conn)
//		//case 5:
//		//	msgJson.MsgMgrObj.SendMsg(201, 5, client.conn)
//		//}
//		//fmt.Println(pr.ModPlayer.Name, ",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘）5.地图")
//	}
//}

func main() {
	//Client := NewTcpClient("116.62.193.144", 8999)
	Client := NewTcpClient("127.0.0.1", 8999)
	Client.start()
}
