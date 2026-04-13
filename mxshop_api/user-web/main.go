package main

import (
	"fmt"
	"mxshop_api/user-web/initialize"
	"go.uber.org/zap"
)

func main() {
	port := 8021

	// 1. 初始化logger
	initialize.InitLogger()

	// 2. 初始化routers
	engine := initialize.Routers()



	zap.S().Debugf("启动服务器，端口: %d", port)
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败: ", err.Error())
	}
}
