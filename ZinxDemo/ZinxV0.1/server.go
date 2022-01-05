package main

import "github.com/lgc202/Zinx/znet"

func main() {
	// 初始化服务器
	s := znet.NewServer("")
	// 启动服务器
	s.Serve()
}
