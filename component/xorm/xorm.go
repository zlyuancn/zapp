/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package xorm

import (
	"errors"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // mssql
	_ "github.com/go-sql-driver/mysql"   // mysql
	_ "github.com/lib/pq"                // postgres
	_ "github.com/mattn/go-sqlite3"      // sqlite
	"xorm.io/xorm"

	"github.com/zlyuancn/zapp/component/conn"
	"github.com/zlyuancn/zapp/core"
)

type Xorm struct {
	app  core.IApp
	conn *conn.Conn
}

func NewXorm(app core.IApp) core.IXormComponent {
	x := &Xorm{
		app:  app,
		conn: conn.NewConn(),
	}
	return x
}

func (x *Xorm) GetXorm(name ...string) *xorm.Engine {
	return x.conn.GetInstance(x.makeClient, name...).(*xorm.Engine)
}

func (x *Xorm) makeClient(name string) (interface{}, error) {
	// 获取配置
	conf, ok := x.app.GetConfig().Config().Components.Xorm[name]
	if !ok {
		return nil, errors.New("试图获取未注册的xorm")
	}

	e, err := xorm.NewEngine(conf.Driver, conf.Source)
	if err != nil {
		return nil, err
	}
	e.SetMaxIdleConns(conf.MaxIdleConns)
	e.SetMaxOpenConns(conf.MaxOpenConns)
	e.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Millisecond)
	return e, nil
}

func (x *Xorm) Close() {
	x.conn.IterInstance(func(name string, instance interface{}) {
		_ = instance.(*xorm.Engine).Close()
	})
}
