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
	if err != nil {
		return
	}
	err = serv.ListenAndAccept()
	if err != nil {
		return
	}
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("http server stopped")
	err = serv.StopServer()
}
