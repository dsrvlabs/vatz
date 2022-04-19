package sdk

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

const (
	registerFuncLimit = 5
)

// Errors
var (
	ErrRegisterMaxLimit = errors.New("too many register functions")
)

// Plugin provides interfaces to manage plugin.
type Plugin interface {
	Start(ctx context.Context, address string, port int) error
	Stop()
	Register(cb func(info, option map[string]*structpb.Value) error) error
}

type plugin struct {
	grpc grpcServer
	ch   chan os.Signal
}

func (p *plugin) Start(ctx context.Context, address string, port int) error {
	log.Println("plugin - Start")

	p.ch = make(chan os.Signal, 1)
	signal.Notify(p.ch, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		_ = <-p.ch

		log.Println("grpcServer - Shutting down")

		p.grpc.Stop()
	}()

	return p.grpc.Start(ctx, address, port)
}

func (p *plugin) Stop() {
	log.Println("plugin - Stop")

	p.ch <- syscall.SIGTERM
}

func (p *plugin) Register(cb func(info, option map[string]*structpb.Value) error) error {
	log.Println("RegisterFeature function")

	if p.grpc.callbacks == nil {
		p.grpc.callbacks = make([]func(map[string]*structpb.Value, map[string]*structpb.Value) error, 0)
	}

	if len(p.grpc.callbacks) == registerFuncLimit {
		return ErrRegisterMaxLimit
	}

	p.grpc.callbacks = append(p.grpc.callbacks, cb)

	return nil
}

// NewPlugin creates new plugin service instance.
func NewPlugin() Plugin {
	return &plugin{
		grpc: grpcServer{},
	}
}
