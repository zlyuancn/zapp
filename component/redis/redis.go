/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/28
   Description :
-------------------------------------------------
*/

package redis

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/zlyuancn/zapp/component/conn"
	"github.com/zlyuancn/zapp/core"
)

type Redis struct {
	app  core.IApp
	conn *conn.Conn
}

func NewRedis(app core.IApp) core.IRedisComponent {
	return &Redis{
		app:  app,
		conn: conn.NewConn(),
	}
}

func (r *Redis) GetRedis(name ...string) redis.UniversalClient {
	return r.conn.GetInstance(r.makeClient, name...).(redis.UniversalClient)
}

func (r *Redis) makeClient(name string) (interface{}, error) {
	// 获取配置
	conf, ok := r.app.GetConfig().Config().Components.Redis[name]
	if !ok {
		return nil, errors.New("试图获取未注册的redis")
	}
	if conf.Address == "" {
		return nil, errors.New("redis的address为空")
	}

	if conf.IsCluster {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(conf.Address, ","),
			Password:     conf.Password,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		}), nil
	}

	return redis.NewClient(&redis.Options{
		Addr:         conf.Address,
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
	}), nil
}

func (r *Redis) Close() {
	r.conn.IterInstance(func(name string, instance interface{}) {
		_ = instance.(redis.UniversalClient).Close()
	})
}
