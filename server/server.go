package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-service-discovery/cluster"
	"go-service-discovery/cluster/config"
	"go-service-discovery/cluster/events"
	"go-service-discovery/prober"
	"log"
	"net"
	"time"
)

// EventTYPE defines types of events that can occur in the cluster
type EventTYPE int

const (
	NEW_JOIN EventTYPE = iota // A new member has joined the cluster
	LEAVE                     // A member has left the cluster
)

// StartServer initializes the TCP server and starts listening for connections
func (server *Server) StartServer() (*Server, error) {
	// Start listening on the specified host and port
	listener, err := net.Listen("tcp", server.Host+":"+server.Port)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return nil, err
	}
	server.TCPListener = listener
	fmt.Println("TCP Server listening on", server.Host+":"+server.Port)
	// Start the health check routine
	go InitiateHealthCheck(server)
	return server, nil
}

// InitiateHealthCheck starts a routine to periodically check the health of each cluster
func InitiateHealthCheck(server *Server) {
	fmt.Println("Initiating health checks ...")
	prob := prober.NewProberService()
	ctx := context.Background()
	timer := time.NewTicker(time.Second * 20)
	defer timer.Stop()
	// Create a ticker that ticks every 10 seconds
	// Ensure the ticker is stopped when the function exits
	go prob.MonitorForFailedChecks(ctx)
	for {
		select {
		case <-timer.C:
			// Lock the server to safely access the cluster details
			server.RLock()
			for _, clusterConfig := range server.ClusterDetails {
				// Perform health check on each cluster in a separate goroutine
				log.Printf("Running Cluster Health Check for cluster %s ... :-  ", clusterConfig.ClusterID)
				go prob.ClusterHealthCheck(ctx, clusterConfig)
			}
			server.RUnlock()
		}
	}
}

// ListenAndAccept accepts incoming connections and handles them
func (server *Server) ListenAndAccept() error {
	for {
		// Accept a new connection
		conn, err := server.TCPListener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return err
		}
		// Handle the connection in a new goroutine
		go server.handleConnection(conn)
	}
}

// handleConnection handles an individual connection to the server
func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when the function exits

	for {
		buf := make([]byte, 1024) // Create a buffer to read data from the connection
		bytesread, err := conn.Read(buf)
		if err != nil || bytesread == 0 {
			log.Printf("Failed to read from connection: %v", err)
			return
		}
		eventType := buf[0]
		switch eventType {
		case 0:
			// Handle a new node joining the cluster
			if server.UpdateClusterConfig(conn, buf, bytesread) {
				return
			}
		case 1:
			// Handle a new node joining the cluster
			if server.RemoveClusterConfig(conn, buf, bytesread) {
				return
			}
		}
	}
}

func (server *Server) UpdateClusterConfig(conn net.Conn, buf []byte, readsize int) bool {
	var NodeDetails *cluster.ClusterMember
	nodeDetails := buf[1:readsize]
	err := json.Unmarshal(nodeDetails, &NodeDetails)
	if err != nil {
		fmt.Printf("Failed to unmarshal node details: %v\n", err)
		return true
	}
	clusters, existing := IdentifyCluster(server, NodeDetails)
	clusters.AddClusterMemberList(NodeDetails)

	if !existing {
		server.ClusterDetails = append(server.ClusterDetails, clusters)
		server.ClusterDetails[0].ClusterSecretKey, _ = generateSecretKey(32)
		go clusters.ListenForBroadcasts()
	}
	clusters.BroadCastChannel <- events.NewClusterEvent(events.EventTYPE(0), *NodeDetails)
	clusterInfo, err := json.Marshal(config.RegistrationRepose{Members: clusters.ClusterMemList, Secret: server.ClusterDetails[0].ClusterSecretKey})
	if err != nil {
		log.Printf("Failed to marshal clusters info: %v\n", err)
	}
	_, err = conn.Write(clusterInfo)
	return false
}

func generateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (server *Server) RemoveClusterConfig(conn net.Conn, buf []byte, readsize int) bool {
	fmt.Println("Removing cluster members .. ")
	var NodeDetails *cluster.ClusterMember
	nodeDetails := buf[1:readsize]
	err := json.Unmarshal(nodeDetails, &NodeDetails)
	if err != nil {
		fmt.Printf("Failed to unmarshal node details: %v\n", err)
		return true
	}
	clusters, existing := IdentifyCluster(server, NodeDetails)
	if existing {
		clusters.BroadCastChannel <- events.NewClusterEvent(events.EventTYPE(1), *NodeDetails)
	}
	nodesInfo, err := json.Marshal(clusters.ClusterMemList)
	if err != nil {
		log.Printf("Failed to marshal clusters info: %v\n", err)
	}
	_, err = conn.Write(nodesInfo)
	return false
}

// StopServer stops the TCP server from listening for new connections
func (server *Server) StopServer() error {
	err := server.TCPListener.Close()
	fmt.Println("TCP Server stopped listening. Server will be stopped")
	if err != nil {
		return err
	}
	return nil
}

// IdentifyCluster finds or creates a cluster configuration for the given node
func IdentifyCluster(s *Server, node *cluster.ClusterMember) (*config.ClusterDetails, bool) {
	s.Lock()
	defer s.Unlock()
	var existing bool = true
	var clusterConfig *config.ClusterDetails
	for _, cluster := range s.ClusterDetails {
		if cluster.ClusterID == node.ClusterID {
			clusterConfig = cluster
			break
		}
	}

	if clusterConfig == nil {
		existing = false
		// Create a new cluster configuration if not found
		clusterConfig = &config.ClusterDetails{
			ClusterID:        node.ClusterID,
			ClusterName:      "",
			ClusterMemList:   make([]*cluster.ClusterMember, 0, 5), // Initialize with a capacity of 5
			BroadCastChannel: make(chan events.ClusterEvent),
		}
	}
	return clusterConfig, existing
}
