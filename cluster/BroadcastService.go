package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func BroadCastEvents(targetMember *ClusterMember, data ClusterMember) {
	body, err := json.Marshal(data)
	if err != nil {
		return
	}
	post, err := http.Post(targetMember.NodeAddr+":"+targetMember.NodePort+"/JOIN", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	if post.StatusCode != http.StatusOK {
		return
	}
}
