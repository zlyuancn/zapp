/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/1/1
   Description :
-------------------------------------------------
*/

package conn

import (
	"sync"

	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/utils"
)

type CreatorFunc = func(name string) (interface{}, error)

// 连接器
type Conn struct {
	wgs map[string]*connWaitGroup
	mx  sync.RWMutex
}

type connWaitGroup struct {
	instance interface{}
	e        error
	wg       sync.WaitGroup
}

func NewConn() *Conn {
	return &Conn{
		wgs: make(map[string]*connWaitGroup),
	}
}

// 获取实例
func (c *Conn) GetInstance(creator CreatorFunc, name ...string) interface{} {
	if len(name) == 0 {
		return c.getInstance(creator, consts.DefaultComponentName)
	}
	return c.getInstance(creator, name[0])
}

func (c *Conn) getInstance(creator CreatorFunc, name string) interface{} {
	c.mx.RLock()
	wg, ok := c.wgs[name]
	c.mx.RUnlock()

	if ok {
		wg.wg.Wait()
		if wg.e != nil {
			logger.Log.Panic(wg.e.Error(), zap.String("name", name))
		}
		return wg.instance
	}

	c.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if wg, ok = c.wgs[name]; ok {
		c.mx.Unlock()
		wg.wg.Wait()
		if wg.e != nil {
			logger.Log.Panic(wg.e.Error(), zap.String("name", name))
		}
		return wg.instance
	}

	// 占位置
	wg = new(connWaitGroup)
	wg.wg.Add(1)
	c.wgs[name] = wg
	c.mx.Unlock()

	var err error
	err = utils.Recover.WarpCall(func() error {
		wg.instance, err = creator(name)
		return err
	})

	// 如果出现错误, 删除占位
	if err != nil {
		wg.e = err
		wg.wg.Done()
		c.mx.Lock()
		delete(c.wgs, name)
		c.mx.Unlock()
		logger.Log.Panic(wg.e.Error(), zap.String("name", name))
	}

	wg.wg.Done()
	return wg.instance
}

// 迭代实例
func (c *Conn) IterInstance(fn func(name string, instance interface{})) {
	c.mx.Lock()
	defer c.mx.Unlock()

	for name, wg := range c.wgs {
		if wg.instance != nil {
			fn(name, wg.instance)
		}
	}
}
