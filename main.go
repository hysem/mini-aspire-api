package main

import (
	"log"
	"time"

	"github.com/hysem/mini-aspire-api/app"
	"github.com/hysem/mini-aspire-api/app/config"
	"github.com/hysem/mini-aspire-api/app/core/logger"
	"go.uber.org/zap"
)

func main() {
	time.Local = time.UTC

	if err := config.Load(); err != nil {
		log.Fatal("failed to load configuration", err)
	}

	logger.Init(config.Current().Log.Level, config.Current().App.Environment)

	zap.L().Debug("loaded configuration", zap.Any("config", config.Current()))

	app.New().Run()
}
