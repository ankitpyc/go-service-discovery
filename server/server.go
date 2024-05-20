package server

import (
	"encoding/json"
	"fmt"
	"go-service-discovery/Discovery"
	"go-service-discovery/cluster"
	"log"
	"net"
	"sync"
	"time"
)

// EventTYPE defines types of events that can occur in the cluster
type EventTYPE int

const (
	NEW_JOIN EventTYPE = iota // A new member has joined the cluster
	LEAVE                     // A member has left the cluster
)

// StartServer initializes the TCP server and starts listening for connections
func (s *Server) StartServer() (*Server, error) {
	// Start listening on the specified host and port
	listener, err := net.Listen("tcp", s.Host+":"+s.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return nil, err
	}
	s.TCPListener = listener
	fmt.Println("TCP Server listening on", s.Host+":"+s.Port)

	// Start the health check routine
	go InitiateHealthCheck(s)

	return s, nil
}

// InitiateHealthCheck starts a routine to periodically check the health of each cluster
func InitiateHealthCheck(s *Server) {
	discoveryService := Discovery.DiscoveryService{}
	timer := time.NewTicker(time.Second * 10) // Create a ticker that ticks every 10 seconds
	defer timer.Stop()                        // Ensure the ticker is stopped when the function exits

	for {
		select {
		case <-timer.C:
			// Lock the server to safely access the cluster details
			s.SSMu.RLock()
			for _, clusterConfig := range s.ClusterDetails {
				// Perform health check on each cluster in a separate goroutine
				go func(cc *cluster.ClusterConfig) {
					discoveryService.ClusterHealthCheck(cc)
				}(clusterConfig)
			}
			s.SSMu.RUnlock()
		}
	}
}

// ListenAndAccept accepts incoming connections and handles them
func (s *Server) ListenAndAccept() error {
	for {
		// Accept a new connection
		conn, err := s.TCPListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return err
		}
		// Handle the connection in a new goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection handles an individual connection to the server
func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function exits

	for {
		buf := make([]byte, 1024) // Create a buffer to read data from the connection
		readsize, err := conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		if readsize == 0 {
			return // Connection closed by client
		}

		fmt.Println(readsize)
		switch buf[0] {
		case 0:
			// Handle a new node joining the cluster
			var NodeDetails cluster.ClusterMember
			nodeDetails := buf[1:readsize]
			err := json.Unmarshal(nodeDetails, &NodeDetails)
			if err != nil {
				fmt.Printf("Failed to unmarshal node details: %v\n", err)
				return
			}

			clusters := IdentifyCluster(server, &NodeDetails)
			clusters.Mut.Lock()
			clusters.AddClusterMemberList(NodeDetails)
			clusters.BroadCastChannel <- *clusters.CreateClusterEvent(1, NodeDetails)
			clusters.Mut.Unlock()
		}
	}
}

// StopServer stops the TCP server from listening for new connections
func (s *Server) StopServer() error {
	err := s.TCPListener.Close()
	fmt.Println("TCP Server stopped listening. Server will be stopped")
	if err != nil {
		return err
	}
	return nil
}

// IdentifyCluster finds or creates a cluster configuration for the given node
func IdentifyCluster(s *Server, node *cluster.ClusterMember) *cluster.ClusterConfig {
	s.SSMu.RLock()
	defer s.SSMu.Unlock()

	var clusterConfig *cluster.ClusterConfig
	for _, cluster := range s.ClusterDetails {
		if cluster.ClusterID == node.ClusterID {
			clusterConfig = cluster
			break
		}
	}

	if clusterConfig == nil {
		// Create a new cluster configuration if not found
		clusterConfig = &cluster.ClusterConfig{
			ClusterID:        node.ClusterID,
			ClusterName:      "",
			ClusterMemList:   make([]*cluster.ClusterMember, 0, 5), // Initialize with a capacity of 5
			BroadCastChannel: make(chan cluster.ClusterEvent),
			Mut:              sync.RWMutex{},
		}
		s.ClusterDetails = append(s.ClusterDetails, clusterConfig)
	}

	return clusterConfig
}
