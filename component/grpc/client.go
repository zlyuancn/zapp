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
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/zlyuancn/zapp/component/grpc/balance/round_robin"
	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
	"github.com/zlyuancn/zapp/utils"
)

var typeOfGrpcClientConn = reflect.TypeOf((*grpc.ClientConn)(nil))

type Client struct {
	app core.IApp
	mx  sync.RWMutex

	connMap    map[string]*Conn
	creatorMap map[string]reflect.Value
}

type Conn struct {
	wg     sync.WaitGroup
	cc     *grpc.ClientConn
	client interface{}
	e      error
}

func NewClient(app core.IApp) core.IGrpcComponent {
	g := &Client{
		app:        app,
		connMap:    make(map[string]*Conn),
		creatorMap: make(map[string]reflect.Value),
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
	defer conn.wg.Done()
	g.connMap[name] = conn

	g.mx.Unlock()

	return g.makeClient(name, conn)
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

func (g *Client) makeClient(name string, conn *Conn) interface{} {
	// 获取配置
	conf, ok := g.app.GetConfig().Config().GrpcClient[name]
	if !ok {
		conn.e = errors.New("试图获取未注册的grpc客户端")
		logger.Log.Panic(conn.e, zap.String("name", name))
	}

	// 获取建造者
	creator, ok := g.creatorMap[name]
	if !ok {
		conn.e = errors.New("未注册grpc客户端建造者")
		logger.Log.Panic(conn.e, zap.String("name", name))
	}

	// 构建cc
	var cc *grpc.ClientConn
	err := utils.Recover.WarpCall(func() error {
		var e error
		cc, e = g.makeConn(name, conf.Registry, conf.Balance, conf.DialTimeout)
		return e
	})
	if err != nil {
		conn.e = err
		g.mx.Lock()
		delete(g.connMap, name)
		g.mx.Unlock()
		logger.Log.Panic(conn.e, zap.String("name", name), zap.String("registry", conf.Registry), zap.String("balance", conf.Balance), zap.Error(err))
	}

	// 构建客户端
	var client interface{}
	err = utils.Recover.WarpCall(func() error {
		client = creator.Call([]reflect.Value{reflect.ValueOf(cc)})[0].Interface()
		return nil
	})
	if err != nil {
		conn.e = err
		g.mx.Lock()
		delete(g.connMap, name)
		g.mx.Unlock()
		logger.Log.Panic(conn.e, zap.String("name", name), zap.String("registry", conf.Registry), zap.String("balance", conf.Balance), zap.Error(err))
	}

	conn.client = client
	conn.cc = cc
	return conn.client
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
