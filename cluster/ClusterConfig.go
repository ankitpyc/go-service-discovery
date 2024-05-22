package cluster

import (
	"fmt"
	"sync"
)

type ClusterConfig struct {
	ClusterMemList   []*ClusterMember
	ClusterName      string
	ClusterID        string
	TotalSize        int
	BroadCastChannel chan ClusterEvent
	Mut              sync.RWMutex
}

func (cc *ClusterConfig) AddClusterMemberList(member *ClusterMember) []*ClusterMember {
	cc.Mut.Lock()
	defer cc.Mut.Unlock()
	member.NodeStatus = "HEALTHY"
	fmt.Println("Discovered Node -> ", member.NodeAddr+":"+member.NodePort)
	cc.ClusterMemList = append(cc.ClusterMemList, member)
	return cc.ClusterMemList
}

func (cc *ClusterConfig) ListenForBroadcasts() {
	for {
		select {
		case event := <-cc.BroadCastChannel:
			if event.ClusterEvent == EventTYPE(0) {
				cc.JoinCluster(event)
			} else if event.ClusterEvent == EventTYPE(1) {
				cc.LeaveCluster(event)
			}
		}
	}
}

func (cc *ClusterConfig) JoinCluster(event ClusterEvent) {
	fmt.Println("New Node Joined Cluster -> ", event.ClusterMember.NodeID)
	for _, clusterMember := range cc.ClusterMemList {
		BroadCastEvents(clusterMember, event.ClusterMember)
	}
}

func (cc *ClusterConfig) LeaveCluster(event ClusterEvent) {
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

func (cc *ClusterConfig) CreateClusterEvent(eventTYPE EventTYPE, details ClusterMember) *ClusterEvent {
	return &ClusterEvent{
		ClusterEvent:  eventTYPE,
		ClusterMember: details,
	}
}
