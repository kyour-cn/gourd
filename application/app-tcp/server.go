package app_tcp

import (
	"fmt"
	"github.com/kyour-cn/guerd/server/tcp"
)

/**
 * This file is part of Guerd.
 *
 * @link     http://guerd.kyour.cn
 * @document http://guerd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/guerd/blob/master/LICENSE
 */

//Http服务配置信息
type TcpConfig struct {
	Enable bool   `toml:"enable"`
	Addr   string `toml:"addr"`
}

func Serve(config *TcpConfig) (err error) {

	server, err := tcp.Listen("tcp", config.Addr)
	if err != nil {
		return err
	}

	tcp.Accept(server, func(conn tcp.Connection) {
		//新的连接
		fmt.Printf("新的连接：(%v) %s\n", conn.Fd, conn.Addr)
	}, func(conn tcp.Connection, buffer []byte) {
		//新的数据接收
		fmt.Println("来自" + conn.Addr + "新的消息：" + string(buffer))

	}, func(conn tcp.Connection) {
		//连接断开
		fmt.Printf("连接断开：(%v) %s\n", conn.Fd, conn.Addr)

	})

	return
}
