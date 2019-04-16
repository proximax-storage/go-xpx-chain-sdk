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
	messageMapper sdk.CosignatureMapper
	handlers      subscribers.Cosignature
	errCh         chan<- error
}

func NewCosignatureHandler(messageMapper sdk.CosignatureMapper, handlers subscribers.Cosignature, errCh chan<- error) *cosignatureHandler {
	return &cosignatureHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
		errCh:         errCh,
	}
}

func (h *cosignatureHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapCosignature(resp)
	if err != nil {
		h.errCh <- errors.Wrap(err, "message mapper error")
		return true
	}

	handlers := h.handlers.GetHandlers(address)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for f := range handlers {
		wg.Add(1)
		go func(f *subscribers.CosignatureHandler) {
			defer wg.Done()

			callFunc := *f
			if rm := callFunc(res); !rm {
				return
			}

			_, err := h.handlers.RemoveHandlers(address, f)
			if err != nil {
				h.errCh <- errors.Wrap(err, "removing handler from storage")
				return
			}
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
