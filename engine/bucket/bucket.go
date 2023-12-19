package bucket

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jhump/protoreflect/desc"
)

var (
	errInvalidField  = errors.New("empty fields are not allowed")
	errInvalidMethod = errors.New("empty method spec")
	errNotExist      = errors.New("cannot find the plugin")
)

var (
	once   sync.Once
	bucket *pluginBucket
)

// PluginDescriptor contains plugin server specs.
type PluginDescriptor struct {
	Address string
	Name    string
	Methods map[string]MethodArgDescriptor
}

func (d PluginDescriptor) GetMethod(methodName string) (*MethodArgDescriptor, error) {
	m, ok := d.Methods[methodName]
	if !ok {
		return nil, errNotExist
	}

	return &m, nil
}

// MethodArgDescriptor contains reflection descriptor to create raw format message.
type MethodArgDescriptor struct {
	InDesc  *desc.MessageDescriptor
	OutDesc *desc.MessageDescriptor
}

// PluginBucket is a storage of plugin's metadata.
type PluginBucket interface {
	Set(newPlugin PluginDescriptor) error
	Get(name string) (*PluginDescriptor, error)
}

type pluginBucket struct {
	plugins map[string]PluginDescriptor
	lock    sync.Mutex
}

func (b *pluginBucket) Set(newPlugin PluginDescriptor) error {
	fmt.Println("Set")

	if newPlugin.Address == "" || newPlugin.Name == "" {
		return errInvalidField
	}

	if len(newPlugin.Methods) == 0 {
		return errInvalidMethod
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	b.plugins[newPlugin.Name] = newPlugin

	return nil
}

func (b *pluginBucket) Get(name string) (*PluginDescriptor, error) {
	fmt.Println("Get")

	p, ok := b.plugins[name]
	if !ok {
		return nil, errNotExist
	}

	return &p, nil
}

func NewBucket() PluginBucket {
	c := make(chan bool, 1)

	once.Do(func() {
		bucket = &pluginBucket{
			plugins: map[string]PluginDescriptor{},
			lock:    sync.Mutex{},
		}

		c <- true
	})

	if bucket == nil {
		_ = <-c
	}

	return bucket
}
