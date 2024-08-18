package subs

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	TopicBlock              Topic = "block"
	TopicConfirmedAdded     Topic = "confirmedAdded"
	TopicUnconfirmedAdded   Topic = "unconfirmedAdded"
	TopicUnconfirmedRemoved Topic = "unconfirmedRemoved"
	TopicStatus             Topic = "status"
	TopicPartialAdded       Topic = "partialAdded"
	TopicPartialRemoved     Topic = "partialRemoved"
	TopicCosignature        Topic = "cosignature"
	TopicDriveState         Topic = "driveState"
)

type Notifier interface {
	Notify(path *Path, payload []byte) error
}

type Publisher struct {
	subsMutex sync.Mutex

	subs map[Topic]Notifier
}

func NewPublisher() *Publisher {
	return &Publisher{
		subsMutex: sync.Mutex{},
		subs:      make(map[Topic]Notifier),
	}
}

func (p *Publisher) AddSubscriber(topic Topic, sub Notifier) error {
	p.subsMutex.Lock()
	defer p.subsMutex.Unlock()
	p.subs[topic] = sub

	return nil
}

func (p *Publisher) Publish(data []byte) error {
	p.subsMutex.Lock()
	defer p.subsMutex.Unlock()

	msgInfo, err := MapMessageInfo(data)
	if err != nil {
		return err
	}
	path := PathFromWsMessageInfo(msgInfo)

	sub, ok := p.subs[path.Topic()]
	if !ok {
		return errors.New("topic not found")
	}

	err = sub.Notify(path, data)
	if err != nil {
		return err
	}

	return nil
}

func MapMessageInfo(m []byte) (*sdk.WsMessageInfo, error) {
	var messageInfoDTO sdk.WsMessageInfoDTO
	if err := json.Unmarshal(m, &messageInfoDTO); err != nil {
		return nil, errors.Wrap(err, "unmarshaling message info data")
	}

	i, err := messageInfoDTO.ToStruct()
	if err != nil {
		println(string(m))
	}

	return i, err
}
