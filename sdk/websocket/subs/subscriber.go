package subs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Mapper[T any] interface {
	Map([]byte) (T, error)
}

type MapperFunc[T any] func(payload []byte) (T, error)

type SubscribersPool[T any] interface {
	Notifier
	NewSubscription(path *Path) (_ <-chan T, id int)
	CloseSubscription(path *Path, id int)
	GetPaths() []string
	HasSubscriptions(path *Path) bool
}

type subscribersPool[T any] struct {
	dataCh chan []byte

	subsPerPathsMutex sync.Mutex
	subsPerPaths      map[string]*subscriptions[T]

	mapper Mapper[T]
}

func NewSubscribersPool[T any](mapper Mapper[T]) SubscribersPool[T] {
	c := &subscribersPool[T]{
		dataCh:            make(chan []byte, 10),
		subsPerPathsMutex: sync.Mutex{},
		subsPerPaths:      make(map[string]*subscriptions[T]),
		mapper:            mapper,
	}

	return c
}

func (c *subscribersPool[T]) Notify(ctx context.Context, path *Path, payload []byte) error {
	v, err := c.mapper.Map(payload)
	if err != nil {
		return err
	}

	go func() {
		c.subsPerPathsMutex.Lock()
		defer c.subsPerPathsMutex.Unlock()

		subs, ok := c.subsPerPaths[path.String()]
		if !ok {
			return
		}

		err := subs.notify(ctx, v)
		if err != nil {
			log.Printf("Cannot notify %s: %s\n", path.String(), err)
		}
	}()

	return nil
}

func (c *subscribersPool[T]) NewSubscription(path *Path) (_ <-chan T, id int) {
	c.subsPerPathsMutex.Lock()
	defer c.subsPerPathsMutex.Unlock()

	subs, ok := c.subsPerPaths[path.String()]
	if !ok {
		subs = newSubscriptions[T]()
		c.subsPerPaths[path.String()] = subs
	}

	return subs.new()
}

func (c *subscribersPool[T]) CloseSubscription(path *Path, id int) {
	c.subsPerPathsMutex.Lock()
	defer c.subsPerPathsMutex.Unlock()

	sub, ok := c.subsPerPaths[path.String()]
	if !ok {
		return
	}
	sub.delete(id)

	if sub.length() == 0 {
		delete(c.subsPerPaths, path.String())
	}
}

func (c *subscribersPool[T]) HasSubscriptions(path *Path) bool {
	c.subsPerPathsMutex.Lock()
	defer c.subsPerPathsMutex.Unlock()

	subs, ok := c.subsPerPaths[path.String()]
	if !ok {
		return false
	}

	return subs.length() > 0
}

func (c *subscribersPool[T]) GetPaths() []string {
	c.subsPerPathsMutex.Lock()
	defer c.subsPerPathsMutex.Unlock()

	paths := make([]string, 0, len(c.subsPerPaths))
	for p := range c.subsPerPaths {
		paths = append(paths, p)
	}

	return paths
}

type subscriptions[T any] struct {
	subsMutex sync.Mutex
	subs      map[int]chan T

	randomizer *rand.Rand
}

func newSubscriptions[T any]() *subscriptions[T] {
	return &subscriptions[T]{
		subsMutex:  sync.Mutex{},
		subs:       make(map[int]chan T),
		randomizer: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (s *subscriptions[T]) new() (_ <-chan T, id int) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	for {
		id = s.randomizer.Intn(10000)
		_, ok := s.subs[id]
		if !ok {
			break
		}
	}

	ch := make(chan T)
	s.subs[id] = ch

	return ch, id
}

func (s *subscriptions[T]) delete(id int) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	delete(s.subs, id)
}

func (s *subscriptions[T]) notify(ctx context.Context, v T) error {

	errCh := make(chan error, len(s.subs))
	wg := sync.WaitGroup{}

	s.subsMutex.Lock()
	for id, sub := range s.subs {
		wg.Add(1)
		go func(id int, sub chan T) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case sub <- v:
			case <-time.After(time.Second * 30):
				close(sub)

				s.subsMutex.Lock()
				delete(s.subs, id)
				s.subsMutex.Unlock()

				errCh <- errors.New(fmt.Sprintf("Close %d subscription because deadline has expired\n", id))
			}
		}(id, sub)
	}
	s.subsMutex.Unlock()

	wg.Wait()
	close(errCh)

	var err error
	for e := range errCh {
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
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
