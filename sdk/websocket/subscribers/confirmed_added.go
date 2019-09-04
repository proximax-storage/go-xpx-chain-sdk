package subscribers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	ConfirmedAddedHandler func(sdk.Transaction) bool
)

func NewConfirmedAdded() ConfirmedAdded {
	return &confirmedAddedImpl{
		subscribers: make(map[string]map[*ConfirmedAddedHandler]struct{}),
	}
}

type ConfirmedAdded interface {
	AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*ConfirmedAddedHandler]struct{}
	GetAddresses() []string
}

type confirmedAddedImpl struct {
	sync.RWMutex
	subscribers map[string]map[*ConfirmedAddedHandler]struct{}
}

func (e *confirmedAddedImpl) AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*ConfirmedAddedHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *confirmedAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error) {
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

func (e *confirmedAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *confirmedAddedImpl) GetHandlers(address *sdk.Address) map[*ConfirmedAddedHandler]struct{} {
	e.RLock()
	defer e.RUnlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *confirmedAddedImpl) GetAddresses() []string {
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
