package handlers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"sync"
)

type Handler interface {
	Handle(*sdk.Address, []byte) bool // Is subscription still necessary ?
}

type cosignatureHandler struct {
	messageProcessor sdk.CosignatureProcessor
	handlers         subscribers.Cosignature
	errCh            chan<- error
}

func NewCosignatureHandler(messageProcessor sdk.CosignatureProcessor, handlers subscribers.Cosignature, errCh chan<- error) *cosignatureHandler {
	return &cosignatureHandler{
		messageProcessor: messageProcessor,
		handlers:         handlers,
		errCh:            errCh,
	}
}

func (h *cosignatureHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageProcessor.ProcessCosignature(resp)
	if err != nil {
		h.errCh <- errors.Wrap(err, "message processor error")
		return true
	}

	handlers := h.handlers.GetHandlers(address)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for f := range handlers {
		wg.Add(1)
		go func(address *sdk.Address, callFuncPtr *subscribers.CosignatureHandler, errCh chan<- error, wg *sync.WaitGroup) {
			defer wg.Done()

			callFunc := *callFuncPtr
			if rm := callFunc(res); !rm {
				return
			}

			_, err := h.handlers.RemoveHandlers(address, callFuncPtr)
			if err != nil {
				errCh <- errors.Wrap(err, "error removing handler from storage")
				return
			}
		}(address, f, h.errCh, &wg)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
