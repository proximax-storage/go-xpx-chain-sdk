package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	PartialRemovedHandler func(*sdk.PartialRemovedInfo) bool
)

func NewPartialRemoved() PartialRemoved {

	return &partialRemovedImpl{
		subscribers: make(map[string]map[*PartialRemovedHandler]struct{}),
	}
}

type PartialRemoved interface {
	AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*PartialRemovedHandler]struct{}
}

type partialRemovedImpl struct {
	sync.RWMutex
	subscribers map[string]map[*PartialRemovedHandler]struct{}
}

func (e *partialRemovedImpl) AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*PartialRemovedHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *partialRemovedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error) {
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

func (e *partialRemovedImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *partialRemovedImpl) GetHandlers(address *sdk.Address) map[*PartialRemovedHandler]struct{} {
	e.RLock()
	defer e.RUnlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}
