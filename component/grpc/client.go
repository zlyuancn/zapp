/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/18
   Description :
-------------------------------------------------
*/

package grpc

import (
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/component/grpc/registry/local"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

type Client struct {
}

func NewClient(app core.IApp) *Client {
	for name, conf := range app.GetConfig().Config().GrpcClient {
		switch conf.Registry {
		case "", local.Schema:
			local.RegistryAddress(name, conf.Address)
		default:
			utils.Fatal("未定义的Grpc注册器", zap.String("Registry", conf.Registry))
		}
	}
	return &Client{}
}
