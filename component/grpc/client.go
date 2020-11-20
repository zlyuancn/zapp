/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/18
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zapp/component/grpc/balance/round_robin"
	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type Client struct {
	app core.IApp
	mx  sync.RWMutex

	connMap map[string]*Conn
}

type Conn struct {
	wg     sync.WaitGroup
	cc     *grpc.ClientConn
	client interface{}
	e      error
}

func NewClient(app core.IApp) core.IGrpcClient {
	g := &Client{
		app:     app,
		connMap: make(map[string]*Conn),
	}

	for name, conf := range app.GetConfig().Config().GrpcClient {
		if conf.Registry == "" {
			conf.Registry = consts.DefaultConfig_GrpcClient_Registry
		}
		if conf.Balance == "" {
			conf.Balance = consts.DefaultConfig_GrpcClient_Balance
		}
		app.GetConfig().Config().GrpcClient[name] = conf

		switch conf.Registry {
		case local.Name:
			local.RegistryAddress(name, conf.Address)
		default:
			logger.Log.Fatal("未定义的Grpc注册器", zap.String("registry", conf.Registry))
		}

		_ = g.getBalance(conf.Balance)
	}

	return g
}

func (g *Client) Close() {
	for _, conn := range g.connMap {
		if conn.cc != nil {
			_ = conn.cc.Close()
		}
	}
}

func (g *Client) GetGrpcClient(name string, creator func(cc *grpc.ClientConn) interface{}) interface{} {
	g.mx.RLock()
	conn, ok := g.connMap[name]
	g.mx.RUnlock()

	if ok {
		conn.wg.Wait()
		if conn.e != nil {
			logger.Log.Panic(conn.e, zap.String("name", name))
		}
		return conn.client
	}

	g.mx.Lock()

	// 再获取一次, 它可能在获取锁的过程中完成了
	if conn, ok = g.connMap[name]; ok {
		g.mx.Unlock()

		conn.wg.Wait()
		if conn.e != nil {
			logger.Log.Panic(conn.e, zap.String("name", name))
		}
		return conn.client
	}

	// 占位置
	conn = new(Conn)
	conn.wg.Add(1)
	g.connMap[name] = conn

	g.mx.Unlock()

	// 获取配置, 如果配置不存在则不需要删除conn, 因为它永远不会有配置了
	conf, ok := g.app.GetConfig().Config().GrpcClient[name]
	if !ok {
		conn.e = errors.New("试图获取未注册的grpc客户端")
		conn.wg.Done()
		logger.Log.Panic(conn.e, zap.String("name", name))
	}

	cc, err := g.makeConn(name, conf.Registry, conf.Balance)
	if err != nil {
		conn.e = errors.New("构建grpc客户端conn失败")
		conn.wg.Done()

		// 删除位置
		g.mx.Lock()
		delete(g.connMap, name)
		g.mx.Unlock()

		logger.Log.Panic(conn.e,
			zap.String("name", name),
			zap.String("registry", conf.Registry),
			zap.String("balance", conf.Balance),
			zap.Error(err))
	}

	conn.cc = cc
	conn.client = creator(conn.cc)
	conn.wg.Done()

	return conn.client
}

func (g *Client) makeConn(name, scheme, balance string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		scheme+":///"+name,
		grpc.WithInsecure(),   // 不安全连接
		g.getBalance(balance), // 均衡器
		grpc.WithBlock(),      // 等待连接
	)
	return conn, err
}

func (g *Client) getBalance(balance string) grpc.DialOption {
	switch balance {
	case round_robin.Name:
		return round_robin.Balance()
	default:
		logger.Log.Fatal("未定义的Grpc均衡器", zap.String("balance", balance))
	}
	return nil
}
