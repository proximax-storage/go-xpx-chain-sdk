package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	ConfirmedAddedHandler func(sdk.Transaction) bool

	ConfirmedAdded interface {
		AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error)
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*ConfirmedAddedHandler
		GetAddresses() []string
	}

	confirmedAddedSubscription struct {
		address  *sdk.Address
		handlers []*ConfirmedAddedHandler
		resultCh chan bool
	}

	confirmedAddedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *confirmedAddedSubscription
		removeSubscriberCh chan *confirmedAddedSubscription
		subscribers        map[string][]*ConfirmedAddedHandler
	}
)

func NewConfirmedAdded() ConfirmedAdded {
	p := &confirmedAddedImpl{
		subscribers:        make(map[string][]*ConfirmedAddedHandler),
		newSubscriberCh:    make(chan *confirmedAddedSubscription),
		removeSubscriberCh: make(chan *confirmedAddedSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *confirmedAddedImpl) addSubscription(s *confirmedAddedSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*ConfirmedAddedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *confirmedAddedImpl) removeSubscription(s *confirmedAddedSubscription) {
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

func (e *confirmedAddedImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *confirmedAddedImpl) AddHandlers(address *sdk.Address, handlers ...ConfirmedAddedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*ConfirmedAddedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &confirmedAddedSubscription{
		address:  address,
		handlers: refHandlers,
	}
	return nil
}

func (e *confirmedAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*ConfirmedAddedHandler) (bool, error) {

	if len(handlers) == 0 {
		return false, nil
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &confirmedAddedSubscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (e *confirmedAddedImpl) HasHandlers(address *sdk.Address) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil
}

func (e *confirmedAddedImpl) GetHandlers(address *sdk.Address) []*ConfirmedAddedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *confirmedAddedImpl) GetAddresses() []string {
	e.Lock()
	defer e.Unlock()
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
