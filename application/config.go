package application

/**
 * This file is part of Guerd.
 *
 * @link     http://guerd.kyour.cn
 * @document http://guerd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/guerd/blob/master/LICENSE
 */

import (
	"github.com/kyour-cn/guerd/application/app-tcp"
	"github.com/kyour-cn/guerd/utils/toml"
	"log"
)

//Http服务配置信息
type HttpConfig struct {
	Enable   bool     `toml:"enable"`
	Addr     []string `toml:"addr,omitempty"`
	WsEnable bool     `toml:"websocket"`
	Path     string   `toml:"path"`
	Index    string   `toml:"index"`
}

type Config struct {
	Http HttpConfig
	Tcp  app_tcp.TcpConfig
}

//读取配置信息
func readConfig(file string) (config *Config) {

	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal(err)
	}

	return
}
