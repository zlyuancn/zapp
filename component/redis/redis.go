/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package redis

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/utils"
)

type Redis struct {
	app core.IApp
	mx  sync.RWMutex

	clientMap map[string]*Client
}

type Client struct {
	wg     sync.WaitGroup
	client redis.UniversalClient
	e      error
}

func NewRedis(app core.IApp) core.IRedisComponent {
	return &Redis{
		app:       app,
		clientMap: make(map[string]*Client),
	}
}

func (r *Redis) GetRedis(name ...string) redis.UniversalClient {
	if len(name) > 0 {
		return r.getClient(name[0])
	}
	return r.getClient(consts.DefaultComponentName)
}

func (r *Redis) getClient(name string) redis.UniversalClient {
	r.mx.RLock()
	client, ok := r.clientMap[name]
	r.mx.RUnlock()

	if ok {
		client.wg.Wait()
		if client.e != nil {
			logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
		}
		return client.client
	}

	r.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if client, ok = r.clientMap[name]; ok {
		r.mx.Unlock()

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
	r.clientMap[name] = client

	r.mx.Unlock()

	return r.makeClient(name, client)
}

func (r *Redis) makeClient(name string, client *Client) redis.UniversalClient {
	// 获取配置
	conf, ok := r.app.GetConfig().Config().Redis[name]
	if !ok {
		client.e = errors.New("试图获取未注册的redis")
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}
	if conf.Address == "" {
		client.e = errors.New("redis的address为空")
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}

	err := utils.Recover.WarpCall(func() error {
		if conf.IsCluster {
			client.client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:        strings.Split(conf.Address, ","),
				Password:     conf.Password,
				PoolSize:     conf.PoolSize,
				ReadTimeout:  time.Duration(conf.ReadTimeout * 1e6),
				WriteTimeout: time.Duration(conf.WriteTimeout * 1e6),
				DialTimeout:  time.Duration(conf.DialTimeout * 1e6),
			})
		} else {
			client.client = redis.NewClient(&redis.Options{
				Addr:         conf.Address,
				Password:     conf.Password,
				DB:           conf.DB,
				PoolSize:     conf.PoolSize,
				ReadTimeout:  time.Duration(conf.ReadTimeout * 1e6),
				WriteTimeout: time.Duration(conf.WriteTimeout * 1e6),
				DialTimeout:  time.Duration(conf.DialTimeout * 1e6),
			})
		}
		return nil
	})
	if err != nil {
		client.e = err
		r.mx.Lock()
		delete(r.clientMap, name)
		r.mx.Unlock()
		logger.Log.Panic(zap.String("name", name), zap.Error(client.e))
	}
	return client.client
}

func (r *Redis) Close() {
	for _, c := range r.clientMap {
		if c.client != nil {
			_ = c.client.Close()
		}
	}
}
