package server

import (
	clusters "go-service-discovery/cluster"
	"net"
	"sync"
)

type Server struct {
	Host           string
	Port           string
	TCPListener    net.Listener
	SSMu           sync.RWMutex
	ClusterDetails []*clusters.ClusterConfig
}

func NewServer(host string, port string) *Server {
	return &Server{
		Host:           host,
		Port:           port,
		ClusterDetails: []*clusters.ClusterConfig{},
		SSMu:           sync.RWMutex{},
	}
}
