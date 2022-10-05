package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func TestRPCs(t *testing.T) {
	rpc := NewRPCService()

	go rpc.Start("127.0.0.1", 19090, 19091)

	time.Sleep(time.Second * 1) // Wait ready

	go func() {
		time.Sleep(time.Millisecond * 100)

		fmt.Println("Call Stop")
		rpc.Stop()
	}()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:19091/v1/plugin_status", nil)
	assert.Nil(t, err)

	cli := http.Client{}
	resp, err := cli.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)

	assert.Nil(t, err)

	respData := map[string]any{}
	err = json.Unmarshal(data, &respData)

	assert.Nil(t, err)
	assert.Contains(t, respData, "status")
	assert.Contains(t, respData, "pluginStatus")
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
