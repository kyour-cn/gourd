package cache

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

// https://www.codingsky.com/doc/2020/3/14/140.html

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

type Item struct {
	Object     interface{} // 真正的数据项
	Expiration int64       // 生存时间
}

// 判断数据项是否已经过期
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

const (
	// 没有过期时间标志 -不过期
	NoExpiration time.Duration = -1

	// 默认的过期时间
	DefaultExpiration time.Duration = 0
)

type Cache struct {
	defaultExpiration time.Duration   // 默认的过期时间
	Items             map[string]Item // 缓存数据项存储在 map 中
	mu                sync.RWMutex    // 读写锁
	gcInterval        time.Duration   // 过期数据项清理周期
	stopGc            chan bool
	File              string
}

// 过期缓存数据项清理
func (c *Cache) gcLoop() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}

// 删除缓存数据项
func (c *Cache) delete(k string) {
	//timeObj := time.Unix(c.Items[k].Expiration/int64(time.Second), 0).Format("2006-01-02 15:04:05")
	//
	//log.Println("到期：", timeObj)
	//log.Printf("删除：%s---%v\n", k, timeObj)
	delete(c.Items, k)
}

// 遍历删除过期数据项
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	//timeObj := time.Unix(now/int64(time.Second), 0).Format("2006-01-02 15:04:05")
	//log.Printf("清理垃圾:%v,Count:%v", timeObj, len(c.Items))

	for k, v := range c.Items {

		//if v.Expiration <= 0 {
		//	log.Printf("不过期数据:%v", v.Object)
		//} else {
		//	log.Printf("Key（%s）距离过期还有%v", k, time.Duration(now-v.Expiration))
		//}

		if v.Expiration > 0 && now > v.Expiration {
			c.delete(k)

		}
	}
	c.mu.Unlock()

}

// 设置缓存数据项，如果数据项存在则覆盖
func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	//timeObj := time.Unix(e/int64(time.Second), 0).Format("2006-01-02 15:04:05")

	//log.Printf("设置值-%s:%v", k, timeObj)
	c.mu.Lock()
	c.Items[k] = Item{
		Object:     v,
		Expiration: e,
	}
	c.mu.Unlock()
}

// 设置数据项, 没有锁操作
func (c *Cache) set(k string, v interface{}, d time.Duration) {

	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	//timeObj := time.Unix(e/int64(time.Second), 0).Format("2006-01-02 15:04:05")

	//log.Printf("设置值-%s:%v,(==count:%v)", k, timeObj, len(c.Items))

	//timeObj := time.Unix(e/int64(time.Second), 0).Format("2006-01-02 15:04:05")
	//
	//log.Println("到期：", timeObj)

	c.Items[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

// 获取数据项，如果找到数据项，还需要判断数据项是否已经过期
func (c *Cache) get(k string) (interface{}, bool) {
	item, found := c.Items[k]
	if !found {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// 添加数据项，如果数据项已经存在，则返回错误
func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {

		c.mu.Unlock()
		return errors.New("Item " + k + " already exists")
	}
	c.set(k, v, d)
	c.mu.Unlock()

	return nil
}

// 获取数据项
func (c *Cache) Get(k string) (interface{}, bool) {

	c.mu.RLock()
	item, found := c.Items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	c.mu.RUnlock()

	return item.Object, true
}

// 替换一个存在的数据项
func (c *Cache) Replace(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return errors.New("Item " + k + " doesn't exist")
	}
	c.set(k, v, d)
	c.mu.Unlock()

	return nil
}

// 删除一个数据项
func (c *Cache) Delete(k string) {

	c.mu.Lock()
	c.delete(k)
	c.mu.Unlock()

}

// 将缓存数据项写入到 io.Writer 中
func (c *Cache) Save(w io.Writer) (err error) {

	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("Error registering item types with Gob library")
		}
	}()

	c.mu.RLock()
	for _, v := range c.Items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.Items)
	c.mu.RUnlock()

	return
}

// 保存数据项到文件中
func (c *Cache) SaveToFile(file string) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = c.Save(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// 从 io.Reader 中读取数据项
func (c *Cache) Load(r io.Reader) error {

	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {

		c.mu.Lock()

		for k, v := range items {
			ov, found := c.Items[k]
			if !found || ov.Expired() {
				c.Items[k] = v
			}
		}
		c.mu.Unlock()

	}
	return err
}

// 从文件中加载缓存数据项
func (c *Cache) LoadFile(file string) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	if err = c.Load(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// 返回缓存数据项的数量
func (c *Cache) Count() int {
	return len(c.Items)
}

// 清空缓存
func (c *Cache) Flush() {

	c.mu.Lock()
	c.Items = map[string]Item{}
	c.mu.Unlock()

}

// 停止过期缓存清理
func (c *Cache) StopGc() {
	c.stopGc <- true
}

// 创建一个缓存系统
func NewCache(defaultExpiration, gcInterval time.Duration, file string) *Cache {
	c := &Cache{
		defaultExpiration: defaultExpiration,
		gcInterval:        gcInterval,
		Items:             map[string]Item{},
		stopGc:            make(chan bool),
		File:              file,
	}
	// 开始启动过期清理 goroutine
	//go c.gcLoop()

	go func() {
		tick := time.NewTicker(c.gcInterval)
		for {
			select {
			case <-tick.C:
				//log.Printf("Cache.Item:%v\n", Obj.Count())

				//定时执行保存数据
				_ = Obj.SaveToFile(c.File)
				Obj.DeleteExpired()
			}
		}
	}()
	return c
}
