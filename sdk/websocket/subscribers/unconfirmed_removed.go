package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	UnconfirmedRemovedHandler func(*sdk.UnconfirmedRemoved) bool

	UnconfirmedRemoved interface {
		AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...UnconfirmedRemovedHandler) error
		RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*UnconfirmedRemovedHandler) bool
		HasHandlers(handle *sdk.TransactionChannelHandle) bool
		GetHandlers(handle *sdk.TransactionChannelHandle) []*UnconfirmedRemovedHandler
		GetHandles() []string
	}

	unconfirmedRemovedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *unconfirmedRemovedSubscription
		removeSubscriberCh chan *unconfirmedRemovedSubscription
		subscribers        map[string][]*UnconfirmedRemovedHandler
	}
	unconfirmedRemovedSubscription struct {
		handle   *sdk.TransactionChannelHandle
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
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*UnconfirmedRemovedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *unconfirmedRemovedImpl) removeSubscription(s *unconfirmedRemovedSubscription) {
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

func (e *unconfirmedRemovedImpl) AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...UnconfirmedRemovedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*UnconfirmedRemovedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &unconfirmedRemovedSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *unconfirmedRemovedImpl) RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*UnconfirmedRemovedHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &unconfirmedRemovedSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *unconfirmedRemovedImpl) HasHandlers(handle *sdk.TransactionChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *unconfirmedRemovedImpl) GetHandlers(handle *sdk.TransactionChannelHandle) []*UnconfirmedRemovedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *unconfirmedRemovedImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
