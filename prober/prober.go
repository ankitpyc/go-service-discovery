package prober

import (
	"context"
	"fmt"
	"go-service-discovery/cluster"
	"net/http"
)

func (prober *ProberService) MonitorForFailedChecks() {
	// Correct way to compare atomic.Int32 with an int
	for {
		select {
		case member := <-prober.FailedChecks:
			fmt.Println("Node is unhealthy")
			if member.ClusterMember.MissedHeartbeats >= 2 {
				member.ClusterMember.NodeStatus = "Unreachable"
				member.ClusterConfig.BroadCastChannel <- cluster.ClusterEvent{ClusterEvent: cluster.EventTYPE(1), ClusterMember: *member.ClusterMember}
				fmt.Println("MissedCount is greater than or equal to 2. Nodes are removed from Cluster Checks")
			} else {
				member.ClusterMember.MissedHeartbeats = member.ClusterMember.MissedHeartbeats + 1
			}
		}
	}
}

func (prober *ProberService) ClusterHealthCheck(config *cluster.ClusterConfig) {

	for _, mem := range config.ClusterMemList {
		if mem.NodeStatus != "Healthy" {
		}
		ctx, _ := context.WithTimeout(context.Background(), prober.TimeOut)
		url := "http://" + mem.NodeAddr + ":" + mem.NodePort + "/health"
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			prober.FailedChecks <- FailedMemConfig{mem, config}
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if resp == nil || err != nil {
			prober.FailedChecks <- FailedMemConfig{mem, config}
		}
		if resp.StatusCode != http.StatusOK {
			prober.FailedChecks <- FailedMemConfig{mem, config}
		}

		defer func() {
			if r := recover(); r != nil {
				prober.FailedChecks <- FailedMemConfig{mem, config}
				fmt.Printf("Probe panicked for member %s: %v\n", mem.NodeAddr, r)
			}
		}()

		fmt.Printf("Node %s is healthy", mem.NodeAddr+":"+mem.NodePort)
		fmt.Println()
	}
}
