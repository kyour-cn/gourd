package app_http

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
	"github.com/kyour-cn/gourd/server/router"
	"github.com/kyour-cn/gourd/server/session"
	_ "github.com/kyour-cn/gourd/server/session/memory"
	"mime"
	"net/http"
	"time"
)

func Serve(config *application.HttpConfig, router *router.Router) (error error) {

	if !config.Enable {
		//不启用
		return
	}

	//初始化Session -暂时放这里
	session.Init()

	//注册中间件
	//router.Use()

	//静态资源mime定义
	_ = mime.AddExtensionType(".js", "text/javascript")
	_ = mime.AddExtensionType(".html", "text/html; charset=utf-8")
	_ = mime.AddExtensionType(".htm", "text/html; charset=utf-8")
	_ = mime.AddExtensionType(".css", "text/css")

	for _, addr := range config.Addr {
		//监听多个地址
		srv := &http.Server{
			Handler: router, //router
			Addr:    addr,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				error = err
				fmt.Printf("Gourd_Error:%v\n", err)
			}

		}()

	}

	return
}
