// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	hdlrs "github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/handlers"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
	"golang.org/x/net/websocket"
	"io"
	"sync"
)

type pathType string

const (
	pathBlock              pathType = "block"
	pathConfirmedAdded     pathType = "confirmedAdded"
	pathUnconfirmedAdded   pathType = "unconfirmedAdded"
	pathUnconfirmedRemoved pathType = "unconfirmedRemoved"
	pathStatus             pathType = "status"
	pathPartialAdded       pathType = "partialAdded"
	pathPartialRemoved     pathType = "partialRemoved"
	pathCosignature        pathType = "cosignature"
)

var (
	unsupportedMessageTypeError = errors.New("unsupported message type")
)

func NewCatapultWebSocketClient(endpoint string) (CatapultWebsocketClient, error) {
	conn, err := websocket.Dial(endpoint, "tcp", "http://localhost")
	if err != nil {
		return nil, err
	}

	var raw []byte
	if err := websocket.Message.Receive(conn, &raw); err != nil {
		return nil, err
	}

	var resp *sdk.WsConnectionResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	return &CatapultWebsocketClientImpl{
		conn: conn,
		UID:  resp.Uid,

		topicHandlers:    make(map[pathType]topicHandler),
		messagePublisher: newMessagePublisher(conn),
		errorsChan:       make(chan error, 1000),
	}, nil
}

type WebsocketClient interface {
	Listen(wg *sync.WaitGroup)
	GetErrorsChan() (chan error, error)
}

type CatapultWebsocketClient interface {
	WebsocketClient

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

	topicHandlers map[pathType]topicHandler

	messagePublisher messagePublisher

	errorsChan chan error
}

func (c *CatapultWebsocketClientImpl) Listen(wg *sync.WaitGroup) {
	var resp []byte
	for {
		err := websocket.Message.Receive(c.conn, &resp)

		if err == io.EOF {
			c.errorsChan <- errors.Wrap(err, "error websocket disconnect ")
			wg.Done()
			return
		}

		if err != nil {
			c.errorsChan <- errors.Wrap(err, "error receiving message from websocket")
			wg.Done()
			return
		}

		messageInfo, err := c.getMessageInfo(resp)
		if err != nil {
			c.errorsChan <- errors.Wrap(err, "error getting address and channel name from websocket message")
			continue
		}

		go c.routeMessage(messageInfo, resp)
	}
}

