package healthcheck

import (
	"errors"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"testing"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	dp "github.com/dsrvlabs/vatz/manager/dispatcher"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/dsrvlabs/vatz/mocks"
)

func TestPluginHealthCheckSuccess(t *testing.T) {
	h := healthChecker{}
	ctx := context.Background()

	// Mock
	mockPluginCli := mocks.MockPluginClient{}
	mockPluginCli.
		On("Verify", ctx, new(emptypb.Empty), []grpc.CallOption(nil)).
		Return(&pluginpb.VerifyInfo{VerifyMsg: "test"}, nil)

	mockDispatcher := dp.MockNotification{}

	// Test
	status, err := h.PluginHealthCheck(ctx, &mockPluginCli, config.Plugin{}, &mockDispatcher)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, tp.AliveStatusUp, status)

	mockPluginCli.AssertExpectations(t)
	mockDispatcher.AssertExpectations(t)
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
		h := healthChecker{}
		ctx := context.Background()

		// Mock
		mockPluginCli := mocks.MockPluginClient{}
		mockPluginCli.
			On("Verify", ctx, new(emptypb.Empty), []grpc.CallOption(nil)).
			Return(test.MockVerifyInfo, test.MockVerifyErr)

		mockDispatcher := dp.MockNotification{}
		mockJSONMsg := tp.ReqMsg{
			FuncName:     "isPluginUp",
			State:        pluginpb.STATE_FAILURE,
			Msg:          "Plugin is DOWN!!",
			Severity:     pluginpb.SEVERITY_CRITICAL,
			ResourceType: "test",
		}
		mockDispatcher.On("SendNotification", mockJSONMsg).Return(nil)

		// Test
		status, err := h.PluginHealthCheck(ctx, &mockPluginCli, config.Plugin{Name: "test"}, &mockDispatcher)

		// Asserts
		assert.Nil(t, err)
		assert.Equal(t, tp.AliveStatusDown, status)

		mockPluginCli.AssertExpectations(t)
		mockDispatcher.AssertExpectations(t)
	}

}

func TestVatzHealthCheck(t *testing.T) {
	// TODO: TBD
}
