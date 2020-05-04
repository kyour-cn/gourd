package middleware

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"log"
	"net/http"
)

type Handler struct {
	http.Handler
}

//响应
type Response interface {
	http.ResponseWriter
}

//请求
type Request struct {
	http.Request
}

//AutoRoute
func (p *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test"))
}

//http中间件
//func HttpMiddleware(w Response, r *Request) (c *HttpContext) {
//
//	sess := session.SessionMiddleware(w, r)
//
//	c = &HttpContext{
//		Response: &w,
//		Request:  r,
//		Session:  &sess,
//	}
//
//	return
//}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
