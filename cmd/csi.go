package cmd

import (
	"github.com/spf13/cobra"
)

var csiCmd = &cobra.Command{
	Use:   "csi",
	Short: "Start renitens csi driver",
	Long:  "Start renitens csi driver",
	RunE:  csiRun,
}

// var csiNodeCmd = &cobra.Command{
// 	Use:   "node",
// 	Short: "Run Renitens node driver",
// 	Long:  "Run Renitens node driver",
// 	Run:   csiNodeRun,
// }

func init() {
	// csiCmd.AddCommand(csiNodeCmd)
	csiCmd.PersistentFlags().StringP("node-id", "n", "", "Pass node id for identification")
}

func csiRun(cmd *cobra.Command, args []string) error {
	// instantiate node service
	// instantiate controller service
	// instantiate driver
	// run driver

	return nil
}

// func csiNodeRun(cmd *cobra.Command, args []string) {
// 	fmt.Println("run node here")
// }
