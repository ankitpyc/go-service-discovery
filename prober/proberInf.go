package prober

import (
	"context"
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/config"
	"time"
)

type FailedMemConfig struct {
	ClusterMember *cluster.ClusterMember
	ClusterConfig *config.ClusterDetails
}

type ProberServiceInf interface {
	ClusterHealthCheck(config *config.ClusterDetails)
}

type ProberService struct {
	TimeOut      time.Duration
	Ctx          context.Context
	FailedChecks chan FailedMemConfig
}

func (prober *ProberService) handleTimeout(mem *cluster.ClusterMember, details *config.ClusterDetails) {
	fmt.Println("port ", mem.NodePort, " has timeout")
}

func NewProberService(option ...Options) *ProberService {
	prober := &ProberService{
		TimeOut:      3 * time.Second,
		Ctx:          context.Background(),
		FailedChecks: make(chan FailedMemConfig),
	}
	for _, opt := range option {
		opt(prober)
	}
	return prober
}
