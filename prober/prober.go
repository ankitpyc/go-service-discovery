package prober

import (
	"context"
	"errors"
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/config"
	"go-service-discovery/cluster/events"
	"net/http"
)

func (prober *ProberService) MonitorForFailedChecks(ctx context.Context) {
	// Correct way to compare atomic.Int32 with an int
	for {
		select {
		case member := <-prober.FailedChecks:
			fmt.Println("Node %s is unhealthy", member.ClusterMember.NodeAddr+":"+member.ClusterMember.NodePort)
			if member.ClusterMember.MissedHeartbeats >= 2 {
				member.ClusterMember.NodeStatus = "Unreachable"
				member.ClusterConfig.BroadCastChannel <- events.ClusterEvent{ClusterEvent: events.EventTYPE(1), ClusterMember: *member.ClusterMember}
				fmt.Println("MissedCount is greater than or equal to 2. Nodes are removed from Cluster Checks")
			} else {
				member.ClusterMember.MissedHeartbeats = member.ClusterMember.MissedHeartbeats + 1
			}
		}
	}
}

func (prober *ProberService) ClusterHealthCheck(ctx context.Context, config *config.ClusterDetails) {
	config.RLock()
	defer config.RUnlock()
	for _, mem := range config.ClusterMemList {
		cont, cancel := context.WithTimeout(ctx, prober.TimeOut)
		defer cancel()
		url := "http://" + mem.NodeAddr + ":" + mem.NodePort + "/Health"
		req, err := http.NewRequestWithContext(cont, "GET", url, nil)
		if err != nil {
			prober.FailedChecks <- FailedMemConfig{mem, config}
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		probeFailed := prober.handleProbeResponse(resp, err, mem, config)
		if probeFailed {
			continue
		}
	}
}

func (prober *ProberService) handleProbeResponse(resp *http.Response, err error, mem *cluster.ClusterMember, config *config.ClusterDetails) bool {
	var probeFailed bool = false
	if errors.Is(err, context.DeadlineExceeded) {
		prober.handleTimeout(mem, config)
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