func (c *CatapultWebsocketClientImpl) AddBlockHandlers(handlers ...subscribers.BlockHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.blockSubscriber == nil {
		c.blockSubscriber = subscribers.NewBlockSubscriber()
	}

	if _, ok := c.topicHandlers[pathBlock]; !ok {
		c.topicHandlers[pathBlock] = topicHandler{
			Handler: hdlrs.NewBlockHandler(sdk.BlockProcessorFn(sdk.ProcessBlock), c.blockSubscriber, c.errorsChan),
			Topic:   topicFormatFn(formatBlockTopic),
		}
	}

	if !c.blockSubscriber.HasHandlers() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathBlock); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.blockSubscriber.AddHandlers(handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddConfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.ConfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.confirmedAddedSubscribers == nil {
		c.confirmedAddedSubscribers = subscribers.NewConfirmedAdded()
	}

	if _, ok := c.topicHandlers[pathConfirmedAdded]; !ok {
		c.topicHandlers[pathConfirmedAdded] = topicHandler{
			Handler: hdlrs.NewConfirmedAddedHandler(sdk.NewConfirmedAddedProcessor(sdk.MapTransaction), c.confirmedAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.confirmedAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathConfirmedAdded, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.confirmedAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.unconfirmedAddedSubscribers == nil {
		c.unconfirmedAddedSubscribers = subscribers.NewUnconfirmedAdded()
	}

	if _, ok := c.topicHandlers[pathUnconfirmedAdded]; !ok {
		c.topicHandlers[pathUnconfirmedAdded] = topicHandler{
			Handler: hdlrs.NewUnconfirmedAddedHandler(sdk.NewUnconfirmedAddedProcessor(sdk.MapTransaction), c.unconfirmedAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.unconfirmedAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedRemovedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.unconfirmedRemovedSubscribers == nil {
		c.unconfirmedRemovedSubscribers = subscribers.NewUnconfirmedRemoved()
	}

	if _, ok := c.topicHandlers[pathUnconfirmedRemoved]; !ok {
		c.topicHandlers[pathUnconfirmedRemoved] = topicHandler{
			Handler: hdlrs.NewUnconfirmedRemovedHandler(sdk.UnconfirmedRemovedProcessorFn(sdk.ProcessUnconfirmedRemoved), c.unconfirmedRemovedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.unconfirmedRemovedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedRemovedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialAddedHandlers(address *sdk.Address, handlers ...subscribers.PartialAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.partialAddedSubscribers == nil {
		c.partialAddedSubscribers = subscribers.NewPartialAdded()
	}

	if _, ok := c.topicHandlers[pathPartialAdded]; !ok {
		c.topicHandlers[pathPartialAdded] = topicHandler{
			Handler: hdlrs.NewPartialAddedHandler(sdk.NewPartialAddedProcessor(sdk.MapTransaction), c.partialAddedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.partialAddedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathPartialAdded, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.partialAddedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialRemovedHandlers(address *sdk.Address, handlers ...subscribers.PartialRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.partialRemovedSubscribers == nil {
		c.partialRemovedSubscribers = subscribers.NewPartialRemoved()
	}

	if _, ok := c.topicHandlers[pathPartialRemoved]; !ok {
		c.topicHandlers[pathPartialRemoved] = topicHandler{
			Handler: hdlrs.NewPartialRemovedHandler(sdk.PartialRemovedProcessorFn(sdk.ProcessPartialRemoved), c.partialRemovedSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.partialRemovedSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathPartialRemoved, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.partialRemovedSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddStatusHandlers(address *sdk.Address, handlers ...subscribers.StatusHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.statusSubscribers == nil {
		c.statusSubscribers = subscribers.NewStatus()
	}

	if _, ok := c.topicHandlers[pathStatus]; !ok {
		c.topicHandlers[pathStatus] = topicHandler{
			Handler: hdlrs.NewStatusHandler(sdk.StatusProcessorFn(sdk.ProcessStatus), c.statusSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.statusSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathStatus, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	err := c.statusSubscribers.AddHandlers(address, handlers...)
	if err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddCosignatureHandlers(address *sdk.Address, handlers ...subscribers.CosignatureHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if c.cosignatureSubscribers == nil {
		c.cosignatureSubscribers = subscribers.NewCosignature()
	}

	if _, ok := c.topicHandlers[pathCosignature]; !ok {
		c.topicHandlers[pathCosignature] = topicHandler{
			Handler: hdlrs.NewCosignatureHandler(sdk.CosignatureProcessorFn(sdk.ProcessCosignature), c.cosignatureSubscribers, c.errorsChan),
			Topic:   topicFormatFn(formatPlainTopic),
		}
	}

	if !c.cosignatureSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathType(fmt.Sprintf("%s/%s", pathCosignature, address.Address))); err != nil {
			return errors.Wrap(err, "error publishing subscribe message into websocket")
		}
	}

	if err := c.cosignatureSubscribers.AddHandlers(address, handlers...); err != nil {
		return errors.Wrap(err, "error adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) GetErrorsChan() (chan error, error) {
	return c.errorsChan, nil
}

func (c *CatapultWebsocketClientImpl) getMessageInfo(m []byte) (*sdk.WsMessageInfo, error) {
	var messageInfoDTO sdk.WsMessageInfoDTO
	if err := json.Unmarshal(m, &messageInfoDTO); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling message info data")
	}

	return messageInfoDTO.ToStruct()
}

func (c *CatapultWebsocketClientImpl) routeMessage(messageInfo *sdk.WsMessageInfo, resp []byte) {
	handler, ok := c.topicHandlers[pathType(messageInfo.ChannelName)]
	if !ok {
		c.errorsChan <- errors.Wrap(unsupportedMessageTypeError, "error getting topic handler from topic handlers storage")
		return
	}

	if ok := handler.Handle(messageInfo.Address, resp); !ok {
		if err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, pathType(handler.Format(messageInfo))); err != nil {
			c.errorsChan <- errors.Wrap(err, "error unsubscribing from topic")
		}
	}

	return
}

type Topic interface {
	Format(info *sdk.WsMessageInfo) pathType
}

type topicFormatFn func(info *sdk.WsMessageInfo) pathType

func (ref topicFormatFn) Format(info *sdk.WsMessageInfo) pathType {
	return ref(info)
}

func formatPlainTopic(info *sdk.WsMessageInfo) pathType {
	return pathType(fmt.Sprintf("%s/%s", pathCosignature, info.Address.Address))
}

func formatBlockTopic(_ *sdk.WsMessageInfo) pathType {
	return pathBlock
}

type topicHandler struct {
	Topic
	hdlrs.Handler
}
