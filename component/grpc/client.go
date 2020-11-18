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
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"

	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type Grpc struct {
	app core.IApp
	mx  sync.RWMutex

	schemeMap map[string]string
	connMap   map[string]*grpc.ClientConn
	clientMap map[string]interface{}
}

func NewClient(app core.IApp) core.IGrpcClient {
	schemeMap := make(map[string]string, len(app.GetConfig().Config().GrpcClient))
	for name, conf := range app.GetConfig().Config().GrpcClient {
		scheme := conf.Registry
		if scheme == "" {
			scheme = local.Schema
		}

		schemeMap[name] = scheme

		switch scheme {
		case local.Schema:
			local.RegistryAddress(name, conf.Address)
		default:
			utils.Fatal("未定义的Grpc注册器", zap.String("Registry", scheme))
		}
	}

	return &Grpc{
		app:       app,
		schemeMap: schemeMap,
		connMap:   make(map[string]*grpc.ClientConn, len(schemeMap)),
		clientMap: make(map[string]interface{}, len(schemeMap)),
	}
}

func (g *Grpc) Close() {
	for _, conn := range g.connMap {
		_ = conn.Close()
	}
}

func (g *Grpc) GetGrpcClient(name string, creator func(cc *grpc.ClientConn) interface{}) interface{} {
	g.mx.RLock()
	client, ok := g.clientMap[name]
	g.mx.RUnlock()

	if ok {
		return client
	}

	g.mx.Lock()
	defer g.mx.Unlock()

	if client, ok = g.clientMap[name]; ok {
		return client
	}

	conn, ok := g.connMap[name]
	if !ok {
		scheme, ok := g.schemeMap[name]
		if !ok {
			utils.Panic("试图获取未注册的grpc客户端", zap.String("name", name))
		}

		cc, err := g.makeConn(name, scheme)
		if err != nil {
			utils.Panic("构建grpc客户端conn失败", zap.String("name", name), zap.String("scheme", scheme), zap.Error(err))
		}
		g.connMap[name], conn = cc, cc
	}

	client = creator(conn)
	g.clientMap[name] = client

	return client
}

func (g *Grpc) makeConn(name, scheme string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		scheme+":///"+name,
		grpc.WithInsecure(),                                                                                    // 不安全连接
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{ "loadBalancingConfig": [{"%v": {}}] }`, roundrobin.Name)), // 轮询
		grpc.WithBlock(),                                                                                       // 等待连接
	)
	return conn, err
}
