package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	CosignatureHandler func(*sdk.SignerInfo) bool

	cosignatureHandlers        map[*CosignatureHandler]struct{}
	cosignatureHandlersStorage struct {
		sync.RWMutex
		data cosignatureHandlers
	}

	cosignatureSubscribers        map[string]*cosignatureHandlersStorage
	cosignatureSubscribersStorage struct {
		sync.RWMutex
		data cosignatureSubscribers
	}
)

func NewCosignature() Cosignature {
	subscribers := &cosignatureSubscribersStorage{
		data: make(cosignatureSubscribers),
	}

	return &cosignatureImpl{
		subscribers: subscribers,
	}
}

type Cosignature interface {
	AddHandlers(address *sdk.Address, handlers ...CosignatureHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) cosignatureHandlers
}

type cosignatureImpl struct {
	subscribers *cosignatureSubscribersStorage
}

func (e *cosignatureImpl) AddHandlers(address *sdk.Address, handlers ...CosignatureHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &cosignatureHandlersStorage{
			data: make(cosignatureHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *cosignatureImpl) RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error) {
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

func (e *cosignatureImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *cosignatureImpl) GetHandlers(address *sdk.Address) cosignatureHandlers {
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
