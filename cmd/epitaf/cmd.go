package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	// Initialize ZAP globally
	var config zap.Config
	if os.Getenv("DEV") == "true" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		if os.Getenv("LOGS_PATH") != "" {
			config.OutputPaths = []string{os.Getenv("LOGS_PATH")}
		}
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)

	if godotenv.Load() != nil {
		zap.S().Warn(".env not found")
	}

	// Set timezone to UTC
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		zap.S().Info("Loaded " + tz + " timezone.")
		if err != nil {
			zap.S().Error(err)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "epitaf",
	Short: "Epitaf.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute cmd
func main() {
	if err := rootCmd.Execute(); err != nil {
		_ = rootCmd.Help()
		os.Exit(1)
	}
}
