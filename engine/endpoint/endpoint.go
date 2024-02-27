package endpoint

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/dsrvlabs/vatz-proto/manager/v2"
	"github.com/dsrvlabs/vatz/engine/bucket"
)

type endpointService struct {
	v2.UnimplementedEndpointServiceServer

	bucket bucket.PluginBucket
}

func (s *endpointService) ListPlugin(ctx context.Context, in *v2.ListPluginRequest) (*v2.ListPluginResponse, error) {
	log.Info().Str("module", "endpoint").Msg("list plugin")

	descs := s.bucket.List()

	metas := make([]*v2.PluginMetadata, len(descs))
	for i, desc := range descs {
		metas[i] = &v2.PluginMetadata{Name: desc.Name}
	}

	resp := &v2.ListPluginResponse{Plugin: metas}
	return resp, nil
}

func (s *endpointService) DetailPlugin(ctx context.Context, in *v2.DetailPluginRequest) (*v2.DetailPluginResponse, error) {
	log.Info().Str("module", "endpoint").Msg("detail plugin")

	desc, err := s.bucket.Get(in.PluginName)
	if err != nil {
		return nil, err
	}

	resp := &v2.DetailPluginResponse{
		PluginName: in.PluginName,
	}

	methods := []*v2.PluginMethod{}

	// TODO: How to serialize?
	for name, desc := range desc.Methods {
		methods = append(methods, &v2.PluginMethod{
			Name: name,
		})
		_ = name
		_ = desc
	}

	resp.Methods = methods

	return resp, nil
}

// TODO: Add gateway

// StartEndpointService starts endpoint service
func StartEndpointService(port int) error {
	log.Info().Str("module", "endpoint").Msg("start endpoint service")

	s := grpc.NewServer()
	v2.RegisterEndpointServiceServer(s, &endpointService{bucket: bucket.NewBucket()})
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