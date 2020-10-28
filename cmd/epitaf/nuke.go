package main

import (
	"github.com/aureleoules/epitaf/db"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nukeCmd)
}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "nuke epitaf",
	Run: func(cmd *cobra.Command, args []string) {
		db.Connect()
		db.Delete()
	},
}
