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
	fmt.Println("Discovered Node -> ", member.NodeAddr+":"+member.NodePort)
	cc.ClusterMemList = append(cc.ClusterMemList, member)
	return cc.ClusterMemList
}

func (cc *ClusterConfig) ListenForBroadcasts() {
	for {
		select {
		case NewMember := <-cc.BroadCastChannel:
			for _, clusterMember := range cc.ClusterMemList {
				BroadCastEvents(clusterMember, NewMember.ClusterMember)
			}
		}
	}
}

func (cc *ClusterConfig) CreateClusterEvent(eventTYPE EventTYPE, details ClusterMember) *ClusterEvent {
	return &ClusterEvent{
		ClusterEvent:  eventTYPE,
		ClusterMember: details,
	}
}
