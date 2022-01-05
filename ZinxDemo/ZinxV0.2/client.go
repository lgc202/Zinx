package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	// 1. 连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Printf("connet failed, err: %s\n", err.Error())
		return
	}

	// 2. 调用write方法
	for {
		_, err := conn.Write([]byte("hello zinx v0.2..."))
		if err != nil {
			fmt.Printf("write to server failed, err: %s\n", err.Error())
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("read from server failed, err: %s\n", err.Error())
			return
		}

		fmt.Printf("server call back: %s, cnt: %d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}
