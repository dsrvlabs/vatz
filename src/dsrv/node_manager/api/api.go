package api

import (
	"context"
	structpb "github.com/golang/protobuf/ptypes/struct"
	manager_pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	"google.golang.org/grpc"
	"log"
	manager_presenter "pilot-manager/src/dsrv/node_manager/manager"
)

var (
	ManagerInstance manager_presenter.Manager
	ExecutableRPC   GrpcService
)

type GrpcService struct {
	managerpb.UnimplementedNodeManagerServer
}

func (s *GrpcService) Execute(ctx context.Context, in *managerpb.ExecuteRequest) (*managerpb.ExecuteResponse, error) {

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:9091", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := manager_pluginpb.NewManagerPluginClient(cc)
	request := &manager_pluginpb.ExecuteRequest{ExecuteInfo: in.TargetInfo, Options: in.Command}
	aresp, _ := client.Execute(context.Background(), request)

	var options = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"Name": &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: "Sample_Name",
				},
			},
		},
	}

	resp := managerpb.ExecuteResponse{
		State:    managerpb.ExecuteResponse_SUCCESS,
		Message:  aresp.Message,
		Protocol: "near",
		Options:  options,
	}

	return &resp, nil
}

func (s *GrpcService) Init(ctx context.Context, in *managerpb.InitRequest) (*managerpb.InitResponse, error) {
	// TODO: Check already running, if not it requires to start plugins
	err := ManagerInstance.Init()
	if err != nil {
		return &managerpb.InitResponse{Result: managerpb.CommandStatus_FAIL}, nil
	}
	return &managerpb.InitResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *GrpcService) End(ctx context.Context, in *managerpb.EndRequest) (*managerpb.EndResponse, error) {
	// TODO: Kill the Process if there's running plugins.
	err := ManagerInstance.End()
	if err != nil {
		return &managerpb.EndResponse{Result: managerpb.CommandStatus_FAIL}, nil
	}
	return &managerpb.EndResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *GrpcService) Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error) {
	// Currently, I do not know whether it requires verifying initialized plugin is up and running.
	return nil, nil
}

func (s *GrpcService) UpdateConfig(ctx context.Context, in *managerpb.UpdateRequest) (*managerpb.UpdateResponse, error) {
	// TODO: Set the Proto for Updatefor

	res, err := manager_presenter.RunManager().UpdateConfig(ctx, in)
	if err != nil {
		return &managerpb.UpdateResponse{}, nil
	}
	return res, nil
}

func init() {
	ManagerInstance = manager_presenter.RunManager()
}
