package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/vatz-proto/manager/v2"
)

func TestHandlerStartStop(t *testing.T) {
	h := NewHandler()

	ctx := context.Background()

	req := &v2.UserRequest{
		Plugin: "snippet.grpc.reflection.HelloService",
		Method: "Hello",
		Fields: []*v2.FieldSpec{
			{
				Name:  "name",
				Type:  "string",
				Value: "rootwarp",
			},
			{
				Name:  "age",
				Type:  "int32",
				Value: "40",
			},
		},
	}

	resp, err := h.SendRequest(ctx, req)
	assert.Nil(t, err)

	fmt.Println(resp)
}
