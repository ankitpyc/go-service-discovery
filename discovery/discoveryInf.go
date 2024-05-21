package discovery

import (
	"context"
	"go-service-discovery/cluster"
	"time"
)

type ProberServiceInf interface {
	ClusterHealthCheck(config *cluster.ClusterConfig)
}

type ProberService struct {
	TimeOut      time.Duration
	Ctx          context.Context
	FailedChecks chan *cluster.ClusterMember
}

func NewProberService(option ...Options) *ProberService {
	prober := &ProberService{
		TimeOut:      3 * time.Second,
		Ctx:          context.Background(),
		FailedChecks: make(chan *cluster.ClusterMember),
	}
	for _, opti := range option {
		opti(prober)
	}
	return prober
}
