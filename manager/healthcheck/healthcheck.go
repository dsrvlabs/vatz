package healthcheck

import (
	"fmt"
	"log"
	"context"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	//var empty = make(map[interface{}]interface{})
	e := new(emptypb.Empty)

	callop := grpc.UseCompressor("test")
	verify, err := client.Verify(context.Background(), e, callop)
	fmt.Println("received verify: %v", verify)
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
