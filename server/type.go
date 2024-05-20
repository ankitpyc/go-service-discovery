package server

import (
	clusters "go-service-discovery/cluster"
	"net"
)

type Server struct {
	Host           string
	Port           string
	TCPListener    net.Listener
	ClusterDetails []*clusters.ClusterConfig
}

func NewServer(host string, port string) *Server {
	return &Server{
		Host:           host,
		Port:           port,
		ClusterDetails: []*clusters.ClusterConfig{},
	}
}
