package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	PartialRemovedHandler func(*sdk.PartialRemovedInfo) bool

	partialRemovedHandlers        map[*PartialRemovedHandler]struct{}
	partialRemovedHandlersStorage struct {
		sync.RWMutex
		data partialRemovedHandlers
	}

	partialRemovedSubscribers        map[string]*partialRemovedHandlersStorage
	partialRemovedSubscribersStorage struct {
		sync.RWMutex
		data partialRemovedSubscribers
	}
)

func NewPartialRemoved() PartialRemoved {
	subscribers := &partialRemovedSubscribersStorage{
		data: make(partialRemovedSubscribers),
	}

	return &partialRemovedImpl{
		subscribers: subscribers,
	}
}

type PartialRemoved interface {
	AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) partialRemovedHandlers
}

type partialRemovedImpl struct {
	subscribers *partialRemovedSubscribersStorage
}

func (e *partialRemovedImpl) AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &partialRemovedHandlersStorage{
			data: make(partialRemovedHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *partialRemovedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error) {
	if len(handlers) == 0 {
		return false, nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if external, ok := e.subscribers.data[address.Address]; !ok || len(external.data) == 0 {
		return false, errors.Wrap(handlersNotFound, "handlers not found in handlers storage")
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		delete(e.subscribers.data[address.Address].data, handlers[i])
	}

	if len(e.subscribers.data[address.Address].data) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *partialRemovedImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *partialRemovedImpl) GetHandlers(address *sdk.Address) partialRemovedHandlers {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	e.subscribers.data[address.Address].RLock()
	defer e.subscribers.data[address.Address].RUnlock()

	h, ok := e.subscribers.data[address.Address]
	if !ok {
		return nil
	}

	return h.data
}
