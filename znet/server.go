package znet

import (
	"fmt"
	"github.com/lgc202/Zinx/utils"
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
	// 当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, IP: %s, Port: %d\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	go func() {
		// 0. 启动work工作池机制
		s.msgHandler.StartWorkerPool()

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

			dealConn := NewConnection(conn, connID, s.msgHandler)
			connID++
			// 启动当前连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Serve() {
	// 启动server服务功能
	s.Start()

	// TODO 做一些启动之后的其它业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("add router successful")
}

// NewServer 初始化server
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		msgHandler: NewMsgHandler(),
	}
}
