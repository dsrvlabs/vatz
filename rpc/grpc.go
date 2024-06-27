package rpc

import (
	"context"

	pb "github.com/dsrvlabs/vatz-proto/rpc/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/dsrvlabs/vatz/manager/healthcheck"
	tp "github.com/dsrvlabs/vatz/types"
)

type grpcService struct {
	pb.UnimplementedVatzRPCServer

	healthChecker healthcheck.HealthCheck
}

func (s *grpcService) PluginStatus(ctx context.Context, in *emptypb.Empty) (*pb.PluginStatusResponse, error) {
	pluginStatus := s.healthChecker.PluginStatus(ctx)

	respStatus := pb.PluginStatusResponse{
		Status:       pb.Status_OK,
		PluginStatus: make([]*pb.PluginStatus, len(pluginStatus)),
	}

	for i, status := range pluginStatus {
		newStatus := pb.PluginStatus{PluginName: status.Plugin.Name}

		if status.IsAlive == tp.AliveStatusUp {
			newStatus.Status = pb.Status_OK
		} else {
			newStatus.Status = pb.Status_FAIL
		}

		respStatus.PluginStatus[i] = &newStatus
	}

	return &respStatus, nil
}
