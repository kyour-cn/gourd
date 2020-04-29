package tcp

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"fmt"
	"net"
	"sync"
)

type TcpServer struct {
	Listen net.Listener
}

type Event struct {
	OnConnect func(conn Connection)
	OnReceive func(conn Connection, buffer []byte)
	OnClose   func(conn Connection)
}

//客户端连接池
type Client struct {
	mu      sync.Mutex
	Conn    []Connection
	Count   uint32 //连接数量
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

//onConnect func(conn Connection), onReceive func(conn Connection, buffer []byte), onClose func(conn Connection)
//封装tcp连接池管理
func Accept(server TcpServer, event Event) {

	for {
		conn, e := server.Listen.Accept()
		if e != nil {
			fmt.Println("Accept Error.")
		}

		//创建tcp连接对象
		tcpConn := &Connection{
			Socket:   conn,
			isClosed: false,
			Fd:       Clients.GetPoolIndex(),
			Addr:     conn.RemoteAddr().String(),
		}

		//将连接加入连接池
		Clients.mu.Lock() //加锁，避免资源争夺
		Clients.Count++   //连接数量
		Clients.Conn = append(Clients.Conn, *tcpConn)
		Clients.mu.Unlock()

		//新的连接
		go Connect(*tcpConn, event.OnConnect, event.OnReceive, event.OnClose)
	}
}

//获取连接池自增Id并加一
func (c *Client) GetPoolIndex() uint32 {

	c.fdIndex++
	return c.fdIndex

}

//新连接
func Connect(conn Connection, onConnect func(conn Connection),
	onReceive func(conn Connection, buffer []byte), onClose func(conn Connection)) {

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

	//将对象从连接池中移除
	for k, v := range Clients.Conn {
		if v == conn {
			Clients.Conn = append(Clients.Conn[:k], Clients.Conn[k+1:]...)
		}
	}

}
