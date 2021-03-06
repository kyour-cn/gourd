package session

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/kyour-cn/gourd/common"
	"github.com/kyour-cn/gourd/server/session/memory"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var manager *Manager
var config Config

//session存储方式接口
//type Provider interface {
//	//初始化一个session，sid根据需要生成后传入
//	SessionInit(sid string) (Session, error)
//	//根据sid,获取session
//	SessionRead(sid string) (Session, error)
//	//销毁session
//	SessionDestroy(sid string) error
//	//回收
//	SessionGC(maxLifeTime int64)
//	//保存到文件
//	SaveToFile(path string) error
//	//从文件加载
//	LoadFile(file string) error
//}

type Manager struct {
	cookieName  string
	lock        sync.Mutex         //互斥锁
	provider    *memory.FromMemory //存储session方式
	maxLifeTime int64              //有效期
}

type Config struct {
	Sessname   string `toml:"sessname"`
	Expiration int    `toml:"expiration"`
	SaveFile   string `toml:"savefile_path"`
	SaveCycle  int    `toml:"savefile_cycle"`
	Enable     bool   `toml:"enable"`
}

//初始化，创建容器
func Init() {

	//获取配置
	_ = common.ReadConfig("session", &config)

	//log.Printf("Session:%v", config)

	if !config.Enable {
		//不启用Session
		return
	}

	//provide, _ := provides["memory"]

	fm := memory.NewFromMemory()

	manager = &Manager{cookieName: config.Sessname, provider: fm, maxLifeTime: int64(config.Expiration)}

	if config.SaveFile != "" {
		err := manager.provider.LoadFile(config.SaveFile)
		if err != nil {
			log.Printf("Session LoadFile Err:%v", err)
		}
	}

	manager.GC()

}

//实例化一个session管理器
//func NewSessionManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
//
//	provide, ok := provides[provideName]
//	if !ok {
//		return nil, fmt.Errorf("session: unknown provide %q ", provideName)
//	}
//
//	return &Manager{cookieName: cookieName, provider: provide, maxLifeTime: maxLifeTime}, nil
//}

//Session中间件
//判断当前请求的cookie中是否存在有效的session，存在返回，否则创建
func GetSession(w http.ResponseWriter, r *http.Request) (session memory.SessionStore) {

	manager.lock.Lock() //加锁
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie.Value == "" {
		//创建一个
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)

		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime),
			Expires:  time.Now().Add(time.Duration(manager.maxLifeTime)),
			//MaxAge和Expires都可以设置cookie持久化时的过期时长，Expires是老式的过期方法，
			// 如果可以，应该使用MaxAge设置过期时间，但有些老版本的浏览器不支持MaxAge。
			// 如果要支持所有浏览器，要么使用Expires，要么同时使用MaxAge和Expires。
		}

		//log.Printf("过期：%v", time.Now().Add(time.Duration(manager.maxLifeTime)))

		http.SetCookie(w, &cookie)

	} else {

		sid, _ := url.QueryUnescape(cookie.Value) //反转义特殊符号
		session, _ = manager.provider.SessionRead(sid)

	}

	return session
}

/*
//注册 由实现Provider接口的结构体调用
func Register(name string, provide memory.FromMemory) {

	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, ok := provides[name]; ok {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provide

	//========初始化
	//获取配置
	//err := common.ReadConfig("session", &config)
	//
	//log.Printf("Session:%v,%v", config, err)
	//
	//manager = &Manager{cookieName: config.Sessname, provider: provide, maxLifeTime: int64(time.Duration(config.Expiration) * time.Second)}
	//
	//go manager.GC()

}

*/

//var provides = make(map[string]memory.FromMemory)

//生成sessionId
func (manager *Manager) sessionId() string {

	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	//加密
	ctx := md5.New()
	ctx.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	return hex.EncodeToString(ctx.Sum(nil))

}

//销毁session 同时删除cookie
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()

		sid, _ := url.QueryUnescape(cookie.Value)
		_ = manager.provider.SessionDestroy(sid)

		expiration := time.Now()

		cookie := http.Cookie{
			Name:     manager.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration,
			MaxAge:   -1}

		http.SetCookie(w, &cookie)
	}
}

//垃圾回收机制
func (manager *Manager) GC() {

	go func() {

		tick := time.NewTicker(time.Duration(config.SaveCycle) * time.Second)
		for {
			select {
			case <-tick.C:
				manager.lock.Lock()

				manager.provider.SessionGC(int64(time.Duration(config.SaveCycle) * time.Second))

				if config.SaveFile != "" {
					_ = manager.provider.SaveToFile(config.SaveFile)
					//log.Printf("保存Session:%v", err)
				}

				manager.lock.Unlock()

			}
		}
	}()

	//manager.lock.Lock()
	//defer manager.lock.Unlock()
	//
	//
	//manager.provider.SessionGC(int64(time.Duration(config.SaveCycle) * time.Second))
	//
	//time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
	//	manager.GC()
	//})

}
