package subscribers

import (
	"sync"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	CosignatureHandler func(*sdk.SignerInfo) bool

	Cosignature interface {
		AddHandlers(address *sdk.Address, handlers ...*CosignatureHandler) error
		RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error)
		HasHandlers(address *sdk.Address) bool
		GetHandlers(address *sdk.Address) []*CosignatureHandler
		GetAddresses() []string
	}
	cosignatureSubscription struct {
		address  *sdk.Address
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

	if _, ok := e.subscribers[s.address.Address]; !ok {
		e.subscribers[s.address.Address] = make([]*CosignatureHandler, 0)
	}
	for i := 0; i < len(s.handlers); i++ {
		e.subscribers[s.address.Address] = append(e.subscribers[s.address.Address], s.handlers[i])
	}
}

func (e *cosignatureImpl) removeSubscription(s *cosignatureSubscription) {

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

func (e *cosignatureImpl) AddHandlers(address *sdk.Address, handlers ...*CosignatureHandler) error {

	if len(handlers) == 0 {
		return nil
	}

	e.newSubscriberCh <- &cosignatureSubscription{
		address:  address,
		handlers: handlers,
	}
	return nil
}

func (e *cosignatureImpl) RemoveHandlers(address *sdk.Address, handlers ...*CosignatureHandler) (bool, error) {

	if len(handlers) == 0 {
		return false, nil
	}

	resCh := make(chan bool)

	e.removeSubscriberCh <- &cosignatureSubscription{
		address:  address,
		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (e *cosignatureImpl) HasHandlers(address *sdk.Address) bool {
	return len(e.subscribers[address.Address]) > 0 && e.subscribers[address.Address] != nil
}

func (e *cosignatureImpl) GetHandlers(address *sdk.Address) []*CosignatureHandler {
	if res, ok := e.subscribers[address.Address]; ok && res != nil {
		return res
	}

	return nil
}

func (e *cosignatureImpl) GetAddresses() []string {
	addresses := make([]string, 0, len(e.subscribers))
	for addr := range e.subscribers {
		addresses = append(addresses, addr)
	}

	return addresses
}
