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
	"strings"

	"github.com/spf13/viper"
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
func NewConfig(defaultConfig *core.Config) core.IConfig {
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
		c.c = new(core.Config)
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

	c.checkDefaultConfig(c.c)

	if *testFlag {
		fmt.Println("配置文件测试成功")
		os.Exit(0)
	}

	return c
}

func (c *configCli) checkDefaultConfig(conf *core.Config) {
	conf.Frame.FreeMemoryInterval = zutils.Ternary.Or(conf.Frame.FreeMemoryInterval, consts.FrameConfig_FreeMemoryInterval).(int)
	conf.Frame.WaitServiceRunTime = zutils.Ternary.Or(conf.Frame.WaitServiceRunTime, consts.FrameConfig_WaitServiceRunTime).(int)
	conf.Frame.ContinueWaitServiceRunTime = zutils.Ternary.Or(conf.Frame.ContinueWaitServiceRunTime, consts.FrameConfig_ContinueWaitServiceRunTime).(int)
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
