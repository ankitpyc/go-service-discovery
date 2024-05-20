package server

import (
	"encoding/json"
	"fmt"
	"go-service-discovery/cluster"
	"log"
	"net"
	"sync"
)

type EventTYPE int

const (
	NEW_JOIN EventTYPE = iota
	LEAVE
)

func (s *Server) StartServer() (*Server, error) {
	listener, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return nil, err
	}
	s.TCPListener = listener
	fmt.Println("TCP Server listening on", s.Host+":"+s.Port)
	return s, nil
}

func (s *Server) ListenAndAccept() error {
	for {
		conn, err := s.TCPListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return err
		}
		go s.handleConnection(conn)

	}

}

func (server *Server) handleConnection(conn net.Conn) {
	for {
		var buf []byte = make([]byte, 1024)
		readsize, err := conn.Read(buf)
		if err != nil {
			defer conn.Close()
		}
		fmt.Println(readsize)
		switch buf[0] {
		case 0:
			var NodeDetails cluster.ClusterMember
			nodeDetails := buf[1:readsize]
			err := json.Unmarshal([]byte(nodeDetails), &NodeDetails)
			if err != nil {
				fmt.Printf("Failed to unmarshal node details: %v\n", err)
				return
			}
			clusters := IdentifyCluster(server, &NodeDetails)
			clusters.Mut.Lock()
			clusters.AddClusterMemberList(NodeDetails)
			clusters.BroadCastChannel <- *clusters.CreateClusterEvent(0, NodeDetails)
			clusters.Mut.Unlock()
		}
	}
}

func (s *Server) StopServer() error {
	err := s.TCPListener.Close()
	fmt.Println("TCP Server stopped listening . Server will be stopped")
	if err != nil {
		return err
	}
	return nil
}

func IdentifyCluster(s *Server, node *cluster.ClusterMember) *cluster.ClusterConfig {
	var clusterConfig *cluster.ClusterConfig = nil
	for _, cluster := range s.ClusterDetails {
		if cluster.ClusterID == node.ClusterID {
			clusterConfig = cluster
		}
	}

	if clusterConfig == nil {
		clusterConfig = &cluster.ClusterConfig{
			ClusterID:        node.ClusterID,
			ClusterName:      "",
			ClusterMemList:   make([]*cluster.ClusterMember, 0, 5),
			BroadCastChannel: make(chan cluster.ClusterEvent),
			Mut:              sync.RWMutex{},
		}
	}
	s.ClusterDetails = append(s.ClusterDetails, clusterConfig)
	return clusterConfig
}
