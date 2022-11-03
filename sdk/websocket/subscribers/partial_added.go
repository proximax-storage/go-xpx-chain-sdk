package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	PartialAdded interface {
		AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...PartialAddedHandler) error
		RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*PartialAddedHandler) bool
		HasHandlers(handle *sdk.CompoundChannelHandle) bool
		GetHandlers(handle *sdk.CompoundChannelHandle) []*PartialAddedHandler
		GetHandles() []string
	}
	PartialAddedHandler func(sdk.Transaction) bool

	partialAddedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *partialAddedSubscription
		removeSubscriberCh chan *partialAddedSubscription
		subscribers        map[string][]*PartialAddedHandler
	}
	partialAddedSubscription struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*PartialAddedHandler
		resultCh chan bool
	}
)

func NewPartialAdded() PartialAdded {

	p := &partialAddedImpl{
		subscribers:        make(map[string][]*PartialAddedHandler),
		newSubscriberCh:    make(chan *partialAddedSubscription),
		removeSubscriberCh: make(chan *partialAddedSubscription),
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

func (e *partialAddedImpl) addSubscription(s *partialAddedSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*PartialAddedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *partialAddedImpl) removeSubscription(s *partialAddedSubscription) {
	e.Lock()
	defer e.Unlock()
	if external, ok := e.subscribers[s.handle.String()]; !ok || len(external) == 0 {
		s.resultCh <- false
	}

	itemCount := len(e.subscribers[s.handle.String()])
	for _, removeHandler := range s.handlers {
		for index, currentHandlers := range e.subscribers[s.handle.String()] {
			if removeHandler == currentHandlers {
				e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()][:index],
					e.subscribers[s.handle.String()][index+1:]...)
			}
		}
	}

	s.resultCh <- itemCount != len(e.subscribers[s.handle.String()])
}

func (e *partialAddedImpl) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...PartialAddedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*PartialAddedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &partialAddedSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *partialAddedImpl) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*PartialAddedHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &partialAddedSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *partialAddedImpl) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *partialAddedImpl) GetHandlers(handle *sdk.CompoundChannelHandle) []*PartialAddedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *partialAddedImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
