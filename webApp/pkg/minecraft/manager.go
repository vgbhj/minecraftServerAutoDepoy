package minecraft

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

func StartDockerContainer() error {
	dir := setting.MinecraftSetting.ServerDir
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("docker-compose output: %s", string(output))
		return err
	}
	log.Printf("Minecraft server started with docker-compose in %s", dir)
	return nil
}

func StopDockerContainer() error {
	dir := setting.MinecraftSetting.ServerDir
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("docker-compose output: %s", string(output))
		return err
	}
	log.Printf("Minecraft server stopped with docker-compose in %s", dir)
	return nil
}

func RestartDockerContainer() error {
	err := StopDockerContainer()
	if err != nil {
		log.Printf("Failed to stop container during restart: %v", err)
		return fmt.Errorf("failed to stop container: %w", err)
	}

	err = StartDockerContainer()
	if err != nil {
		log.Printf("Failed to start container during restart: %v", err)
		return fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Minecraft server restarted successfully")
	return nil
}

// @Summary Get current server configuration
// @Description Returns currently selected server type and version
// @Tags minecraft
// @Produce json
// @Success 200 {object} ServerInfo
// @Router /api/v1/minecraft/current [get]
func GetCurrentConfig(c *gin.Context) {
	config := GetServerInfo()
	c.JSON(http.StatusOK, config)
}
