package minecraft

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vgbhj/minecraftServerAutoDepoy/pkg/setting"
)

type ValidationResult struct {
	Valid       bool
	Version     string
	Build       string
	Message     string
	DownloadURL string
}

var minecraftData = struct {
	Versions []string
	Cores    map[string][]string
}{
	Versions: []string{"1.20.1", "1.19.4", "1.18.2", "1.17.1", "1.16.5"},
	Cores: map[string][]string{
		"Vanilla": {"Official"},
		"Paper":   {"Latest", "Stable"},
		"Spigot":  {"Latest", "Recommended"},
		"Forge":   {"Latest", "Recommended"},
		"Fabric":  {"Latest", "Stable"},
	},
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
		"versions": minecraftData.Versions,
		"cores":    minecraftData.Cores,
	}
	c.JSON(http.StatusOK, response)
}

func validateVersion(serverType, version string) (ValidationResult, error) {
	var (
		url      string
		buildNum string
		finalURL string
		message  string
		valid    bool
	)

	switch strings.ToLower(serverType) {
	case "vanilla":
		if versionCheck(version, ">=", "1.0") && versionCheck(version, "<", "1.2") {
			url = fmt.Sprintf("http://files.betacraft.uk/server-archive/release/%s/%s.jar", version, version)
		} else {
			url = fmt.Sprintf("https://mcversions.net/download/%s", version)
		}

		if resp, err := http.Head(url); err == nil && resp.StatusCode == http.StatusOK {
			finalURL = url
			valid = true
		}

	case "paper":
		paperAPI := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", version)
		resp, err := http.Get(paperAPI)
		if err != nil {
			return ValidationResult{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var data struct {
				Builds []int `json:"builds"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				return ValidationResult{}, err
			}

			if len(data.Builds) > 0 {
				buildNum = fmt.Sprintf("%d", data.Builds[len(data.Builds)-1])
				finalURL = fmt.Sprintf("%s/builds/%s/downloads/paper-%s-%s.jar", paperAPI, buildNum, version, buildNum)
				valid = true
			}
		}

	case "forge":
		forgeURL := fmt.Sprintf("https://files.minecraftforge.net/net/minecraftforge/forge/index_%s.html", version)
		resp, err := http.Get(forgeURL)
		if err != nil {
			return ValidationResult{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return ValidationResult{}, err
			}

			re := regexp.MustCompile(`href="([^"]+\.jar)"`)
			matches := re.FindStringSubmatch(string(body))
			if len(matches) > 1 {
				finalURL = matches[1]
				valid = true
			}
		}

	default:
		return ValidationResult{}, errors.New("unsupported server type")
	}

	if !valid {
		message = fmt.Sprintf("Version '%s' not found or cannot be downloaded", version)
	}

	return ValidationResult{
		Valid:       valid,
		Version:     version,
		Build:       buildNum,
		Message:     message,
		DownloadURL: finalURL,
	}, nil
}

func versionCheck(version, comparator, targetVersion string) bool {
	// Я пока не понял че там нахуй написано в том auto-mc с проверкой версий
	// поэтому пока заглушку
	return true
}

// VersionSelectionRequest запрос на выбор версии
type VersionSelectionRequest struct {
	Version    string `json:"version" binding:"required"`
	CoreType   string `json:"core_type" binding:"required"`
	CoreOption string `json:"core_option" binding:"required"`
}

// @Summary Select Minecraft version and core
// @Description Saves selected Minecraft version and server type
// @Tags minecraft
// @Accept json
// @Produce json
// @Param input body VersionSelectionRequest true "Selection data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/minecraft/select [post]
func HandleVersionSelection(c *gin.Context) {
	var req VersionSelectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Selection saved successfully",
		"server": gin.H{
			"type":    req.CoreType,
			"version": req.Version,
		},
	})
}

func GetMinecraftJarVersion(jarPath string) (string, error) {
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "version.json" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			var data struct {
				Id string `json:"id"`
			}
			if err := json.NewDecoder(rc).Decode(&data); err != nil {
				return "", err
			}
			return data.Id, nil
		}
		if f.Name == "META-INF/MANIFEST.MF" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			re := regexp.MustCompile(`Implementation-Version: ([^\s]+)`)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				return matches[1], nil
			}
		}
	}

	return "", errors.New("version info not found in jar")
}

// @Summary      Download Minecraft server jar
// @Description  Downloads the selected Minecraft server jar and replaces /opt/minecraft-server/server.jar
// @Tags         minecraft
// @Accept       json
// @Produce      json
// @Param        serverType  query  string  true  "Server type (e.g. paper, vanilla, forge)"
// @Param        version     query  string  true  "Minecraft version"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/minecraft/download [post]
func DownloadAndReplaceServerJar(c *gin.Context) {
	serverType := c.Query("serverType")
	version := c.Query("version")

	if serverType == "" || version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "serverType and version are required"})
		return
	}

	validation, err := validateVersion(serverType, version)
	if err != nil || !validation.Valid || validation.DownloadURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find download URL for this version"})
		return
	}

	// Download the jar
	resp, err := http.Get(validation.DownloadURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download jar"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download jar, bad status"})
		return
	}

	targetDir := setting.MinecraftSetting.ServerDir
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target directory"})
		return
	}

	tmpFile := filepath.Join(targetDir, "server_new.jar")
	out, err := os.Create(tmpFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create jar file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write jar file"})
		return
	}

	// Replace old server.jar with new one
	finalJar := filepath.Join(targetDir, "server.jar")
	if err := os.Rename(tmpFile, finalJar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replace server.jar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Server jar downloaded and replaced successfully",
		"serverType": serverType,
		"version":    version,
		"jarPath":    finalJar,
	})
}
