package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	UnconfirmedRemovedHandler func(*sdk.UnconfirmedRemoved) bool

	unconfirmedRemovedHandlers        map[*UnconfirmedRemovedHandler]struct{}
	unconfirmedRemovedHandlersStorage struct {
		sync.RWMutex
		data unconfirmedRemovedHandlers
	}

	unconfirmedRemovedSubscribers        map[string]*unconfirmedRemovedHandlersStorage
	unconfirmedRemovedSubscribersStorage struct {
		sync.RWMutex
		data unconfirmedRemovedSubscribers
	}
)

func NewUnconfirmedRemoved() UnconfirmedRemoved {
	subscribers := &unconfirmedRemovedSubscribersStorage{
		data: make(unconfirmedRemovedSubscribers),
	}

	return &unconfirmedRemovedImpl{
		subscribers: subscribers,
	}
}

type UnconfirmedRemoved interface {
	AddHandlers(address *sdk.Address, handlers ...UnconfirmedRemovedHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedRemovedHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) unconfirmedRemovedHandlers
}

type unconfirmedRemovedImpl struct {
	subscribers *unconfirmedRemovedSubscribersStorage
}

func (e *unconfirmedRemovedImpl) AddHandlers(address *sdk.Address, handlers ...UnconfirmedRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &unconfirmedRemovedHandlersStorage{
			data: make(unconfirmedRemovedHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *unconfirmedRemovedImpl) RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedRemovedHandler) (bool, error) {
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

func (e *unconfirmedRemovedImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *unconfirmedRemovedImpl) GetHandlers(address *sdk.Address) unconfirmedRemovedHandlers {
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
