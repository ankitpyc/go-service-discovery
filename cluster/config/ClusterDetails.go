package config

import (
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/broadcast"
	"go-service-discovery/cluster/events"
	"sync"
)

type ClusterDetails struct {
	ClusterMemList   []*cluster.ClusterMember
	ClusterName      string
	ClusterID        string
	TotalSize        int
	BroadCastChannel chan events.ClusterEvent
	Mut              sync.RWMutex
}

func (cc *ClusterDetails) AddClusterMemberList(member *cluster.ClusterMember) []*cluster.ClusterMember {
	cc.Mut.Lock()
	defer cc.Mut.Unlock()
	member.NodeStatus = "HEALTHY"
	fmt.Println("Discovered Node -> ", member.NodeAddr+":"+member.NodePort)
	cc.ClusterMemList = append(cc.ClusterMemList, member)
	return cc.ClusterMemList
}

func (cc *ClusterDetails) ListenForBroadcasts() {
	for {
		select {
		case event := <-cc.BroadCastChannel:
			if event.ClusterEvent == events.EventTYPE(0) {
				cc.JoinCluster(event)
			} else if event.ClusterEvent == events.EventTYPE(1) {
				cc.LeaveCluster(event)
			}
		}
	}
}

func (cc *ClusterDetails) JoinCluster(event events.ClusterEvent) {
	fmt.Println("New Node Joined Cluster -> ", event.ClusterMember.NodeID)
	for _, clusterMember := range cc.ClusterMemList {
		broadcast.BroadCastEvents(clusterMember, event.ClusterMember)
	}
}

func (cc *ClusterDetails) LeaveCluster(event events.ClusterEvent) {
	fmt.Println("Node Removed Cluster -> ", event.ClusterMember.NodeID)

	cc.Mut.Lock()
	nodeId := event.ClusterMember.NodeID
	defer cc.Mut.Unlock()
	for i, mem := range cc.ClusterMemList {
		if mem.NodeID == nodeId {
			// Remove the member from the slice
			cc.ClusterMemList = append(cc.ClusterMemList[:i], cc.ClusterMemList[i+1:]...)
			return
		}
	}
}
