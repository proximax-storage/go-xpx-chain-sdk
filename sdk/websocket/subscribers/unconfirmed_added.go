package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	UnconfirmedAddedHandler func(sdk.Transaction) bool
	UnconfirmedAdded        interface {
		AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...UnconfirmedAddedHandler) error
		RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*UnconfirmedAddedHandler) bool
		HasHandlers(handle *sdk.TransactionChannelHandle) bool
		GetHandlers(handle *sdk.TransactionChannelHandle) []*UnconfirmedAddedHandler
		GetHandles() []string
	}

	unconfirmedAddedImpl struct {
		sync.Mutex
		newSubscriberCh    chan *unconfirmedAddedSubscription
		removeSubscriberCh chan *unconfirmedAddedSubscription
		subscribers        map[string][]*UnconfirmedAddedHandler
	}
	unconfirmedAddedSubscription struct {
		handle   *sdk.TransactionChannelHandle
		handlers []*UnconfirmedAddedHandler
		resultCh chan bool
	}
)

func NewUnconfirmedAdded() UnconfirmedAdded {

	p := &unconfirmedAddedImpl{
		subscribers:        make(map[string][]*UnconfirmedAddedHandler),
		newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
		removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *unconfirmedAddedImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *unconfirmedAddedImpl) addSubscription(s *unconfirmedAddedSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*UnconfirmedAddedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *unconfirmedAddedImpl) removeSubscription(s *unconfirmedAddedSubscription) {
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

func (e *unconfirmedAddedImpl) AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...UnconfirmedAddedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*UnconfirmedAddedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &unconfirmedAddedSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *unconfirmedAddedImpl) RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*UnconfirmedAddedHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &unconfirmedAddedSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *unconfirmedAddedImpl) HasHandlers(handle *sdk.TransactionChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *unconfirmedAddedImpl) GetHandlers(handle *sdk.TransactionChannelHandle) []*UnconfirmedAddedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *unconfirmedAddedImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
