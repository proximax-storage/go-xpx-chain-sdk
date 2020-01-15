package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

type Handler interface {
	Handle(*sdk.Address, []byte) bool // Is subscription still necessary ?
}

type cosignatureHandler struct {
	messageMapper sdk.CosignatureMapper
	handlers      subscribers.Cosignature
}

func NewCosignatureHandler(messageMapper sdk.CosignatureMapper, handlers subscribers.Cosignature) *cosignatureHandler {
	return &cosignatureHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

func (h *cosignatureHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapCosignature(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapper error"))
	}

	handlers := h.handlers.GetHandlers(address)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for _, f := range handlers {
		wg.Add(1)
		go func(f *subscribers.CosignatureHandler) {
			defer wg.Done()

			callFunc := *f
			if rm := callFunc(res); !rm {
				return
			}

			h.handlers.RemoveHandlers(address, f)
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
