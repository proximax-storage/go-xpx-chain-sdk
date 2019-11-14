// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	hdlrs "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/handlers"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

const pathWS = "ws"

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
	driveState             Path = "driveState"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

func NewClient(ctx context.Context, cfg *sdk.Config) (CatapultClient, error) {
	ctx, cancelFunc := context.WithCancel(ctx)

	conn, uid, err := connect(cfg)
	if err != nil {
		return nil, err
	}

	topicHandlers := make(topicHandlers)
	messagePublisher := newMessagePublisher(conn)
	messageRouter := NewRouter(uid, messagePublisher, topicHandlers)

	return &CatapultWebsocketClientImpl{
		config:     cfg,
		conn:       conn,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		connectFn:  connect,
		UID:        uid,

		blockSubscriber:               subscribers.NewBlock(),
		confirmedAddedSubscribers:     subscribers.NewConfirmedAdded(),
		unconfirmedAddedSubscribers:   subscribers.NewUnconfirmedAdded(),
		unconfirmedRemovedSubscribers: subscribers.NewUnconfirmedRemoved(),
		partialAddedSubscribers:       subscribers.NewPartialAdded(),
		partialRemovedSubscribers:     subscribers.NewPartialRemoved(),
		statusSubscribers:             subscribers.NewStatus(),
		cosignatureSubscribers:        subscribers.NewCosignature(),
		driveStateSubscribers:         subscribers.NewDriveState(),

		topicHandlers:    topicHandlers,
		messageRouter:    messageRouter,
		messagePublisher: messagePublisher,
	}, nil
}

type Client interface {
	io.Closer

	Listen()
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
	AddDriveStateHandlers(address *sdk.Address, handlers ...subscribers.DriveStateHandler) error
}

type CatapultWebsocketClientImpl struct {
	config     *sdk.Config
	conn       *websocket.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
	UID        string

	blockSubscriber               subscribers.Block
	confirmedAddedSubscribers     subscribers.ConfirmedAdded
	unconfirmedAddedSubscribers   subscribers.UnconfirmedAdded
	unconfirmedRemovedSubscribers subscribers.UnconfirmedRemoved
	partialAddedSubscribers       subscribers.PartialAdded
	partialRemovedSubscribers     subscribers.PartialRemoved
	statusSubscribers             subscribers.Status
	cosignatureSubscribers        subscribers.Cosignature
	driveStateSubscribers         subscribers.DriveState

	topicHandlers    TopicHandlersStorage
	messageRouter    Router
	messagePublisher MessagePublisher

	connectFn        func(cfg *sdk.Config) (*websocket.Conn, string, error)
	alreadyListening bool
}

func (c *CatapultWebsocketClientImpl) Listen() {
	if c.alreadyListening {
		return
	}

	c.alreadyListening = true

	messagesChan := make(chan []byte)

	go func() {
		defer c.cancelFunc()

	ReadMessageLoop:
		for {
			_, resp, e := c.conn.ReadMessage()
			if e != nil {
				if _, ok := e.(*net.OpError); ok {
					// Stop ReadMessage goroutine if user called Close function for websocket client
					return
				}

				if _, ok := e.(*websocket.CloseError); ok {
					// Start websocket reconnect processing if connection was closed
					for range time.NewTicker(c.config.WsReconnectionTimeout).C {
						if err := c.reconnect(); err != nil {
							continue
						}

						continue ReadMessageLoop
					}
				}

				return
			}

			messagesChan <- resp
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			if c.conn != nil {
				if err := c.conn.Close(); err != nil {
					panic(err)
				}
				c.conn = nil
			}
			return
		case msg := <-messagesChan:
			go c.messageRouter.RouteMessage(msg)
		}
	}
}

