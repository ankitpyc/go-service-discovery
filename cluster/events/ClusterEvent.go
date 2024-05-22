package events

import "go-service-discovery/cluster"

type EventTYPE int

const (
	NEW_JOIN EventTYPE = iota
	LEAVE
)

type ClusterEvent struct {
	ClusterEvent  EventTYPE
	ClusterMember cluster.ClusterMember
}

func NewClusterEvent(eventType EventTYPE, member cluster.ClusterMember) ClusterEvent {
	return ClusterEvent{
		ClusterEvent:  eventType,
		ClusterMember: member,
	}
}
