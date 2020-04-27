package tcp

/**
 * This file is part of Guerd.
 *
 * @link     http://guerd.kyour.cn
 * @document http://guerd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/guerd/blob/master/LICENSE
 */

import (
	"fmt"
	"net"
	"sync"
)

type TcpServer struct {
	Listen net.Listener
}

//客户端连接池
type Client struct {
	mu      sync.Mutex
	conn    []Connection
	count   uint32 //连接数量
	fdIndex uint32 // fd自动递增
}

var Clients Client

//监听tcp端口
func Listen(network string, addr string) (server TcpServer, err error) {

	listen, err := net.Listen(network, addr)
	if err != nil {
		return server, err
	}

	server.Listen = listen

	return
}

//封装tcp连接池管理
func Accept(server TcpServer, onConnect func(conn Connection), onReceive func(conn Connection, buffer []byte), onClose func(conn Connection)) {

	for {
		conn, e := server.Listen.Accept()
		if e != nil {
			fmt.Println("Accept Error.")
		}

		//创建tcp连接对象
		tcpConn := &Connection{
			Socket:   conn,
			isClosed: false,
			Fd:       Clients.fdIndex,
			Addr:     conn.RemoteAddr().String(),
		}

		//新的连接
		go Connect(*tcpConn, onConnect, onReceive, onClose)
	}
}

//新连接
func Connect(conn Connection, onConnect func(conn Connection), onReceive func(conn Connection, buffer []byte), onClose func(conn Connection)) {

	//将连接加入连接池
	Clients.mu.Lock() //加锁，避免资源争夺
	Clients.fdIndex++ //自增ID
	Clients.count++   //连接数量
	Clients.conn = append(Clients.conn, conn)
	Clients.mu.Unlock()

	//新的连接
	onConnect(conn)

	buffer := make([]byte, 1024)

	//监听客户端消息
	for {

		n, err := conn.Socket.Read(buffer)
		if err != nil {
			break
		}

		//接收到数据回调
		onReceive(conn, buffer[0:n])
	}

	//连接关闭回调
	onClose(conn)

	//客户端断开
	_ = conn.Socket.Close()

}
