package executor

import (
	"context"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type mockPluginClient struct {
	mock.Mock
}

func (c *mockPluginClient) Verify(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pluginpb.VerifyInfo, error) {
	// TODO: Not using yet.
	return nil, nil
}

func (c *mockPluginClient) Execute(ctx context.Context, in *pluginpb.ExecuteRequest, opts ...grpc.CallOption) (*pluginpb.ExecuteResponse, error) {
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
