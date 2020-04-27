package main

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
	"github.com/gorilla/websocket"
	"github.com/kyour-cn/Guerd/server/ws"
	"net/http"
	"strconv"
)

//ws服务器示例
func main() {

	fmt.Println("启动...")

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {

		//封装ws握手，维护连接池
		ws.Handshake(writer, request,
			func(conn ws.Connection) {
				//新连接
				fmt.Println("新的连接：" + strconv.Itoa(conn.Fd))
				err := conn.Send(websocket.TextMessage, []byte("Hello "+strconv.Itoa(conn.Fd)))
				if err != nil {
					fmt.Println("发送消息失败")
				}
			}, func(conn ws.Connection, msg ws.Message) {
				//收到消息
				//fmt.Println("收到：" + string(msg.Data))

				//回复消息
				err := conn.Send(msg.Type, []byte("收到:"+string(msg.Data)))
				if err != nil {
					fmt.Println("回复消息失败")
				}

				//根据fd发送消息
				//_ = ws.SendWsMessage(fd, "来自"+strconv.Itoa(conn.Fd)+"的消息")

			}, func(conn ws.Connection) {
				//新连接
				fmt.Println("连接" + strconv.Itoa(conn.Fd) + "关闭...")
			})
	})
	_ = http.ListenAndServe("0.0.0.0:8100", nil)

}
