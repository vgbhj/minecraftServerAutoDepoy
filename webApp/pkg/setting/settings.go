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

type Minecraft struct {
	ServerDir    string `env:"MINECRAFT_SERVER_DIR,required"`
	ServerIP     string `env:"MINECRAFT_SERVER_IP,required"`
	RconPort     string `env:"MINECRAFT_RCON_PORT,required"`
	RconPassword string `env:"MINECRAFT_RCON_PASSWORD,required"`
}

var (
	AppSetting       = &App{}
	ServerSetting    = &Server{}
	MinecraftSetting = &Minecraft{}
)

func Setup() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	if err := env.Parse(AppSetting); err != nil {
		log.Fatalf("Failed to parse AppSetting: %v", err)
	}
	if err := env.Parse(ServerSetting); err != nil {
		log.Fatalf("Failed to parse ServerSetting: %v", err)
	}
	if err := env.Parse(MinecraftSetting); err != nil {
		log.Fatalf("Failed to parse MinecraftSetting: %v", err)
	}

	log.Printf("Config: mode=%s port=%d minecraft_dir=%s\n", ServerSetting.RunMode, ServerSetting.HttpPort, MinecraftSetting.ServerDir)
}
