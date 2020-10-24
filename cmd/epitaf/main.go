package main

import (
	"github.com/aureleoules/epitaf/api"
	"github.com/aureleoules/epitaf/db"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)
	godotenv.Load()
}

func main() {
	// db.Init()

	go db.Connect()
	api.Serve()
}
