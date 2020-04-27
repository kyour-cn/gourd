package app_http

/**
 * This file is part of Guerd.
 *
 * @link     http://guerd.kyour.cn
 * @document http://guerd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/guerd/blob/master/LICENSE
 */

import (
	"net/http"
)

//Http服务配置信息
type HttpConfig struct {
	Enable   bool     `toml:"enable"`
	Addr     []string `toml:"addr,omitempty"`
	WsEnable bool     `toml:"websocket"`
}

func Serve(config HttpConfig) (err error) {

	if !config.Enable {
		//不启用
		return
	}

	//文件服务
	http.Handle("/", http.FileServer(http.Dir("./public")))

	//ws服务
	//http.HandleFunc("/ws", websocketHandle)

	//创建监听端口
	err = http.ListenAndServe(":8100", nil)

	return
}
