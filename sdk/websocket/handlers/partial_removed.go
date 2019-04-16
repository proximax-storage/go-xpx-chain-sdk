package handlers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"sync"
)

func NewPartialRemovedHandler(messageMapper sdk.PartialRemovedMapper, handlers subscribers.PartialRemoved, errCh chan<- error) *partialRemovedHandler {
	return &partialRemovedHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
		errCh:         errCh,
	}
}

type partialRemovedHandler struct {
	messageMapper sdk.PartialRemovedMapper
	handlers      subscribers.PartialRemoved
	errCh         chan<- error
}

func (h *partialRemovedHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapPartialRemoved(resp)
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
		go func(f *subscribers.PartialRemovedHandler) {
			defer wg.Done()

			callFunc := *f

			if rm := callFunc(res); !rm {
				return
			}

			_, err = h.handlers.RemoveHandlers(address, f)
			if err != nil {
				h.errCh <- errors.Wrap(err, "removing handler from storage")
				return
			}

			return
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
