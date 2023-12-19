package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketAddNewError(t *testing.T) {
	tests := []struct {
		Desc      string
		InDesc    PluginDescriptor
		ExpectErr error
	}{
		{
			Desc:      "Empty fields are not allowed",
			InDesc:    PluginDescriptor{},
			ExpectErr: errInvalidField,
		},
		{
			Desc: "Empty method is not allowed",
			InDesc: PluginDescriptor{
				Address: "dummy",
				Name:    "dummy",
			},
			ExpectErr: errInvalidMethod,
		},
	}

	b := NewBucket()

	for _, test := range tests {
		err := b.Set(test.InDesc)
		assert.Equal(t, err, test.ExpectErr)
	}
}

func TestBucketAddNew(t *testing.T) {
	b := NewBucket()

	p := PluginDescriptor{
		Address: "localhost:9090",
		Name:    "HelloService",
		Methods: map[string]MethodArgDescriptor{
			"SayHello": {},
			"Greeting": {},
		},
	}

	// Set first
	err := b.Set(p)

	assert.Nil(t, err)

	// Then, Get
	pluginFromBucket, err := b.Get("HelloService")

	assert.Nil(t, err)
	assert.NotNil(t, pluginFromBucket)
	assert.Equal(t, p.Address, pluginFromBucket.Address)
	assert.Equal(t, p.Name, pluginFromBucket.Name)

	// Check existing method
	m, err := pluginFromBucket.GetMethod("SayHello")

	assert.Nil(t, err)
	assert.NotNil(t, m)

	m, err = pluginFromBucket.GetMethod("Greeting")

	assert.Nil(t, err)
	assert.NotNil(t, m)

	// check non-existing method
	m, err = pluginFromBucket.GetMethod("NonExistMethod")

	assert.Nil(t, m)
	assert.Equal(t, err, errNotExist)
}

func TestBucketSingleton(t *testing.T) {
	b1 := NewBucket()
	b2 := NewBucket()

	assert.Equal(t, b1, b2)
}
