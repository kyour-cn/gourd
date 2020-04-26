package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

//type WsServer struct {
//	http http.Server
//}

//Websocket连接池
type Client struct {
	mu      sync.Mutex
	conn    []*Connection
	count   int //连接数量
	fdIndex int // fd自动递增
}

var Clients Client

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端读写消息
type Message struct {
	Type int
	Data []byte
}

/**
 * 封装ws握手，维护连接池
 * @param resp      http.ResponseWriter
 * @param req       http.Request
 * @param onOpen    OpenEvent    新连接事件
 * @param onMessage MessageEvent 收到消息事件
 * @param onClose   CloseEvent   连接断开事件
 */
func Handshake(
	resp http.ResponseWriter, req *http.Request, onOpen func(conn Connection), onMessage func(conn Connection, msg Message), onClose func(conn Connection)) {

	// 应答客户端告知升级连接为websocket
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	wsConn := &Connection{
		WsSocket:  wsSocket,
		closeChan: make(chan byte),
		isClosed:  false,
		Fd:        Clients.fdIndex,
	}

	//将连接加入连接池
	Clients.mu.Lock() //加锁，避免资源争夺
	Clients.fdIndex++ //自增ID
	Clients.count++   //连接数量
	Clients.conn = append(Clients.conn, wsConn)
	Clients.mu.Unlock()

	// 处理器
	go wsConn.procLoop(onOpen, onMessage, onClose)
}

//维护新的连接
func (wsConn *Connection) procLoop(onOpen func(conn Connection), onMessage func(conn Connection, msg Message), onClose func(conn Connection)) {

	//fd := strconv.Itoa(wsConn.Fd)
	//fmt.Println("新的连接:" + fd)

	//新连接事件
	onOpen(*wsConn)

	// 启动一个gouroutine发送心跳
	/*
		go func() {
			for {
				time.Sleep(time.Duration(heart_time))
				if err := wsConn.Send(websocket.TextMessage, []byte("heartbeat from server")); err != nil {
					fmt.Println("heartbeat fail")
					wsConn.wsClose()
					break
				}
			}
		}()
	*/

	// 这是一个同步处理模型（只是一个例子），如果希望并行处理可以每个请求一个gorutine，注意控制并发goroutine的数量!!!
	for {

		// 读一个message
		msgType, data, err := wsConn.WsSocket.ReadMessage()
		if err != nil {
			//连接断开
			onClose(*wsConn)
			_ = wsConn.Close()
			break
		}
		msg := &Message{
			msgType,
			data,
		}

		//新消息回调
		onMessage(*wsConn, *msg)

	}
}

//根据fd取得ws连接
func FindClient(fd int) Connection {

	var c Connection
	for k, v := range Clients.conn {
		if v.Fd == fd {
			c = *Clients.conn[k]
		}
	}
	return c

}

//根据fd查到连接发送消息
func SendWsMessage(fd int, msg string) error {
	conn := FindClient(fd)

	err := conn.Send(websocket.TextMessage, []byte(msg))
	if err != nil {
		return errors.New("websocket closed")
	}

	return nil
}
