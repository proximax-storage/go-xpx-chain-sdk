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
	pathStateStatement     Path = "stateStatement"
	pathPublicKeyStatement Path = "publicKeyStatement"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

type (
	// Subscribe path
	Path string

	CatapultWebsocketClientImpl struct {
		ctx        context.Context
		cancelFunc context.CancelFunc

		UID    string
		config *sdk.Config

		conn *websocket.Conn

		blockSubscriber               subscribers.Block
		statusSubscribers             subscribers.Status
		cosignatureSubscribers        subscribers.Cosignature
		driveStateSubscribers         subscribers.DriveState
		partialAddedSubscribers       subscribers.PartialAdded
		partialRemovedSubscribers     subscribers.PartialRemoved
		confirmedAddedSubscribers     subscribers.ConfirmedAdded
		unconfirmedAddedSubscribers   subscribers.UnconfirmedAdded
		unconfirmedRemovedSubscribers subscribers.UnconfirmedRemoved
		stateStatementSubscribers     subscribers.Receipt
		publicKeyStatementSubscribers subscribers.Receipt

		messageRouter    Router
		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher

		// connectionStatusCh chan bool
		listenCh     chan bool            // channel for manage current listen status for connection
		reconnectCh  chan *websocket.Conn // channel for connection with we will close, and open new connection
		connectionCh chan *websocket.Conn // channel for new opened connection

		connectFn func(cfg *sdk.Config) (*websocket.Conn, string, error)
	}

	Client interface {
		io.Closer

		Listen()
	}

	CatapultClient interface {
		Client

		Config() *sdk.Config
		AddBlockHandlers(handlers ...subscribers.BlockHandler) error
		AddConfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.ConfirmedAddedHandler) error
		AddUnconfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedAddedHandler) error
		AddUnconfirmedRemovedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedRemovedHandler) error
		AddPartialAddedHandlers(address *sdk.Address, handlers ...subscribers.PartialAddedHandler) error
		AddPartialRemovedHandlers(address *sdk.Address, handlers ...subscribers.PartialRemovedHandler) error
		AddConfirmedAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.ConfirmedAddedHandler) error
		AddUnconfirmedAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.UnconfirmedAddedHandler) error
		AddUnconfirmedRemovedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.UnconfirmedRemovedHandler) error
		AddPartialAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.PartialAddedHandler) error
		AddPartialRemovedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.PartialRemovedHandler) error
		AddStatusHandlers(address *sdk.Address, handlers ...subscribers.StatusHandler) error
		AddCosignatureHandlers(address *sdk.Address, handlers ...subscribers.CosignatureHandler) error
		AddDriveStateHandlers(address *sdk.Address, handlers ...subscribers.DriveStateHandler) error
		AddPublicKeyStatementHandlers(entityType *sdk.EntityType, handlers ...subscribers.ReceiptHandler) error
		AddStateStatementHandlers(entityType *sdk.EntityType, handlers ...subscribers.ReceiptHandler) error
	}
)

func NewClient(ctx context.Context, cfg *sdk.Config) (CatapultClient, error) {
	ctx, cancelFunc := context.WithCancel(ctx)

	socketClient := &CatapultWebsocketClientImpl{
		ctx:        ctx,
		cancelFunc: cancelFunc,

		config: cfg,

		blockSubscriber:               subscribers.NewBlock(),
		statusSubscribers:             subscribers.NewStatus(),
		cosignatureSubscribers:        subscribers.NewCosignature(),
		driveStateSubscribers:         subscribers.NewDriveState(),
		partialAddedSubscribers:       subscribers.NewPartialAdded(),
		partialRemovedSubscribers:     subscribers.NewPartialRemoved(),
		confirmedAddedSubscribers:     subscribers.NewConfirmedAdded(),
		unconfirmedAddedSubscribers:   subscribers.NewUnconfirmedAdded(),
		unconfirmedRemovedSubscribers: subscribers.NewUnconfirmedRemoved(),

		topicHandlers: &topicHandlers{h: make(topicHandlersMap)},

		listenCh:     make(chan bool),
		reconnectCh:  make(chan *websocket.Conn),
		connectionCh: make(chan *websocket.Conn),

		connectFn: connect,
	}

	go socketClient.handleSignal()

	if err := socketClient.initNewConnection(); err != nil {
		return socketClient, err
	}

	return socketClient, nil
}

func (c *CatapultWebsocketClientImpl) Listen() {

	c.listenCh <- true

	select {
	case <-c.ctx.Done():
		c.closeConnection(c.conn)
	}
}

func (c *CatapultWebsocketClientImpl) Close() error {
	c.cancelFunc()
	return nil
}

