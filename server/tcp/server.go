package tcp

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
	conn    []*Connection
	count   int //连接数量
	fdIndex int // fd自动递增
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
func Accept(server TcpServer, onOpen func(conn Connection)) {

	for {
		conn, e := server.Listen.Accept()
		if e != nil {
			fmt.Println("Accept Error.")
		}
		tcpConn := &Connection{
			Socket:    conn,
			closeChan: make(chan byte),
			isClosed:  false,
			Fd:        Clients.fdIndex,
		}

		//新的连接
		go Connect(*tcpConn)
	}
}

//新连接
func Connect(conn Connection) {

	buffer := make([]byte, 1024)

	//监听客户端消息
	for {
		n, err := conn.Socket.Read(buffer)

		if err != nil {
			fmt.Println("Read Error.")
			break
		}

		clientMsg := string(buffer[0:n])
		fmt.Printf("收到%v的消息:%s\n", conn.Socket.RemoteAddr(), clientMsg)

		if clientMsg != "im off" {
			_, _ = conn.Socket.Write([]byte("已阅:" + clientMsg))
		} else {
			_, _ = conn.Socket.Write([]byte("bye!"))
			break
		}
	}

	//客户端断开
	_ = conn.Socket.Close()
	fmt.Printf("客户端%s已经断开连接\n", conn.Socket.RemoteAddr())

}

//消息接收处理
func Receive() {

}

//连接关闭处理
func Close() {

}
