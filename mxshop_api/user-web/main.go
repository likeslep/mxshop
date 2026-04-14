package main

import (
	"fmt"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化logger
	initialize.InitLogger()

	// 2. 初始化routers
	engine := initialize.Routers()

	// 3. 初始化配置文件
	initialize.InitConfig()

	port := global.ServerConfig.Port

	zap.S().Debugf("启动服务器，端口: %d", port)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败: ", err.Error())
	}
}
