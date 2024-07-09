package gcp

import (
	"cloud.google.com/go/logging"
	"context"
	tp "github.com/dsrvlabs/vatz/types"
	"google.golang.org/api/option"
)

type GCPLoggingClient struct {
	client *logging.Client
	logger *logging.Logger
}

func NewGCPLoggingClient(ctx context.Context, projectID string, ClientOption option.ClientOption) (*GCPLoggingClient, error) {
	var client *logging.Client

	if ClientOption == nil {
		noOptionClient, err := logging.NewClient(ctx, projectID)
		if err != nil {
			return nil, err
		}
		client = noOptionClient
	} else {

		withOptionClient, err := logging.NewClient(ctx, projectID, ClientOption)
		if err != nil {
			return nil, err
		}
		client = withOptionClient
	}

	return &GCPLoggingClient{
		client: client,
		logger: client.Logger(tp.MonitoringIdentifier),
	}, nil
}

func (c *GCPLoggingClient) Log(entry LogEntry) error {
	c.logger.Log(logging.Entry{
		Payload:   entry.Payload,
		Severity:  entry.Severity,
		Labels:    entry.Labels,
		Timestamp: entry.Timestamp,
	})
	return nil
}

func (c *GCPLoggingClient) Close() error {
	return c.client.Close()
}
