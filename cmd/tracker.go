package cmd

import (
	"time"

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
	ts := time.Now().Unix() - int64(100)

	chk1 := &co.CheckpointInfo{
		Namespace: "namespace1",
		Timestamp: ts,
		Pod:       "example1",
		Container: "container1",
	}
	chk2 := &co.CheckpointInfo{
		Namespace: "namespace2",
		Timestamp: ts,
		Pod:       "example2",
		Container: "container2",
	}

	mockClient.SetNodes([]*co.Node{
		{Name: "node1", IP: "192.168.1.1", Checkpoints: []*co.CheckpointInfo{chk1}},
		{Name: "node2", IP: "192.168.1.2", Checkpoints: []*co.CheckpointInfo{chk1, chk2}},
		{Name: "node3", IP: "192.168.1.3", Checkpoints: []*co.CheckpointInfo{}},
		{Name: "node4", IP: "192.168.1.4", Checkpoints: []*co.CheckpointInfo{chk2}},
	})
	mockClient.SetPodLocation("node2")
	trk := tracker.NewTrackerService(store, mockClient)
	trk.Run()
	return nil
}
