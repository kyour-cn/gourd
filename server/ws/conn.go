package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

// Websocket连接
type Connection struct {
	Fd       uint32
	WsSocket *websocket.Conn // 底层websocket
	mutex    sync.Mutex      // 避免重复关闭管道
	isClosed bool
	//closeChan chan byte // 关闭通知
}

//ws连接发送消息
func (wsConn *Connection) Send(msgType int, data []byte) error {

	if err := wsConn.WsSocket.WriteMessage(msgType, data); err != nil {
		return err
	}
	return nil
}

//关闭ws连接
func (wsConn *Connection) Close() error {
	wsConn.isClosed = true
	err := wsConn.WsSocket.Close()
	return err

}
