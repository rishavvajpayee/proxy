package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProxyTargetUrl string
}

var AppConfig Config

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	AppConfig = Config{
		ProxyTargetUrl: os.Getenv("PROXYTARGETURL"),
	}

	return nil
}
