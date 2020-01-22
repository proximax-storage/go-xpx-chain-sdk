package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	UnconfirmedRemovedHandler func(*sdk.UnconfirmedRemoved) bool

	UnconfirmedRemoved interface {
		AddHandlers(address *sdk.Address, handlers ...UnconfirmedRemovedHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedRemovedHandler) bool
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*UnconfirmedRemovedHandler
		GetAddresses() []string
	}

	unconfirmedRemovedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *unconfirmedRemovedSubscription
		removeSubscriberCh chan *unconfirmedRemovedSubscription
		subscribers        map[string][]*UnconfirmedRemovedHandler
	}
	unconfirmedRemovedSubscription struct {
		address  *sdk.Address
		handlers []*UnconfirmedRemovedHandler
		resultCh chan bool
	}
)

func NewUnconfirmedRemoved() UnconfirmedRemoved {

	p := &unconfirmedRemovedImpl{
		subscribers:        make(map[string][]*UnconfirmedRemovedHandler),
		newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
		removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *unconfirmedRemovedImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *unconfirmedRemovedImpl) addSubscription(s *unconfirmedRemovedSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*UnconfirmedRemovedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *unconfirmedRemovedImpl) removeSubscription(s *unconfirmedRemovedSubscription) {
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

func (e *unconfirmedRemovedImpl) AddHandlers(address *sdk.Address, handlers ...UnconfirmedRemovedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*UnconfirmedRemovedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &unconfirmedRemovedSubscription{
		address:  address,
		handlers: refHandlers,
	}
	return nil
}

func (e *unconfirmedRemovedImpl) RemoveHandlers(address *sdk.Address, handlers ...*UnconfirmedRemovedHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &unconfirmedRemovedSubscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *unconfirmedRemovedImpl) HasHandlers(address *sdk.Address) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil
}

func (e *unconfirmedRemovedImpl) GetHandlers(address *sdk.Address) []*UnconfirmedRemovedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *unconfirmedRemovedImpl) GetAddresses() []string {
	e.Lock()
	defer e.Unlock()
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
