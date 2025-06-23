package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Ошибка загрузки файла", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("serverjar")
	if err != nil {
		http.Error(w, "Файл не выбран", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll(uploadDir, os.ModePerm)
	dstPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Ошибка копирования файла", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Ядро сервера %s успешно загружено и развернуто!", handler.Filename)
}
