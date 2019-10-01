package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	PartialAdded interface {
		AddHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error)
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*PartialAddedHandler
		GetAddresses() []string
	}
	PartialAddedHandler func(*sdk.AggregateTransaction) bool

	partialAddedImpl struct {
		newSubscriberCh    chan *subscription
		removeSubscriberCh chan *subscription
		subscribers        map[string][]*PartialAddedHandler
	}
	subscription struct {
		address  *sdk.Address
		handlers []*PartialAddedHandler
		resultCh chan bool
	}
)

func NewPartialAdded() PartialAdded {

	p := &partialAddedImpl{
		subscribers:        make(map[string][]*PartialAddedHandler),
		newSubscriberCh:    make(chan *subscription),
		removeSubscriberCh: make(chan *subscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *partialAddedImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *partialAddedImpl) addSubscription(s *subscription) {

	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*PartialAddedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *partialAddedImpl) removeSubscription(s *subscription) {

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

func (e *partialAddedImpl) AddHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	e.newSubscriberCh <- &subscription{
		address:  address,
		handlers: handlers,
	}
	return nil
}

func (e *partialAddedImpl) RemoveHandlers(address *sdk.Address, handlers ...*PartialAddedHandler) (bool, error) {
	if len(handlers) == 0 {
		return false, nil
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &subscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (e *partialAddedImpl) HasHandlers(address *sdk.Address) bool {

	if len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil {
		return true
	}

	return false
}

func (e *partialAddedImpl) GetHandlers(address *sdk.Address) []*PartialAddedHandler {

	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *partialAddedImpl) GetAddresses() []string {
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
