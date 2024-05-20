package Discovery

import (
	"fmt"
	"go-service-discovery/cluster"
	"net/http"
)

func (service *DiscoveryService) ClusterHealthCheck(config *cluster.ClusterConfig) {
	for _, mem := range config.ClusterMemList {
		resp, err := http.Get("http://" + mem.NodeAddr + ":" + mem.NodePort + "/health")
		if err != nil {
			fmt.Println(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Node %s is unhealthy", mem.NodeAddr+":"+mem.NodePort)
			fmt.Println()
			return
		}
		fmt.Printf("Node %s is healthy", mem.NodeAddr+":"+mem.NodePort)
		fmt.Println()
	}
}
