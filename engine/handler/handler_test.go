package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerStartStop(t *testing.T) {
	h := NewHandler()

	err := h.Start()
	assert.Nil(t, err)

	err = h.Stop()
	assert.Nil(t, err)
}
