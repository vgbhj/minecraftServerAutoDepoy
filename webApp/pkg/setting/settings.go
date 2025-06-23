package setting

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type App struct {
	JwtSecret string `env:"JWT_SECRET,required"`
}

type Server struct {
	RunMode  string `env:"RUN_MODE" envDefault:"release"`
	HttpPort int    `env:"HTTP_PORT" envDefault:"8000"`
}

var (
	AppSetting    = &App{}
	ServerSetting = &Server{}
)

func Setup() {
	// 1️⃣ Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// 2️⃣ Парсим env-переменные в структуры
	if err := env.Parse(AppSetting); err != nil {
		log.Fatalf("Failed to parse AppSetting: %v", err)
	}
	if err := env.Parse(ServerSetting); err != nil {
		log.Fatalf("Failed to parse ServerSetting: %v", err)
	}

	log.Printf("Config: mode=%s port=%d\n", ServerSetting.RunMode, ServerSetting.HttpPort)
}
