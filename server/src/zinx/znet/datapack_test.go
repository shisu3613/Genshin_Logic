package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//单元测试
//只是负责测试datapack拆包和封包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	//1.创建socket TCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server Listen err :", err)
		return
	}

	//创建一个go负责处理客户端发送的业务
	go func() {
		//2.从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept err :", err)
				return
			}

			go func(net.Conn) {
				//处理客户端的请求
				//————————》是一个unPack的过程《————————————//
				//定义一个拆包的对象
				dp := NewDataPack()
				for {
					//第一次读把head读出来
					HeadData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, HeadData)
					if err != nil {
						fmt.Println("read head err :", err)
						return
					}
					//将HeadData解包
					msgHead, err := dp.Unpack(HeadData)
					if err != nil {
						fmt.Println("sever unpack err :", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//说明由数据的，需要进行第二次读取
						//第二次读把data读出来
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//跟据datalen的长度再次从IO流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("跟据datalen的长度再次从IO流中读取失败:", err)
							return
						}
						//完整的一个消息已经读取完毕了
						fmt.Println("--》recedata(ID,LEN,DATA):", msg.Id, msg.DataLen, string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	/*模拟客户端*/
	conn, error := net.Dial("tcp", "127.0.0.1:7777")
	if error != nil {
		fmt.Println("client start err,", error)
		return
	}
	//创建一个封包对象dp
	dp := NewDataPack()

	//创建第一个Msg
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	//打包
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack err :", err)
		return
	}

	//创建第二个Msg
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	//打包
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack err :", err)
		return
	}
	fmt.Println("sendData1:", sendData1)
	fmt.Println("sendData2:", sendData2)
	sendData1 = append(sendData1, sendData2...)
	fmt.Println("sendData:", sendData1)
	//将两个包放在一起模拟粘包过程，一次性发送给服务器端
	if _, err := conn.Write(sendData1); err != nil {
		return
	}
	select {}
}
