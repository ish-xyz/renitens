package tracker

import (
	"fmt"
	"sort"
	"strconv"

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

func needToExclude(exclusions []string, nodeName string) bool {
	exclude := false
	for _, ex := range exclusions {
		if ex == nodeName {
			exclude = true
			break
		}
	}

	return exclude
}

func findLeastUsedNodes(nodes []*co.NodeInfo, amount int, exclusions []string) []*co.NodeInfo {
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
			if needToExclude(exclusions, n.Name) {
				continue
			}

			filteredNodes = append(filteredNodes, n)
			if len(filteredNodes) == amount {
				return filteredNodes
			}
		}
	}

	return filteredNodes
}

func (t *TrackerService) serveCheckpoints(c *gin.Context) {
	dump := t.store.Dump()
	c.JSON(200, dump)
}

func (t *TrackerService) serveFindHostsForNewCheckpoint(c *gin.Context) {
	// get amount from parameters
	amount, err := strconv.Atoi(c.Query("amount"))
	if err != nil {
		c.JSON(400, map[string]string{"error": "failed to conver amount into int"})
		return
	}
	exclusions := c.QueryArray("exclusions")

	allNodes, err := t.co.GetNodes()
	if err != nil {
		c.JSON(500, map[string]string{"error": "failed to retrieve nodes"})
		return
	}

	filteredNodes := findLeastUsedNodes(allNodes, amount, exclusions)
	if len(filteredNodes) < amount {
		c.JSON(404, map[string]string{"error": "not enough nodes available"})
		return
	}

	c.JSON(200, filteredNodes)
}

func (t *TrackerService) Run() {

	// Setup Gin router
	router := gin.Default()
	router.GET("/api/v1/checkpoints", t.serveCheckpoints)
	router.GET("/api/v1/find", t.serveFindHostsForNewCheckpoint)
	router.Run("localhost:8080")

	// For testing purposes
	_ = t.reconcileCheckpoints()
	_, _ = t.findHostsForNewCheckpoint(1, []string{})
}
