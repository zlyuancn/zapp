/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/shima-park/agollo"
	"github.com/spf13/viper"
	"github.com/zlyuancn/zlog"
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type configCli struct {
	v     *viper.Viper
	c     *core.Config
	files []string
}

// 解析配置
//
// 多个配置文件如果存在同配置分片则只识别最后的分片
func NewConfig(appName string, defaultConfig *core.Config) core.IConfig {
	c := &configCli{
		v:     viper.New(),
		c:     defaultConfig,
		files: []string{},
	}

	confText := flag.String("c", "", "配置文件,多个文件用逗号隔开,同名配置分片会完全覆盖之前的分片")
	testFlag := flag.Bool("t", false, "测试配置文件")
	flag.Parse()

	if *confText == "" && defaultConfig == nil { // 如果命令行没有指定配置文件 且 没有主动设置配置
		*confText = consts.DefaultConfig_ConfigFiles
		fmt.Printf("未指定配置文件, 将使用 %s 配置文件\n", consts.DefaultConfig_ConfigFiles)
	}
	if *confText != "" { // 不管是命令行指定的还是由于没有主动设置配置选择默认配置文件
		files := strings.Split(*confText, ",")
		log := zlog.DefaultConfig
		c.c = newConfig()
		c.files = files
		for _, file := range files {
			vp := viper.New()
			vp.SetConfigFile(file)
			if err := vp.ReadInConfig(); err != nil {
				logger.Log.Fatal("配置文件加载失败", zap.String("file", file), zap.Error(err))
			}
			for k, v := range vp.AllSettings() {
				c.v.SetDefault(k, v)
			}
		}

		if err := c.v.Unmarshal(c.c); err != nil {
			logger.Log.Fatal("配置解析失败", zap.Strings("files", files), zap.Error(err))
		}
	}

	c.checkDefaultConfig(appName, c.c)

	if *testFlag {
		fmt.Println("配置文件测试成功")
		os.Exit(0)
	}

	return c
}

func newConfig() *core.Config {
	return &core.Config{
		Log: zlog.DefaultConfig,
	}
}

type ApolloConfig struct {
	Address              string // apollo-api地址, 多个地址用英文逗号连接
	AppId                string // 应用名
	AccessKey            string // 验证key, 优先级高于基础认证
	AuthBasicUser        string // 基础认证用户名
	AuthBasicPassword    string // 基础认证密码
	Cluster              string // 集群名
	AlwaysLoadFromRemote bool   // 总是从远程获取, 如果为false, 在远程加载失败时从备份文件加载
	BackupFile           string // 备份文件名
}

// 从apollo中获取配置结构
func GetConfigFromApollo(conf *ApolloConfig, nameSpaces ...Namespace) (*core.Config, error) {
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
	data := make(map[string]interface{}, len(nameSpaces))
	for _, name := range nameSpaces {
		d := apolloClient.GetNameSpace(string(name))
		if len(d) == 0 {
			return nil, fmt.Errorf("命名空间[%s]的数据为空", name)
		}
		data[strings.ReplaceAll(string(name), "_", "")] = d
	}

	// 反序列化
	v := viper.New()
	if err = v.MergeConfigMap(data); err != nil {
		return nil, fmt.Errorf("合并配置失败: %s", err)
	}

	out := newConfig()
	if err = v.Unmarshal(out); err != nil {
		return nil, fmt.Errorf("反序列化失败: %s", err)
	}

	return out, nil
}

func (c *configCli) checkDefaultConfig(appName string, conf *core.Config) {
	conf.Frame.FreeMemoryInterval = zutils.Ternary.Or(conf.Frame.FreeMemoryInterval, consts.FrameConfig_FreeMemoryInterval).(int)
	conf.Frame.WaitServiceRunTime = zutils.Ternary.Or(conf.Frame.WaitServiceRunTime, consts.FrameConfig_WaitServiceRunTime).(int)
	conf.Frame.ContinueWaitServiceRunTime = zutils.Ternary.Or(conf.Frame.ContinueWaitServiceRunTime, consts.FrameConfig_ContinueWaitServiceRunTime).(int)
	if conf.Log.Name == "" {
		conf.Log.Name = appName
	}
}

func (c *configCli) Config() *core.Config {
	return c.c
}

func (c *configCli) ConfigFiles() []string {
	return c.files
}

func (c *configCli) Parse(outPtr interface{}) error {
	return c.v.Unmarshal(outPtr)
}

func (c *configCli) ParseShard(shard string, outPtr interface{}) error {
	if !c.v.IsSet(shard) {
		return fmt.Errorf("分片<%s>不存在", shard)
	}
	return c.v.UnmarshalKey(shard, outPtr)
}

func (c *configCli) GetViper() *viper.Viper {
	return c.v
}