func (c *CatapultWebsocketClientImpl) AddBlockHandlers(handlers ...subscribers.BlockHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathBlock) {
		c.topicHandlers.SetTopicHandler(pathBlock, &TopicHandler{
			Handler: hdlrs.NewBlockHandler(sdk.BlockMapperFn(sdk.MapBlock), c.blockSubscriber),
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
			Handler: hdlrs.NewConfirmedAddedHandler(sdk.NewConfirmedAddedMapper(sdk.MapTransaction), c.confirmedAddedSubscribers),
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
			Handler: hdlrs.NewUnconfirmedAddedHandler(sdk.NewUnconfirmedAddedMapper(sdk.MapTransaction), c.unconfirmedAddedSubscribers),
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
			Handler: hdlrs.NewUnconfirmedRemovedHandler(sdk.UnconfirmedRemovedMapperFn(sdk.MapUnconfirmedRemoved), c.unconfirmedRemovedSubscribers),
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
			Handler: hdlrs.NewPartialAddedHandler(sdk.NewPartialAddedMapper(sdk.MapTransaction), c.partialAddedSubscribers),
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
			Handler: hdlrs.NewPartialRemovedHandler(sdk.PartialRemovedMapperFn(sdk.MapPartialRemoved), c.partialRemovedSubscribers),
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
			Handler: hdlrs.NewStatusHandler(sdk.StatusMapperFn(sdk.MapStatus), c.statusSubscribers),
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
			Handler: hdlrs.NewCosignatureHandler(sdk.CosignatureMapperFn(sdk.MapCosignature), c.cosignatureSubscribers),
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


func (c *CatapultWebsocketClientImpl) AddDriveStateHandlers(address *sdk.Address, handlers ...subscribers.DriveStateHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(driveState) {
		c.topicHandlers.SetTopicHandler(driveState, &TopicHandler{
			Handler: hdlrs.NewDriveStateHandler(sdk.DriveStateMapperFn(sdk.MapDriveState), c.driveStateSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.driveStateSubscribers.HasHandlers(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", driveState, address.Address))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.driveStateSubscribers.AddHandlers(address, handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) reconnect() error {

	conn, uid, err := c.connectFn(c.config)
	if err != nil {
		return err
	}

	c.conn = conn
	c.UID = uid

	c.messagePublisher.SetConn(conn)
	c.messageRouter.SetUid(uid)

	if c.blockSubscriber.HasHandlers() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s", pathBlock))); err != nil {
			return err
		}
	}

	for _, value := range c.confirmedAddedSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.cosignatureSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, value))); err != nil {
			return err
		}
	}

	for _, value := range c.partialAddedSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.partialRemovedSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, value))); err != nil {
			return err
		}
	}

	for _, value := range c.statusSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, value))); err != nil {
			return err
		}
	}

	for _, value := range c.unconfirmedAddedSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.unconfirmedRemovedSubscribers.GetAddresses() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, value))); err != nil {
			return err
		}
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}

		c.conn = nil
	}

	c.cancelFunc()

	c.alreadyListening = false

	c.blockSubscriber = nil
	c.confirmedAddedSubscribers = nil
	c.unconfirmedAddedSubscribers = nil
	c.unconfirmedRemovedSubscribers = nil
	c.partialAddedSubscribers = nil
	c.partialRemovedSubscribers = nil
	c.statusSubscribers = nil
	c.cosignatureSubscribers = nil

	c.topicHandlers = nil

	return nil
}

func connect(cfg *sdk.Config) (*websocket.Conn, string, error) {
	var conn *websocket.Conn
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(convertToWsUrl(cfg.UsedBaseUrl).String(), nil)
	if err != nil {
		for _, u := range cfg.BaseURLs {

			if u == cfg.UsedBaseUrl {
				continue
			}

			conn, _, err = websocket.DefaultDialer.Dial(convertToWsUrl(u).String(), nil)
			if err != nil {
				continue
			}

			cfg.UsedBaseUrl = u
			break
		}
	}

	if conn == nil {
		return nil, "", err
	}

	resp := new(wsConnectionResponse)
	if err = conn.ReadJSON(resp); err != nil {
		return nil, "", err
	}

	return conn, resp.Uid, nil
}

func convertToWsUrl(url *url.URL) *url.URL {
	copyUrl := *url
	copyUrl.Scheme = "ws" // always ws
	copyUrl.Path = pathWS
	return &copyUrl
}
