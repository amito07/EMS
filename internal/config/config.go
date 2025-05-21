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

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"development"`
	storagePath string 
	ShutdownTimeout int `yaml:"shutdown_timeout" env-required:"true" env-default:"5"`
	HTTPServer `yaml:"http_server"`
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
