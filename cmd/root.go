package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "renitens",
	Short: "renitens root command",
	Run:   rootRun,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(csiCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println("here I am")
	return
}
