package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

const (
	minecraftManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

func downloadLatestVanillaServerJar(dest string) error {
	// 1. Получаем manifest с версиями
	resp, err := http.Get(minecraftManifestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type Version struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	type Manifest struct {
		Latest   map[string]string `json:"latest"`
		Versions []Version         `json:"versions"`
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return err
	}

	latestReleaseID := manifest.Latest["release"]
	var latestVersionURL string
	for _, v := range manifest.Versions {
		if v.ID == latestReleaseID {
			latestVersionURL = v.URL
			break
		}
	}
	if latestVersionURL == "" {
		return fmt.Errorf("latest release version not found in manifest")
	}

	// 2. Получаем ссылку на server.jar для этой версии
	resp, err = http.Get(latestVersionURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type Downloads struct {
		Server struct {
			URL string `json:"url"`
		} `json:"server"`
	}
	type VersionInfo struct {
		Downloads Downloads `json:"downloads"`
	}
	var versionInfo VersionInfo
	if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
		return err
	}
	serverJarURL := versionInfo.Downloads.Server.URL
	if serverJarURL == "" {
		return fmt.Errorf("server.jar url not found")
	}

	// 3. Скачиваем server.jar
	resp, err = http.Get(serverJarURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func copyFileIfNotExists(src, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	}
	return nil
}

func Setup() {
	dir := setting.MinecraftSetting.ServerDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		log.Printf("Created directory: %s", dir)
	}

	serverJarPath := filepath.Join(dir, "server.jar")
	if _, err := os.Stat(serverJarPath); os.IsNotExist(err) {
		log.Printf("server.jar not found, downloading latest vanilla Minecraft server...")
		if err := downloadLatestVanillaServerJar(serverJarPath); err != nil {
			log.Fatalf("Failed to download server.jar: %v", err)
		}
		log.Printf("Downloaded latest vanilla server.jar to %s", serverJarPath)
	}

	assetsDir := filepath.Join(".", "assets")
	dockerfileSrc := filepath.Join(assetsDir, "Dockerfile")
	dockerComposeSrc := filepath.Join(assetsDir, "docker-compose.yaml")
	dockerfileDst := filepath.Join(dir, "Dockerfile")
	dockerComposeDst := filepath.Join(dir, "docker-compose.yaml")

	if err := copyFileIfNotExists(dockerfileSrc, dockerfileDst); err != nil {
		log.Printf("Failed to copy Dockerfile: %v", err)
	} else {
		log.Printf("Dockerfile checked/copied to %s", dockerfileDst)
	}
	if err := copyFileIfNotExists(dockerComposeSrc, dockerComposeDst); err != nil {
		log.Printf("Failed to copy docker-compose.yaml: %v", err)
	} else {
		log.Printf("docker-compose.yaml checked/copied to %s", dockerComposeDst)
	}
}
