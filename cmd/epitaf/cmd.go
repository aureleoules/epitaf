package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	// Initialize ZAP globally
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		zap.S().Panic(err)
	}
	zap.ReplaceGlobals(logger)

	godotenv.Load()
}

var rootCmd = &cobra.Command{
	Use:   "epitaf",
	Short: "Epitaf.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute cmd
func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Help()
		os.Exit(1)
	}
}
