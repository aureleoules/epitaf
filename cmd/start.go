package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aureleoules/epitaf/api"
	"github.com/aureleoules/epitaf/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

func handleInterrupt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		zap.S().Info("Caught Ctrl+C signal")

		db.Close()
		os.Exit(0)
	}()
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Epitaf",
	Run: func(cmd *cobra.Command, args []string) {
		handleInterrupt()

		zap.S().Info("Starting epitaf...")

		db.Connect()
		defer db.Close()

		api.Serve()
	},
}
