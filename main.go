package main

import (
	"log"
	"time"

	plugin_grpc "pilot-manager/grpc"
	"pilot-manager/manager"
)

const (
	servName = "Node Manager (Pilot)"
)

func main() {
	log.Println("Starting ...", servName)

	manager.RunManager()
	err := manager.RunManager().Start()
	if err != nil {
		log.Panic(err)
	}

	plugin_grpc.StartServer()

	for {

		time.Sleep(time.Second * 10)
	}
}
