package sdk

import (
	"github.com/pkg/errors"
)

type (
	blockHandler              func(*BlockInfo) bool
	confirmedAddedHandler     func(Transaction) bool
	unconfirmedAddedHandler   func(Transaction) bool
	unconfirmedRemovedHandler func(*UnconfirmedRemoved) bool
	statusHandler             func(*StatusInfo) bool
	partialAddedHandler       func(*AggregateTransaction) bool
	partialRemovedHandler     func(*PartialRemovedInfo) bool
	cosignatureHandler        func(*SignerInfo) bool
)

type (
	blockHandlers              map[*blockHandler]struct{}
	confirmedAddedHandlers     map[*confirmedAddedHandler]struct{}
	unconfirmedAddedHandlers   map[*unconfirmedAddedHandler]struct{}
	unconfirmedRemovedHandlers map[*unconfirmedRemovedHandler]struct{}
	statusHandlers             map[*statusHandler]struct{}
	partialAddedHandlers       map[*partialAddedHandler]struct{}
	partialRemovedHandlers     map[*partialRemovedHandler]struct{}
	cosignatureHandlers        map[*cosignatureHandler]struct{}
)

type (
	confirmedAddedSubscribers     map[string]confirmedAddedHandlers
	unconfirmedAddedSubscribers   map[string]unconfirmedAddedHandlers
	unconfirmedRemovedSubscribers map[string]unconfirmedRemovedHandlers
	statusSubscribers             map[string]statusHandlers
	partialAddedSubscribers       map[string]partialAddedHandlers
	partialRemovedSubscribers     map[string]partialRemovedHandlers
	cosignatureSubscribers        map[string]cosignatureHandlers
)

var (
	subscriptionNotFoundError = errors.New("subscription not found")
	handlersNotFound          = errors.New("handlers not found")
)

func newEventSubscribers() eventsSubscribers {
	return &eventsSubscribersImpl{
		blockHandlers:                 make(blockHandlers),
		unconfirmedAddedSubscribers:   make(unconfirmedAddedSubscribers),
		confirmedAddedSubscribers:     make(confirmedAddedSubscribers),
		unconfirmedRemovedSubscribers: make(unconfirmedRemovedSubscribers),
		partialAddedSubscribers:       make(partialAddedSubscribers),
		partialRemovedSubscribers:     make(partialRemovedSubscribers),
		statusSubscribers:             make(statusSubscribers),
		cosignatureSubscribers:        make(cosignatureSubscribers),
	}
}

type eventsSubscribers interface {
	AddBlockHandlers(handlers ...blockHandler) error
	AddConfirmedAddedHandlers(address *Address, handlers ...confirmedAddedHandler) error
	AddUnconfirmedAddedHandlers(address *Address, handlers ...unconfirmedAddedHandler) error
	AddUnconfirmedRemovedHandlers(address *Address, handlers ...unconfirmedRemovedHandler) error
	AddPartialAddedHandlers(address *Address, handlers ...partialAddedHandler) error
	AddPartialRemovedHandlers(address *Address, handlers ...partialRemovedHandler) error
	AddStatusInfoHandlers(address *Address, handlers ...statusHandler) error
	AddCosignatureHandlers(address *Address, handlers ...cosignatureHandler) error

	RemoveBlockHandlers(handlers ...*blockHandler) (bool, error)
	RemoveConfirmedAddedHandlers(address *Address, handlers ...*confirmedAddedHandler) (bool, error)
	RemoveUnconfirmedAddedHandlers(address *Address, handlers ...*unconfirmedAddedHandler) (bool, error)
	RemoveUnconfirmedRemovedHandlers(address *Address, handlers ...*unconfirmedRemovedHandler) (bool, error)
	RemovePartialAddedHandlers(address *Address, handlers ...*partialAddedHandler) (bool, error)
	RemovePartialRemovedHandlers(address *Address, handlers ...*partialRemovedHandler) (bool, error)
	RemoveStatusInfoHandlers(address *Address, handlers ...*statusHandler) (bool, error)
	RemoveCosignatureHandlers(address *Address, handlers ...*cosignatureHandler) (bool, error)

	IsBlockSubscribed() bool
	IsAddConfirmedAddedSubscribed(address *Address) bool
	IsAddUnconfirmedAddedSubscribed(address *Address) bool
	IsAddUnconfirmedRemovedSubscribed(address *Address) bool
	IsAddPartialAddedSubscribed(address *Address) bool
	IsAddPartialRemovedSubscribed(address *Address) bool
	IsAddStatusInfoSubscribed(address *Address) bool
	IsAddCosignatureSubscribed(address *Address) bool

	GetBlockHandlers() blockHandlers
	GetConfirmedAddedHandlers(address *Address) (confirmedAddedHandlers, error)
	GetUnconfirmedAddedHandlers(address *Address) (unconfirmedAddedHandlers, error)
	GetUnconfirmedRemovedHandlers(address *Address) (unconfirmedRemovedHandlers, error)
	GetPartialAddedHandlers(address *Address) (partialAddedHandlers, error)
	GetPartialRemovedHandlers(address *Address) (partialRemovedHandlers, error)
	GetStatusInfoHandlers(address *Address) (statusHandlers, error)
	GetCosignatureHandlers(address *Address) (cosignatureHandlers, error)
}

