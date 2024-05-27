package config

import "go-service-discovery/cluster"

type RegistrationRepose struct {
	Secret  string                   `json:"Secret"`
	Members []*cluster.ClusterMember `json:"Members"`
}
