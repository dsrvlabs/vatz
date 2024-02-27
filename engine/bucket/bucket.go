package bucket

import (
	"errors"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/rs/zerolog/log"
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
	List() []*PluginDescriptor
}

type pluginBucket struct {
	plugins map[string]PluginDescriptor
	lock    sync.Mutex
}

func (b *pluginBucket) Set(newPlugin PluginDescriptor) error {
	log.Info().Str("module", "bucket").Msgf("set %s", newPlugin.Name)

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
	log.Info().Str("module", "bucket").Msgf("get %s", name)

	p, ok := b.plugins[name]
	if !ok {
		return nil, errNotExist
	}

	return &p, nil
}

func (b *pluginBucket) List() []*PluginDescriptor {
	descs := []*PluginDescriptor{}

	for k := range b.plugins {
		desc := b.plugins[k]
		descs = append(descs, &desc)
	}

	return descs
}

func NewBucket() PluginBucket {
	c := make(chan bool, 1)

	once.Do(func() {
		log.Info().Str("module", "bucket").Msg("new bucket")

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
