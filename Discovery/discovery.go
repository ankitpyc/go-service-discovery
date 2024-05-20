package Discovery

import (
	"go-service-discovery/cluster"
	"net"
	"time"
)

func (service *DiscoveryService) ClusterHealthCheck(config *cluster.ClusterConfig) {
	for _, mem := range config.ClusterMemList {
		dial, err := net.Dial("tcp", mem.NodeAddr+":"+mem.NodePort)
		if err != nil {
			return
		}
		err = dial.SetDeadline(time.Now().Add(100 * time.Millisecond))
		if err != nil {
			return
		}
	}
}
