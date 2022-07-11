package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func NewPartialAddedHandler(messageMapper sdk.PartialAddedMapper, handlers subscribers.PartialAdded) *partialAddedHandler {
	return &partialAddedHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

type partialAddedHandler struct {
	messageMapper sdk.PartialAddedMapper
	handlers      subscribers.PartialAdded
}

func (h *partialAddedHandler) Handle(handle *sdk.TransactionChannelHandle, resp []byte) bool {
	res, err := h.messageMapper.MapPartialAdded(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapper error"))
	}

	handlers := h.handlers.GetHandlers(handle)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for _, f := range handlers {
		wg.Add(1)
		go func(f *subscribers.PartialAddedHandler) {
			defer wg.Done()

			callFunc := *f

			if rm := callFunc(res); !rm {
				return
			}

			h.handlers.RemoveHandlers(handle, f)
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(handle)
}
