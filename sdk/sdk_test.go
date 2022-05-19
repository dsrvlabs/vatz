package sdk

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestStartStop(t *testing.T) {
	ctx := context.Background()

	p := NewPlugin("unittest")

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
	tests := []struct {
		MockResp CallResponse
		MockErr  error
	}{
		{
			MockResp: CallResponse{
				State:   pb.STATE_SUCCESS,
				Message: "hello world",
			},
			MockErr: nil,
		},
		{
			MockResp: CallResponse{
				State: pb.STATE_FAILURE,
			},
			MockErr: errors.New("dummy error"),
		},
	}

	ctx := context.Background()

	p := plugin{}

	for _, test := range tests {
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
			Return(test.MockResp, test.MockErr)

		// Test
		err := p.Register(mockCallback.DummyCall1)
		resp, err := p.grpc.Execute(ctx, &req)

		// Asserts
		assert.Nil(t, err)
		if test.MockErr == nil {
			assert.Equal(t, test.MockResp.State, resp.State)
			assert.Equal(t, test.MockResp.Message, resp.Message)
		} else {
			assert.Equal(t, test.MockResp.State, resp.State)
			assert.Equal(t, test.MockErr.Error(), resp.Message)
		}

		mockCallback.AssertExpectations(t)
	}
}
