package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	UnconfirmedAddedHandler func(sdk.Transaction) bool

	unconfirmedAddedHandlers        map[*UnconfirmedAddedHandler]struct{}
	unconfirmedAddedHandlersStorage struct {
		sync.RWMutex
		data unconfirmedAddedHandlers
	}

	unconfirmedAddedSubscribers        map[string]*unconfirmedAddedHandlersStorage
	unconfirmedAddedSubscribersStorage struct {
		sync.RWMutex
		data unconfirmedAddedSubscribers
	}
)

func NewUnconfirmedAdded() UnconfirmedAdded {
	subscribers := &unconfirmedAddedSubscribersStorage{
		data: make(unconfirmedAddedSubscribers),
	}

	return &unconfirmedAddedImpl{
		subscribers: subscribers,
	}
}

type UnconfirmedAdded interface {
	AddHandlers(address *sdk.Address, handlers ...UnconfirmedAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) unconfirmedAddedHandlers
}

type unconfirmedAddedImpl struct {
	subscribers *unconfirmedAddedSubscribersStorage
}

func (e *unconfirmedAddedImpl) AddHandlers(address *sdk.Address, handlers ...UnconfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &unconfirmedAddedHandlersStorage{
			data: make(unconfirmedAddedHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *unconfirmedAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedAddedHandler) (bool, error) {
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

func (e *unconfirmedAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *unconfirmedAddedImpl) GetHandlers(address *sdk.Address) unconfirmedAddedHandlers {
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
