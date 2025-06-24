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

type ServerInfo struct {
	Type    string
	Version string
}

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

func GetAvailableVersions() ([]string, map[string][]string, error) {
	// Здесь может быть логика получения версий из:
	// - Локального кэша
	// - Внешнего API (например, Mojang API)
	// - Базы данных
	// - Конфигурационного файла

	versions := []string{"1.20.1", "1.19.4", "1.18.2"}
	cores := map[string][]string{
		"Vanilla": {"Official"},
		"Paper":   {"Latest", "Stable"},
		"Forge":   {"Recommended"},
	}

	return versions, cores, nil
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

// HandleVersionSelection обрабатывает выбор версии
func HandleVersionSelection(c *gin.Context) {
	var req VersionSelectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, что выбранная версия существует
	validVersion := false
	for _, v := range minecraftData.Versions {
		if v == req.Version {
			validVersion = true
			break
		}
	}

	if !validVersion {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	// Проверка, что выбранное ядро существует
	coreOptions, exists := minecraftData.Cores[req.CoreType]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid core type"})
		return
	}

	// Проверка, что выбранная опция ядра существует
	validOption := false
	for _, opt := range coreOptions {
		if opt == req.CoreOption {
			validOption = true
			break
		}
	}

	if !validOption {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid core option"})
		return
	}

	// Здесь можно добавить логику обработки выбора (сохранение в БД и т.д.)

	c.JSON(http.StatusOK, gin.H{
		"message": "Selection successful",
		"selection": gin.H{
			"version":     req.Version,
			"core_type":   req.CoreType,
			"core_option": req.CoreOption,
		},
	})
}
