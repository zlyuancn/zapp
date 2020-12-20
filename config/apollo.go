/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/12/18
   Description :
-------------------------------------------------
*/

package config

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/shima-park/agollo"
	"github.com/spf13/viper"
	"github.com/zlyuancn/zutils"

	"github.com/zlyuancn/zapp/consts"
)

type Namespace string

const (
	FrameNamespace              Namespace = "frame"
	LogNamespace                          = "log"
	ApiServiceNamespace                   = "api_service"
	GrpcServiceNamespace                  = "grpc_service"
	CronServiceNamespace                  = "cron_service"
	MysqlBinlogServiceNamespace           = "mysql_binlog_service"
	GrpcClientNamespace                   = "grpc_client"
	XormNamespace                         = "xorm"
	RedisNamespace                        = "redis"
	ES7Namespace                          = "es7"
	CacheNamespace                        = "cache"
)

type ApolloConfig struct {
	Address              string      // apollo-api地址, 多个地址用英文逗号连接
	AppId                string      // 应用名
	AccessKey            string      // 验证key, 优先级高于基础认证
	AuthBasicUser        string      // 基础认证用户名
	AuthBasicPassword    string      // 基础认证密码
	Cluster              string      // 集群名
	AlwaysLoadFromRemote bool        // 总是从远程获取, 在远程加载失败时不会从备份文件加载
	BackupFile           string      // 备份文件名
	Namespaces           []Namespace // 要加载的命名空间, 一个命名空间相当于一个配置组
}

// 从viper构建apollo配置
func makeApolloConfigFromViper(vi *viper.Viper) (*ApolloConfig, error) {
	var conf ApolloConfig
	err := vi.UnmarshalKey(consts.ConfigGroupName_Apollo, &conf)
	return &conf, err
}

// 从apollo中获取配置构建viper
func makeViperFromApollo(conf *ApolloConfig) (*viper.Viper, error) {
	// 构建选项
	opts := []agollo.Option{
		agollo.AutoFetchOnCacheMiss(),                                       // 当本地缓存中namespace不存在时，尝试去apollo缓存接口去获取
		agollo.Cluster(zutils.Ternary.Or(conf.Cluster, "default").(string)), // 集群名
	}
	if !conf.AlwaysLoadFromRemote {
		opts = append(opts, agollo.FailTolerantOnBackupExists()) // 从服务获取数据失败时从备份文件加载
	}
	if conf.BackupFile != "" {
		opts = append(opts, agollo.BackupFile(conf.BackupFile))
	} else if runtime.GOOS == "windows" {
		opts = append(opts, agollo.BackupFile("/nul"))
	} else {
		opts = append(opts, agollo.BackupFile("/dev/null"))
	}

	// 验证方式
	if conf.AccessKey != "" {
		opts = append(opts, agollo.AccessKey(conf.AccessKey))
	} else if conf.AuthBasicUser != "" {
		opts = append(opts,
			agollo.WithClientOptions(
				agollo.WithAccessKey("basic "+zutils.Crypto.Base64Encode(conf.AuthBasicUser+":"+conf.AuthBasicPassword)),
				agollo.WithSignatureFunc(func(ctx *agollo.SignatureContext) agollo.Header {
					return agollo.Header{"authorization": ctx.AccessKey}
				}),
			))
	}

	// 构建apollo客户端
	apolloClient, err := agollo.New(conf.Address, conf.AppId, opts...)
	if err != nil {
		return nil, fmt.Errorf("初始化agollo失败: %s", err)
	}

	// 加载数据
	data := make(map[string]interface{}, len(conf.Namespaces))
	for _, name := range conf.Namespaces {
		d := apolloClient.GetNameSpace(string(name))
		if len(d) == 0 {
			return nil, fmt.Errorf("命名空间[%s]的数据为空", name)
		}
		data[strings.ReplaceAll(string(name), "_", "")] = map[string]interface{}(d)
	}

	// 构建viper
	vi := viper.New()
	if err = vi.MergeConfigMap(data); err != nil {
		return nil, fmt.Errorf("合并配置失败: %s", err)
	}
	return vi, nil
}
