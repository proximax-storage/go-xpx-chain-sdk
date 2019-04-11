package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	PartialAddedHandler func(*sdk.AggregateTransaction) bool

	partialAddedHandlers        map[*PartialAddedHandler]struct{}
	partialAddedHandlersStorage struct {
		sync.RWMutex
		data partialAddedHandlers
	}

	partialAddedSubscribers        map[string]*partialAddedHandlersStorage
	partialAddedSubscribersStorage struct {
		sync.RWMutex
		data partialAddedSubscribers
	}
)

func NewPartialAdded() PartialAdded {

	subscribers := &partialAddedSubscribersStorage{
		data: make(partialAddedSubscribers),
	}

	return &partialAddedImpl{
		subscribers: subscribers,
	}
}

type PartialAdded interface {
	AddHandlers(address *sdk.Address, handlers ...PartialAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) partialAddedHandlers
}

type partialAddedImpl struct {
	subscribers *partialAddedSubscribersStorage
}

func (e *partialAddedImpl) AddHandlers(address *sdk.Address, handlers ...PartialAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &partialAddedHandlersStorage{
			data: make(partialAddedHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *partialAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error) {
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

func (e *partialAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *partialAddedImpl) GetHandlers(address *sdk.Address) partialAddedHandlers {
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
