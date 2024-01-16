package endpoint

import (
	"context"
	"fmt"
	"testing"

	"github.com/dsrvlabs/vatz/engine/bucket"

	"github.com/dsrvlabs/vatz-proto/manager/v2"
	"github.com/stretchr/testify/assert"
)

func init() {
	b := bucket.NewBucket()

	_ = b.Set(
		bucket.PluginDescriptor{
			Address: "localhost:9090",
			Name:    "service-1",
			Methods: map[string]bucket.MethodArgDescriptor{
				"SayHello": {},
				"Greeting": {},
			},
		},
	)

	_ = b.Set(
		bucket.PluginDescriptor{
			Address: "localhost:9091",
			Name:    "service-2",
			Methods: map[string]bucket.MethodArgDescriptor{
				"SayHello": {},
				"Greeting": {},
			},
		},
	)
}

func TestEndpointListPlugin(t *testing.T) {
	s := endpointService{bucket: bucket.NewBucket()}

	ctx := context.Background()
	req := v2.ListPluginRequest{}
	resp, err := s.ListPlugin(ctx, &req)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.Plugin))

	metaNames := make([]string, len(resp.Plugin))
	for i, p := range resp.Plugin {
		metaNames[i] = p.Name
	}

	assert.Contains(t, metaNames, "service-1")
	assert.Contains(t, metaNames, "service-2")
	assert.NotContains(t, metaNames, "non-exist")
}

func TestEndpointDetailPlugin(t *testing.T) {
	s := endpointService{bucket: bucket.NewBucket()}

	ctx := context.Background()
	req := v2.DetailPluginRequest{
		PluginName: "service-1",
	}

	resp, err := s.DetailPlugin(ctx, &req)

	assert.Nil(t, err)

	fmt.Println(resp)
}
