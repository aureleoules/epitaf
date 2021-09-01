package cmd

import (
	"fmt"

	"github.com/aureleoules/epitaf/db"
	"github.com/aureleoules/epitaf/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(updateClassesCmd)
}

var updateClassesCmd = &cobra.Command{
	Use:   "update-classes",
	Short: "update classes of users",
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()
		users, err := models.GetUsers()
		if err != nil {
			zap.S().Panic(err)
		}
		fmt.Println(len(users))
		for _, u := range users {
			user, err := models.PrepareUser(u.Email)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			err = models.UpdateUser(&models.UpdateUserReq{
				Login:     u.Login,
				Promotion: int(user.Promotion.Int64Value()),
				Class:     user.Class.String(),
				Region:    user.Region.String(),
				Semester:  user.Semester.String(),
			})

			if err != nil {
				zap.S().Error(err)
				continue
			}

			fmt.Println(user.Class)
		}
	},
}
