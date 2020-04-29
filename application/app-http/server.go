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
	"github.com/gorilla/mux"
	"github.com/kyour-cn/gourd/application"
	"net/http"
	"time"
)

func Serve(config *application.HttpConfig, router *mux.Router) (error error) {

	if !config.Enable {
		//不启用
		return
	}

	for _, addr := range config.Addr {
		//监听多个地址
		srv := &http.Server{
			Handler: router,
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
