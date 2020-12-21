package cmd

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(registerAdminCmd)
}

var registerAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "register admin user",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()
		defer db.Close()

		user := models.Admin{
			Name:     args[0],
			Login:    args[1],
			Email:    args[2],
			Password: args[3],
		}

		realm, err := models.GetRealmBySlug(args[4])
		if err != nil {
			zap.S().Error(err)
			return
		}
		user.RealmID = realm.UUID

		user.HashPassword()

		err = user.Insert()
		if err != nil {
			zap.S().Error(err)
			return
		}
	},
}
