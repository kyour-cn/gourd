package guerd

/**
 * This file is part of Guerd.
 *
 * @link     http://guerd.kyour.cn
 * @document http://guerd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/guerd/blob/master/LICENSE
 */

import (
	app_http "github.com/kyour-cn/guerd/application/app-http"
	"github.com/kyour-cn/guerd/application/app-tcp"
	"github.com/kyour-cn/guerd/utils/toml"
	"log"
)

type Config struct {
	Http app_http.HttpConfig
	Tcp  app_tcp.TcpConfig
}

//读取配置信息
func readConfig(file string) (config *Config) {

	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatal(err)
	}

	return
}
