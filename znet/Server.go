package znet

import (
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"net"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
}

func (s Server) Start() {
	fmt.Printf("[Start] Server Listening at IP: %s, Port: %d\n", s.IP, s.Port)
	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		fmt.Println("start Zinx server succ")

		// 3. 阻塞等待客户端连接， 处理读写请求
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("listen failed, err: %s", err.Error())
				continue
			}

			// 已经与客户端建立连接， 做一些基本的回写业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Printf("read data failed, err: %s\n", err.Error())
						continue
					}

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Printf("write failed, err: %s\n", err.Error())
						continue
					}
				}
			}()
		}
	}()
}

func (s Server) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s Server) Serve() {
	// 启动server服务功能
	s.Start()

	// TODO 做一些启动之后的其它业务

	// 阻塞状态
	select {}
}

// NewServer 初始化server
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}
