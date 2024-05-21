package discovery

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
			if member.MissedHeartbeats >= 2 {
				member.NodeStatus = "Unreachable"
				fmt.Println("MissedCount is greater than or equal to 2")
			} else {
				member.MissedHeartbeats = member.MissedHeartbeats + 1
			}
		}
	}
}

func (prober *ProberService) ClusterHealthCheck(config *cluster.ClusterConfig) {

	for _, mem := range config.ClusterMemList {
		ctx, _ := context.WithTimeout(context.Background(), prober.TimeOut)
		url := "http://" + mem.NodeAddr + ":" + mem.NodeAddr + "/health"
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			prober.FailedChecks <- mem
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if resp.StatusCode != http.StatusOK {
			prober.FailedChecks <- mem
		}
		fmt.Printf("Node %s is healthy", mem.NodeAddr+":"+mem.NodePort)
		fmt.Println()
	}
}
