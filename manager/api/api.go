package api

import (
	"context"
	managerpb "github.com/hqueue/vatz-secret/dist/proto/vatz/manager/v1"
)

var (
	ExecutableRPC GrpcService
)

type GrpcService struct {
	managerpb.UnimplementedManagerServer
}

func (s *GrpcService) Execute(ctx context.Context, in *managerpb.ExecuteRequest) (*managerpb.ExecuteResponse, error) {

	//opts := grpc.WithInsecure()
	//cc, err := grpc.Dial("localhost:9091", opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer cc.Close()
	//
	//client := pluginpb.NewPluginClient(cc)
	//request := &pluginpb.ExecuteRequest{ExecuteInfo: in.TargetInfo, Options: in.Command}
	//aresp, _ := client.Execute(context.Background(), request)

	//data := make(map[interface{}]interface{})
	//
	//resp := managerpb.ExecuteResponse{
	//	State:    managerpb.ExecuteResponse_SUCCESS,
	//	Message:  "Empty Result",
	//	Protocol: "Near",
	//	Options:  data,
	//}

	return nil, nil
}

func (s *GrpcService) Init(ctx context.Context, in *managerpb.InitRequest) (*managerpb.InitResponse, error) {

	return &managerpb.InitResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *GrpcService) End(ctx context.Context, in *managerpb.EndRequest) (*managerpb.EndResponse, error) {
	return &managerpb.EndResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *GrpcService) Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error) {
	return nil, nil
}

func (s *GrpcService) UpdateConfig(ctx context.Context, in *managerpb.UpdateRequest) (*managerpb.UpdateResponse, error) {
	return nil, nil
}
