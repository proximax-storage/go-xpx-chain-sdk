package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	BlockHandler func(*sdk.BlockInfo) bool
)

func NewBlock() Block {
	return &blockSubscriberImpl{
		handlers: make(map[*BlockHandler]struct{}),
	}
}

type Block interface {
	AddHandlers(handlers ...BlockHandler) error
	RemoveHandlers(handlers ...*BlockHandler) (bool, error)
	HasHandlers() bool
	GetHandlers() map[*BlockHandler]struct{}
}

type blockSubscriberImpl struct {
	sync.RWMutex
	handlers map[*BlockHandler]struct{}
}

func (s *blockSubscriberImpl) AddHandlers(handlers ...BlockHandler) error {
	s.Lock()
	defer s.Unlock()

	if s.handlers == nil || len(handlers) == 0 {
		return nil
	}

	for i := 0; i < len(handlers); i++ {
		s.handlers[&handlers[i]] = struct{}{}
	}

	return nil
}

func (s *blockSubscriberImpl) RemoveHandlers(handlers ...*BlockHandler) (bool, error) {
	s.Lock()
	defer s.Unlock()

	if s.handlers == nil || len(handlers) == 0 {
		return false, errors.Wrap(handlersNotFound, "handlers not found in handlers storage")
	}

	for i := 0; i < len(handlers); i++ {
		delete(s.handlers, handlers[i])
	}

	if len(s.handlers) > 0 {
		return false, nil
	}

	return true, nil
}

func (s *blockSubscriberImpl) HasHandlers() bool {
	s.RLock()
	defer s.RUnlock()
	return len(s.handlers) > 0
}

func (s *blockSubscriberImpl) GetHandlers() map[*BlockHandler]struct{} {
	s.RLock()
	defer s.RUnlock()

	return s.handlers
}
