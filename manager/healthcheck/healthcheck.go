package healthcheck

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
)

type healthCheck struct {
}

func (h healthCheck) HealthCheck() (string, error) {
	fmt.Println("this is sending Notification")
	//TODO: Set Client to Plugin
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:9091", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	var client = pluginpb.NewPluginClient(cc)
	var empty = make(map[interface{}]interface{})
	verify, err := client.Verify(empty)
	if err != nil {
		return "", err
	}
	return "UP", nil
}

type HealthCheck interface {
	HealthCheck() (string, error)
}

func NewHealthChecker() HealthCheck {
	return &healthCheck{}
}
