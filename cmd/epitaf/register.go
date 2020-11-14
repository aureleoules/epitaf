package main

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(registerAdminCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		user, err := models.PrepareUser(args[0])
		if err != nil {
			zap.S().Error(err)
			return
		}

		err = user.Insert()
		if err != nil {
			zap.S().Error(err)
			return
		}
	},
}

var registerAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "register admin user",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()

		user := models.User{
			Name:      args[0],
			Login:     args[1],
			Email:     args[2],
			Class:     "",
			Promotion: 0,
			Teacher:   true,
			Region:    "",
			Semester:  "",
		}

		err := user.Insert()
		if err != nil {
			zap.S().Error(err)
			return
		}
	},
}
