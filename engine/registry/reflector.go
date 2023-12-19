package registry

import (
	"context"
	"fmt"

	"github.com/dsrvlabs/vatz/engine/bucket"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"github.com/rs/zerolog/log"
)

type PluginReflector interface {
	Query(context.Context, string, string, int) error
}

type reflectionHandler struct{}

func (r *reflectionHandler) Query(ctx context.Context, name, address string, port int) error {
	log.Info().Str("module", "registry").Msgf("query %s %s:%d", name, address, port)

	pluginAddress := fmt.Sprintf("%s:%d", address, port)
	cred := insecure.NewCredentials()
	conn, err := grpc.Dial(pluginAddress, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Err(err)
		return err
	}

	defer conn.Close()

	reflectCli := grpc_reflection_v1.NewServerReflectionClient(conn)
	reflectInfoCli, err := reflectCli.ServerReflectionInfo(ctx)

	listReq := grpc_reflection_v1.ServerReflectionRequest_ListServices{}
	reflectReq := grpc_reflection_v1.ServerReflectionRequest{
		Host:           pluginAddress,
		MessageRequest: &listReq,
	}

	err = reflectInfoCli.Send(&reflectReq)
	if err != nil {
		log.Err(err)
		return err
	}

	reflectResp, err := reflectInfoCli.Recv()
	if err != nil {
		log.Err(err)
		return err
	}

	listServiceResp := reflectResp.GetListServicesResponse()
	services := listServiceResp.GetService()

	var findService *grpc_reflection_v1.ServiceResponse
	for _, service := range services {
		if service.Name == name {
			findService = service
			break
		}
	}

	if findService == nil {
		return fmt.Errorf("cannot find %s", name)
	}

	log.Info().Str("module", "registry").Msgf("service found %s", findService.Name)

	// List functions
	grpcReflectCli := grpcreflect.NewClientAuto(ctx, conn)
	serviceDesc, err := grpcReflectCli.ResolveService(name)
	if err != nil {
		log.Err(err)
		return err
	}

	pDesc := bucket.PluginDescriptor{
		Address: pluginAddress,
		Name:    name,
		Methods: map[string]bucket.MethodArgDescriptor{},
	}

	methods := serviceDesc.GetMethods()
	for _, method := range methods {
		log.Debug().Str("module", "registry").Msgf("method %s", method.GetName())

		pDesc.Methods[method.GetName()] = bucket.MethodArgDescriptor{
			InDesc:  method.GetInputType(),
			OutDesc: method.GetOutputType(),
		}
	}

	b := bucket.NewBucket()
	err = b.Set(pDesc)

	return nil
}

func NewPluginReflector() PluginReflector {
	return &reflectionHandler{}
}
