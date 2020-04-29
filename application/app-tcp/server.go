package app_tcp

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"github.com/kyour-cn/gourd/application"
	"github.com/kyour-cn/gourd/server/tcp"
)

func Serve(config *application.TcpConfig, event *tcp.Event) (err error) {

	server, err := tcp.Listen("tcp", config.Addr)
	if err != nil {
		return err
	}

	tcp.Accept(server, *event)

	return
}
