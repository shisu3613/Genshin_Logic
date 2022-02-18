package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"server/zinx/ziface"
)

/*
存储一切有关zinx框架的全局参数，供其他模块使用
*/

type GlobalObj struct {
	/*
		server
	*/
	TcpServer ziface.Iserver //当前zinx全局的server对象
	Host      string         //当前服务器主机的IP
	TcpPort   int            //当前监听的端口号
	Name      string         //当前服务器的名称

	/*zinx*/
	Version        string //当前版本
	MaxConn        int    //当前允许的最大连接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值

	//任务池消息队列模块参数
	WorkerPoolSize    uint32 //当前业务池子中的worker的数量
	MaxWorkerPoolSize uint32 //框架允许用户最多开辟多少个worker
	MaxTaskQueueLen   uint32 //当前框架内消息队列的最大长度
}

// GlobalObject 定义一个全局的对外的对象
var GlobalObject *GlobalObj

// Reload zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	//path, _ := os.Getwd()
	//fmt.Println("when init the GlobalObj, the path is ", path)
	data, err := ioutil.ReadFile("json/zinx.json")
	//将json文件解析
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
	if GlobalObject.WorkerPoolSize > GlobalObject.MaxWorkerPoolSize {
		GlobalObject.WorkerPoolSize = GlobalObject.MaxWorkerPoolSize
		fmt.Println("The max workerPoolSize in this frame is", GlobalObject.MaxWorkerPoolSize, "reset the WorkerPoolSize as ", GlobalObject.MaxWorkerPoolSize)
	}
}

//提供一个init方法初始化当前globalobject
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		TcpServer:         nil,
		Host:              "",
		TcpPort:           8999,
		Name:              "ZinxServerApp",
		Version:           "V0.8",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		MaxWorkerPoolSize: 4096,
		MaxTaskQueueLen:   1024,
	}
	GlobalObject.Reload()
}
