package znet

import (
	"errors"
	"fmt"
	"github.com/lgc202/Zinx/ziface"
	"io"
	"net"
)

type Connection struct {
	// 当前连接的socket tcp套接字
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool

	// 消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandler

	// 无缓冲管道， 用于读、写两个goroutine直接的消息通信
	msgChan chan []byte
}

func (c *Connection) Start() {
	fmt.Println("Conn Start(), ConnID = ", c.ConnID)
	// 1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	// 2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	for {
		select {
		case <-c.ExitChan:
			// 得到退出消息，不再阻塞
			return
		}
	}
	// TODO 启动从当前连接的写业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn Sopt(), ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	// 通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitChan <- true

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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	// 写回客户端
	c.msgChan <- msg
	return nil
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()
		// 读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head failed, err ", err)
			c.ExitChan <- true
			return
		}

		// 拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitChan <- true
			return
		}

		// 根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitChan <- true
				return
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中,找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data failed, err: ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			// conn已经关闭
			return
		}
	}
}
