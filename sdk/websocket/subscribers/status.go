package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	StatusHandler func(*sdk.StatusInfo) bool

	Status interface {
		AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...StatusHandler) error
		RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*StatusHandler) bool
		HasHandlers(handle *sdk.CompoundChannelHandle) bool
		GetHandlers(handle *sdk.CompoundChannelHandle) []*StatusHandler
		GetHandles() []string
	}

	statusImpl struct {
		sync.Mutex
		newSubscriberCh    chan *statusSubscription
		removeSubscriberCh chan *statusSubscription
		subscribers        map[string][]*StatusHandler
	}
	statusSubscription struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*StatusHandler
		resultCh chan bool
	}
)

func NewStatus() Status {

	p := &statusImpl{
		subscribers:        make(map[string][]*StatusHandler),
		newSubscriberCh:    make(chan *statusSubscription),
		removeSubscriberCh: make(chan *statusSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *statusImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *statusImpl) addSubscription(s *statusSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*StatusHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *statusImpl) removeSubscription(s *statusSubscription) {
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

func (e *statusImpl) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...StatusHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*StatusHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &statusSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *statusImpl) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*StatusHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &statusSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *statusImpl) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *statusImpl) GetHandlers(handle *sdk.CompoundChannelHandle) []*StatusHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *statusImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
