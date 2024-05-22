package config

import (
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/events"
)

type ClusterConfigInf interface {
	AddClusterMemberList(member *cluster.ClusterMember) []*cluster.ClusterMember
	ListenForBroadcasts()
	JoinCluster(event events.ClusterEvent)
	LeaveCluster(event events.ClusterEvent)
}
