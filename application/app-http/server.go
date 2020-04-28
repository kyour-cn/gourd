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
	"github.com/gorilla/mux"
	"github.com/kyour-cn/guerd/application"
	"net/http"
	"time"
)

func Serve(config *application.HttpConfig, router *mux.Router) (err error) {

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
			err = srv.ListenAndServe()
		}()

	}

	return
}
