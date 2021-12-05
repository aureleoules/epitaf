package cmd

import (
	"fmt"

	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(apiKeyCmd)
}

var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Generates a new api key",
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		newAPIKey := ""
		for i := 0; i < 8; i++ {
			newAPIKey += shortid.MustGenerate()
		}
		newAPIKey = newAPIKey[:64]

		err := models.InsertAPIKey(newAPIKey)
		if err != nil {
			zap.S().Error(err)
			return
		}

		fmt.Printf("\033[32m%s\033[0m\n", newAPIKey)
	},
}
