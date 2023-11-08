package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	//"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	//"google.golang.org/protobuf/encoding/protojson"
	//"google.golang.org/protobuf/reflect/protodesc"
	// "google.golang.org/protobuf/types/descriptorpb"

	// "github.com/rootwarp/snippets/golang/grpc/reflection/proto/agent"
	agent "github.com/dsrvlabs/vatz-proto/manager/v2"
)

var (
	ErrServiceNotFound = errors.New("cannot find service")
)

type reflectionHandler struct {
	ServiceSpecs map[string]serviceMeta
}

type serviceMeta struct {
	Address   string
	Functions map[string]funcMeta
}

type funcMeta struct {
	InDesc  *desc.MessageDescriptor
	OutDesc *desc.MessageDescriptor
}

func (r *reflectionHandler) Query(ctx context.Context, name, address string, port int) error {
	fmt.Println("Query")

	host := fmt.Sprintf("%s:%d", address, port)
	cred := insecure.NewCredentials()
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(cred))
	if err != nil {
		return err
	}

	defer conn.Close()

	reflectCli := grpc_reflection_v1.NewServerReflectionClient(conn)
	reflectInfoCli, err := reflectCli.ServerReflectionInfo(ctx)

	listReq := grpc_reflection_v1.ServerReflectionRequest_ListServices{}
	reflectReq := grpc_reflection_v1.ServerReflectionRequest{
		Host:           host,
		MessageRequest: &listReq,
	}

	err = reflectInfoCli.Send(&reflectReq)
	if err != nil {
		return err
	}

	reflectResp, err := reflectInfoCli.Recv()
	if err != nil {
		return err
	}

	listServiceResp := reflectResp.GetListServicesResponse()
	services := listServiceResp.GetService()

	var findService *grpc_reflection_v1.ServiceResponse
	for _, service := range services {
		fmt.Println("service", service.Name)
		if service.Name == name {
			findService = service
		}
	}

	if findService == nil {
		return fmt.Errorf("cannot find %s", name)
	}

	fmt.Println("found", findService)

	// List functions
	grpcReflectCli := grpcreflect.NewClientAuto(ctx, conn)
	serviceDesc, err := grpcReflectCli.ResolveService(name)
	if err != nil {
		return err
	}

	methods := serviceDesc.GetMethods()

	r.ServiceSpecs[name] = serviceMeta{
		Address:   fmt.Sprintf("%s:%d", address, port),
		Functions: map[string]funcMeta{},
	}

	for _, method := range methods {
		fmt.Println("*****", method.GetName())

		in := method.GetInputType()
		out := method.GetOutputType()

		meta := funcMeta{
			InDesc:  in,
			OutDesc: out,
		}
		r.ServiceSpecs[name].Functions[method.GetName()] = meta
	}

	return nil
}

func (s *reflectionHandler) Invoke(ctx context.Context, serviceName, funcName string) (any, error) {
	fmt.Println("Invoke")

	service, ok := r.ServiceSpecs[serviceName]
	if !ok {
		return nil, ErrServiceNotFound
	}

	fMeta := service.Functions[funcName]
	fmt.Println(fMeta)

	// TODO: We should fill here dynamically
	inSpec := fMeta.InDesc

	newInMsg := dynamicpb.NewMessage(inSpec.UnwrapMessage())

	nameField := inSpec.FindFieldByName("name").UnwrapField()
	ageField := inSpec.FindFieldByName("age").UnwrapField()
	newInMsg.Set(nameField, protoreflect.ValueOfString("rootwarp"))
	newInMsg.Set(ageField, protoreflect.ValueOfInt32(40))

	outSpec := fMeta.OutDesc
	newOutMsg := dynamicpb.NewMessage(outSpec.UnwrapMessage())

	// Invoke
	cred := insecure.NewCredentials()
	conn, err := grpc.Dial(service.Address, grpc.WithTransportCredentials(cred))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer conn.Close()

	funcURL := fmt.Sprintf("%s/%s", serviceName, funcName)
	err = conn.Invoke(ctx, funcURL, newInMsg, newOutMsg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newOutMsg, nil
}

type registrationServer struct {
	agent.UnimplementedRegistrationServiceServer
}

func (s *registrationServer) RegisterPlugin(ctx context.Context, in *agent.RegisterRequest) (*agent.RegisterResponse, error) {
	fmt.Println("Register", in.Name, in.Address, in.Port)

	err := r.Query(ctx, in.Name, in.Address, int(in.Port))
	if err != nil {
		return nil, err
	}

	// TODO: Testing purpose.
	r.Invoke(ctx, in.Name, "Hello")

	return &agent.RegisterResponse{
		Msg: fmt.Sprintf("%s - %s:%d", in.Name, in.Address, in.Port),
	}, nil
}

var r *reflectionHandler

func main() {
	fmt.Println("Start server")

	go func() {
		r = &reflectionHandler{
			ServiceSpecs: map[string]serviceMeta{},
		}
	}()

	s := grpc.NewServer()
	agent.RegisterRegistrationServiceServer(s, &registrationServer{})
	reflection.Register(s)

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
