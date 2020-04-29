package application

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"github.com/kyour-cn/gourd/utils/toml"
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

//Tcp服务配置信息
type TcpConfig struct {
	Enable bool   `toml:"enable"`
	Addr   string `toml:"addr"`
}

type Config struct {
	Http HttpConfig
	Tcp  TcpConfig
}

//读取配置信息
func readConfig(file string) (config *Config) {

	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal(err)
	}

	return
}
