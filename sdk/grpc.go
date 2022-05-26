package sdk

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type grpcServer struct {
	pb.UnimplementedPluginServer

	srv       *grpc.Server
	callbacks []func(map[string]*structpb.Value, map[string]*structpb.Value) (CallResponse, error)
}

// Verify returns liveness.
func (s *grpcServer) Verify(context.Context, *emptypb.Empty) (*pb.VerifyInfo, error) {
	return &pb.VerifyInfo{
		VerifyMsg: "OK",
	}, nil
}

// Execute runs plugin features.
func (s *grpcServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Println("PluginServer.Execute")

	resp := &pb.ExecuteResponse{
		ResourceType: PluginName,
	}

	for _, f := range s.callbacks {
		var (
			executeInfo map[string]*structpb.Value
			option      map[string]*structpb.Value
		)

		if req.GetExecuteInfo() != nil {
			executeInfo = req.GetExecuteInfo().GetFields()
		}

		if req.GetOptions() != nil {
			option = req.GetOptions().GetFields()
		}

		callResp, err := f(executeInfo, option)

		resp.AlertType = callResp.AlertTypes
		if err != nil {
			resp.Severity = pb.SEVERITY_ERROR
			resp.State = pb.STATE_FAILURE
			resp.Message = err.Error()
		} else {
			resp.Severity = callResp.Severity
			resp.State = callResp.State
			resp.Message = callResp.Message
		}
	}

	return resp, nil
}

// Start starts gRPC service.
func (s *grpcServer) Start(ctx context.Context, address string, port int) error {
	log.Println("grpcServer - Start")

	c, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Println(err)
		return err
	}

	s.srv = grpc.NewServer()

	pb.RegisterPluginServer(s.srv, s)
	reflection.Register(s.srv)

	return s.srv.Serve(c)
}

func (s *grpcServer) Stop() {
	log.Println("grpcServer - Stop")

	if s.srv != nil {
		s.srv.GracefulStop()
	}
}
