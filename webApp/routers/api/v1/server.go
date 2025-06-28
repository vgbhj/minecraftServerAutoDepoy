package v1

import (
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/minecraft"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

type ServerInfo struct {
	Version string `json:"version"`
}

type ServerStatus struct {
	Status string `json:"status"` // "UP" или "DOWN"
}

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

// @Summary Get current server configuration
// @Description Returns currently selected server type and version
// @Tags minecraft
// @Produce json
// @Success 200 {object} ServerInfo
// @Router /api/v1/minecraft/current [get]
func GetCurrentVersion(c *gin.Context) {
	jarPath := setting.MinecraftSetting.ServerDir + "server.jar"
	if jarPath == "" {
		c.JSON(500, gin.H{
			"error": "MINECRAFT_JAR_PATH is not set in environment",
		})
		return
	}

	version, err := minecraft.GetMinecraftJarVersion(jarPath)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to get Minecraft server version",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"version": version,
	})
}

// @Summary Get Minecraft server status
// @Description Returns Minecraft server status based on Docker container state
// @Tags minecraft
// @Produce json
// @Success 200 {object} ServerStatus
// @Router /api/v1/server/status [get]
func GetServerStatus(c *gin.Context) {
	isUp, err := minecraft.IsDockerContainerRunning("minecraft-server")
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to get server status",
			"details": err.Error(),
		})
		return
	}

	status := "DOWN"
	if isUp {
		status = "UP"
	}

	c.JSON(200, ServerStatus{Status: status})
}

// @Summary Get host server IP
// @Description Returns the host server's external IP address
// @Tags minecraft
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/server/ip [get]
func GetServerIP(c *gin.Context) {
	cmd := exec.Command("sh", "-c", `ip route get 8.8.8.8 | grep -oP 'src \K[\d.]+'`)
	out, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to get server IP",
			"details": err.Error(),
		})
		return
	}
	ip := strings.TrimSpace(string(out))
	if ip == "" {
		c.JSON(500, gin.H{
			"error": "Could not determine server IP",
		})
		return
	}
	c.JSON(200, gin.H{
		"ip": ip,
	})
}
