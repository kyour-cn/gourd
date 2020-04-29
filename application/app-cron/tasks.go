package app_cron

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"github.com/kyour-cn/gourd/utils/cron"
	"github.com/kyour-cn/gourd/utils/toml"
	"log"
	"strings"
)

//任务池
var Tasks map[string]func()

var Cron *cron.Cron

type Config struct {
	Rule   []string `toml:"rule"`
	Enable bool     `toml:"enable"`
}

func Register(name string, tf func()) {

	Tasks = make(map[string]func())

	Tasks[name] = tf

}

//直接添加任务
func AddFunc(spec string, cmd func()) error {

	if Cron == nil {
		Cron = cron.NewWithSeconds()
	}

	_, err := Cron.AddFunc(spec, cmd)

	return err

}

//启动定时任务
func Start() {

	var config Config

	congFile := "./conf/cron.conf"

	//获取配置
	if _, err := toml.DecodeFile(congFile, &config); err != nil {
		log.Printf("Crontab配置错误(%s):%v", congFile, err)
	}

	if !config.Enable {
		//未开启
		return
	}

	if Cron == nil {
		Cron = cron.NewWithSeconds()
	}

	//遍历
	for _, rule := range config.Rule {

		arr := strings.Split(rule, "=>")

		arr[0] = strings.TrimSpace(arr[0])
		arr[1] = strings.TrimSpace(arr[1])

		_, ok := Tasks[arr[0]]
		if ok {
			_, _ = Cron.AddFunc(arr[1], Tasks[arr[0]])
		}

	}

	//启动任务
	Cron.Start()

}
