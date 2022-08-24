package rpc

import (
	"net"

	vatzpb "github.com/dsrvlabs/vatz-proto/rpc/v1"
	"github.com/dsrvlabs/vatz/manager/healthcheck"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// VatzRPC provides RPC interfaces.
type VatzRPC interface {
	Start() error
	Stop()
}

type rpcService struct {
	vatzRPCService vatzpb.VatzRPCServer
	grpcServer     *grpc.Server
}

func (s *rpcService) Start() error {
	log.Info().Str("module", "rpc").Msg("start rpc server")

	l, err := net.Listen("tcp", "127.0.0.1:19090")
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	s.vatzRPCService = &grpcService{
		healthChecker: healthcheck.GetHealthChecker(),
	}

	vatzpb.RegisterVatzRPCServer(s.grpcServer, s.vatzRPCService)
	reflection.Register(s.grpcServer)

	return s.grpcServer.Serve(l)
}

func (s *rpcService) Stop() {
	s.grpcServer.Stop()
}

// NewRPCService creates new rpc server instance.
func NewRPCService() VatzRPC {
	return &rpcService{}
}
