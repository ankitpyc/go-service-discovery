package Discovery

import "go-service-discovery/cluster"

type DiscoveryInf interface {
	ClusterHealthCheck(config *cluster.ClusterConfig)
}

type DiscoveryService struct {
}
