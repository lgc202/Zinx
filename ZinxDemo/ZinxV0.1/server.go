package main

import "github.com/lgc202/Zinx/znet"

func main() {
	// 初始化服务器
	s := znet.NewServer("ZinxV0.1")
	// 启动服务器
	s.Serve()
}
