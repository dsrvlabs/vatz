package main

import (
	"log"
	"time"

	plugin_grpc "pilot-manager/grpc"
	"pilot-manager/manager"
)

const (
	servName = "Sample Node manager"
)

func main() {
	log.Println("Start ", servName)

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
