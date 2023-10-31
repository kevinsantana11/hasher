package cluster

import (
	"slashslinging/hasher/models/server"
)

type Cluster struct {
	servers []server.Server
}

func New(count int) *Cluster {
	servers := make([]server.Server, 0)
	for i := 0; i < count; i++ {
		servers = append(servers, server.New(i))
	}
	return &Cluster{servers}
}

func (clus Cluster) Servers() []server.Server {
	return clus.servers
}

func (clus Cluster) GetServer(id int) (server *server.Server) {
	for _, s := range clus.servers {
		if s.Id() == id {
			server = &s
		}
	}
	return
}

func (clus *Cluster) Add() int {
	servId := clus.servers[len(clus.servers)-1].Id() + 1
	clus.servers = append(clus.servers, server.New(servId))
	return servId
}

func (clus *Cluster) Del(id int) (delServer server.Server) {
	newArray := make([]server.Server, 0)

	for _, server := range clus.servers {
		if server.Id() != id {
			newArray = append(newArray, server)
		} else {
			delServer = server
		}
	}

	clus.servers = newArray
	return
}
