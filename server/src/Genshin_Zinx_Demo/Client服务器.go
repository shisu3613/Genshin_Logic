package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"server/msgJson"
	"server/zinx/znet"
)

type TcpClientServer struct {
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

// NewTcpClientSever 创立新连接
func NewTcpClientSever(ip string, port int) *TcpClientServer {
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		panic(err)
	}

	client := &TcpClientServer{
		conn:            conn,
		PID:             0,
		isOnline:        make(chan bool),
		BackToMainLogic: make(chan struct{}),
	}
	return client
}

func (client *TcpClientServer) start() {
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

func (client *TcpClientServer) DoMsg(msg *znet.Message) {
	switch msg.Id {
	case 0:
		//打印返回信息
		client.PrintMsg(msg)
	case 1:
		sycnID := new(msgJson.SyncPID)
		_ = json.Unmarshal(msg.Data, sycnID)
		client.PID = sycnID.PID
		fmt.Println("模拟用户创建成功OK------开始测试")
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

func (client *TcpClientServer) PrintMsg(msg *znet.Message) {
	var response string
	_ = json.Unmarshal(msg.Data, &response)
	fmt.Println(response)
}

func main() {
	Client := NewTcpClientSever("116.62.193.144", 8999)
	//Client := NewTcpClient(""127.0.0.1", 8999)
	Client.start()
}
