package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	PartialRemovedHandler func(*sdk.PartialRemovedInfo) bool

	PartialRemoved interface {
		AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error)
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*PartialRemovedHandler
		GetAddresses() []string
	}

	partialRemovedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *partialRemovedSubscription
		removeSubscriberCh chan *partialRemovedSubscription
		subscribers        map[string][]*PartialRemovedHandler
	}
	partialRemovedSubscription struct {
		address  *sdk.Address
		handlers []*PartialRemovedHandler
		resultCh chan bool
	}
)

func NewPartialRemoved() PartialRemoved {

	p := &partialRemovedImpl{
		subscribers:        make(map[string][]*PartialRemovedHandler),
		newSubscriberCh:    make(chan *partialRemovedSubscription),
		removeSubscriberCh: make(chan *partialRemovedSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *partialRemovedImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *partialRemovedImpl) addSubscription(s *partialRemovedSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*PartialRemovedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *partialRemovedImpl) removeSubscription(s *partialRemovedSubscription) {
	e.Lock()
	defer e.Unlock()
	if external, ok := e.subscribers[s.address.Address]; !ok || len(external) == 0 {
		s.resultCh <- false
	}

	itemCount := len(e.subscribers[s.address.Address])
	for _, removeHandler := range s.handlers {
		for index, currentHandlers := range e.subscribers[s.address.Address] {
			if removeHandler == currentHandlers {
				e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address][:index],
					e.subscribers[s.address.Address][index+1:]...)
			}
		}
	}

	s.resultCh <- itemCount != len(e.subscribers[s.address.Address])
}

func (e *partialRemovedImpl) AddHandlers(address *sdk.Address, handlers ...PartialRemovedHandler) error {
	e.Lock()
	defer e.Unlock()
	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*PartialRemovedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &partialRemovedSubscription{
		address:  address,
		handlers: refHandlers,
	}
	return nil
}

func (e *partialRemovedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialRemovedHandler) (bool, error) {
	if len(handlers) == 0 {
		return false, nil
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &partialRemovedSubscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (e *partialRemovedImpl) HasHandlers(address *sdk.Address) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil
}

func (e *partialRemovedImpl) GetHandlers(address *sdk.Address) []*PartialRemovedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *partialRemovedImpl) GetAddresses() []string {
	e.Lock()
	defer e.Unlock()
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
