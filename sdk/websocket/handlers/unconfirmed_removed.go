package handlers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"sync"
)

func NewUnconfirmedRemovedHandler(messageProcessor sdk.UnconfirmedRemovedProcessor, handlers subscribers.UnconfirmedRemoved, errCh chan<- error) *unconfirmedRemovedHandler {
	return &unconfirmedRemovedHandler{
		messageProcessor: messageProcessor,
		handlers:         handlers,
		errCh:            errCh,
	}
}

type unconfirmedRemovedHandler struct {
	messageProcessor sdk.UnconfirmedRemovedProcessor
	handlers         subscribers.UnconfirmedRemoved
	errCh            chan<- error
}

func (h *unconfirmedRemovedHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageProcessor.ProcessUnconfirmedRemoved(resp)
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
		go func(callFuncPtr *subscribers.UnconfirmedRemovedHandler, errCh chan<- error, wg *sync.WaitGroup) {
			defer wg.Done()

			callFunc := *callFuncPtr

			if rm := callFunc(res); !rm {
				return
			}

			_, err = h.handlers.RemoveHandlers(address, callFuncPtr)
			if err != nil {
				errCh <- errors.Wrap(err, "error removing handler from storage")
				return
			}

			return
		}(f, h.errCh, &wg)
	}

	wg.Wait()

	return h.handlers.HasHandlers(address)
}
