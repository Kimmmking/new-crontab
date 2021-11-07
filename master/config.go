package master

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	ApiPort         int      `mapstructure:"apiPort"`
	ApiReadTimeout  int      `mapstructure:"apiReadTimeout"`
	ApiWriteTimeout int      `mapstructure:"apiWriteTimeout"`
	EtcdEndpoints   []string `mapstructure:"etcdEndpoints"`
	EtcdDialTimeout int      `mapstructure:"etcdDialTimeout"`
	WebRoot         string   `mapstructure:"webroot"`
}

var (
	G_config *Config
)

func setDefaultValue() {
	viper.SetDefault("", "")
}

func InitConfig(filename string) (err error) {

	// 设置默认值
	setDefaultValue()

	// 1. 指定配置文件
	viper.SetConfigFile(filename)
	// 2. 读取配置信息
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// 3. 将读取的配置信息保存至全局变量Conf
	if err = viper.Unmarshal(G_config); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	// 4. 监控配置文件的变化
	viper.WatchConfig()
	// 5. 配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改")
		if err = viper.Unmarshal(G_config); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})

	return
}
