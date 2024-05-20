package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func BroadCastEvents(targetMember *ClusterMember, data ClusterMember) {
	fmt.Println("BroadCastEvents Join Events from  ", data.NodeAddr+":"+data.NodePort)
	fmt.Println("BroadCastEvents Join Events to  ", targetMember.NodeAddr+":"+targetMember.NodePort)
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
