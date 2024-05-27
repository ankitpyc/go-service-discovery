package cluster

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
	MissedHeartbeats int32
	NodeStatus       string
	GrpcPort         string
}

func NewClusterMember(nodeType string, nodeID string, nodeAddr string, nodePort string) *ClusterMember {
	return &ClusterMember{
		NodeType: FOLLOWER,
		NodeID:   nodeID,
		NodeAddr: nodeAddr,
		NodePort: nodePort,
	}
}
