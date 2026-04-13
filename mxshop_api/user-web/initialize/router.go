package initialize

import (
	"mxshop_api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	engine := gin.Default()
	ApiGroup := engine.Group("/v1")
	router.InitUserRouter(ApiGroup)
	
	
	
	
	return engine
}
