package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dsrvlabs/vatz-proto/manager/v2"
	"github.com/dsrvlabs/vatz/engine/bucket"
	"github.com/jhump/protoreflect/desc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type handlerService struct {
	v2.UnimplementedRequestHandlerServer
}

func (s *handlerService) SendRequest(ctx context.Context, in *v2.UserRequest) (*v2.UserResponse, error) {
	log.Info().Str("module", "handler").Msgf("send request %s:%s", in.GetPlugin(), in.GetMethod())

	// Find relevant proto message.
	b := bucket.NewBucket()
	pDesc, err := b.Get(in.GetPlugin())
	if err != nil {
		return nil, err
	}

	mDesc, err := pDesc.GetMethod(in.GetMethod())
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", pDesc)

	// Convert into Protobuf message.
	inMsg, err := s.convert(mDesc.InDesc, in.Fields)
	outMsg := dynamicpb.NewMessage(mDesc.OutDesc.UnwrapMessage())

	// Invoke
	cred := insecure.NewCredentials()
	conn, err := grpc.Dial(pDesc.Address, grpc.WithTransportCredentials(cred))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer conn.Close()

	funcURL := fmt.Sprintf("%s/%s", pDesc.Name, in.GetMethod())
	err = conn.Invoke(ctx, funcURL, inMsg, outMsg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// TODO: How to serialize?
	// TODO: And How to parse this from client side?
	// fmt.Printf("resp %+v\n", outMsg)

	d, err := protojson.Marshal(outMsg)

	// fmt.Printf("Marshal %s, %+v\n", string(d), err)

	anyMsg, err := anypb.New(outMsg)
	if err != nil {
		return nil, err
	}

	resp := v2.UserResponse{
		Plugin:    in.Plugin,
		Method:    in.Method,
		Result:    anyMsg,
		StrResult: string(d),
	}
	return &resp, nil
}

func (s *handlerService) convert(d *desc.MessageDescriptor, fields []*v2.FieldSpec) (*dynamicpb.Message, error) {
	// TODO: Define return type

	newMsg := dynamicpb.NewMessage(d.UnwrapMessage())

	for _, spec := range fields {
		fDesc := d.FindFieldByName(spec.Name).UnwrapField()

		if spec.Type == "string" {
			newMsg.Set(fDesc, protoreflect.ValueOfString(spec.Value))
		} else if spec.Type == "int32" {
			intValue, err := strconv.ParseInt(spec.Value, 10, 64)
			if err != nil {
				return nil, err
			}

			newMsg.Set(fDesc, protoreflect.ValueOfInt32(int32(intValue)))
		} else {
			return nil, errors.New("invalid type")
		}
	}

	return newMsg, nil
}

func NewHandler() v2.RequestHandlerServer {
	return &handlerService{}
}
