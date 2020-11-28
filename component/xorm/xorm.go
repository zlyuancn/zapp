/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package xorm

import (
	"errors"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // mssql
	_ "github.com/go-sql-driver/mysql"   // mysql
	_ "github.com/lib/pq"                // postgres
	_ "github.com/mattn/go-sqlite3"      // sqlite
	"go.uber.org/zap"
	"xorm.io/xorm"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/utils"
)

type Xorm struct {
	app core.IApp
	mx  sync.RWMutex

	engineMap map[string]*Engine
}

type Engine struct {
	wg     sync.WaitGroup
	engine *xorm.Engine
	e      error
}

func NewXorm(app core.IApp) core.IXormComponent {
	x := &Xorm{
		app:       app,
		engineMap: make(map[string]*Engine),
	}
	return x
}

func (x *Xorm) GetXorm(name ...string) *xorm.Engine {
	if len(name) > 0 {
		return x.getEngine(name[0])
	}
	return x.getEngine(consts.DefaultComponentName)
}

func (x *Xorm) getEngine(name string) *xorm.Engine {
	x.mx.RLock()
	engine, ok := x.engineMap[name]
	x.mx.RUnlock()

	if ok {
		engine.wg.Wait()
		if engine.e != nil {
			logger.Log.Panic(zap.String("name", name), zap.Error(engine.e))
		}
		return engine.engine
	}

	x.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if engine, ok = x.engineMap[name]; ok {
		x.mx.Unlock()

		engine.wg.Wait()
		if engine.e != nil {
			logger.Log.Panic(zap.String("name", name), zap.Error(engine.e))
		}
		return engine.engine
	}

	// 占位置
	engine = new(Engine)
	engine.wg.Add(1)
	defer engine.wg.Done()
	x.engineMap[name] = engine

	x.mx.Unlock()

	return x.makeClient(name, engine)
}

func (x *Xorm) makeClient(name string, engine *Engine) *xorm.Engine {
	// 获取配置
	conf, ok := x.app.GetConfig().Config().Xorm[name]
	if !ok {
		engine.e = errors.New("试图获取未注册的xorm")
		logger.Log.Panic(zap.String("name", name), zap.Error(engine.e))
	}

	err := utils.Recover.WarpCall(func() error {
		e, err := xorm.NewEngine(conf.Driver, conf.Source)
		if err != nil {
			return err
		}
		e.SetMaxIdleConns(conf.MaxIdleConns)
		e.SetMaxOpenConns(conf.MaxOpenConns)
		e.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Millisecond)
		engine.engine = e
		return nil
	})
	if err != nil {
		engine.e = err
		x.mx.Lock()
		delete(x.engineMap, name)
		x.mx.Unlock()
		logger.Log.Panic(zap.String("name", name), zap.Error(engine.e))
	}
	return engine.engine
}

func (x *Xorm) Close() {
	for _, e := range x.engineMap {
		if e.engine != nil {
			_ = e.engine.Close()
		}
	}
}
