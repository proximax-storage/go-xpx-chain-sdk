package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func NewReceiptHandler(messageMapper sdk.ReceiptMapper, handlers subscribers.Receipt) *receiptHandler {
	return &receiptHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

type receiptHandler struct {
	messageMapper sdk.ReceiptMapper
	handlers      subscribers.Receipt
}

func (h *receiptHandler) Handle(handle *sdk.CompoundChannelHandle, resp []byte) bool {
	res, err := h.messageMapper.MapReceipt(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapping"))
	}

	handlers := h.handlers.GetHandlers(handle)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for _, f := range handlers {
		wg.Add(1)
		go func(f *subscribers.ReceiptHandler) {
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
