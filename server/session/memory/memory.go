package memory

import (
	"container/list"
	"encoding/gob"
	"errors"
	"github.com/kyour-cn/gourd/server/session"
	"io"
	"os"
	"sync"
	"time"
)

var pder = &FromMemory{list: list.New()}

func init() {

	pder.sessions = make(map[string]*list.Element, 0)
	//注册  memory 调用的时候一定要一致
	session.Register("memory", pder)

}

//session实现
type SessionStore struct {
	sid              string                      //session id 唯一标示
	LastAccessedTime time.Time                   //最后访问时间
	value            map[interface{}]interface{} //session 里面存储的值
}

//设置
func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	_ = pder.SessionUpdate(st.sid)
	return nil
}

//获取session
func (st *SessionStore) Get(key interface{}) interface{} {
	_ = pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

//获取全部数据
func (st *SessionStore) GetAll() interface{} {

	return st.value

}

//删除
func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	_ = pder.SessionUpdate(st.sid)
	return nil
}

//清空Session
func (st *SessionStore) Clear() error {
	return pder.SessionDestroy(st.sid)
}

//获取session'Id
func (st *SessionStore) SessionID() string {
	return st.sid
}

//session来自内存 实现
type FromMemory struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做 gc
}

func (frommemory *FromMemory) SessionInit(sid string) (session.Session, error) {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, LastAccessedTime: time.Now(), value: v}
	element := frommemory.list.PushBack(newsess)
	frommemory.sessions[sid] = element
	return newsess, nil
}

func (frommemory *FromMemory) SessionRead(sid string) (session.Session, error) {
	if element, ok := frommemory.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := frommemory.SessionInit(sid)
		return sess, err
	}
}

func (frommemory *FromMemory) SessionDestroy(sid string) error {
	if element, ok := frommemory.sessions[sid]; ok {
		delete(frommemory.sessions, sid)
		frommemory.list.Remove(element)
		return nil
	}
	return nil
}

func (frommemory *FromMemory) SessionGC(maxLifeTime int64) {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	for {
		element := frommemory.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).LastAccessedTime.Unix() + maxLifeTime) <
			time.Now().Unix() {
			frommemory.list.Remove(element)
			delete(frommemory.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (frommemory *FromMemory) SessionUpdate(sid string) error {
	frommemory.lock.Lock()
	defer frommemory.lock.Unlock()
	if element, ok := frommemory.sessions[sid]; ok {
		element.Value.(*SessionStore).LastAccessedTime = time.Now()
		frommemory.list.MoveToFront(element)
		return nil
	}
	return nil
}

// 将缓存数据项写入到 io.Writer 中
func (fm *FromMemory) save(w io.Writer) (err error) {

	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("Error registering item types with Gob library")
		}
	}()

	fm.lock.Lock()
	for _, v := range fm.sessions {
		gob.Register(v.Value)
	}
	err = enc.Encode(&fm.sessions)
	fm.lock.Unlock()

	return
}

// 保存数据项到文件中
func (fm *FromMemory) SaveToFile(file string) error {

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = fm.save(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// 从 io.Reader 中读取数据项
func (c *FromMemory) load(r io.Reader) error {

	dec := gob.NewDecoder(r)
	items := map[string]*list.Element{}
	err := dec.Decode(&items)
	if err == nil {

		c.lock.Lock()

		for k, v := range items {
			_, ok := c.sessions[k]
			if !ok {
				c.sessions[k] = v
			}
		}
		c.lock.Unlock()

	}
	return err
}

// 从文件中加载缓存数据项
func (c *FromMemory) LoadFile(file string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	if err = c.load(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
