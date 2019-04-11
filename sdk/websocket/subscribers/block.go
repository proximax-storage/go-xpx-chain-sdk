package subscribers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"sync"
)

type (
	BlockHandler  func(*sdk.BlockInfo) bool
	blockHandlers map[*BlockHandler]struct{}

	blockHandlersStorage struct {
		sync.RWMutex
		data blockHandlers
	}
)

func NewBlockSubscriber() Block {
	handlers := blockHandlersStorage{}
	handlers.data = make(map[*BlockHandler]struct{})

	return &blockSubscriberImpl{
		handlers: handlers,
	}
}

type Block interface {
	AddHandlers(handlers ...BlockHandler) error
	RemoveHandlers(handlers ...*BlockHandler) (bool, error)
	HasHandlers() bool
	GetHandlers() blockHandlers
}

type blockSubscriberImpl struct {
	handlers blockHandlersStorage
}

func (s *blockSubscriberImpl) AddHandlers(handlers ...BlockHandler) error {
	s.handlers.Lock()
	defer s.handlers.Unlock()

	for i := 0; i < len(handlers); i++ {
		s.handlers.data[&handlers[i]] = struct{}{}
	}

	return nil
}

func (s *blockSubscriberImpl) RemoveHandlers(handlers ...*BlockHandler) (bool, error) {
	s.handlers.Lock()
	defer s.handlers.Unlock()

	if s.handlers.data == nil {
		return false, errors.Wrap(handlersNotFound, "handlers not found in handlers storage")
	}

	for i := 0; i < len(handlers); i++ {
		delete(s.handlers.data, handlers[i])
	}

	if len(s.handlers.data) > 0 {
		return false, nil
	}

	s.handlers.data = nil

	return true, nil
}

func (s *blockSubscriberImpl) HasHandlers() bool {
	s.handlers.RLock()
	defer s.handlers.RUnlock()

	return len(s.handlers.data) > 0
}

func (s *blockSubscriberImpl) GetHandlers() blockHandlers {
	s.handlers.RLock()
	defer s.handlers.RUnlock()

	return s.handlers.data
}
