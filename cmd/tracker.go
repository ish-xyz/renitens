package cmd

import (
	co "github.com/ish-xyz/renitens/pkg/containerorchestrator"
	"github.com/ish-xyz/renitens/pkg/storage"
	"github.com/ish-xyz/renitens/pkg/tracker"
	"github.com/spf13/cobra"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Start renitens tracker service",
	Long:  "Start renitens tracker service",
	RunE:  trackerSvcRun,
}

// var csiNodeCmd = &cobra.Command{
// 	Use:   "node",
// 	Short: "Run Renitens node driver",
// 	Long:  "Run Renitens node driver",
// 	Run:   csiNodeRun,
// }

func init() {
	// csiCmd.AddCommand(csiNodeCmd)
	//trackerCmd.PersistentFlags().StringP("node-id", "n", "", "Pass node id for identification")
	_ = ""
}

func trackerSvcRun(cmd *cobra.Command, args []string) error {

	// just for testing, remove later
	store := storage.NewMemStorage()
	mockClient := co.NewMockCOClient()
	mockClient.SetNodes([]*co.Node{
		{Name: "node1", IP: "192.168.1.1", Checkpoints: []string{"x", "y", "z"}},
		{Name: "node2", IP: "192.168.1.2", Checkpoints: []string{"x"}},
		{Name: "node3", IP: "192.168.1.3", Checkpoints: []string{"x", "y", "z", "n", "u"}},
		{Name: "node4", IP: "192.168.1.4", Checkpoints: []string{}},
	})
	mockClient.SetPodLocation("node2")
	t := tracker.NewTrackerService(store, mockClient)
	t.Run()
	return nil
}
