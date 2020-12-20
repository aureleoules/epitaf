package cmd

import (
	"github.com/spf13/cobra"
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
		panic("not implemented")
	},
}

var registerAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "register admin user",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		panic("not implemented")
	},
}
