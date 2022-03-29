package executor

import (
	"fmt"
	"log"
	"context"
	"google.golang.org/protobuf/types/known/structpb"
	pluginpb "github.com/xellos00/dk-yuba-proto/dist/proto/vatz/plugin/v1"
	//"vatz/manager/config"
)

var (
	executorInstance Executor
	EManager         executor_manager
)

func init() {
	executorInstance = NewExecutor()
}

type executor_manager struct {
}

func (s *executor_manager) Execute() error {
	fmt.Println("this is Execute call from Manager ")

	//TODO: get dial address, protocol from config
	targetMap := map[string]interface{}{
		"addr": "localhost:9091",
		"protocol": "vatz-plugin-cosmos",
	}

	target, err := structpb.NewStruct(targetMap)
	if err != nil {
		log.Fatalf("failed to check target structpb: %v", err)
	}

	//TODO: get target plugin info and commands from config
	commandMap := map[string]interface{}{
		"command": "getBlockHeight",
		"params": "",
	}

	commands, err := structpb.NewStruct(commandMap)
	if err != nil {
		log.Fatalf("failed to check command structpb: %v", err)
	}

	req := &pluginpb.ExecuteRequest{
		ExecuteInfo: target,
		Options: commands,
	}

	resp, _ := executorInstance.Execute(context.Background(), req)

	//TODO: handle return response
	log.Printf("received state :%v", resp.GetState())

	return nil
}
