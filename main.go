package main

import (
	"go-service-discovery/server"
	"log"
	"os"
	"os/signal"
)

const (
	Host = "127.0.0.1"
	Port = "2212"
)

func main() {
	serv, err := server.NewServer(Host, Port).StartServer()
	if serv == nil || err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	err = serv.ListenAndAccept()
	if err != nil {
		log.Fatalf("Failed while listening for connections: %v", err)
	}
	handleStopServer(err, serv)
}

func handleStopServer(err error, serv *server.Server) {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("Server is shutting down")
	err = serv.StopServer()
}
