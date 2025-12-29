package tracker

import (
	"fmt"
	"math"
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

func findLeastUsedNodes(nodes []*co.NodeInfo, amount int, exclusions []string) ([]*co.NodeInfo, error) {
	nodesBuckets := map[int][]*co.NodeInfo{}
	for _, node := range nodes {
		nodesBuckets[len(node.Checkpoints)] = append(nodesBuckets[len(node.Checkpoints)], node)
	}

	sortedBucketsKeys := make([]int, 0, len(nodesBuckets))
	for k := range nodesBuckets {
		sortedBucketsKeys = append(sortedBucketsKeys, k)
	}
	sort.Slice(sortedBucketsKeys, func(i, j int) bool {
		return sortedBucketsKeys[i] < sortedBucketsKeys[j]
	})

	size := int(math.Min(float64(amount), float64(len(nodes))))
	filteredNodes := make([]*co.NodeInfo, size)
	c := 0
	for _, i := range sortedBucketsKeys {
		for _, n := range nodesBuckets[i] {
			if needToExclude(exclusions, n.Name) {
				continue
			}

			filteredNodes[c] = n
			c += 1
			if c >= size {
				return filteredNodes, nil
			}
		}
	}

	return nil, fmt.Errorf("cannot find enough nodes")
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

	filteredNodes, err := findLeastUsedNodes(allNodes, amount, exclusions)
	if err != nil {
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

}
