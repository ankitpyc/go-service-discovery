package server

import (
	clusters "go-service-discovery/cluster"
	"net"
	"sync"
	"time"
)

type ServerInf interface {
	StartServer() (*Server, error)
	ListenAndAccept() error
	handleConnection(conn net.Conn)
	StopServer() error
	UpdateClusterConfig(conn net.Conn, buf []byte, readsize int) bool
}

// Server represents a TCP server that manages cluster details and handles connections.
type Server struct {
	Host           string
	Port           string
	TCPListener    net.Listener
	SSMu           sync.RWMutex
	ClusterDetails []*clusters.ClusterConfig
	TimeOut        time.Duration
}

// NewServer initializes a new Server instance with the given host and port.
// It also initializes an empty list of cluster configurations and an RWMutex for synchronization.
func NewServer(host string, port string) *Server {
	return &Server{
		Host:           host,
		Port:           port,
		ClusterDetails: []*clusters.ClusterConfig{}, // Initialize an empty slice for cluster configurations
		SSMu:           sync.RWMutex{},              // Initialize the RWMutex
	}
}
