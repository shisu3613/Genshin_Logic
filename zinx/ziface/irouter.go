package ziface

/*
	路由抽象接口
	路由里的数据大欧式IRequest请求
*/

type IRouter interface {
	// PreHandle 处理业务之前的钩子方法
	PreHandle(request IRequest)
	// Handler main handling way
	Handler(request IRequest)
	// PostHandler 之后的方法
	PostHandler(request IRequest)
}
