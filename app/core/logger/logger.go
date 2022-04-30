package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const production = "prod"

// Init the logger. This will replace the global logger
// Use zap.L() to get the global logger
func Init(level string, environment string) {
	var logLevel zapcore.Level
	if err := logLevel.Set(level); err != nil {
		log.Fatalf("invalid log level: %v", level)
	}

	logConfig := zap.NewDevelopmentConfig()
	if environment == production {
		logConfig = zap.NewProductionConfig()
	}
	logConfig.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("failed to config logger:%v", err)
	}

	zap.ReplaceGlobals(logger)
}
