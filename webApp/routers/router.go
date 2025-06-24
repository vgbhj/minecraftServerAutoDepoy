package routers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vgbhj/minecraftServerAutoDepoy/docs"
	v1 "github.com/vgbhj/minecraftServerAutoDepoy/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiv1 := r.Group("/api/v1")

	apiv1.Use()
	{
		apiv1.POST("/server/start", v1.StartServer)
		apiv1.POST("/server/stop", v1.StopServer)
		apiv1.POST("/server/restart", v1.RestartServer)
		apiv1.GET("/minecraft/versions", v1.GetAvailableVersions)
	}
	return r
}
