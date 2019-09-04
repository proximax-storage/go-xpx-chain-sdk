package subscribers

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	StatusHandler func(*sdk.StatusInfo) bool
)

func NewStatus() Status {
	return &statusImpl{
		subscribers: make(map[string]map[*StatusHandler]struct{}),
	}
}

type Status interface {
	AddHandlers(address *sdk.Address, handlers ...StatusHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*StatusHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*StatusHandler]struct{}
	GetAddresses() []string
}

type statusImpl struct {
	sync.RWMutex
	subscribers map[string]map[*StatusHandler]struct{}
}

func (e *statusImpl) AddHandlers(address *sdk.Address, handlers ...StatusHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*StatusHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *statusImpl) RemoveHandlers(address *sdk.Address, handlers ...*StatusHandler) (bool, error) {
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

func (e *statusImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *statusImpl) GetHandlers(address *sdk.Address) map[*StatusHandler]struct{} {
	e.Lock()
	defer e.Unlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *statusImpl) GetAddresses() []string {
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
