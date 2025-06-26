package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

	// Валидация входных данных
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Конфигурация SSH
	config := &ssh.ClientConfig{
		User: req.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(req.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Подключение к серверу
	client, err := ssh.Dial("tcp", req.ServerIP+":22", config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "SSH connection failed",
			Details: err.Error(),
		})
		return
	}
	defer client.Close()

	// Выполнение команд
	commands := []string{
		"sudo pacman -Syu --noconfirm",
		"sudo pacman -S docker docker-compose git --noconfirm",
		"sudo systemctl start docker",
		"sudo systemctl enable docker",
		"git clone https://github.com/vgbhj/minecraftServerAutoDepoy.git /opt/webApp",
		"cd /opt/webApp/webApp && sudo go run main.go",
	}

	var output strings.Builder
	success := true

	for _, cmd := range commands {
		session, err := client.NewSession()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to create SSH session",
				Details: err.Error(),
			})
			return
		}

		var cmdOutput bytes.Buffer
		session.Stdout = &cmdOutput
		session.Stderr = &cmdOutput

		err = session.Run(cmd)
		session.Close()

		output.WriteString(fmt.Sprintf("$ %s\n%s\n", cmd, cmdOutput.String()))

		if err != nil {
			output.WriteString(fmt.Sprintf("Error: %v\n", err))
			success = false
			break
		}
	}

	if !success {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Deployment failed",
			Details: output.String(),
		})
		return
	}

	c.JSON(http.StatusOK, DeploymentResponse{
		Message: "Deployment completed successfully",
		Output:  output.String(),
	})
}

// Структуры для Swagger документации
type DeploymentRequest struct {
	ServerIP string `json:"server_ip" example:"192.168.1.100" binding:"required,ip"`
	Username string `json:"username" example:"admin" binding:"required"`
	Password string `json:"password" example:"securepassword" binding:"required"`
}

type DeploymentResponse struct {
	Message string `json:"message" example:"Deployment completed successfully"`
	Output  string `json:"output" example:"$ sudo pacman -Syu..."`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"SSH connection failed"`
	Details string `json:"details" example:"dial tcp 192.168.1.100:22: connect: connection refused"`
}
