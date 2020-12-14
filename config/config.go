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
func NewConfig() core.IConfig {
	confText := flag.String("c", "", "配置文件,多个文件用逗号隔开,同名配置分片会完全覆盖之前的分片")
	testFlag := flag.Bool("t", false, "测试配置文件")
	flag.Parse()

	var files []string
	if *confText == "" {
		*confText = consts.DefaultConfig_ConfigFiles
		fmt.Printf("未指定配置文件, 将使用 %s 配置文件\n", consts.DefaultConfig_ConfigFiles)
	}
	files = strings.Split(*confText, ",")

	c := &configCli{
		v: viper.New(),
		c: &core.Config{
			Debug: false,
		},
		files: files,
	}

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

	c.c.Frame.FreeMemoryInterval = zutils.Ternary.Or(c.c.Frame.FreeMemoryInterval, consts.FrameConfig_FreeMemoryInterval).(int)
	c.c.Frame.WaitServiceRunTime = zutils.Ternary.Or(c.c.Frame.WaitServiceRunTime, consts.FrameConfig_WaitServiceRunTime).(int)
	c.c.Frame.ContinueWaitServiceRunTime = zutils.Ternary.Or(c.c.Frame.ContinueWaitServiceRunTime, consts.FrameConfig_ContinueWaitServiceRunTime).(int)

	if *testFlag {
		fmt.Println("配置文件测试成功")
		os.Exit(0)
	}

	if err := c.v.Unmarshal(c.c); err != nil {
		logger.Log.Fatal("配置解析失败", zap.Strings("files", files), zap.Error(err))
	}
	return c
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
