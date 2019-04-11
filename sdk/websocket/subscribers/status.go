package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	StatusHandler func(*sdk.StatusInfo) bool

	statusHandlers        map[*StatusHandler]struct{}
	statusHandlersStorage struct {
		sync.RWMutex
		data statusHandlers
	}

	statusSubscribers        map[string]*statusHandlersStorage
	statusSubscribersStorage struct {
		sync.RWMutex
		data statusSubscribers
	}
)

func NewStatus() Status {
	subscribers := &statusSubscribersStorage{
		data: make(statusSubscribers),
	}

	return &statusImpl{
		subscribers: subscribers,
	}
}

type Status interface {
	AddHandlers(address *sdk.Address, handlers ...StatusHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*StatusHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) statusHandlers
}

type statusImpl struct {
	subscribers *statusSubscribersStorage
}

func (e *statusImpl) AddHandlers(address *sdk.Address, handlers ...StatusHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	if _, ok := e.subscribers.data[address.Address]; !ok {
		e.subscribers.data[address.Address] = &statusHandlersStorage{
			data: make(statusHandlers),
		}
	}

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	for i := 0; i < len(handlers); i++ {
		e.subscribers.data[address.Address].data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *statusImpl) RemoveHandlers(address *sdk.Address, handlers ...*StatusHandler) (bool, error) {
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

func (e *statusImpl) HasHandlers(address *sdk.Address) bool {
	e.subscribers.RLock()
	defer e.subscribers.RUnlock()

	_, ok := e.subscribers.data[address.Address]
	return ok
}

func (e *statusImpl) GetHandlers(address *sdk.Address) statusHandlers {
	e.subscribers.Lock()
	defer e.subscribers.Unlock()

	e.subscribers.data[address.Address].Lock()
	defer e.subscribers.data[address.Address].Unlock()

	h, ok := e.subscribers.data[address.Address]
	if !ok {
		return nil
	}

	return h.data
}
