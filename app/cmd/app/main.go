package main

import (
	"emails_extractor/app/internal/app"
	"emails_extractor/app/internal/config"
	logger "emails_extractor/app/pkg/logger/zap"
)

func main() {
	cfg := config.GetConfig()

	logger.GetLogger().Info("create app")
	a := app.NewApp(cfg)

	logger.GetLogger().Info("extracting emails from files...")
	if err := a.ExtractEmails(); err != nil {
		logger.GetLogger().Fatal(err.Error())
	}
	logger.GetLogger().Info("done")
}
