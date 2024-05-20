package cluster

func (cc *ClusterConfig) AddClusterMemberList(member ClusterMember) []*ClusterMember {
	cc.ClusterMemList = append(cc.ClusterMemList, &member)
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
