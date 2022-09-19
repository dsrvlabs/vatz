package rpc

import (
	"context"
	"net"
	"net/http"

	vatzpb "github.com/dsrvlabs/vatz-proto/rpc/v1"
	"github.com/dsrvlabs/vatz/manager/healthcheck"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	// TODO: Address should be configurable.
	grpcAddress = "localhost:19090"
	httpAddress = "localhost:19091"
)

// VatzRPC provides RPC interfaces.
type VatzRPC interface {
	Start() error
	Stop()
}

type rpcService struct {
	ctx    context.Context
	cancel context.CancelFunc

	vatzRPCService vatzpb.VatzRPCServer
	grpcServer     *grpc.Server
	httpServer     *http.Server
}

func (s *rpcService) Start() error {
	log.Info().Str("module", "rpc").Msg("start rpc server")

	errChan := make(chan error, 2)

	go func(errChan chan<- error) {
		log.Info().Str("module", "rpc").Msg("start gRPC server")

		l, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Info().Str("module", "rpc").Err(err)
			errChan <- err
			return
		}

		s.grpcServer = grpc.NewServer()
		s.vatzRPCService = &grpcService{
			healthChecker: healthcheck.GetHealthChecker(),
		}

		vatzpb.RegisterVatzRPCServer(s.grpcServer, s.vatzRPCService)
		reflection.Register(s.grpcServer)

		err = s.grpcServer.Serve(l)
		if err != nil {
			log.Info().Str("module", "rpc").Err(err)
			errChan <- err
			return
		}
	}(errChan)

	go func(errChan chan<- error) {
		log.Info().Str("module", "rpc").Msg("start gRPC gateway server")

		mux := runtime.NewServeMux()

		// FIXME: Remove hardcoded address
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		err := vatzpb.RegisterVatzRPCHandlerFromEndpoint(s.ctx, mux, grpcAddress, opts)
		if err != nil {
			log.Info().Str("module", "rpc").Err(err)
			errChan <- err
			return
		}

		s.httpServer = &http.Server{
			Addr:    httpAddress,
			Handler: mux,
		}

		err = s.httpServer.ListenAndServe()
		if err != nil {
			log.Info().Str("module", "rpc").Err(err)
			errChan <- err
			return
		}
	}(errChan)

	err := <-errChan

	log.Info().Str("module", "rpc").Err(err)

	return err
}

func (s *rpcService) Stop() {
	log.Info().Str("module", "rpc").Msg("cancel")
	defer s.cancel()

	if s.httpServer != nil {
		s.httpServer.Shutdown(s.ctx)
	}

	log.Info().Str("module", "rpc").Msg("stop")
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
}

// NewRPCService creates new rpc server instance.
func NewRPCService() VatzRPC {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &rpcService{
		ctx:    ctx,
		cancel: cancel,
	}
}
