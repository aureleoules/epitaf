package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tokenCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "token of user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		panic("not implemented")
	},
}
