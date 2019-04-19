// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	hdlrs "github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/handlers"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"io"
	"sync"
)

type Path string

const (
	pathBlock              Path = "block"
	pathConfirmedAdded     Path = "confirmedAdded"
	pathUnconfirmedAdded   Path = "unconfirmedAdded"
	pathUnconfirmedRemoved Path = "unconfirmedRemoved"
	pathStatus             Path = "status"
	pathPartialAdded       Path = "partialAdded"
	pathPartialRemoved     Path = "partialRemoved"
	pathCosignature        Path = "cosignature"
)

var (
	unsupportedMessageTypeError = errors.New("unsupported message type")
)

func NewClient(endpoint string) (CatapultClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp *wsConnectionResponse
	err = conn.ReadJSON(&resp)
	if err != nil {
		return nil, err
	}

	topicHandlers := make(topicHandlers)
	errorsChan := make(chan error, 1000)
	messagePublisher := newMessagePublisher(conn)
	messageRouter := NewRouter(resp.Uid, messagePublisher, topicHandlers, errorsChan)

	return &CatapultWebsocketClientImpl{
		conn: conn,
		UID:  resp.Uid,

		blockSubscriber:               subscribers.NewBlock(),
		confirmedAddedSubscribers:     subscribers.NewConfirmedAdded(),
		unconfirmedAddedSubscribers:   subscribers.NewUnconfirmedAdded(),
		unconfirmedRemovedSubscribers: subscribers.NewUnconfirmedRemoved(),
		partialAddedSubscribers:       subscribers.NewPartialAdded(),
		partialRemovedSubscribers:     subscribers.NewPartialRemoved(),
		statusSubscribers:             subscribers.NewStatus(),
		cosignatureSubscribers:        subscribers.NewCosignature(),

		topicHandlers:    topicHandlers,
		messageRouter:    messageRouter,
		messagePublisher: messagePublisher,
		errorsChan:       errorsChan,
	}, nil
}

type Client interface {
	Listen(wg *sync.WaitGroup)
	GetErrorsChan() chan error
}

type CatapultClient interface {
	Client

	AddBlockHandlers(handlers ...subscribers.BlockHandler) error
	AddConfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.ConfirmedAddedHandler) error
	AddUnconfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedAddedHandler) error
	AddUnconfirmedRemovedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedRemovedHandler) error
	AddPartialAddedHandlers(address *sdk.Address, handlers ...subscribers.PartialAddedHandler) error
	AddPartialRemovedHandlers(address *sdk.Address, handlers ...subscribers.PartialRemovedHandler) error
	AddStatusHandlers(address *sdk.Address, handlers ...subscribers.StatusHandler) error
	AddCosignatureHandlers(address *sdk.Address, handlers ...subscribers.CosignatureHandler) error
}

type CatapultWebsocketClientImpl struct {
	conn *websocket.Conn
	UID  string

	blockSubscriber               subscribers.Block
	confirmedAddedSubscribers     subscribers.ConfirmedAdded
	unconfirmedAddedSubscribers   subscribers.UnconfirmedAdded
	unconfirmedRemovedSubscribers subscribers.UnconfirmedRemoved
	partialAddedSubscribers       subscribers.PartialAdded
	partialRemovedSubscribers     subscribers.PartialRemoved
	statusSubscribers             subscribers.Status
	cosignatureSubscribers        subscribers.Cosignature

	topicHandlers    TopicHandlersStorage
	messageRouter    Router
	messagePublisher MessagePublisher

	errorsChan chan error
}

func (c *CatapultWebsocketClientImpl) Listen(wg *sync.WaitGroup) {
	for {
		_, resp, err := c.conn.ReadMessage()

		if err == io.EOF {
			c.errorsChan <- errors.Wrap(err, "websocket disconnect ")
			wg.Done()
			return
		}

		if err != nil {
			c.errorsChan <- errors.Wrap(err, "receiving message from websocket")
			wg.Done()
			return
		}

		go c.messageRouter.RouteMessage(resp)
	}
}

func (c *CatapultWebsocketClientImpl) AddBlockHandlers(handlers ...subscribers.BlockHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathBlock) {
		c.topicHandlers.SetTopicHandler(pathBlock, &TopicHandler{
			Handler: hdlrs.NewBlockHandler(sdk.BlockMapperFn(sdk.MapBlock), c.blockSubscriber, c.errorsChan),
			Topic:   topicFormatFn(formatBlockTopic),
		})
	}

	if !c.blockSubscriber.HasHandlers() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathBlock); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.blockSubscriber.AddHandlers(handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddConfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.ConfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathConfirmedAdded) {
		c.topicHandlers.SetTopicHandler(pathConfirmedAdded, &TopicHandler{
			Handler: hdlrs.NewConfirmedAddedHandler(sdk.NewConfirmedAddedMapper(sdk.MapTransaction), c.confirmedAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.confirmedAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.confirmedAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathUnconfirmedAdded) {
		c.topicHandlers.SetTopicHandler(pathUnconfirmedAdded, &TopicHandler{
			Handler: hdlrs.NewUnconfirmedAddedHandler(sdk.NewUnconfirmedAddedMapper(sdk.MapTransaction), c.unconfirmedAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.unconfirmedAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedRemovedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathUnconfirmedRemoved) {
		c.topicHandlers.SetTopicHandler(pathUnconfirmedRemoved, &TopicHandler{
			Handler: hdlrs.NewUnconfirmedRemovedHandler(sdk.UnconfirmedRemovedMapperFn(sdk.MapUnconfirmedRemoved), c.unconfirmedRemovedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.unconfirmedRemovedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedRemovedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialAddedHandlers(address *sdk.Address, handlers ...subscribers.PartialAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathPartialAdded) {
		c.topicHandlers.SetTopicHandler(pathPartialAdded, &TopicHandler{
			Handler: hdlrs.NewPartialAddedHandler(sdk.NewPartialAddedMapper(sdk.MapTransaction), c.partialAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.partialAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.partialAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialRemovedHandlers(address *sdk.Address, handlers ...subscribers.PartialRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathPartialRemoved) {
		c.topicHandlers.SetTopicHandler(pathPartialRemoved, &TopicHandler{
			Handler: hdlrs.NewPartialRemovedHandler(sdk.PartialRemovedMapperFn(sdk.MapPartialRemoved), c.partialRemovedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.partialRemovedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.partialRemovedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddStatusHandlers(address *sdk.Address, handlers ...subscribers.StatusHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathStatus) {
		c.topicHandlers.SetTopicHandler(pathStatus, &TopicHandler{
			Handler: hdlrs.NewStatusHandler(sdk.StatusMapperFn(sdk.MapStatus), c.statusSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.statusSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.statusSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddCosignatureHandlers(address *sdk.Address, handlers ...subscribers.CosignatureHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathCosignature) {
		c.topicHandlers.SetTopicHandler(pathCosignature, &TopicHandler{
			Handler: hdlrs.NewCosignatureHandler(sdk.CosignatureMapperFn(sdk.MapCosignature), c.cosignatureSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.cosignatureSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.cosignatureSubscribers.AddHandlers(address, handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) GetErrorsChan() chan error {
	return c.errorsChan
}
