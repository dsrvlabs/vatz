package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := GetVersion()

	assert.Equal(t, "dev-N/A", ver)
}
