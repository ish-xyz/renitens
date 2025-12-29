package tracker

import (
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	co "github.com/ish-xyz/renitens/pkg/containerorchestrator"
	"github.com/ish-xyz/renitens/pkg/storage"
)

type TrackerService struct {
	store storage.StorageClient
	co    co.ContainerOrchestratorClient
}

func NewTrackerService(store storage.StorageClient, coClient co.ContainerOrchestratorClient) *TrackerService {
	return &TrackerService{
		store: store,
		co:    coClient,
	}
}

func (t *TrackerService) exposeCheckpoints(c *gin.Context) {
	dump := t.store.Dump()
	c.JSON(200, dump)
}

func (t *TrackerService) reconcileCheckpoints() error {

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

func findLeastUsedNodes(nodes []*co.NodeInfo, amount int) []*co.NodeInfo {
	nodesSort := map[int][]*co.NodeInfo{}
	for _, node := range nodes {
		nodesSort[len(node.Checkpoints)] = append(nodesSort[len(node.Checkpoints)], node)
	}

	keys := make([]int, 0, len(nodesSort))
	for k := range nodesSort {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	filteredNodes := make([]*co.NodeInfo, amount)
	for i := 0; i < amount; i++ {
		_ = nodesSort[keys[i]]
		for _, n := range nodesSort[keys[i]] {
			filteredNodes = append(filteredNodes, n)
			if len(filteredNodes) == amount {
				return filteredNodes
			}
		}
	}

	return filteredNodes
}

func (t *TrackerService) findHostsForNewCheckpoint(amount int, exclude []string) ([]*co.NodeInfo, error) {
	nodes, err := t.co.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes from CO: %v", err)
	}

	return findLeastUsedNodes(nodes, amount), nil
}

func (t *TrackerService) Run() {

	// Setup Gin router
	router := gin.Default()
	router.GET("/api/v1/checkpoints", t.exposeCheckpoints)
	router.Run("localhost:8080")

	// For testing purposes
	_ = t.reconcileCheckpoints()
	_, _ = t.findHostsForNewCheckpoint(1, []string{})
}
