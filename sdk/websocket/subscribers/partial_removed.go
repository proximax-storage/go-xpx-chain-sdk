package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	PartialRemovedHandler func(*sdk.PartialRemovedInfo) bool

	PartialRemoved interface {
		AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...PartialRemovedHandler) error
		RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*PartialRemovedHandler) bool
		HasHandlers(handle *sdk.CompoundChannelHandle) bool
		GetHandlers(handle *sdk.CompoundChannelHandle) []*PartialRemovedHandler
		GetHandles() []string
	}

	partialRemovedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *partialRemovedSubscription
		removeSubscriberCh chan *partialRemovedSubscription
		subscribers        map[string][]*PartialRemovedHandler
	}
	partialRemovedSubscription struct {
		handle   *sdk.CompoundChannelHandle
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
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*PartialRemovedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *partialRemovedImpl) removeSubscription(s *partialRemovedSubscription) {
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

func (e *partialRemovedImpl) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...PartialRemovedHandler) error {
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
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *partialRemovedImpl) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*PartialRemovedHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &partialRemovedSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *partialRemovedImpl) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *partialRemovedImpl) GetHandlers(handle *sdk.CompoundChannelHandle) []*PartialRemovedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *partialRemovedImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
