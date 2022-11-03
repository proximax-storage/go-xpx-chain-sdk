package subscribers

import (
	"sync"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	CosignatureHandler func(*sdk.SignerInfo) bool

	Cosignature interface {
		AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...CosignatureHandler) error
		RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*CosignatureHandler) bool
		HasHandlers(handle *sdk.CompoundChannelHandle) bool
		GetHandlers(handle *sdk.CompoundChannelHandle) []*CosignatureHandler
		GetHandles() []string
	}
	cosignatureSubscription struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*CosignatureHandler
		resultCh chan bool
	}
	cosignatureImpl struct {
		sync.RWMutex
		subscribers        map[string][]*CosignatureHandler
		newSubscriberCh    chan *cosignatureSubscription
		removeSubscriberCh chan *cosignatureSubscription
	}
)

func NewCosignature() Cosignature {

	p := &cosignatureImpl{
		subscribers:        make(map[string][]*CosignatureHandler),
		newSubscriberCh:    make(chan *cosignatureSubscription),
		removeSubscriberCh: make(chan *cosignatureSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *cosignatureImpl) addSubscription(s *cosignatureSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*CosignatureHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *cosignatureImpl) removeSubscription(s *cosignatureSubscription) {
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

func (e *cosignatureImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *cosignatureImpl) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...CosignatureHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*CosignatureHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &cosignatureSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *cosignatureImpl) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*CosignatureHandler) bool {

	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)

	e.removeSubscriberCh <- &cosignatureSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *cosignatureImpl) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *cosignatureImpl) GetHandlers(handle *sdk.CompoundChannelHandle) []*CosignatureHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *cosignatureImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
