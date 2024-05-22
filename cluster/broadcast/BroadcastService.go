package broadcast

import (
	"bytes"
	"encoding/json"
	"go-service-discovery/cluster"
	"net/http"
)

type BroadcastInf interface {
	BroadCastEvents(targetMember *cluster.ClusterMember, data cluster.ClusterMember)
}

func BroadCastEvents(targetMember *cluster.ClusterMember, data cluster.ClusterMember) {
	if targetMember.NodePort == data.NodePort {
		return
	}
	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	url := "http://" + targetMember.NodeAddr + ":" + targetMember.NodePort + "/Join"
	post, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	if post.StatusCode != http.StatusOK {
		return
	}
}
