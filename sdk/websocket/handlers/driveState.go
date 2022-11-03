package handlers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

type driveStateHandler struct {
	messageMapper sdk.DriveStateMapper
	handlers      subscribers.DriveState
}

func NewDriveStateHandler(messageMapper sdk.DriveStateMapper, handlers subscribers.DriveState) *driveStateHandler {
	return &driveStateHandler{
		messageMapper: messageMapper,
		handlers:      handlers,
	}
}

func (h *driveStateHandler) Handle(handle *sdk.CompoundChannelHandle, resp []byte) bool {
	res, err := h.messageMapper.MapDriveState(resp)
	if err != nil {
		panic(errors.Wrap(err, "message mapper error"))
	}

	handlers := h.handlers.GetHandlers(handle.Address)
	if len(handlers) == 0 {
		return true
	}

	var wg sync.WaitGroup

	for _, f := range handlers {
		wg.Add(1)
		go func(f *subscribers.DriveStateHandler) {
			defer wg.Done()

			callFunc := *f
			if rm := callFunc(res); !rm {
				return
			}

			h.handlers.RemoveHandlers(handle.Address, f)
		}(f)
	}

	wg.Wait()

	return h.handlers.HasHandlers(handle.Address)
}
