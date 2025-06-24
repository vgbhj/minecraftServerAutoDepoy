// тестовый файлик
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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
	// Здесь должна быть логика сравнения версий (например, "1.18.2" > "1.17.1")
	// Для простоты возвращаем true.
	return true
}

func main() {
	serv1 := ServerInfo{
		Type:    "paper",
		Version: "1.20.1",
	}
	fmt.Println(validateVersion(serv1))
}
