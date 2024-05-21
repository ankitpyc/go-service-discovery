package cluster

import (
	"sync/atomic"
)

type NODETYPE int

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

func NewClusterMember(nodeType string, nodeID string, nodeAddr string, nodePort string) *ClusterMember {
	return &ClusterMember{
		NodeType: FOLLOWER,
		NodeID:   nodeID,
		NodeAddr: nodeAddr,
		NodePort: nodePort,
	}
}