type eventsSubscribersImpl struct {
	blockHandlers                 blockHandlers
	confirmedAddedSubscribers     confirmedAddedSubscribers
	unconfirmedAddedSubscribers   unconfirmedAddedSubscribers
	unconfirmedRemovedSubscribers unconfirmedRemovedSubscribers
	partialAddedSubscribers       partialAddedSubscribers
	partialRemovedSubscribers     partialRemovedSubscribers
	statusSubscribers             statusSubscribers
	cosignatureSubscribers        cosignatureSubscribers
}

func (e *eventsSubscribersImpl) AddBlockHandlers(handler ...blockHandler) error {
	for i := 0; i < len(handler); i++ {
		e.blockHandlers[&handler[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddConfirmedAddedHandlers(address *Address, handlers ...confirmedAddedHandler) error {
	if e.confirmedAddedSubscribers[address.Address] == nil {
		e.confirmedAddedSubscribers[address.Address] = make(confirmedAddedHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.confirmedAddedSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddUnconfirmedAddedHandlers(address *Address, handlers ...unconfirmedAddedHandler) error {
	if e.unconfirmedAddedSubscribers[address.Address] == nil {
		e.unconfirmedAddedSubscribers[address.Address] = make(unconfirmedAddedHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.unconfirmedAddedSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddUnconfirmedRemovedHandlers(address *Address, handlers ...unconfirmedRemovedHandler) error {
	if e.unconfirmedRemovedSubscribers[address.Address] == nil {
		e.unconfirmedRemovedSubscribers[address.Address] = make(unconfirmedRemovedHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.unconfirmedRemovedSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddPartialAddedHandlers(address *Address, handlers ...partialAddedHandler) error {
	if e.partialAddedSubscribers[address.Address] == nil {
		e.partialAddedSubscribers[address.Address] = make(partialAddedHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.partialAddedSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddPartialRemovedHandlers(address *Address, handlers ...partialRemovedHandler) error {
	if e.partialRemovedSubscribers[address.Address] == nil {
		e.partialRemovedSubscribers[address.Address] = make(partialRemovedHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.partialRemovedSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddStatusInfoHandlers(address *Address, handlers ...statusHandler) error {
	if e.statusSubscribers[address.Address] == nil {
		e.statusSubscribers[address.Address] = make(statusHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.statusSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) AddCosignatureHandlers(address *Address, handlers ...cosignatureHandler) error {
	if e.cosignatureSubscribers[address.Address] == nil {
		e.cosignatureSubscribers[address.Address] = make(cosignatureHandlers)
	}

	for i := 0; i < len(handlers); i++ {
		e.cosignatureSubscribers[address.Address][&handlers[i]] = struct{}{}
	}

	return nil
}

func (e *eventsSubscribersImpl) RemoveBlockHandlers(handlers ...*blockHandler) (bool, error) {
	if e.blockHandlers == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.blockHandlers, handlers[i])
	}

	if len(e.blockHandlers) > 0 {
		return false, nil
	}

	e.blockHandlers = nil

	return true, nil
}

func (e *eventsSubscribersImpl) RemoveConfirmedAddedHandlers(address *Address, handlers ...*confirmedAddedHandler) (bool, error) {
	if e.confirmedAddedSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.confirmedAddedSubscribers[address.Address], handlers[i])
	}

	if len(e.confirmedAddedSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemoveUnconfirmedAddedHandlers(address *Address, handlers ...*unconfirmedAddedHandler) (bool, error) {
	if e.unconfirmedAddedSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.unconfirmedAddedSubscribers[address.Address], handlers[i])
	}

	if len(e.unconfirmedAddedSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemoveUnconfirmedRemovedHandlers(address *Address, handlers ...*unconfirmedRemovedHandler) (bool, error) {
	if e.unconfirmedRemovedSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.unconfirmedRemovedSubscribers[address.Address], handlers[i])
	}

	if len(e.unconfirmedRemovedSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemovePartialAddedHandlers(address *Address, handlers ...*partialAddedHandler) (bool, error) {
	if e.partialAddedSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.partialAddedSubscribers[address.Address], handlers[i])
	}

	if len(e.partialAddedSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemovePartialRemovedHandlers(address *Address, handlers ...*partialRemovedHandler) (bool, error) {
	if e.partialRemovedSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.partialRemovedSubscribers[address.Address], handlers[i])
	}

	if len(e.partialRemovedSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemoveStatusInfoHandlers(address *Address, handlers ...*statusHandler) (bool, error) {
	if e.statusSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.statusSubscribers[address.Address], handlers[i])
	}

	if len(e.statusSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) RemoveCosignatureHandlers(address *Address, handlers ...*cosignatureHandler) (bool, error) {
	if e.cosignatureSubscribers[address.Address] == nil {
		return false, handlersNotFound
	}

	for i := 0; i < len(handlers); i++ {
		delete(e.cosignatureSubscribers[address.Address], handlers[i])
	}

	if len(e.cosignatureSubscribers[address.Address]) > 0 {
		return false, nil
	}

	return true, nil
}

func (e *eventsSubscribersImpl) IsBlockSubscribed() bool {
	return len(e.blockHandlers) > 0
}

func (e *eventsSubscribersImpl) IsAddConfirmedAddedSubscribed(address *Address) bool {
	_, ok := e.confirmedAddedSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddUnconfirmedAddedSubscribed(address *Address) bool {
	_, ok := e.unconfirmedAddedSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddUnconfirmedRemovedSubscribed(address *Address) bool {
	_, ok := e.unconfirmedRemovedSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddPartialAddedSubscribed(address *Address) bool {
	_, ok := e.partialAddedSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddPartialRemovedSubscribed(address *Address) bool {
	_, ok := e.partialRemovedSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddStatusInfoSubscribed(address *Address) bool {
	_, ok := e.statusSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) IsAddCosignatureSubscribed(address *Address) bool {
	_, ok := e.cosignatureSubscribers[address.Address]
	return ok
}

func (e *eventsSubscribersImpl) GetBlockHandlers() blockHandlers {
	return e.blockHandlers
}

func (e *eventsSubscribersImpl) GetConfirmedAddedHandlers(address *Address) (confirmedAddedHandlers, error) {
	h, ok := e.confirmedAddedSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetUnconfirmedAddedHandlers(address *Address) (unconfirmedAddedHandlers, error) {
	h, ok := e.unconfirmedAddedSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetUnconfirmedRemovedHandlers(address *Address) (unconfirmedRemovedHandlers, error) {
	h, ok := e.unconfirmedRemovedSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetPartialAddedHandlers(address *Address) (partialAddedHandlers, error) {
	h, ok := e.partialAddedSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetPartialRemovedHandlers(address *Address) (partialRemovedHandlers, error) {
	h, ok := e.partialRemovedSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetStatusInfoHandlers(address *Address) (statusHandlers, error) {
	h, ok := e.statusSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

func (e *eventsSubscribersImpl) GetCosignatureHandlers(address *Address) (cosignatureHandlers, error) {
	h, ok := e.cosignatureSubscribers[address.Address]
	if !ok {
		return nil, subscriptionNotFoundError
	}

	return h, nil
}

//func (s *eventsSubscribersImpl) GetBlockChannel() (chan *BlockInfo, error) {
//	if s.blockChannel == nil {
//		return nil, channelDoesNotExistsError
//	}
//
//	return s.blockChannel, nil
//}
//
//func (s *eventsSubscribersImpl) GetConfirmedAddedChannel(address *Address) (chan Transaction, error) {
//	ch, ok := s.confirmedAddedChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetUnconfirmedAddedChannel(address *Address) (chan Transaction, error) {
//	ch, ok := s.unconfirmedAddedChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetUnconfirmedRemovedChannel(address *Address) (chan *UnconfirmedRemoved, error) {
//	ch, ok := s.unconfirmedRemovedChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetPartialAddedChannel(address *Address) (chan *AggregateTransaction, error) {
//	ch, ok := s.partialAddedChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetPartialRemovedChannel(address *Address) (chan *PartialRemovedInfo, error) {
//	ch, ok := s.partialRemovedChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetStatusInfoChannel(address *Address) (chan *StatusInfo, error) {
//	ch, ok := s.statusInfoChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) GetCosignatureChannel(address *Address) (chan *SignerInfo, error) {
//	ch, ok := s.cosignatureChannels[address.Address]
//
//	if !ok {
//		return nil, channelDoesNotExistsError
//	}
//
//	return ch, nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeBlock() error {
//	if s.blockChannel == nil {
//		return channelDoesNotExistsError
//	}
//
//	close(s.blockChannel)
//	s.blockChannel = nil
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeConfirmedAdded(address *Address) error {
//	ch, ok := s.confirmedAddedChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.confirmedAddedChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeUnconfirmedAdded(address *Address) error {
//	ch, ok := s.unconfirmedAddedChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.unconfirmedAddedChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeUnconfirmedRemoved(address *Address) error {
//	ch, ok := s.unconfirmedRemovedChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.unconfirmedRemovedChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribePartialAdded(address *Address) error {
//	ch, ok := s.partialAddedChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.partialAddedChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribePartialRemoved(address *Address) error {
//	ch, ok := s.partialRemovedChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.partialRemovedChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeStatusInfo(address *Address) error {
//	ch, ok := s.statusInfoChannels[address.Address]
//
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.statusInfoChannels, address.Address)
//
//	return nil
//}
//
//func (s *eventsSubscribersImpl) UnsubscribeCosignature(address *Address) error {
//	ch, ok := s.cosignatureChannels[address.Address]
//	if !ok {
//		return channelDoesNotExistsError
//	}
//
//	close(ch)
//	delete(s.cosignatureChannels, address.Address)
//
//	return nil
//}
