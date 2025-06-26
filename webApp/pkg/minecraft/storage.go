package minecraft

import "sync"

var (
	currentServerInfo ServerInfo
	mu                sync.RWMutex
)

type ServerInfo struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

// SetServerInfo сохраняет выбранную конфигурацию сервера
func SetServerInfo(serverType, version string) {
	mu.Lock()
	defer mu.Unlock()
	currentServerInfo = ServerInfo{
		Type:    serverType,
		Version: version,
	}
}

// GetServerInfo возвращает текущую конфигурацию сервера
func GetServerInfo() ServerInfo {
	mu.RLock()
	defer mu.RUnlock()
	return currentServerInfo
}
