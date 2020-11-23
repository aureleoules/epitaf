package cmd

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize epitaf",
	Run: func(cmd *cobra.Command, args []string) {
		db.Init()

		err := models.InjectSQLSchemas()
		if err != nil {
			zap.S().Fatal(err)
		}
	},
}
