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
	Use:     "apikey",
	Example: "epitaf apikey <label>",
	Short:   "Generates a new api key",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		newAPIKey := models.APIKey{
			Token: "",
			Label: args[0],
		}

		for i := 0; i < 8; i++ {
			newAPIKey.Token += shortid.MustGenerate()
		}
		newAPIKey.Token = newAPIKey.Token[:64]

		err := newAPIKey.Insert()
		if err != nil {
			zap.S().Error(err)
			return
		}

		fmt.Printf("\033[32m%s => %s\033[0m\n", newAPIKey.Label, newAPIKey.Token)
	},
}
