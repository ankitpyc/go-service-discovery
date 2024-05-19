package cluster

type EventTYPE int

const (
	NEW_JOIN EventTYPE = iota
	LEAVE
)

type ClusterEvent struct {
	ClusterEvent  EventTYPE
	ClusterMember ClusterMember
}

func NewClusterEvent(eventType EventTYPE, member ClusterMember) *ClusterEvent {
	return &ClusterEvent{
		ClusterEvent:  eventType,
		ClusterMember: member,
	}
}
