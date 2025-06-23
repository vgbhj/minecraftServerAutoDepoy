package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates = template.Must(template.ParseGlob(filepath.Join("internal", "templates", "*.html")))

func Render(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
