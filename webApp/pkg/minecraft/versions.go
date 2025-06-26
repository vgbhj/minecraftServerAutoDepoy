package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
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

func validateVersion(serverInfo ServerInfo) (ValidationResult, error) {
	var (
		url      string
		buildNum string
		finalURL string
		message  string
		valid    bool
	)

	switch strings.ToLower(serverInfo.Type) {
	case "vanilla":
		// Для Vanilla (официальный сервер Mojang)
		if versionCheck(serverInfo.Version, ">=", "1.0") && versionCheck(serverInfo.Version, "<", "1.2") {
			url = fmt.Sprintf("http://files.betacraft.uk/server-archive/release/%s/%s.jar", serverInfo.Version, serverInfo.Version)
		} else {
			url = fmt.Sprintf("https://mcversions.net/download/%s", serverInfo.Version)
		}

		// Проверяем доступность URL
		if resp, err := http.Head(url); err == nil && resp.StatusCode == http.StatusOK {
			finalURL = url
			valid = true
		}

	case "paper":
		// Для PaperMC
		paperAPI := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", serverInfo.Version)
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
				finalURL = fmt.Sprintf("%s/builds/%s/downloads/paper-%s-%s.jar", paperAPI, buildNum, serverInfo.Version, buildNum)
				valid = true
			}
		}

	case "forge":
		// Для Forge
		forgeURL := fmt.Sprintf("https://files.minecraftforge.net/net/minecraftforge/forge/index_%s.html", serverInfo.Version)
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

			// Парсим HTML для поиска ссылки на JAR
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
		message = fmt.Sprintf("Version '%s' not found or cannot be downloaded", serverInfo.Version)
	}

	return ValidationResult{
		Valid:       valid,
		Version:     serverInfo.Version,
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

	// Сохраняем выбранную конфигурацию
	SetServerInfo(req.CoreType, req.Version)

	c.JSON(http.StatusOK, gin.H{
		"message": "Selection saved successfully",
		"server": gin.H{
			"type":    req.CoreType,
			"version": req.Version,
		},
	})
}
