package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func NewBlockHandler(messageMapper sdk.BlockMapper, handlers subscribers.Block) *blockHandler {
	return &blockHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

type blockHandler struct {
	messageMapper sdk.BlockMapper
	handlers      subscribers.Block
}

func (h *blockHandler) Handle(address *sdk.Address, resp []byte) bool {
	res, err := h.messageMapper.MapBlock(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapping"))
	}

	handlers := h.handlers.GetHandlers()
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for _, f := range handlers {
		wg.Add(1)
		go func(f *subscribers.BlockHandler) {
			defer wg.Done()

			callFunc := *f

			if rm := callFunc(res); !rm {
				return
			}

			h.handlers.RemoveHandlers(f)
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers()
}
