/*
-------------------------------------------------
   Author :       Zhang Fan
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

	"github.com/zlyuancn/zapp/core"
	"github.com/zlyuancn/zapp/utils"
)

const (
	DefaultConfigFile         = "./configs/default.toml,./configs/local.toml"
	DefaultFreeMemoryInterval = 120000 // 默认清理内存间隔时间(毫秒)
)

type configCli struct {
	v     *viper.Viper
	c     *core.Config
	files []string
}

func NewConfig() core.IConfig {
	confText := flag.String("c", "", "配置文件,多个文件用逗号隔开,同名配置分片会完全覆盖之前的分片")
	testConfig := flag.Bool("t", false, "测试配置文件")
	flag.Parse()

	var files []string
	if *confText == "" {
		*confText = DefaultConfigFile
		fmt.Printf("未指定配置文件, 将使用 %s 配置文件\n", DefaultConfigFile)
	}
	files = strings.Split(*confText, ",")

	c := &configCli{
		v: viper.New(),
		c: &core.Config{
			Debug:              false,
			FreeMemoryInterval: DefaultFreeMemoryInterval,
		},
		files: files,
	}

	for _, file := range files {
		vp := viper.New()
		vp.SetConfigFile(file)
		utils.FailOnErrorf(vp.ReadInConfig(), "配置文件<%s>加载失败", file)
		for k, v := range vp.AllSettings() {
			c.v.SetDefault(k, v)
		}
	}

	if *testConfig {
		fmt.Println("配置文件测试成功")
		os.Exit(0)
	}

	utils.FailOnError(c.v.Unmarshal(c.c), "配置解析失败")
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
