package znet

import (
	"server/zinx/ziface"
)

// BaseRouter 定义baseRouter的初衷是实现router时候继承base基类，根据需求对基类进行重写
type BaseRouter struct{}

// PreHandle 这里之所以Base Router的方法都为空，是因为有的Router不需要之中的一些方法，比如Post Handler
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handler(request ziface.IRequest) {}

func (br *BaseRouter) PostHandler(request ziface.IRequest) {}
