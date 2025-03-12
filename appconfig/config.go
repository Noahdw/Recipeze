package appconfig

import (
	"log/slog"
	"os"
)

const (
	EnvDev     = "dev"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

type AppConfig struct {
	URL       string
	FromEmail string
}

var Config AppConfig

func Initialize() {
	env := os.Getenv("APP_ENV")

	switch env {
	case EnvDev:
		Config.URL = "http://localhost:8080"
		Config.FromEmail = os.Getenv("DEV_EMAIL")
	case EnvProd:
		Config.URL = "TODO.com"
	default:
		slog.Error("Invalid env set for app config")
	}
}
