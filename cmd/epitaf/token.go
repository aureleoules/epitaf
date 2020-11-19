package main

import (
	"github.com/aureleoules/epitaf/api"
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(tokenCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "token of user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		user, err := models.GetUserByEmail(args[0])
		if err != nil {
			zap.S().Error(err)
			return
		}

		token, t, err := api.AuthMiddleware().TokenGenerator(user)
		if err != nil {
			zap.S().Error(err)
			return
		}
		zap.S().Info("Generated JWT: ", token)
		zap.S().Info("Expires at", t.String())
	},
}
