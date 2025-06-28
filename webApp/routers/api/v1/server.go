package v1

import (
	"bufio"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// @Summary Get Minecraft server properties
// @Description Returns all fields from server.properties as JSON
// @Tags minecraft
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/server/properties [get]
func GetServerProperties(c *gin.Context) {
	propsPath := setting.MinecraftSetting.ServerDir + "server.properties"
	file, err := os.Open(propsPath)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to open server.properties",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		sepIdx := strings.Index(line, "=")
		if sepIdx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:sepIdx])
		val := strings.TrimSpace(line[sepIdx+1:])
		properties[key] = val
	}
	if err := scanner.Err(); err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to read server.properties",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, properties)
}

// @Summary Update Minecraft server properties
// @Description Updates fields in server.properties from JSON body
// @Tags minecraft
// @Accept json
// @Produce json
// @Param properties body map[string]string true "Properties to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/server/properties [put]
func UpdateServerProperties(c *gin.Context) {
	propsPath := setting.MinecraftSetting.ServerDir + "server.properties"

	var updates map[string]string
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid JSON body",
		})
		return
	}

	file, err := os.Open(propsPath)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to open server.properties",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' || strings.Index(line, "=") < 0 {
			lines = append(lines, line)
			continue
		}
		sepIdx := strings.Index(line, "=")
		key := strings.TrimSpace(line[:sepIdx])
		if val, ok := updates[key]; ok {
			lines = append(lines, key+"="+val)
			delete(updates, key)
		} else {
			lines = append(lines, line)
		}
	}
	for key, val := range updates {
		lines = append(lines, key+"="+val)
	}

	if err := os.WriteFile(propsPath, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		c.JSON(500, gin.H{
			"error":   "Failed to write server.properties",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "server.properties updated",
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConsoleStream(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	cmd := exec.Command("docker", "logs", "-f", "minecraft-server")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err := cmd.Start(); err != nil {
		return
	}
	defer cmd.Process.Kill()

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
