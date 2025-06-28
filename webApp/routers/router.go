package routers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vgbhj/minecraftServerAutoDepoy/docs"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/minecraft"
	v1 "github.com/vgbhj/minecraftServerAutoDepoy/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiv1 := r.Group("/api/v1")

	apiv1.Use()
	{
		apiv1.POST("/server/start", v1.StartServer)
		apiv1.POST("/server/stop", v1.StopServer)
		apiv1.POST("/server/restart", v1.RestartServer)
		apiv1.GET("/minecraft/versions", minecraft.GetAvailableVersions)
		apiv1.POST("/minecraft/select", minecraft.HandleVersionSelection)
		apiv1.GET("/minecraft/current", v1.GetCurrentVersion)
		apiv1.GET("/server/status", v1.GetServerStatus)
		apiv1.GET("/server/ip", v1.GetServerIP)
		apiv1.GET("/server/properties", v1.GetServerProperties)
		apiv1.PUT("/server/properties", v1.UpdateServerProperties)
		apiv1.GET("/console/stream", v1.ConsoleStream)
		apiv1.POST("/console/rcon", v1.SendRconCommand)
	}
	distPath := filepath.Join(".", "frontend", "dist")
	r.StaticFS("/assets", http.Dir(filepath.Join(distPath, "assets")))
	r.StaticFile("/", filepath.Join(distPath, "index.html"))
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(distPath, "index.html"))
	})

	return r
}
