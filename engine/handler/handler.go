package handler

import "context"

// InvokeRequest provides detail arguments to call plugin function.
// TODO: Detail properties are not fixed yet.
type InvokeRequest struct {
	PluginAddress string

	// Warning: Below properties are TBD
	PluginName string
	Function   string
	Args       any
}

// InvokeResponse contains results of InvokeRequest.
type InvokeResponse struct {
	// Warning: TBD
	Error error
}

// RequestHandler provides interface to send plugin call request form relay.
type RequestHandler interface {
	Start() error
	Stop() error

	Request(request InvokeRequest) InvokeResponse
}

type basicHandler struct {
	workers     []handlerWorker
	chanEnqueue chan InvokeRequest
	chanResult  chan InvokeResponse
}

func (h *basicHandler) Start() error {
	h.workers = make([]handlerWorker, 10)
	for i := 0; i < len(h.workers); i++ {
		ctx, cancel := context.WithCancel(context.Background())

		h.workers[i] = handlerWorker{
			ctx:    ctx,
			cancel: cancel,
		}

		go h.workers[i].run(h.chanEnqueue, h.chanResult)
	}

	return nil
}

func (h *basicHandler) Stop() error {
	for _, worker := range h.workers {
		worker.stop()
	}

	return nil
}

func (h *basicHandler) Request(request InvokeRequest) InvokeResponse {
	// TODO:
	// enqueue to worker.
	// get response
	// check timeout
	// return
	return InvokeResponse{}
}

type handlerWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *handlerWorker) request(request InvokeRequest) InvokeResponse {
	// TODO: Handler request
	return InvokeResponse{}
}

func (w *handlerWorker) run(in <-chan InvokeRequest, out chan<- InvokeResponse) error {
	for {
		select {
		case req := <-in:
			resp := w.request(req)
			out <- resp
		case <-w.ctx.Done():
			return nil
		}
	}
}

func (w *handlerWorker) stop() error {
	w.cancel()
	return nil
}

// NewHandler creates a new RequestHandler instance.
func NewHandler() RequestHandler {
	handler := &basicHandler{}

	return handler
}
