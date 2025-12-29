package tracker

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ish-xyz/renitens/pkg/config"
	co "github.com/ish-xyz/renitens/pkg/containerorchestrator"
	"github.com/ish-xyz/renitens/pkg/storage"
)

type TrackerService struct {
	store storage.StorageClient
	co    co.ContainerOrchestratorClient
}

func (t *TrackerService) ExposeCheckpoints(c *gin.Context) {
	dump := t.store.Dump()
	c.JSON(200, dump)
}

func (t *TrackerService) ReconcileCheckpoints() error {

	nodes, err := t.co.GetNodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes from CO: %v", err)
	}
	for _, node := range nodes {
		// connect to agent on node.IP
		// get list of checkpoints on that node
		// for each checkpoint, update centralized index in t.store
		_ = node
	}
	return nil
}

func (t *TrackerService) Serve() {

	// just for testing, remove later
	mockClient := t.co.(*co.MockClient)
	mockClient.SetNodes([]*co.NodeInfo{
		{Name: "node1", IP: "192.168.1.68"},
	})
	t.co = mockClient

	// Setup Gin router
	router := gin.Default()
	router.GET("/api"+config.API_VERSION+"/checkpoints", t.ExposeCheckpoints)

	router.Run("localhost:8080")
}
