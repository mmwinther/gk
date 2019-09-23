package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GK",
	Long:  `GK's current version'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GK v0.0.6")
	},
}