package registry

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	agent "github.com/dsrvlabs/vatz-proto/manager/v2"
)

type registrationServer struct {
	agent.UnimplementedRegistrationServiceServer
	reflector PluginReflector
}

func (s *registrationServer) RegisterPlugin(ctx context.Context, in *agent.RegisterRequest) (*agent.RegisterResponse, error) {
	log.Info().Str("module", "registry").Msg("register plugin")

	log.Info().Str("module", "registry").Msgf("request %s:%d", in.Address, in.Port)

	err := s.reflector.Query(ctx, in.Name, in.Address, int(in.Port))
	if err != nil {
		return nil, err
	}

	return &agent.RegisterResponse{
		Msg: fmt.Sprintf("%s - %s:%d", in.Name, in.Address, in.Port),
	}, nil
}

// StartRegistrationService starts registration service
func StartRegistrationService(port int) error {
	log.Info().Str("module", "registry").Msg("start registration service")

	s := grpc.NewServer()
	agent.RegisterRegistrationServiceServer(
		s, &registrationServer{reflector: NewPluginReflector()})

	reflection.Register(s)

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Err(err)
		return err
	}

	if err := s.Serve(l); err != nil {
		log.Err(err)
		return err
	}

	return nil
}
