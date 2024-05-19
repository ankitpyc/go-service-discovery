package cluster

type NODETYPE int

const (
	FOLLOWER NODETYPE = iota
	LEADER
	CANDIDATE
)

type ClusterMember struct {
	NodeType  NODETYPE
	NodeID    string
	NodeAddr  string
	NodePort  string
	ClusterID string
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
