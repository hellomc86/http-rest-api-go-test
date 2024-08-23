package main

import (
	"log/slog"
	"os"

	"http-rest-api-go/internal/app/apiserver"
		
	"http-rest-api-go/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Todo: init config
	config := config.MustLoad()

	log := setupLogger(config.Env)	

	if err := apiserver.Start(config, log); err != nil {
		log.Error(err.Error())
	}

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
