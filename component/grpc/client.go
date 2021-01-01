/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/11/18
   Description :
-------------------------------------------------
*/

package grpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zapp/component/conn"
	"github.com/zlyuancn/zapp/component/grpc/balance/round_robin"
	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

var typeOfGrpcClientConn = reflect.TypeOf((*grpc.ClientConn)(nil))

type Client struct {
	app        core.IApp
	conn       *conn.Conn
	creatorMap map[string]reflect.Value
}

type Conn struct {
	cc     *grpc.ClientConn
	client interface{}
}

func NewClient(app core.IApp) core.IGrpcComponent {
	g := &Client{
		app:        app,
		conn:       conn.NewConn(),
		creatorMap: make(map[string]reflect.Value),
	}

	configs := app.GetConfig().Config().Components.GrpcClient
	for name, conf := range configs {
		if conf.Registry == "" {
			conf.Registry = consts.DefaultConfig_GrpcClient_Registry
		}
		if conf.Balance == "" {
			conf.Balance = consts.DefaultConfig_GrpcClient_Balance
		}
		configs[name] = conf

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
	g.conn.IterInstance(func(name string, instance interface{}) {
		_ = instance.(*Conn).cc.Close()
	})
}

func (g *Client) RegistryGrpcClientCreator(name string, creator interface{}) {
	createType := reflect.TypeOf(creator)
	if createType.Kind() != reflect.Func {
		logger.Log.Fatal("grpc客户端建造者必须是函数")
		return
	}

	if createType.NumIn() != 1 {
		logger.Log.Fatal("grpc客户端建造者入参为1个")
		return
	}

	if !createType.In(0).AssignableTo(typeOfGrpcClientConn) {
		logger.Log.Fatal("grpc客户端建造者入参类型必须是 *grpc.ClientConn")
		return
	}

	if createType.NumOut() != 1 {
		logger.Log.Fatal("grpc客户端建造者必须有一个返回值")
		return
	}

	g.creatorMap[name] = reflect.ValueOf(creator)
}

func (g *Client) GetGrpcClient(name string) interface{} {
	return g.conn.GetInstance(g.makeClient, name).(*Conn).client
}

func (g *Client) makeConn(name, scheme, balance string, timeout int) (*grpc.ClientConn, error) {
	if timeout <= 0 {
		timeout = consts.DefaultConfig_GrpcClient_DialTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	return grpc.DialContext(ctx,
		scheme+":///"+name,
		grpc.WithInsecure(),   // 不安全连接
		g.getBalance(balance), // 均衡器
		grpc.WithBlock(),      // 等待连接
	)
}

func (g *Client) makeClient(name string) (interface{}, error) {
	// 获取配置
	conf, ok := g.app.GetConfig().Config().Components.GrpcClient[name]
	if !ok {
		return nil, errors.New("试图获取未注册的grpc客户端")
	}

	// 获取建造者
	creator, ok := g.creatorMap[name]
	if !ok {
		return nil, errors.New("未注册grpc客户端建造者")
	}

	cc, err := g.makeConn(name, conf.Registry, conf.Balance, conf.DialTimeout)
	if err != nil {
		return nil, fmt.Errorf("make conn error: name=%s, registry=%s, balance=%s, err: %s", name, conf.Registry, conf.Balance, err)
	}

	client := creator.Call([]reflect.Value{reflect.ValueOf(cc)})[0].Interface()

	return &Conn{
		cc:     cc,
		client: client,
	}, nil
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
