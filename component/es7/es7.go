/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package es7

import (
	"context"
	"errors"
	"strings"
	"time"

	elastic7 "github.com/olivere/elastic/v7"

	"github.com/zlyuancn/zapp/component/conn"
	"github.com/zlyuancn/zapp/core"
)

type ES7 struct {
	app  core.IApp
	conn *conn.Conn
}

func NewES7(app core.IApp) core.IES7Component {
	return &ES7{
		app:  app,
		conn: conn.NewConn(),
	}
}

func (e *ES7) GetES7(name ...string) *elastic7.Client {
	return e.conn.GetInstance(e.makeClient, name...).(*elastic7.Client)
}

func (e *ES7) makeClient(name string) (interface{}, error) {
	// 获取配置
	conf, ok := e.app.GetConfig().Config().Components.ES7[name]
	if !ok {
		return nil, errors.New("试图获取未注册的es7")
	}
	if conf.Address == "" {
		return nil, errors.New("es7的address为空")
	}

	opts := []elastic7.ClientOptionFunc{
		elastic7.SetURL(strings.Split(conf.Address, ",")...),
		elastic7.SetSniff(conf.Sniff),
		elastic7.SetHealthcheck(conf.Healthcheck == nil || *conf.Healthcheck),
		elastic7.SetGzip(conf.GZip),
	}
	if conf.UserName != "" || conf.Password != "" {
		opts = append(opts, elastic7.SetBasicAuth(conf.UserName, conf.Password))
	}
	if conf.Retry > 0 {
		ticks := make([]int, conf.Retry)
		for i := 0; i < conf.Retry; i++ {
			ticks[i] = conf.RetryInterval
		}
		opts = append(opts, elastic7.SetRetrier(elastic7.NewBackoffRetrier(elastic7.NewSimpleBackoff(ticks...))))
	}

	ctx := context.Background()
	if conf.DialTimeout > 0 {
		c, cancel := context.WithTimeout(ctx, time.Duration(conf.DialTimeout)*time.Millisecond)
		defer cancel()
		ctx = c
	}

	return elastic7.DialContext(ctx, opts...)
}

func (e *ES7) Close() {
	e.conn.IterInstance(func(name string, instance interface{}) {
		instance.(*elastic7.Client).Stop()
	})
}
