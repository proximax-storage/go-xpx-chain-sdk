package subscribers

import (
	"sync"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	DriveStateHandler func(*sdk.DriveStateInfo) bool

	DriveState interface {
		AddHandlers(address *sdk.Address, handlers ...DriveStateHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*DriveStateHandler) (bool, error)
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*DriveStateHandler
		GetAddresses() []string
	}
	driveStateSubscription struct {
		address  *sdk.Address
		handlers []*DriveStateHandler
		resultCh chan bool
	}
	driveStateImpl struct {
		sync.RWMutex
		subscribers        map[string][]*DriveStateHandler
		newSubscriberCh    chan *driveStateSubscription
		removeSubscriberCh chan *driveStateSubscription
	}
)

func NewDriveState() DriveState {

	p := &driveStateImpl{
		subscribers:        make(map[string][]*DriveStateHandler),
		newSubscriberCh:    make(chan *driveStateSubscription),
		removeSubscriberCh: make(chan *driveStateSubscription),
	}
	go p.handleNewSubscription()
	return p
}

func (e *driveStateImpl) addSubscription(s *driveStateSubscription) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*DriveStateHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *driveStateImpl) removeSubscription(s *driveStateSubscription) {
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

func (e *driveStateImpl) handleNewSubscription() {
	for {
		select {
		case s := <-e.newSubscriberCh:
			e.addSubscription(s)
		case s := <-e.removeSubscriberCh:
			e.removeSubscription(s)
		}
	}
}

func (e *driveStateImpl) AddHandlers(address *sdk.Address, handlers ...DriveStateHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	refHandlers := make([]*DriveStateHandler, len(handlers))
	for i, h := range handlers {
		refHandlers[i] = &h
	}

	e.newSubscriberCh <- &driveStateSubscription{
		address:  address,
		handlers: refHandlers,
	}
	return nil
}

func (e *driveStateImpl) RemoveHandlers(address *sdk.Address, handlers ...*DriveStateHandler) (bool, error) {

	if len(handlers) == 0 {
		return false, nil
	}

	resCh := make(chan bool)

	e.removeSubscriberCh <- &driveStateSubscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (e *driveStateImpl) HasHandlers(address *sdk.Address) bool {
	e.Lock()
	defer e.Unlock()
	return len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil
}

func (e *driveStateImpl) GetHandlers(address *sdk.Address) []*DriveStateHandler {
	e.Lock()
	defer e.Unlock()
	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *driveStateImpl) GetAddresses() []string {
	e.Lock()
	defer e.Unlock()
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
