package subscribers

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	BlockHandler func(*sdk.BlockInfo) bool

	Block interface {
		AddHandlers(handlers ...*BlockHandler) error
		RemoveHandlers(handlers ...*BlockHandler) (bool, error)
		HasHandlers() bool
		GetHandlers() []*BlockHandler
	}

	blockSubscriberImpl struct {
		newSubscriberCh    chan *blockSubscription
		removeSubscriberCh chan *blockSubscription
		handlers           []*BlockHandler
	}

	blockSubscription struct {
		handlers []*BlockHandler
		resultCh chan bool
	}
)

func NewBlock() Block {
	s := &blockSubscriberImpl{
		handlers:           make([]*BlockHandler, 0),
		newSubscriberCh:    make(chan *blockSubscription),
		removeSubscriberCh: make(chan *blockSubscription),
	}
	go s.handleNewSubscription()
	return s
}

func (s *blockSubscriberImpl) addSubscription(b *blockSubscription) {

	for i := 0; i < len(b.handlers); i++ {
		s.handlers = append(s.handlers, b.handlers[i])
	}
}

func (s *blockSubscriberImpl) removeSubscription(b *blockSubscription) {

	if s.handlers == nil || len(b.handlers) == 0 {
		b.resultCh <- true
	}

	itemCount := len(s.handlers)
	for _, removeHandler := range s.handlers {
		for index, currentHandlers := range b.handlers {
			if removeHandler == currentHandlers {
				s.handlers = append(s.handlers[:index],
					s.handlers[index+1:]...)
			}
		}
	}

	b.resultCh <- itemCount != len(s.handlers)
}

func (s *blockSubscriberImpl) handleNewSubscription() {
	for {
		select {
		case c := <-s.newSubscriberCh:
			s.addSubscription(c)
		case c := <-s.removeSubscriberCh:
			s.removeSubscription(c)
		}
	}
}
func (s *blockSubscriberImpl) AddHandlers(handlers ...*BlockHandler) error {

	if s.handlers == nil || len(handlers) == 0 {
		return nil
	}
	s.newSubscriberCh <- &blockSubscription{
		handlers: handlers,
	}

	return nil
}

func (s *blockSubscriberImpl) RemoveHandlers(handlers ...*BlockHandler) (bool, error) {

	resCh := make(chan bool)
	s.removeSubscriberCh <- &blockSubscription{

		handlers: handlers,
		resultCh: resCh,
	}

	return <-resCh, nil
}

func (s *blockSubscriberImpl) HasHandlers() bool {
	return len(s.handlers) > 0
}

func (s *blockSubscriberImpl) GetHandlers() []*BlockHandler {
	return s.handlers
}
