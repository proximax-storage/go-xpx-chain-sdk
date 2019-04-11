package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	ConfirmedAddedHandler func(sdk.Transaction) bool

	ConfirmedAddedHandlers        map[*ConfirmedAddedHandler]struct{}
	confirmedAddedHandlersStorage struct {
		sync.RWMutex
		data ConfirmedAddedHandlers
	}

	confirmedAddedSubscribers        map[string]*confirmedAddedHandlersStorage
	confirmedAddedSubscribersStorage struct {
		sync.RWMutex
		data confirmedAddedSubscribers
	}
)

func NewConfirmedAdded() ConfirmedAdded {

	subscribers := &confirmedAddedSubscribersStorage{
		data: make(map[string]*confirmedAddedHandlersStorage),
	}

	return &confirmedAddedImpl{
		subscribers: subscribers,
	}
}

type ConfirmedAdded interface {
	AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) ConfirmedAddedHandlers
}

type confirmedAddedImpl struct {
	subscribers *confirmedAddedSubscribersStorage
}

func (e *confirmedAddedImpl) AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &confirmedAddedHandlersStorage{
			data: make(ConfirmedAddedHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *confirmedAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error) {
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

func (e *confirmedAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *confirmedAddedImpl) GetHandlers(address *sdk.Address) ConfirmedAddedHandlers {
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
