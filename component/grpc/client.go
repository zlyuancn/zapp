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
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zapp/component/grpc/balance/round_robin"
	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type Client struct {
	app core.IApp
	mx  sync.RWMutex

	connMap   map[string]*grpc.ClientConn
	clientMap map[string]interface{}
}

func NewClient(app core.IApp) core.IGrpcClient {
	g := &Client{
		app:       app,
		connMap:   make(map[string]*grpc.ClientConn),
		clientMap: make(map[string]interface{}),
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
			utils.Fatal("未定义的Grpc注册器", zap.String("registry", conf.Registry))
		}

		_ = g.getBalance(conf.Balance)
	}

	return g
}

func (g *Client) Close() {
	for _, conn := range g.connMap {
		_ = conn.Close()
	}
}

func (g *Client) GetGrpcClient(name string, creator func(cc *grpc.ClientConn) interface{}) interface{} {
	g.mx.RLock()
	client, ok := g.clientMap[name]
	g.mx.RUnlock()

	if ok {
		return client
	}

	g.mx.Lock()
	defer g.mx.Unlock()

	// todo 这里等待时间过长, 考虑加个占位符优化下

	if client, ok = g.clientMap[name]; ok {
		return client
	}

	conn, ok := g.connMap[name]
	if !ok {
		conf, ok := g.app.GetConfig().Config().GrpcClient[name]
		if !ok {
			utils.Panic("试图获取未注册的grpc客户端", zap.String("name", name))
		}

		cc, err := g.makeConn(name, conf.Registry, conf.Balance)
		if err != nil {
			utils.Panic("构建grpc客户端conn失败",
				zap.String("name", name),
				zap.String("registry", conf.Registry),
				zap.String("balance", conf.Balance),
				zap.Error(err))
		}
		g.connMap[name], conn = cc, cc
	}

	client = creator(conn)
	g.clientMap[name] = client

	return client
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
		utils.Fatal("未定义的Grpc均衡器", zap.String("balance", balance))
	}
	return nil
}
