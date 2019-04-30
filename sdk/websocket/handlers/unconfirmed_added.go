package handlers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"sync"
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

func (h *unconfirmedAddedHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapUnconfirmedAdded(resp)
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
		go func(f *subscribers.UnconfirmedAddedHandler) {
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
