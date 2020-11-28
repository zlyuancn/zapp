/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package es7

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	elastic7 "github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/utils"
)

type ES7 struct {
	app core.IApp
	mx  sync.RWMutex

	clientMap map[string]*Client
}

type Client struct {
	wg     sync.WaitGroup
	client *elastic7.Client
	e      error
}

func NewES7(app core.IApp) core.IES7Component {
	return &ES7{
		app:       app,
		clientMap: make(map[string]*Client),
	}
}

func (e *ES7) GetES7(name ...string) *elastic7.Client {
	if len(name) > 0 {
		return e.getClient(name[0])
	}
	return e.getClient(consts.DefaultComponentName)
}

func (e *ES7) getClient(name string) *elastic7.Client {
	e.mx.RLock()
	client, ok := e.clientMap[name]
	e.mx.RUnlock()

	if ok {
		client.wg.Wait()
		if client.e != nil {
			logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
		}
		return client.client
	}

	e.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if client, ok = e.clientMap[name]; ok {
		e.mx.Unlock()

		client.wg.Wait()
		if client.e != nil {
			logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
		}
		return client.client
	}

	// 占位置
	client = new(Client)
	client.wg.Add(1)
	defer client.wg.Done()
	e.clientMap[name] = client

	e.mx.Unlock()

	return e.makeClient(name, client)
}

func (e *ES7) makeClient(name string, client *Client) *elastic7.Client {
	// 获取配置
	conf, ok := e.app.GetConfig().Config().ES7[name]
	if !ok {
		client.e = errors.New("试图获取未注册的es7")
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}
	if conf.Address == "" {
		client.e = errors.New("es7的address为空")
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}

	err := utils.Recover.WarpCall(func() error {
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
			c, cancel := context.WithTimeout(ctx, time.Duration(conf.DialTimeout*1e6))
			defer cancel()
			ctx = c
		}

		c, err := elastic7.DialContext(ctx, opts...)
		if err != nil {
			return fmt.Errorf("连接失败: %s", err)
		}
		client.client = c
		return nil
	})
	if err != nil {
		client.e = err
		e.mx.Lock()
		delete(e.clientMap, name)
		e.mx.Unlock()
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}
	return client.client
}

func (e *ES7) Close() {
	for _, c := range e.clientMap {
		if c.client != nil {
			c.client.Stop()
		}
	}
}
