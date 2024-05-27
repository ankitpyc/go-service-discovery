package config

import (
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/broadcast"
	"go-service-discovery/cluster/events"
	"sync"
)

type ClusterDetails struct {
	sync.RWMutex
	ClusterMemList   []*cluster.ClusterMember
	ClusterName      string
	ClusterID        string
	TotalSize        int
	BroadCastChannel chan events.ClusterEvent
}

func (cluster *ClusterDetails) AddClusterMemberList(member *cluster.ClusterMember) []*cluster.ClusterMember {
	cluster.Lock()
	defer cluster.Unlock()
	member.NodeStatus = "HEALTHY"
	fmt.Println("Discovered Node -> ", member.NodeAddr+":"+member.NodePort)
	cluster.ClusterMemList = append(cluster.ClusterMemList, member)
	return cluster.ClusterMemList
}

func (cluster *ClusterDetails) ListenForBroadcasts() {
	for {
		select {
		case event := <-cluster.BroadCastChannel:
			if event.ClusterEvent == events.EventTYPE(0) {
				cluster.JoinCluster(event)
			} else if event.ClusterEvent == events.EventTYPE(1) {
				cluster.LeaveCluster(event)
			}
		}
	}
}

func (cluster *ClusterDetails) JoinCluster(event events.ClusterEvent) {
	fmt.Println("New Node Joined Cluster -> ", event.ClusterMember.NodeID)
	for _, clusterMember := range cluster.ClusterMemList {
		broadcast.BroadCastEvents(clusterMember, event.ClusterMember)
	}
}

func (cluster *ClusterDetails) LeaveCluster(event events.ClusterEvent) {
	fmt.Println("Node Removed Cluster -> ", event.ClusterMember.NodePort)

	cluster.Lock()
	nodeId := event.ClusterMember.NodePort
	defer cluster.Unlock()
	for i, mem := range cluster.ClusterMemList {
		if mem.NodePort == nodeId {
			// Remove the member from the slice
			cluster.ClusterMemList = append(cluster.ClusterMemList[:i], cluster.ClusterMemList[i+1:]...)
		}
	}
	fmt.Println("Nodes Remaining.. ")
	for _, clusterMember := range cluster.ClusterMemList {
		broadcast.BroadLeaveCastEvents(clusterMember, event.ClusterMember)
	}
}
