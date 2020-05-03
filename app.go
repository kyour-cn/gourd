package gourd

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"fmt"
	"github.com/kyour-cn/gourd/application"
	app_cron "github.com/kyour-cn/gourd/application/app-cron"
	app_http "github.com/kyour-cn/gourd/application/app-http"
	app_tcp "github.com/kyour-cn/gourd/application/app-tcp"
	"github.com/kyour-cn/gourd/common"
	"github.com/kyour-cn/gourd/server/router"
	"github.com/kyour-cn/gourd/server/tcp"
	"github.com/kyour-cn/gourd/utils/cache"
	"github.com/kyour-cn/gourd/utils/gut"
	"time"
)

const version = "0.1"

var logo = `
_____/\\\\\\\\\\\\_________________________________________________/\\\__
 ___/\\\//////////_________________________________________________\/\\\__
  __/\\\____________________________________________________________\/\\\__
   _\/\\\____/\\\\\\\_____/\\\\\_____/\\\____/\\\__/\\/\\\\\\\_______\/\\\__
    _\/\\\___\/////\\\___/\\\///\\\__\/\\\___\/\\\_\/\\\/////\\\_/\\\\\\\\\__
     _\/\\\_______\/\\\__/\\\__\//\\\_\/\\\___\/\\\_\/\\\___\///_/\\\////\\\__
      _\/\\\_______\/\\\_\//\\\__/\\\__\/\\\___\/\\\_\/\\\_______\/\\\__\/\\\__
       _\//\\\\\\\\\\\\/___\///\\\\\/___\//\\\\\\\\\__\/\\\_______\//\\\\\\\/\\_
        __\////////////_______\/////______\/////////___\///_________\///////\//__

                             SERVER INFORMATION(v%s)
  *********************************************************************************
  * Http | Ws | Enabled：%v Listen: %v
  * TCP       | Enabled：%v Listen: %s
  *********************************************************************************
`

type Application struct {
	Name     string
	Debug    bool
	ConfPath string
	Router   router.Router
	TcpEvent tcp.Event
	Config   application.Config
}

//创建Application
func NewApp() Application {

	app := Application{
		Name:     "Gourd App",
		Debug:    false,
		ConfPath: "./conf/app.conf",
	}

	return app
}

//获取版本
func (app *Application) getVersion() string {
	return version
}

//设定配置文件
func (app *Application) ConfigPath(path string) {

	common.SetConfigPath(path)
	app.ConfPath = path

}

//启动
func (app *Application) Serve() {

	var errors []error

	var config application.Config

	//获取配置
	_ = common.ReadConfig("app", &config)

	//控制台输出logo
	fmt.Printf(logo, version, config.Http.Enable, config.Http.Addr, config.Tcp.Enable, config.Tcp.Addr)

	//创建运行目录
	_, err := gut.Mkdir("./runtime")
	if err != nil {
		errors = append(errors, err)
	}

	//启动http\ws服务
	if config.Http.Enable {
		go func() {
			err := app_http.Serve(&config.Http, &app.Router)
			if err != nil {
				errors = append(errors, err)
			}
		}()
	}

	//启动Tcp服务
	if config.Tcp.Enable {
		go func() {
			err := app_tcp.Serve(&config.Tcp, &app.TcpEvent)
			if err != nil {
				errors = append(errors, err)
			}
		}()
	}

	//初始化缓存
	cache.Init()

	//启动Crontab任务
	app_cron.Start()

	//阻塞应用 - 每一秒检查是否有报错并输出
	tick := time.NewTicker(time.Duration(1) * time.Second)
	for {
		select {
		case <-tick.C:
			if len(errors) > 0 {
				for i, e := range errors {
					//移除这个数据
					errors = append(errors[:i], errors[i+1:]...)
					fmt.Printf("Gourd_Error:%v\n", e)
				}
			}
		}
	}

}

//设定Http路由
func (app *Application) HttpRoute(LoadRouter func() (route *router.Router)) {

	var config application.Config

	_ = common.ReadConfig("app", &config)

	//判断Http开启才会获取路由
	if config.Http.Enable {
		app.Router = *LoadRouter()
	}
}

//设定Tcp事件
func (app *Application) RegistTcp(e tcp.Event) {

	var config application.Config

	_ = common.ReadConfig("app", &config)

	//判断tcp开启
	if config.Tcp.Enable {
		app.TcpEvent = e
	}

}
