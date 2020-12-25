/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2020/7/2
   Description :
-------------------------------------------------
*/

package config

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/zlyuancn/zutils"
	"go.uber.org/zap"

	"github.com/zlyuancn/zapp/consts"
	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/logger"
)

type configCli struct {
	vi *viper.Viper
	c  *core.Config
}

func newConfig() *core.Config {
	conf := &core.Config{
		Frame: core.FrameConfig{
			Debug: true,
		},
	}
	return conf
}

// 解析配置
//
// 配置来源优先级 命令行 > WithViper > WithConfig > WithFiles > WithApollo > 默认配置文件
// 注意: 多个配置文件如果存在同配置分片会智能合并, 完全相同的配置节点以最后的文件为准, 从apollo拉取的配置优先级最高
func NewConfig(appName string, opts ...Option) core.IConfig {
	opt := newOptions()
	for _, o := range opts {
		o(opt)
	}

	confText := flag.String("c", "", "配置文件,多个文件用逗号隔开,同名配置分片会完全覆盖之前的分片")
	testFlag := flag.Bool("t", false, "测试配置文件")
	flag.Parse()

	var vi *viper.Viper
	var err error
	if *confText != "" { // 命令行
		files := strings.Split(*confText, ",")
		vi, err = makeViperFromFile(files)
		if err != nil {
			logger.Log.Fatal("从命令指定文件构建viper失败", zap.Strings("files", files), zap.Error(err))
		}
	} else if opt.vi != nil { // WithViper
		vi = opt.vi
	} else if opt.conf != nil { // WithConfig
		vi, err = makeViperFromStruct(opt.conf)
		if err != nil {
			logger.Log.Fatal("从配置结构构建viper失败", zap.Any("config", opt.conf), zap.Error(err))
		}
	} else if len(opt.files) > 0 { // WithFiles
		vi, err = makeViperFromFile(opt.files)
		if err != nil {
			logger.Log.Fatal("从用户指定文件构建viper失败", zap.Strings("files", opt.files), zap.Error(err))
		}
	} else if opt.apolloConfig != nil { // WithApollo
		vi, err = makeViperFromApollo(opt.apolloConfig)
		if err != nil {
			logger.Log.Fatal("从apollo构建viper失败", zap.Any("apolloConfig", opt.apolloConfig), zap.Error(err))
		}
	} else { // 默认
		files := strings.Split(consts.DefaultConfig_ConfigFiles, ",")
		logger.Log.Debug("使用默认配置文件", zap.Strings("files", files))
		vi, err = makeViperFromFile(files)
		if err != nil {
			logger.Log.Fatal("从默认配置文件构建viper失败", zap.Strings("files", files), zap.Error(err))
		}
	}

	// 如果从viper中发现了apollo配置
	if vi.IsSet(consts.ConfigGroupName_Apollo) {
		apolloConf, err := makeApolloConfigFromViper(vi)
		if err != nil {
			logger.Log.Fatal("解析apollo配置失败", zap.Error(err))
		}
		newVi, err := makeViperFromApollo(apolloConf)
		if err != nil {
			logger.Log.Fatal("从apollo构建viper失败", zap.Any("apolloConfig", apolloConf), zap.Error(err))
		}
		if err = vi.MergeConfigMap(newVi.AllSettings()); err != nil {
			logger.Log.Fatal("合并apollo配置失败", zap.Error(err))
		}
	}

	c := &configCli{
		vi: vi,
		c:  newConfig(),
	}
	if err = vi.Unmarshal(c.c); err != nil {
		logger.Log.Fatal("配置解析失败", zap.Error(err))
	}

	c.checkDefaultConfig(appName, c.c)

	if *testFlag {
		fmt.Println("配置文件测试成功")
		os.Exit(0)
	}

	return c
}

// 从文件构建viper
func makeViperFromFile(files []string) (*viper.Viper, error) {
	vi := viper.New()
	for _, file := range files {
		vp := viper.New()
		vp.SetConfigFile(file)
		if err := vp.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("配置文件'%s'加载失败: %s", file, err)
		}
		if err := vi.MergeConfigMap(vp.AllSettings()); err != nil {
			return nil, fmt.Errorf("合并配置文件'%s'失败: %s", file, err)
		}
	}
	return vi, nil
}

// 从结构体构建viper
func makeViperFromStruct(a interface{}) (*viper.Viper, error) {
	bs, err := jsoniter.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("编码失败: %s", err)
	}

	vi := viper.New()
	vi.SetConfigType("json")
	err = vi.ReadConfig(bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("数据解析失败: %s", err)
	}
	return vi, nil
}

func (c *configCli) checkDefaultConfig(appName string, conf *core.Config) {
	conf.Frame.FreeMemoryInterval = zutils.Ternary.Or(conf.Frame.FreeMemoryInterval, consts.FrameConfig_FreeMemoryInterval).(int)
	conf.Frame.WaitServiceRunTime = zutils.Ternary.Or(conf.Frame.WaitServiceRunTime, consts.FrameConfig_WaitServiceRunTime).(int)
	conf.Frame.ContinueWaitServiceRunTime = zutils.Ternary.Or(conf.Frame.ContinueWaitServiceRunTime, consts.FrameConfig_ContinueWaitServiceRunTime).(int)
}

func (c *configCli) Config() *core.Config {
	return c.c
}

func (c *configCli) Parse(outPtr interface{}) error {
	return c.vi.Unmarshal(outPtr)
}

func (c *configCli) ParseShard(shard string, outPtr interface{}) error {
	if !c.vi.IsSet(shard) {
		return fmt.Errorf("分片<%s>不存在", shard)
	}
	return c.vi.UnmarshalKey(shard, outPtr)
}

func (c *configCli) GetViper() *viper.Viper {
	return c.vi
}
