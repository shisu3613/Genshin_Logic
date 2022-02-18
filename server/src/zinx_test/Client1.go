package main

import (
	"fmt"
	"io"
	"net"
	"server/zinx/znet"
	"time"
)

//模拟客户端

func main() {
	fmt.Println("client1 start...")
	//连接客户端
	time.Sleep(time.Second * 1)

	conn, error := net.Dial("tcp", "127.0.0.1:8999")
	if error != nil {
		fmt.Println("client start err,", error)
		return
	}

	for {
		//尝试发送封印包的msg消息
		dp := znet.NewDataPack()
		data := []byte("ZinxV0.6 client Test Message")
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, data))
		if err != nil {
			fmt.Println("Client Pack err:", err)
			return
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("Conn Write error:", err)
			return
		}

		//服务器应该给我们回复一个Msg数据，MsgID：200 pingpingping
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
			fmt.Println("————————>Recv Server Msg : ID =", msg.Id, ",Len = ", msg.DataLen, ",data = ", string(msg.Data))
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
