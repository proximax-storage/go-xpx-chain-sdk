package subscribers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	UnconfirmedAddedHandler func(sdk.Transaction) bool
)

func NewUnconfirmedAdded() UnconfirmedAdded {
	return &unconfirmedAddedImpl{
		subscribers: make(map[string]map[*UnconfirmedAddedHandler]struct{}),
	}
}

type UnconfirmedAdded interface {
	AddHandlers(address *sdk.Address, handlers ...UnconfirmedAddedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedAddedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*UnconfirmedAddedHandler]struct{}
	GetAddresses() []string
}

type unconfirmedAddedImpl struct {
	sync.RWMutex
	subscribers map[string]map[*UnconfirmedAddedHandler]struct{}
}

func (e *unconfirmedAddedImpl) AddHandlers(address *sdk.Address, handlers ...UnconfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *unconfirmedAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedAddedHandler) (bool, error) {
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

func (e *unconfirmedAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *unconfirmedAddedImpl) GetHandlers(address *sdk.Address) map[*UnconfirmedAddedHandler]struct{} {
	e.RLock()
	defer e.RUnlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *unconfirmedAddedImpl) GetAddresses() []string {
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
