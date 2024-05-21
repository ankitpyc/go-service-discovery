package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type BroadcastInf interface {
	BroadCastEvents(targetMember *ClusterMember, data ClusterMember)
}

func BroadCastEvents(targetMember *ClusterMember, data ClusterMember) {
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
