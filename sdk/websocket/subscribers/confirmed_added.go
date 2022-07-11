package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	ConfirmedAddedHandler func(sdk.Transaction) bool

	ConfirmedAdded interface {
		AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...ConfirmedAddedHandler) error
		RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*ConfirmedAddedHandler) bool
		HasHandlers(handle *sdk.TransactionChannelHandle) bool
		GetHandlers(handle *sdk.TransactionChannelHandle) []*ConfirmedAddedHandler
		GetHandles() []string
	}

	confirmedAddedSubscription struct {
		handle   *sdk.TransactionChannelHandle
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
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*ConfirmedAddedHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *confirmedAddedImpl) removeSubscription(s *confirmedAddedSubscription) {
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

func (e *confirmedAddedImpl) AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...ConfirmedAddedHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*ConfirmedAddedHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &confirmedAddedSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *confirmedAddedImpl) RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*ConfirmedAddedHandler) bool {

	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &confirmedAddedSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *confirmedAddedImpl) HasHandlers(handle *sdk.TransactionChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *confirmedAddedImpl) GetHandlers(handle *sdk.TransactionChannelHandle) []*ConfirmedAddedHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *confirmedAddedImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
