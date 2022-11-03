package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"sync"
)

type (
	ReceiptHandler func(receipt *sdk.AnonymousReceipt) bool
	Receipt        interface {
		AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...ReceiptHandler) error
		RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*ReceiptHandler) bool
		HasHandlers(handle *sdk.CompoundChannelHandle) bool
		GetHandlers(handle *sdk.CompoundChannelHandle) []*ReceiptHandler
		GetHandles() []string
	}

	receiptImpl struct {
		sync.Mutex
		newSubscriberCh    chan *receiptSubscription
		removeSubscriberCh chan *receiptSubscription
		subscribers        map[string][]*ReceiptHandler
	}
	receiptSubscription struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*ReceiptHandler
		resultCh chan bool
	}
)

func NewReceipt() Receipt {

	p := &receiptImpl{
		subscribers:        make(map[string][]*ReceiptHandler),
		newSubscriberCh:    make(chan *receiptSubscription),
		removeSubscriberCh: make(chan *receiptSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *receiptImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *receiptImpl) addSubscription(s *receiptSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.handle.String()]; !ok {
		e.subscribers[s.handle.String()] = make([]*ReceiptHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.handle.String()] = append(e.subscribers[s.handle.String()], s.handlers[i])
	}
}

func (e *receiptImpl) removeSubscription(s *receiptSubscription) {
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

func (e *receiptImpl) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...ReceiptHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*ReceiptHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &receiptSubscription{
		handle:   handle,
		handlers: refHandlers,
	}
	return nil
}

func (e *receiptImpl) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*ReceiptHandler) bool {
	if len(handlers) == 0 {
		return false
	}

	resCh := make(chan bool)
	e.removeSubscriberCh <- &receiptSubscription{
		handle:   handle,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh
}

func (e *receiptImpl) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[handle.String()]) > 0 && e.subscribers[handle.String()] != nil
}

func (e *receiptImpl) GetHandlers(handle *sdk.CompoundChannelHandle) []*ReceiptHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[handle.String()]; ok && res != nil {
		return res
	}

	return nil
}

func (e *receiptImpl) GetHandles() []string {
	e.Lock()
	defer e.Unlock()
	handles := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		handles = append(handles, addr)
	}

	return handles
}
