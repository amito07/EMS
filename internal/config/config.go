package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"ems_user"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"ems_password"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"ems_db"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"development"`
	storagePath string 
	ShutdownTimeout int `yaml:"shutdown_timeout" env-required:"true" env-default:"5"`
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "Path to the config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("No config file provided")
		}

	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	} 

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg
}
