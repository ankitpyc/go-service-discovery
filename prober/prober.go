package prober

import (
	"context"
	"errors"
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/config"
	"go-service-discovery/cluster/events"
	"net/http"
	"sync/atomic"
)

func (prober *ProberService) MonitorForFailedChecks(ctx context.Context) {
	for {
		select {
		case member := <-prober.FailedChecks:
			fmt.Printf("Node %s is unhealthy\n", member.ClusterMember.NodeAddr+":"+member.ClusterMember.NodePort)
			if atomic.LoadInt32(&member.ClusterMember.MissedHeartbeats) >= 2 {
				member.ClusterMember.NodeStatus = "Unreachable"
				member.ClusterConfig.BroadCastChannel <- events.ClusterEvent{ClusterEvent: events.EventTYPE(1), ClusterMember: *member.ClusterMember}
				fmt.Println("MissedCount is greater than or equal to 2. Nodes are removed from Cluster Checks")
			} else {
				atomic.AddInt32(&member.ClusterMember.MissedHeartbeats, 1)
			}
		case <-ctx.Done():
			fmt.Println("MonitorForFailedChecks stopped:", ctx.Err())
			return
		}
	}
}

func (prober *ProberService) ClusterHealthCheck(ctx context.Context, config *config.ClusterDetails) {
	config.RLock()
	defer config.RUnlock()
	for _, mem := range config.ClusterMemList {
		checkCtx, cancel := context.WithTimeout(ctx, prober.TimeOut)
		defer cancel()

		url := "http://" + mem.NodeAddr + ":" + mem.NodePort + "/Health"
		req, err := http.NewRequestWithContext(checkCtx, "GET", url, nil)
		if err != nil {
			prober.FailedChecks <- FailedMemConfig{mem, config}
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		probeFailed := prober.handleProbeResponse(resp, err, mem, config)
		if probeFailed {
			continue
		}
		fmt.Printf("Node %s:%s is healthy\n", mem.NodeAddr, mem.NodePort)
	}
}

func (prober *ProberService) handleProbeResponse(resp *http.Response, err error, mem *cluster.ClusterMember, config *config.ClusterDetails) bool {
	var probeFailed bool = false

	if errors.Is(err, context.DeadlineExceeded) {
		prober.handleTimeout(mem, config)
		probeFailed = true
		return probeFailed
	}

	if resp == nil || err != nil {
		prober.FailedChecks <- FailedMemConfig{mem, config}
		probeFailed = true
		return probeFailed
	}

	if resp.StatusCode != http.StatusOK {
		prober.FailedChecks <- FailedMemConfig{mem, config}
		probeFailed = true
	}
	return probeFailed
}

func (prober *ProberService) handleTimeout(mem *cluster.ClusterMember, config *config.ClusterDetails) {
	fmt.Printf("Health check for node %s:%s timed out\n", mem.NodeAddr, mem.NodePort)
	prober.FailedChecks <- FailedMemConfig{mem, config}
}
