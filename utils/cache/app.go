package cache

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

import (
	"github.com/kyour-cn/gourd/common"
	"github.com/kyour-cn/gourd/utils/gut"
	"log"
	"time"
)

var Obj Cache

var isInit bool

type Config struct {
	//是否启用
	Enable bool `toml:"enable"`
	//缓存默认过期时间（秒） ，传递-1时取此值
	Expiration int32 `toml:"expiration"`
	//是否实时输出到缓存文件
	RealtimeSave bool `toml:"realtime_save"`
	//缓存文件保存周期（秒），realtime_save=false 才会生效
	SavefileCycle int32 `toml:"savefile_cycle"`
	//缓存文件保存路径
	SavefilePath string `toml:"savefile_path"`
}

var config Config

//初始化缓存
func Init() {

	//判断是否未初始化
	if !isInit {

		err := common.ReadConfig("cache", &config)
		if err != nil {
			log.Printf("Cache 配置错误：%v\n", err)
		}

		if !config.Enable {
			//不启用缓存
			return
		}

		//创建缓存对象
		Obj = *NewCache(time.Duration(config.Expiration)*time.Second, time.Duration(config.SavefileCycle)*time.Second, config.SavefilePath)

		exist, err := gut.PathExists(config.SavefilePath)
		if err != nil {
			log.Printf("Cache PathExists Error：%v\n", err)
		}
		if exist {
			//读取文件
			err = Obj.LoadFile(config.SavefilePath)
			if err != nil {
				log.Printf("Cache Error:%v\n", err)
			}
		}

		//已初始化
		isInit = true

		//实时保存关闭-定时保存
		//if !config.RealtimeSave {
		//log.Printf("Cache.Init \n", Obj.Items)
		//time.AfterFunc(5*time.Minute, func() {
		//
		//	log.Printf("expired")
		//
		//})
		//time.Duration(config.SavefileCycle)*time.Second

		/*
			go func() {
				tick := time.NewTicker(time.Duration(config.SavefileCycle) * time.Second)
				for {
					select {
					case <-tick.C:
						//log.Printf("Cache.Item:%v\n", Obj.Count())

						//定时执行保存数据
						_ = Obj.SaveToFile(config.SavefilePath)
						Obj.DeleteExpired()
					}
				}
			}()

		*/

		/*
			time.AfterFunc(time.Minute, func() {
				for {

					//time.Sleep(time.Duration(config.SavefileCycle) * time.Second)

					log.Printf("Cache.Item:%v\n", len(Obj.Items))

					//定时执行保存数据
					_ = Obj.SaveToFile(config.SavefilePath)
				}
			})

		*/
		//}
	}

}

//设置数据
func Set(k string, v interface{}, d time.Duration) {

	Init()

	Obj.Set(k, v, d)

	//如果实时保存
	if config.RealtimeSave {
		err := Obj.SaveToFile(config.SavefilePath)
		if err != nil {
			log.Fatalf("Cache SaveToFile Error:%v", err)
		}
	}

}

//获取缓存数据
func Get(key string) (interface{}, bool) {

	Init()

	return Obj.Get(key)

}

func Add(k string, v interface{}, d time.Duration) (err error) {

	Init()

	err = Obj.Add(k, v, d)

	//如果实时保存
	if config.RealtimeSave {
		err = Obj.SaveToFile(config.SavefilePath)
		if err != nil {
			log.Printf("Cache Add Error:%v", err)
		}
	}

	return

}
