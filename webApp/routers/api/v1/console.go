package v1

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
	"github.com/gorilla/websocket"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

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

type RconRequest struct {
	Command string `json:"command"`
}

func SendRconCommand(c *gin.Context) {
	var req RconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	addr := setting.MinecraftSetting.ServerIP + ":" + setting.MinecraftSetting.RconPort
	password := setting.MinecraftSetting.RconPassword

	conn, err := rcon.Dial(addr, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RCON", "details": err.Error()})
		return
	}
	defer conn.Close()

	resp, err := conn.Execute(req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute command", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": resp})
}
