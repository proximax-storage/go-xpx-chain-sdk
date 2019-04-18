package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	CosignatureHandler func(*sdk.SignerInfo) bool
)

func NewCosignature() Cosignature {

	return &cosignatureImpl{
		subscribers: make(map[string]map[*CosignatureHandler]struct{}),
	}
}

type Cosignature interface {
	AddHandlers(address *sdk.Address, handlers ...CosignatureHandler) error
	RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error)
	HasHandlers(address *sdk.Address) bool
	GetHandlers(address *sdk.Address) map[*CosignatureHandler]struct{}
}

type cosignatureImpl struct {
	sync.RWMutex
	subscribers map[string]map[*CosignatureHandler]struct{}
}

func (e *cosignatureImpl) AddHandlers(address *sdk.Address, handlers ...CosignatureHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if _, ok := e.subscribers[address.Address]; !ok {
		e.subscribers[address.Address] = make(map[*CosignatureHandler]struct{})
	}

	for i := 0; i < len(handlers); i++ {
		e.subscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *cosignatureImpl) RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error) {
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

func (e *cosignatureImpl) HasHandlers(address *sdk.Address) bool {
	e.RLock()
	defer e.RUnlock()

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *cosignatureImpl) GetHandlers(address *sdk.Address) map[*CosignatureHandler]struct{} {
	e.RLock()
	defer e.RUnlock()

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}
