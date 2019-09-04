package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func NewConfirmedAddedHandler(messageMapper sdk.ConfirmedAddedMapper, handlers subscribers.ConfirmedAdded) *confirmedAddedHandler {
	return &confirmedAddedHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

type confirmedAddedHandler struct {
	messageMapper sdk.ConfirmedAddedMapper
	handlers      subscribers.ConfirmedAdded
}

func (h *confirmedAddedHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapConfirmedAdded(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapper error"))
	}

	handlers := h.handlers.GetHandlers(address)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for f := range handlers {
		wg.Add(1)
		go func(f *subscribers.ConfirmedAddedHandler) {
			defer wg.Done()

			callFunc := *f

			if rm := callFunc(res); !rm {
				return
			}

			_, err = h.handlers.RemoveHandlers(address, f)
			if err != nil {
				panic(errors.Wrap(err, "removing handler from storage"))
			}

			return
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
