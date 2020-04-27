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
	"fmt"
	"github.com/kyour-cn/guerd/application/app-http"
	app_tcp "github.com/kyour-cn/guerd/application/app-tcp"
	"time"
)

const version = "0.1"

var logo = `
_____/\\\\\\\\\\\\___________________________________________________/\\\__        
 ___/\\\//////////___________________________________________________\/\\\__       
  __/\\\______________________________________________________________\/\\\__      
   _\/\\\____/\\\\\\\__/\\\____/\\\_____/\\\\\\\\___/\\/\\\\\\\________\/\\\__     
    _\/\\\___\/////\\\_\/\\\___\/\\\___/\\\/////\\\_\/\\\/////\\\__/\\\\\\\\\__    
     _\/\\\_______\/\\\_\/\\\___\/\\\__/\\\\\\\\\\\__\/\\\___\///__/\\\////\\\__   
      _\/\\\_______\/\\\_\/\\\___\/\\\_\//\\///////___\/\\\________\/\\\__\/\\\__  
       _\//\\\\\\\\\\\\/__\//\\\\\\\\\___\//\\\\\\\\\\_\/\\\________\//\\\\\\\/\\_ 
        __\////////////_____\/////////_____\//////////__\///__________\///////\//__

                             SERVER INFORMATION(v%s)
  *********************************************************************************
  * Http      | Enabled：true Listen: %v
  * WebSocket | Enabled：true Listen: (@Http.Listen)
  * TCP       | Enabled：true Listen: %s
  *********************************************************************************
`

type Application struct {
	Name     string
	Debug    bool
	ConfPath string
}

//创建Application
func NewApp() Application {

	app := Application{
		Name:     "Guerd App",
		Debug:    false,
		ConfPath: "./app.conf",
	}

	return app
}

//获取版本
func (app *Application) getVersion() string {
	return version
}

//设定配置文件
func (app *Application) ConfigFile(path string) {
	app.ConfPath = path

}

//启动
func (app *Application) Serve() {

	var errors []error

	//获取配置
	config := readConfig(app.ConfPath)
	fmt.Printf("%v", config)

	fmt.Printf(logo, version, config.Http.Addr, config.Tcp.Addr)

	//启动http服务
	go func() {
		err := app_http.Serve(config.Http)
		if err != nil {
			errors = append(errors, err)
		}
	}()

	//启动Tcp服务
	go func() {
		err := app_tcp.Serve(config.Tcp)
		if err != nil {
			errors = append(errors, err)
		}
	}()

	//每一秒检查是否有报错并输出
	for {
		time.Sleep(1000)
		if len(errors) > 0 {
			for i, e := range errors {
				//移除这个数据
				errors = append(errors[:i], errors[i+1:]...)
				fmt.Printf("Error:%v\n", e)
			}
		}
	}

}
