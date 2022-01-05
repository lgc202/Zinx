package znet

import (
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"net"
)

type Connection struct {
	// 当前连接的socket tcp套接字
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 当前连接所绑定的处理业务的方法
	handleAPI ziface.HandleFunc

	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool
}

func (c *Connection) Start() {
	fmt.Println("Conn Start(), ConnID = ", c.ConnID)
	// 启动从当前连接的读业务
	c.StartReader()
	// TODO 启动从当前连接的写业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn Sopt(), ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func NewConnection(conn *net.TCPConn, connID uint32, handleAPI ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: handleAPI,
		isClosed:  true,
		ExitChan:  make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中， 目前最大是512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("read data from client failed, err: %s\n", err.Error())
			continue
		}

		// 调用当前连接所绑定的handleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID ", c.ConnID, "handle err ", err)
			break
		}
	}
}
