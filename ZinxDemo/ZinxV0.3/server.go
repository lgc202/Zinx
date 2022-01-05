package main

import (
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"github.com/lgc202/Zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("Call Router PreHandle err ", err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("pong...pong...pong\n"))
	if err != nil {
		fmt.Println("Call Router Handle err ", err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("Call Router PostHandle err ", err)
		return
	}
}

func main() {
	// 初始化服务器
	s := znet.NewServer("ZinxV0.3")
	// 添加自定义的router
	s.AddRouter(&PingRouter{})
	// 启动服务器
	s.Serve()
}
