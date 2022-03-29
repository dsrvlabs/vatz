package executor

import (
	"context"
	pluginpb "github.com/hqueue/vatz-secret/dist/proto/vatz/plugin/v1"
)

type executor struct {
}

type Executor interface {
	Execute(ctx context.Context, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error)
}

func (v executor) Execute(ctx context.Context, in *pluginpb.ExecuteRequest) (*pluginpb.ExecuteResponse, error) {
	return nil, nil
}

func NewExecutor() Executor {
	return &executor{}
}
