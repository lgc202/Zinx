package ziface

import "net"

type IConnection interface {
	// Start 启动连接
	Start()

	// Stop 停止连接
	Stop()

	// GetTcpConnection 获取当前连接的绑定socket conn
	GetTcpConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态IP port
	RemoteAddr() net.Addr

	// Send 发送数据
	Send([]byte) error
}

// HandleFunc 处理连接业务
type HandleFunc func(*net.TCPConn, []byte, int) error
