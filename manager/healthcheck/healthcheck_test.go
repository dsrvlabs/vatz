package healthcheck

import (
	"errors"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	"testing"

	tp "github.com/dsrvlabs/vatz/manager/types"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/dsrvlabs/vatz/mocks"
)

func TestPluginHealthCheckSuccess(t *testing.T) {
	h := healthChecker{
		pluginStatus: map[string]tp.PluginStatus{},
	}
	ctx := context.Background()

	// Mock
	mockPluginCli := mocks.MockPluginClient{}
	mockPluginCli.
		On("Verify", ctx, new(emptypb.Empty), []grpc.CallOption(nil)).
		Return(&pluginpb.VerifyInfo{VerifyMsg: "test"}, nil)

	var mockDispatchers = []dp.Dispatcher{}
	status, err := h.PluginHealthCheck(ctx, &mockPluginCli, config.Plugin{Name: "dummy"}, mockDispatchers)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, tp.AliveStatusUp, status)

	statuses := h.PluginStatus(ctx)

	assert.Equal(t, 1, len(statuses))
	assert.Equal(t, "dummy", statuses[0].Plugin.Name)
	assert.Equal(t, tp.AliveStatusUp, statuses[0].IsAlive)

	mockPluginCli.AssertExpectations(t)
}

func TestPluginHealthCheckFailed(t *testing.T) {
	tests := []struct {
		Desc           string
		MockVerifyInfo *pluginpb.VerifyInfo
		MockVerifyErr  error
	}{
		{
			Desc:           "VerifyInfo is nil",
			MockVerifyInfo: nil,
			MockVerifyErr:  nil,
		},
		{
			Desc:           "Error VerifyInfo",
			MockVerifyInfo: &pluginpb.VerifyInfo{},
			MockVerifyErr:  errors.New("temporal error occurred"),
		},
	}

	for _, test := range tests {
		h := healthChecker{
			pluginStatus: map[string]tp.PluginStatus{},
		}
		ctx := context.Background()

		// Mock
		mockPluginCli := mocks.MockPluginClient{}
		mockPluginCli.
			On("Verify", ctx, new(emptypb.Empty), []grpc.CallOption(nil)).
			Return(test.MockVerifyInfo, test.MockVerifyErr)

		mockDispatcher := dp.MockDispatcher{}
		var mockDispatchers []dp.Dispatcher

		mockJSONMsg := tp.ReqMsg{
			FuncName:     "isPluginUp",
			State:        pluginpb.STATE_FAILURE,
			Msg:          "Plugin is DOWN!!",
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: "test",
		}
		mockDispatcher.On("SendNotification", mockJSONMsg).Return(nil)

		// Test
		status, err := h.PluginHealthCheck(ctx, &mockPluginCli, config.Plugin{Name: "test"}, mockDispatchers)

		// Asserts
		assert.Nil(t, err)
		assert.Equal(t, tp.AliveStatusDown, status)

		statuses := h.PluginStatus(ctx)

		assert.Equal(t, 1, len(statuses))
		assert.Equal(t, "test", statuses[0].Plugin.Name)
		assert.Equal(t, tp.AliveStatusDown, statuses[0].IsAlive)

		mockPluginCli.AssertExpectations(t)
	}
}

func TestVatzHealthCheck(t *testing.T) {
	// TODO: TBD
}
