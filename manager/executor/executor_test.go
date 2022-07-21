package executor

import (
	"context"
	"fmt"
	"testing"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/manager/config"
	notif "github.com/dsrvlabs/vatz/manager/notification"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestExecutorSuccess(t *testing.T) {
	const (
		testMethodName = "is_up"
		testPluginName = "unittest_plugin"
	)

	tests := []struct {
		Desc          string
		TestExecResp  *pluginpb.ExecuteResponse
		TestNotifInfo notif.NotifyInfo
	}{
		{
			Desc: "No Alert",
			TestExecResp: &pluginpb.ExecuteResponse{
				State:    pluginpb.STATE_SUCCESS,
				Severity: pluginpb.SEVERITY_UNKNOWN,
			},
			TestNotifInfo: notif.NotifyInfo{
				Plugin:   testPluginName,
				Method:   testMethodName,
				State:    pluginpb.STATE_SUCCESS,
				Severity: pluginpb.SEVERITY_UNKNOWN,
			},
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		cfgPlugin := config.Plugin{
			Name: testPluginName,
			ExecutableMethods: []struct {
				Name string `yaml:"method_name"`
			}{
				{testMethodName},
			},
		}

		// Mocks.
		mockExeInfo, err := structpb.NewStruct(map[string]interface{}{
			"execute_method": testMethodName,
		})

		mockOpts, err := structpb.NewStruct(map[string]interface{}{
			"plugin_name": testPluginName,
		})

		exeReq := pluginpb.ExecuteRequest{
			ExecuteInfo: mockExeInfo,
			Options:     mockOpts,
		}

		mockClient := mockPluginClient{}
		mockClient.On("Execute", ctx, &exeReq, []grpc.CallOption(nil)).Return(test.TestExecResp, nil)

		mockNotif := notif.MockNotification{}
		dispatchManager = &mockNotif

		if test.TestNotifInfo.State != pluginpb.STATE_SUCCESS {
			dummyMsg := notif.ReqMsg{
				FuncName:     testMethodName,
				State:        pluginpb.STATE_FAILURE,
				Msg:          "No response from Plugin",
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: testPluginName,
			}
			mockNotif.On("SendNotification", dummyMsg).Return(nil)
		}

		// Test
		e := executor{
			status: map[string]bool{},
		}

		err = e.Execute(ctx, &mockClient, cfgPlugin)

		fmt.Println("Status", e.status)

		// Asserts
		mockClient.AssertExpectations(t)
		mockNotif.AssertExpectations(t)

		assert.Nil(t, err)
		assert.True(t, e.status[testMethodName])
	}
}

func TestExecutorFailure(t *testing.T) {
	const (
		testMethodName = "is_up"
		testPluginName = "unittest_plugin"
	)

	tests := []struct {
		Desc           string
		MockPrevStatus bool
		TestExecResp   *pluginpb.ExecuteResponse
		TestNotifInfo  notif.NotifyInfo
		ExpectReqMsg   notif.ReqMsg
	}{
		{
			Desc:           "Alert ERROR",
			MockPrevStatus: false,
			TestExecResp: &pluginpb.ExecuteResponse{
				State:    pluginpb.STATE_FAILURE,
				Severity: pluginpb.SEVERITY_ERROR,
			},
			TestNotifInfo: notif.NotifyInfo{
				Plugin:   testPluginName,
				Method:   testMethodName,
				State:    pluginpb.STATE_FAILURE,
				Severity: pluginpb.SEVERITY_ERROR,
			},
			ExpectReqMsg: notif.ReqMsg{
				FuncName:     testMethodName,
				State:        pluginpb.STATE_FAILURE,
				Msg:          "No response from Plugin",
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: testPluginName,
			},
		},
		{
			Desc:           "Alert Critical",
			MockPrevStatus: false,
			TestExecResp: &pluginpb.ExecuteResponse{
				State:    pluginpb.STATE_FAILURE,
				Severity: pluginpb.SEVERITY_CRITICAL,
				Message:  "test execute msg",
			},
			TestNotifInfo: notif.NotifyInfo{
				Plugin:     testPluginName,
				Method:     testMethodName,
				State:      pluginpb.STATE_FAILURE,
				Severity:   pluginpb.SEVERITY_CRITICAL,
				ExecuteMsg: "test execute msg",
			},
			ExpectReqMsg: notif.ReqMsg{
				FuncName:     testMethodName,
				State:        pluginpb.STATE_FAILURE,
				Msg:          "test execute msg",
				Severity:     pluginpb.SEVERITY_CRITICAL,
				ResourceType: testPluginName,
			},
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		cfgPlugin := config.Plugin{
			Name: testPluginName,
			ExecutableMethods: []struct {
				Name string `yaml:"method_name"`
			}{
				{testMethodName},
			},
		}

		// Mocks.
		mockExeInfo, err := structpb.NewStruct(map[string]interface{}{
			"execute_method": testMethodName,
		})

		mockOpts, err := structpb.NewStruct(map[string]interface{}{
			"plugin_name": testPluginName,
		})

		exeReq := pluginpb.ExecuteRequest{
			ExecuteInfo: mockExeInfo,
			Options:     mockOpts,
		}

		mockClient := mockPluginClient{}
		mockClient.On("Execute", ctx, &exeReq, []grpc.CallOption(nil)).Return(test.TestExecResp, nil)

		mockNotif := notif.MockNotification{}
		dispatchManager = &mockNotif

		mockNotif.On("SendNotification", test.ExpectReqMsg).Return(nil)

		// Test
		e := executor{
			status:    map[string]bool{},
			isSending: map[string]bool{},
		}

		err = e.Execute(ctx, &mockClient, cfgPlugin)

		fmt.Println("Status", e.status)

		// Asserts
		mockClient.AssertExpectations(t)
		mockNotif.AssertExpectations(t)

		assert.Nil(t, err)
		assert.False(t, e.status[testMethodName])
	}
}
