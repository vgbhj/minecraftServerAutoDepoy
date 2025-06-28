package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	pkg "github.com/vgbhj/minecraftServerAutoDepoy/pkg/server"
	"golang.org/x/crypto/ssh"
)

// @Summary Deploy Minecraft server
// @Description Deploys the Minecraft server web application on the target server via SSH
// @Tags Deployment
// @Accept json
// @Produce json
// @Param input body DeploymentRequest true "SSH connection details"
// @Success 200 {object} DeploymentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/server/deploy [post]
func Deployment(c *gin.Context) {
	var req DeploymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	config := &ssh.ClientConfig{
		User: req.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(req.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", req.ServerIP+":22", config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SSH connection failed",
			Details: err.Error(),
		})
		return
	}
	defer client.Close()

	commands := []string{
		"curl -s -O https://raw.githubusercontent.com/vgbhj/minecraftServerAutoDepoy/refs/heads/main/deployApp/install.sh",
		"chmod +x install.sh",
		"sudo ./install.sh > /tmp/minecraft_deploy.log 2>&1",
		"cat /tmp/minecraft_deploy.log", // Изменено с tail на cat для получения полного лога
	}

	output, err := pkg.DeployCommands(client, commands)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Deployment failed",
			Details: fmt.Sprintf("%v\nLast 500 chars:\n%s", err, output),
		})
		return
	}

	// Parse the script output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var message, adminPanel string

	// Идем с конца лога, так как нужные нам строки обычно в конце
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.HasPrefix(line, "Admin panel") {
			adminPanel = line
		} else if strings.HasPrefix(line, "Done!") {
			message = line
		}

		// Если нашли обе строки, можно выйти из цикла
		if message != "" && adminPanel != "" {
			break
		}
	}
	// Fallback to generic response
	c.JSON(http.StatusOK, DeploymentResponse{
		Message: "Deployment completed successfully",
		Output:  output, // Возвращаем полный вывод, если не распознали формат
	})
}

type DeploymentRequest struct {
	ServerIP string `json:"server_ip" example:"194.87.76.29" binding:"required,ip"`
	Username string `json:"username" example:"root" binding:"required"`
	Password string `json:"password" example:"securepassword" binding:"required"`
}

type DeploymentResponse struct {
	Message string `json:"message" example:"Deployment completed successfully"`
	Output  string `json:"output" example:"Admin panel is available at: http://192.168.1.100:8080"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"SSH connection failed"`
	Details string `json:"details" example:"dial tcp 192.168.1.100:22: connect: connection refused"`
}
