package initialize

import (
	"fmt"
	"mxshop_api/user-web/global"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)


func InitConfig() {

	v := viper.New()
	v.SetConfigFile("config-debug.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	fmt.Println(global.ServerConfig)
	zap.S().Info("配置信息：%v", global.ServerConfig)
	fmt.Printf("%v", v.Get("name"))

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("配置文件发生变化：%v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		fmt.Println(global.ServerConfig)
	})
}
