package cmd

import (
	"fmt"

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
	Use:  "login",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		user, err := models.GetUserByEmail(args[0])
		if err != nil {
			zap.S().Error(err)
			return
		}

		token, _, err := api.AuthMiddleware().TokenGenerator(user)
		if err != nil {
			zap.S().Error(err)
			return
		}

		fmt.Printf("\033[32mhttp://localhost:3000/login?token=%s\033[0m\n", token)
	},
}
