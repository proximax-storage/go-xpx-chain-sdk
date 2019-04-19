package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	PartialAddedHandler func(*sdk.AggregateTransaction) bool
)

func NewPartialAdded() PartialAdded {

	return &partialAddedImpl{
		subscribers: make(map[string]map[*PartialAddedHandler]struct{}),
	}
}

type PartialAdded interface {
	AddHandlers(address *sdk.Address, handlers ...PartialAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*PartialAddedHandler]struct{}
}

type partialAddedImpl struct {
	sync.RWMutex
	subscribers map[string]map[*PartialAddedHandler]struct{}
}

func (e *partialAddedImpl) AddHandlers(address *sdk.Address, handlers ...PartialAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*PartialAddedHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *partialAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error) {
	if len(handlers) == 0 {
		return false, nil
	}

	e.Lock()
	defer e.Unlock()

	if external, ok := e.subscribers[address.Address]; !ok || len(external) == 0 {
		return false, errors.Wrap(handlersNotFound, "handlers not found in handlers storage")
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.subscribers[address.Address], handlers[i])
	}

	if len(e.subscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *partialAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *partialAddedImpl) GetHandlers(address *sdk.Address) map[*PartialAddedHandler]struct{} {
	e.RLock()
	defer e.RUnlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}
