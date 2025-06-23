package minecraft

import (
	"log"
	"os/exec"

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
