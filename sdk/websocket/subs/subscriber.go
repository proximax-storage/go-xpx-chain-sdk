package subs

import (
	"log"
	"sync"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type Mapper[T any] interface {
	Map([]byte) (T, error)
}

type MapperFunc[T any] func(payload []byte) (T, error)

type SubscribersPool[T any] interface {
	Notifier
	NewSubscription(address string) (_ <-chan T, id int)
	CloseSubscription(address string, id int)
	GetAddresses() []string
	HasSubscriptions(address string) bool
}

type subscribersPool[T any] struct {
	dataCh chan []byte

	subsPerAddressMutex sync.Mutex
	subsPerAddress      map[string]*subscriptions[T]

	mapper Mapper[T]
}

func NewSubscribersPool[T any](mapper Mapper[T]) SubscribersPool[T] {
	c := &subscribersPool[T]{
		dataCh:              make(chan []byte, 10),
		subsPerAddressMutex: sync.Mutex{},
		subsPerAddress:      make(map[string]*subscriptions[T]),
		mapper:              mapper,
	}

	return c
}

func (c *subscribersPool[T]) Notify(address *sdk.Address, payload []byte) error {
	v, err := c.mapper.Map(payload)
	if err != nil {
		return err
	}

	go func() {
		c.subsPerAddressMutex.Lock()
		subs, ok := c.subsPerAddress[address.Address]
		c.subsPerAddressMutex.Unlock()
		if !ok {
			return
		}
		subs.notify(v)
	}()

	return nil
}

func (c *subscribersPool[T]) NewSubscription(address string) (_ <-chan T, id int) {
	c.subsPerAddressMutex.Lock()
	subs, ok := c.subsPerAddress[address]
	if !ok {
		subs = newSubscriptions[T]()
		c.subsPerAddress[address] = subs
	}
	c.subsPerAddressMutex.Unlock()

	return subs.new()
}

func (c *subscribersPool[T]) CloseSubscription(address string, id int) {
	c.subsPerAddressMutex.Lock()
	sub, ok := c.subsPerAddress[address]
	c.subsPerAddressMutex.Unlock()
	if !ok {
		return
	}
	sub.delete(id)

	if sub.length() == 0 {
		delete(c.subsPerAddress, address)
	}
}

func (c *subscribersPool[T]) HasSubscriptions(address string) bool {
	c.subsPerAddressMutex.Lock()
	defer c.subsPerAddressMutex.Unlock()
	subs, ok := c.subsPerAddress[address]
	if !ok {
		return false
	}

	return subs.length() > 0
}

func (c *subscribersPool[T]) GetAddresses() []string {
	c.subsPerAddressMutex.Lock()
	defer c.subsPerAddressMutex.Unlock()

	addresses := make([]string, 0, len(c.subsPerAddress))
	for addr := range c.subsPerAddress {
		addresses = append(addresses, addr)
	}

	return addresses
}

type subscriptions[T any] struct {
	subsMutex sync.Mutex
	subs      map[int]chan T
}

func newSubscriptions[T any]() *subscriptions[T] {
	return &subscriptions[T]{
		subsMutex: sync.Mutex{},
		subs:      make(map[int]chan T),
	}
}

func (s *subscriptions[T]) new() (_ <-chan T, id int) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	ch := make(chan T)
	s.subs[0] = ch

	return ch, id
}

func (s *subscriptions[T]) delete(id int) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	delete(s.subs, id)
}

func (s *subscriptions[T]) notify(v T) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	for id, sub := range s.subs {
		select {
		case sub <- v:
		case <-time.After(time.Second * 30):
			log.Println("deadline")

			close(sub)
			delete(s.subs, id)
		}
	}
}

func (s *subscriptions[T]) getAll() map[int]chan T {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	res := make(map[int]chan T)
	for i, sub := range s.subs {
		res[i] = sub
	}

	return res
}

func (s *subscriptions[T]) length() int {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	return len(s.subs)
}
