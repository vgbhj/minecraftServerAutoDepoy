package main

import (
	"log"
	"net/http"

	"github.com/vgbhj/minecraftServerAutoDepoy/webApp/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/upload", handlers.Upload).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