func (c *CatapultWebsocketClientImpl) Config() *sdk.Config {
	return c.config
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

func (c *CatapultWebsocketClientImpl) AddStateStatementHandlers(entityType *sdk.EntityType, handlers ...subscribers.ReceiptHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	handle := sdk.NewCompoundChannelHandleFromEntityType(*entityType)
	if !c.topicHandlers.HasHandler(pathStateStatement) {
		c.topicHandlers.SetTopicHandler(pathStateStatement, &TopicHandler{
			Handler: hdlrs.NewReceiptHandler(sdk.ReceiptMapperFn(sdk.MapReceipt), c.stateStatementSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.stateStatementSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathStateStatement); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.stateStatementSubscribers.AddHandlers(handle, handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPublicKeyStatementHandlers(entityType *sdk.EntityType, handlers ...subscribers.ReceiptHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	handle := sdk.NewCompoundChannelHandleFromEntityType(*entityType)
	if !c.topicHandlers.HasHandler(pathPublicKeyStatement) {
		c.topicHandlers.SetTopicHandler(pathPublicKeyStatement, &TopicHandler{
			Handler: hdlrs.NewReceiptHandler(sdk.ReceiptMapperFn(sdk.MapReceipt), c.publicKeyStatementSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.publicKeyStatementSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathPublicKeyStatement); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.publicKeyStatementSubscribers.AddHandlers(handle, handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddConfirmedAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.ConfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathConfirmedAdded) {
		c.topicHandlers.SetTopicHandler(pathConfirmedAdded, &TopicHandler{
			Handler: hdlrs.NewConfirmedAddedHandler(sdk.NewConfirmedAddedMapper(sdk.MapTransaction, c.config.GenerationHash), c.confirmedAddedSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.confirmedAddedSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.confirmedAddedSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddConfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.ConfirmedAddedHandler) error {
	return c.AddConfirmedAddedHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.UnconfirmedAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathUnconfirmedAdded) {
		c.topicHandlers.SetTopicHandler(pathUnconfirmedAdded, &TopicHandler{
			Handler: hdlrs.NewUnconfirmedAddedHandler(sdk.NewUnconfirmedAddedMapper(sdk.MapTransaction, c.config.GenerationHash), c.unconfirmedAddedSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.unconfirmedAddedSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedAddedSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}
func (c *CatapultWebsocketClientImpl) AddUnconfirmedAddedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedAddedHandler) error {
	return c.AddUnconfirmedAddedHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedRemovedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.UnconfirmedRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathUnconfirmedRemoved) {
		c.topicHandlers.SetTopicHandler(pathUnconfirmedRemoved, &TopicHandler{
			Handler: hdlrs.NewUnconfirmedRemovedHandler(sdk.UnconfirmedRemovedMapperFn(sdk.MapUnconfirmedRemoved), c.unconfirmedRemovedSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.unconfirmedRemovedSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.unconfirmedRemovedSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedRemovedHandlers(address *sdk.Address, handlers ...subscribers.UnconfirmedRemovedHandler) error {
	return c.AddUnconfirmedRemovedHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddPartialAddedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.PartialAddedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathPartialAdded) {
		c.topicHandlers.SetTopicHandler(pathPartialAdded, &TopicHandler{
			Handler: hdlrs.NewPartialAddedHandler(sdk.NewPartialAddedMapper(sdk.MapTransaction, c.config.GenerationHash), c.partialAddedSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.partialAddedSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.partialAddedSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialAddedHandlers(address *sdk.Address, handlers ...subscribers.PartialAddedHandler) error {
	return c.AddPartialAddedHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddPartialRemovedHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.PartialRemovedHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathPartialRemoved) {
		c.topicHandlers.SetTopicHandler(pathPartialRemoved, &TopicHandler{
			Handler: hdlrs.NewPartialRemovedHandler(sdk.PartialRemovedMapperFn(sdk.MapPartialRemoved), c.partialRemovedSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.partialRemovedSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.partialRemovedSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialRemovedHandlers(address *sdk.Address, handlers ...subscribers.PartialRemovedHandler) error {
	return c.AddPartialRemovedHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddStatusHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.StatusHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathStatus) {
		c.topicHandlers.SetTopicHandler(pathStatus, &TopicHandler{
			Handler: hdlrs.NewStatusHandler(sdk.StatusMapperFn(sdk.MapStatus), c.statusSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.statusSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	err := c.statusSubscribers.AddHandlers(handle, handlers...)
	if err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddStatusHandlers(address *sdk.Address, handlers ...subscribers.StatusHandler) error {
	return c.AddStatusHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
}

func (c *CatapultWebsocketClientImpl) AddCosignatureHandlersByHandle(handle *sdk.CompoundChannelHandle, handlers ...subscribers.CosignatureHandler) error {
	if len(handlers) == 0 {
		return nil
	}

	if !c.topicHandlers.HasHandler(pathCosignature) {
		c.topicHandlers.SetTopicHandler(pathCosignature, &TopicHandler{
			Handler: hdlrs.NewCosignatureHandler(sdk.CosignatureMapperFn(sdk.MapCosignature), c.cosignatureSubscribers),
			Topic:   topicFormatFn(formatPlainTopic),
		})
	}

	if !c.cosignatureSubscribers.HasHandlers(handle) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, handle.String()))); err != nil {
			return errors.Wrap(err, "publishing subscribe message into websocket")
		}
	}

	if err := c.cosignatureSubscribers.AddHandlers(handle, handlers...); err != nil {
		return errors.Wrap(err, "adding handlers functions into handlers storage")
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddCosignatureHandlers(address *sdk.Address, handlers ...subscribers.CosignatureHandler) error {
	return c.AddCosignatureHandlersByHandle(sdk.NewCompoundChannelHandleFromAddress(address), handlers...)
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

func (c *CatapultWebsocketClientImpl) handleSignal() {
	for {
		select {
		case conn := <-c.connectionCh:
			c.conn = conn
		case conn := <-c.reconnectCh:
			c.closeConnection(conn)
			go func() {
				c.listenCh <- false
			}()

		case <-c.listenCh:

			if c.conn == nil {
				err := c.initNewConnection()
				if err != nil {
					fmt.Println("websocket: connection is failed. Try again after wait period")
					select {
					case <-time.NewTicker(c.config.WsReconnectionTimeout).C:
						go func() {
							c.listenCh <- true
						}()
						continue
					}
				}
			}

			err := c.updateHandlers()
			if err != nil {
				fmt.Println("websocket: update handles is failed. Try again after timeout period")
				select {
				case <-time.NewTicker(c.config.WsReconnectionTimeout).C:
					continue
				}

			}
			fmt.Println(fmt.Sprintf("websocket: connection established: %s", c.config.UsedBaseUrl.String()))
			c.startListener()
		}
	}
}

func (c *CatapultWebsocketClientImpl) removeHandlers() {
	c.blockSubscriber = nil
	c.confirmedAddedSubscribers = nil
	c.unconfirmedAddedSubscribers = nil
	c.unconfirmedRemovedSubscribers = nil
	c.partialAddedSubscribers = nil
	c.partialRemovedSubscribers = nil
	c.statusSubscribers = nil
	c.cosignatureSubscribers = nil
	c.driveStateSubscribers = nil

	c.topicHandlers = nil
}

func (c *CatapultWebsocketClientImpl) closeConnection(conn *websocket.Conn) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			fmt.Println(fmt.Sprintf("websocket: disconnection error: %s", err))
		}
	}
	c.conn = nil
}

func (c *CatapultWebsocketClientImpl) startListener() {
	for {
		_, resp, e := c.conn.ReadMessage()
		if e != nil {
			if _, ok := e.(*net.OpError); ok {
				// Stop ReadMessage if user called Close function for websocket client
				return
			}

			if _, ok := e.(*websocket.CloseError); ok {
				go func() {
					c.reconnectCh <- c.conn
				}()
				return
			}
		}

		c.messageRouter.RouteMessage(resp)
	}
}

func (c *CatapultWebsocketClientImpl) initNewConnection() error {
	conn, uid, err := c.connectFn(c.config)
	if err != nil {
		return err
	}

	c.UID = uid
	c.conn = conn

	messagePublisher := newMessagePublisher(c.conn)
	messageRouter := NewRouter(c.UID, messagePublisher, c.topicHandlers)

	c.messageRouter = messageRouter
	c.messagePublisher = messagePublisher

	return nil
}

func (c *CatapultWebsocketClientImpl) updateHandlers() error {

	if c.topicHandlers.HasHandler(pathBlock) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s", pathBlock))); err != nil {
			return err
		}
	}

	for _, value := range c.confirmedAddedSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.cosignatureSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, value))); err != nil {
			return err
		}
	}

	for _, value := range c.partialAddedSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.partialRemovedSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, value))); err != nil {
			return err
		}
	}

	for _, value := range c.statusSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, value))); err != nil {
			return err
		}
	}

	for _, value := range c.unconfirmedAddedSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, value))); err != nil {
			return err
		}
	}

	for _, value := range c.unconfirmedRemovedSubscribers.GetHandles() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, value))); err != nil {
			return err
		}
	}

	return nil
}

func connect(cfg *sdk.Config) (*websocket.Conn, string, error) {
	var conn *websocket.Conn
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(newWSUrl(cfg.UsedBaseUrl).String(), nil)
	if err != nil {
		for _, u := range cfg.BaseURLs {

			if u == cfg.UsedBaseUrl {
				continue
			}

			conn, _, err = websocket.DefaultDialer.Dial(newWSUrl(u).String(), nil)
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

func newWSUrl(url url.URL) *url.URL {
	if "https" == url.Scheme {
		url.Scheme = "wss"
		url.Path = "wss"
	} else {
		url.Scheme = "ws"
		url.Path = "ws"
	}

	return &url
}
