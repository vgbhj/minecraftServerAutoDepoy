package routers

import (
	"net/http"
	"path/filepath"

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
		apiv1.POST("/server/deploy", v1.Deployment)
	}
	distPath := filepath.Join(".", "frontend", "dist")
	r.StaticFS("/assets", http.Dir(filepath.Join(distPath, "assets")))
	r.StaticFile("/", filepath.Join(distPath, "index.html"))
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(distPath, "index.html"))
	})
	return r
}
