package relay

import (
	"github.com/dsrvlabs/vatz/engine/handler"
)

type RelayRequest struct{}

type RelayResponse struct{}

type RequestRelay interface {
	Start() error
	Stop() error

	Request() error
}

type basicRelay struct{}

func (r *basicRelay) Start() error {
	return nil
}

func (r *basicRelay) Stop() error {
	return nil
}

// TODO:
// Need arguments
// How to parallize.
func (r *basicRelay) Request() error {
	info, err := r.fetchPluginInfo()
	if err != nil {
		return err
	}
	_ = info

	// TODO: Send received request to handler
	h := handler.NewHandler() // TODO: Hmm..shoudl be singleton?

	invokeRequest := handler.InvokeRequest{} // TODO: Properties should be filled.

	resp := h.Request(invokeRequest)

	_ = resp
	//

	return nil
}

// TODO: the structure of plugininfo should be defined at "plugin service bucket"
func (r *basicRelay) fetchPluginInfo() (any, error) {
	return nil, nil
}

func NewRelay() RequestRelay {
	r := &basicRelay{}
	return r
}
