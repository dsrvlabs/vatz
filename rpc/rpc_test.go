package rpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/dsrvlabs/vatz-proto/rpc/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"github.com/dsrvlabs/vatz/mocks"
)

func TestTerminateRPC(t *testing.T) {
	rpc := NewRPCService()

	go func() {
		time.Sleep(time.Millisecond * 100)
		rpc.Stop()
	}()

	err := rpc.Start()

	assert.Nil(t, err)
}

func TestPluginStatus(t *testing.T) {
	ctx := context.Background()

	// Mocks
	mockHealthCheck := mocks.HealthCheck{}

	mockHealthCheck.
		On("PluginStatus", ctx).
		Return([]tp.PluginStatus{
			{
				Plugin:  config.Plugin{Name: "plugin1"},
				IsAlive: tp.AliveStatusUp,
			},
			{
				Plugin:  config.Plugin{Name: "plugin2"},
				IsAlive: tp.AliveStatusDown,
			},
		})

	// Tests
	srv := grpcService{healthChecker: &mockHealthCheck}
	resp, err := srv.PluginStatus(ctx, &emptypb.Empty{})

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.PluginStatus))

	assert.Equal(t, "plugin1", resp.PluginStatus[0].PluginName)
	assert.Equal(t, pb.Status_OK, resp.PluginStatus[0].GetStatus())

	assert.Equal(t, "plugin2", resp.PluginStatus[1].PluginName)
	assert.Equal(t, pb.Status_FAIL, resp.PluginStatus[1].GetStatus())

	mockHealthCheck.AssertExpectations(t)
}
