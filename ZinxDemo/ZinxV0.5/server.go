package main

import (
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"github.com/lgc202/Zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写pong...pong...pong
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("pong...pong...pong\n"))
	if err != nil {
		fmt.Println("Call Router Handle err ", err)
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
