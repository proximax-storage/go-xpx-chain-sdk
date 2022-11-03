package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func NewUnconfirmedAddedHandler(messageMapper sdk.UnconfirmedAddedMapper, handlers subscribers.UnconfirmedAdded) *unconfirmedAddedHandler {
	return &unconfirmedAddedHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

type unconfirmedAddedHandler struct {
	messageMapper sdk.UnconfirmedAddedMapper
	handlers      subscribers.UnconfirmedAdded
	errCh         chan<- error
}

func (h *unconfirmedAddedHandler) Handle(handle *sdk.CompoundChannelHandle, resp []byte) bool {
	res, err := h.messageMapper.MapUnconfirmedAdded(resp)
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
		go func(f *subscribers.UnconfirmedAddedHandler) {
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
