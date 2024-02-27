package main

import (
	"errors"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	agent "github.com/dsrvlabs/vatz-proto/manager/v2"
	"github.com/dsrvlabs/vatz/engine/endpoint"
	"github.com/dsrvlabs/vatz/engine/handler"
	"github.com/dsrvlabs/vatz/engine/registry"
	"github.com/dsrvlabs/vatz/utils"
)

func init() {
	utils.InitializeLogger()
}

// TODO: Use??
var ErrServiceNotFound = errors.New("cannot find service")

func main() {
	log.Info().Str("module", "main").Msg("start server")

	// TODO: Refactoring here.
	go func() {
		h := handler.NewHandler()

		s := grpc.NewServer()
		agent.RegisterRequestHandlerServer(s, h)

		l, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}

		if err := s.Serve(l); err != nil {
			panic(err)
		}
	}()

	go registry.StartRegistrationService(8080)
	go endpoint.StartEndpointService(9091)

	for {
		time.Sleep(10 * time.Second)
	}
}
