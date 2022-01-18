package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	structpb "github.com/golang/protobuf/ptypes/struct"
	manager_pluginpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/plugin"
	managerpb "github.com/xellos00/silver-bentonville/dist/proto/dsrv/api/node_manager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"pilot-manager/manager"
)

const (
	grpcPort = 9090
)

var (
	ManagerInstance manager.Manager
)

type grpcService struct {
	managerpb.UnimplementedNodeManagerServer
}

func (s *grpcService) Execute(ctx context.Context, in *managerpb.ExecuteRequest) (*managerpb.ExecuteResponse, error) {

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

func (s *grpcService) Init(ctx context.Context, in *managerpb.InitRequest) (*managerpb.InitResponse, error) {
	// TODO: Check already running.

	err := ManagerInstance.Start()
	if err != nil {
		return &managerpb.InitResponse{Result: managerpb.CommandStatus_FAIL}, nil
	}

	return &managerpb.InitResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *grpcService) End(ctx context.Context, in *managerpb.EndRequest) (*managerpb.EndResponse, error) {
	// TODO: Check running.

	err := ManagerInstance.Stop()
	if err != nil {
		return &managerpb.EndResponse{Result: managerpb.CommandStatus_FAIL}, nil
	}

	return &managerpb.EndResponse{Result: managerpb.CommandStatus_SUCCESS}, nil
}

func (s *grpcService) Verify(ctx context.Context, in *managerpb.VerifyRequest) (*managerpb.VerifyInfo, error) {
	// TODO: Update config and refresh service.
	// Check how to verify this connection or API call is valid.
	return nil, nil
}

func (s *grpcService) UpdateConfig(ctx context.Context, in *managerpb.UpdateRequest) (*managerpb.UpdateResponse, error) {
	// TODO: Update config and refresh service.
	return nil, nil
}

func init() {
	ManagerInstance = manager.RunManager()
}

// StartServer try to start grpc service.
func StartServer() error {
	s := grpc.NewServer()
	serv := grpcService{}

	managerpb.RegisterNodeManagerServer(s, &serv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", grpcPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
		return err
	}

	log.Println("listen ", addr)

	go func() {
		if err := s.Serve(l); err != nil {
			log.Panic(err)
		}
	}()

	return nil
}
