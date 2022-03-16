package sdk

import (
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/rootwarp/vatz-plugin-sdk/plugin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestStartStop(t *testing.T) {
	ctx := context.Background()

	p := NewPlugin()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		log.Println("Start")
		_ = p.Start(ctx, "0.0.0.0", 9091)

		log.Println("Bye")
		wg.Done()
	}()

	time.Sleep(time.Second * 1)

	p.Stop()

	wg.Wait()
}

func TestVerify(t *testing.T) {
	p := plugin{}

	ctx := context.Background()
	info, err := p.grpc.Verify(ctx, &emptypb.Empty{})

	assert.Nil(t, err)
	assert.Equal(t, "OK", info.VerifyMsg)
}

func TestInvokeCallback(t *testing.T) {
	ctx := context.Background()

	p := plugin{}

	req := pb.ExecuteRequest{
		ExecuteInfo: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"function": structpb.NewStringValue("testfunc"),
			},
		},
	}

	mockCallback := mockFuncs{}

	mockCallback.
		On("DummyCall1", req.GetExecuteInfo().GetFields(), req.GetOptions().GetFields()).
		Return(nil)

	// Test
	err := p.Register(mockCallback.DummyCall1)
	resp, err := p.grpc.Execute(ctx, &req)

	assert.Nil(t, err)
	assert.Equal(t, pb.ExecuteResponse_SUCCESS, resp.State)

	mockCallback.AssertExpectations(t)
}
