package prober

import (
	"context"
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

func NewProberService(option ...Options) *ProberService {
	prober := &ProberService{
		TimeOut:      3 * time.Second,
		Ctx:          context.Background(),
		FailedChecks: make(chan FailedMemConfig),
	}
	for _, opti := range option {
		opti(prober)
	}
	return prober
}
