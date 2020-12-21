package cmd

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(createRealmCmd)
}

var createRealmCmd = &cobra.Command{
	Use:   "create-realm",
	Short: "Create realm",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()
		defer db.Close()

		realm := models.Realm{
			Name: args[0],
			Slug: args[1],
		}
		err := realm.Insert()
		if err != nil {
			zap.S().Error(err)
		}
	},
}
