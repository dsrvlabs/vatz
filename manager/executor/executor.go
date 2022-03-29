package executor

import (
	"log"
	"context"
	"google.golang.org/grpc"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
)

type executor struct {
}

type Executor interface {
	Execute(ctx context.Context, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error)
}

func (v executor) Execute(ctx context.Context, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	log.Printf("executor call plugin")

	opts := grpc.WithInsecure()
	//TODO: get addr from in.GetExecuteInfo()
	var address = "localhost:9091"
	conn, err := grpc.Dial(address, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pluginpb.NewPluginClient(conn)

	//TODO: execute request to plugin
	resp, _ := client.Execute(context.Background(), in)

	//TODO: handle response from plugin
	//alert by resp.state, resp.severity
	log.Printf("received state :%v", resp.GetState())

	return nil, nil
}

func NewExecutor() Executor {
	return &executor{}
}
