package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Renitens version",
	Long:  "Get Renitens version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("-- Renitens " + version)
	},
}
