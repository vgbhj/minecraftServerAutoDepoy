package v1

import (
	"net/http"

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
			"error": "Failed to start Minecraft server",
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

// @Summary      Get available Minecraft versions and cores
// @Description  Returns list of available Minecraft versions and server cores
// @Tags         minecraft
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Returns versions and cores"
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/minecraft/versions [get]
func GetAvailableVersions(c *gin.Context) {
	response := map[string]interface{}{
		"versions": []string{"1.20.1", "1.19.4", "1.18.2"},
		"cores": map[string][]string{
			"Vanilla": {"Official"},
			"Paper":   {"Latest", "Stable"},
			"Forge":   {"Recommended"},
		},
	}
	c.JSON(http.StatusOK, response)
}
