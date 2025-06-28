package minecraft

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

const podmanBin = "/usr/local/bin/podman"
const podmanComposeBin = "/usr/local/bin/podman-compose"

func StartDockerContainer() error {
	dir := setting.MinecraftSetting.ServerDir
	cmd := exec.Command(podmanComposeBin, "up", "-d")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PATH=/usr/local/bin:"+os.Getenv("PATH"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("podman-compose output: %s", string(output))
		return err
	}
	log.Printf("Minecraft server started with podman-compose in %s", dir)
	return nil
}

func StopDockerContainer() error {
	dir := setting.MinecraftSetting.ServerDir
	cmd := exec.Command(podmanComposeBin, "down")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PATH=/usr/local/bin:"+os.Getenv("PATH"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("podman-compose output: %s", string(output))
		return err
	}
	log.Printf("Minecraft server stopped with podman-compose in %s", dir)
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

func IsDockerContainerRunning(containerName string) (bool, error) {
	cmd := exec.Command(podmanBin, "ps", "--filter", "name="+containerName, "--filter", "status=running", "--format", "{{.Names}}")
	cmd.Env = append(os.Environ(), "PATH=/usr/local/bin:"+os.Getenv("PATH"))
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), containerName), nil
}
