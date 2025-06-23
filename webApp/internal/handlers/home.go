package handlers

import (
	"net/http"

	"github.com/vgbhj/minecraftServerAutoDepoy/webApp/internal/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Загрузка ядра Minecraft сервера",
	}
	templates.Render(w, "home.html", data)
}
