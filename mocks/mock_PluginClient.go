package mocks

import (
	"context"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockPluginClient is mock object for grpc client of VATZ.
type MockPluginClient struct {
	mock.Mock
}

func (c *MockPluginClient) Verify(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pluginpb.VerifyInfo, error) {
	// TODO: Not using yet.
	ret := c.Called(ctx, in, opts)

	var r0 *pluginpb.VerifyInfo
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, []grpc.CallOption) (*pluginpb.VerifyInfo, error)); ok {
		r0, r1 = rf(ctx, in, opts)
	} else {
		r0 = ret.Get(0).(*pluginpb.VerifyInfo)
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (c *MockPluginClient) Execute(ctx context.Context, in *pluginpb.ExecuteRequest, opts ...grpc.CallOption) (*pluginpb.ExecuteResponse, error) {
	ret := c.Called(ctx, in, opts)

	var r0 *pluginpb.ExecuteResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pluginpb.ExecuteRequest, []grpc.CallOption) (*pluginpb.ExecuteResponse, error)); ok {
		r0, r1 = rf(ctx, in, opts)
	} else {
		r0 = ret.Get(0).(*pluginpb.ExecuteResponse)
		r1 = ret.Error(1)
	}
	return r0, r1
}
