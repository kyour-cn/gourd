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
	"net"
	"sync"
)

// Tcp连接
type Connection struct {
	Fd       uint32
	Addr     string
	Socket   net.Conn   // 底层websocket
	mutex    sync.Mutex // 避免重复关闭管道
	isClosed bool
}

//ws连接获取消息
func (conn *Connection) Read(buffer []byte) error {
	_, err := conn.Socket.Read(buffer)

	return err
}

//ws连接发送消息
func (conn *Connection) Write(data []byte) error {
	_, err := conn.Socket.Write(data)

	return err
}

//关闭socket连接
func (conn *Connection) Close() error {

	return conn.Socket.Close()

}
