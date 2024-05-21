package cl

import "sync/atomic"

type NODETYPE int

type ClusterEventTYPE int

// Exported enum values
const (
	EventAdd ClusterEventTYPE = iota
	EventRemove
	// Add other event types as necessary
)
const (
	FOLLOWER NODETYPE = iota
	LEADER
	CANDIDATE
)

type ClusterMember struct {
	NodeType         NODETYPE
	NodeID           string
	NodeAddr         string
	NodePort         string
	ClusterID        string
	MissedHeartbeats int
	MissedCount      atomic.Int32
}

type ClusterEvent struct {
	ClusterEvent  ClusterEventTYPE
	ClusterMember *ClusterMember
}

type ClusterConfig struct {
	ClusterMemList   []*ClusterMember
	ClusterName      string
	ClusterID        string
	TotalSize        int
	BroadCastChannel chan ClusterEvent
}

func NewClusterMember(nodeType string, nodeID string, nodeAddr string, nodePort string) *ClusterMember {
	return &ClusterMember{
		NodeType: FOLLOWER,
		NodeID:   nodeID,
		NodeAddr: nodeAddr,
		NodePort: nodePort,
	}
}
func NewClusterEvent(eventType ClusterEventTYPE, member *ClusterMember) *ClusterEvent {
	return &ClusterEvent{
		ClusterEvent:  eventType,
		ClusterMember: member,
	}
}
