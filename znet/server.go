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
	fmt.Printf("[Start %s] Server Listening at IP: %s, Port: %d\n", s.Name, s.IP, s.Port)
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
		var connID uint32

		// 3. 阻塞等待客户端连接， 处理读写请求
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("listen failed, err: %s", err.Error())
				continue
			}

			dealConn := NewConnection(conn, connID, CallBackToClient)
			connID++
			// 启动当前连接业务处理
			go dealConn.Start()
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

// CallBackToClient 用来处理业务逻辑, 目前先写死， 后面应该由用户自己指定
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		return fmt.Errorf("CallBackToClient failed, err: %s", err.Error())
	}
	return nil
}
