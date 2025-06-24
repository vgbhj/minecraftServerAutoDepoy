package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/minecraft"
)

// @Summary      Запуск сервера Minecraft
// @Description  Запускает сервер Minecraft через docker-compose
// @Tags         minecraft
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/server/start [post]
func StartServer(c *gin.Context) {
	if err := minecraft.StartDockerContainer(); err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to start Minecraft server",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Minecraft server started successfully",
	})
}

// @Summary      Остановка сервера Minecraft
// @Description  Останавливает сервер Minecraft через docker-compose
// @Tags         minecraft
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/server/stop [post]
func StopServer(c *gin.Context) {
	if err := minecraft.StopDockerContainer(); err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to stop Minecraft server",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Minecraft server stopped successfully",
	})
}

// @Summary      Перезагрузка сервера Minecraft
// @Description  Перезагружает сервер Minecraft через docker-compose
// @Tags         minecraft
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/server/restart [post]
func RestartServer(c *gin.Context) {
	if err := minecraft.RestartDockerContainer(); err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to restart Minecraft server",
			"details": err.Error(),
		})
		return
	}
}
