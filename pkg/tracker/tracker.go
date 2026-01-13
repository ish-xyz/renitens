package tracker

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

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

		// get a list of pods and checkpoints that should be there
		_ = node
	}

	// get all pods
	// check if checkpoints exists already
	// if it doesn't schedule
	// add to desired state

	// nodex -> checkpoints
	return nil
}

func findLeastUsedNodes(nodes []*co.Node, amount int, exclusions []string) ([]string, error) {
	nodesBuckets := map[int][]*co.Node{}
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
	filteredNodes := make([]string, size)
	c := 0
	for _, i := range sortedBucketsKeys {
		for _, n := range nodesBuckets[i] {
			if contains(exclusions, n.Name) {
				continue
			}

			filteredNodes[c] = n.Name
			c += 1
			if c >= size {
				return filteredNodes, nil
			}
		}
	}

	return nil, fmt.Errorf("cannot find enough nodes")
}

func (t *TrackerService) serveGetCheckpoints(c *gin.Context) {
	dump := t.store.Dump()
	c.JSON(200, dump)
}

func (t *TrackerService) serveScheduleCheckpoint(c *gin.Context) {

	// get params
	namespace := c.Param("namespace")
	pod := c.Param("pod")
	minDelta, err := strconv.Atoi(c.Param("minDelta"))
	if err != nil {
		c.JSON(500, "cannot convert minDelta parameter into an integer")
		return
	}
	replicas, err := strconv.Atoi(c.Param("replicas"))
	if err != nil {
		c.JSON(500, "failed to convert 'replicas' parameter into an integer")
		return
	}

	// schedule checkpoints
	excludeNode, err := t.co.GetPodLocation(namespace, pod)
	if err != nil {
		c.JSON(500, "failed to determine pod location")
		return
	}

	nodes, err := t.co.GetNodes()
	if err != nil {
		c.JSON(500, "failed to list nodes")
		return
	}

	var existingNodes []string
	now := time.Now().Unix()
	for _, node := range nodes {
		for _, chk := range node.Checkpoints {
			if chk.Namespace != namespace || chk.Pod != pod {
				continue
			}
			if (now - chk.Timestamp) < int64(minDelta) {
				existingNodes = append(existingNodes, node.Name)
				break
			}
		}
	}

	newReplicas := replicas - len(existingNodes)
	if newReplicas > 0 {
		candidateNodes, err := findLeastUsedNodes(nodes, newReplicas, []string{excludeNode})
		if err != nil {
			c.JSON(500, "failed to list nodes")
			return
		}
		existingNodes = append(existingNodes, candidateNodes...)
	}

	c.JSON(200, fmtData(existingNodes))
}

func (t *TrackerService) Run() {

	// Setup Gin router
	router := gin.Default()
	router.GET("/api/v1/checkpoints", t.serveGetCheckpoints)
	router.GET("/api/v1/schedule/:namespace/:pod/:replicas/:minDelta", t.serveScheduleCheckpoint)
	router.Run("localhost:8080")

	// For testing purposes
	_ = t.reconcileCheckpoints()

}
