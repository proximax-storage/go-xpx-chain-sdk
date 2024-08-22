// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subs"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	topicBlock              subs.Topic = "block"
	topicConfirmedAdded     subs.Topic = "confirmedAdded"
	topicUnconfirmedAdded   subs.Topic = "unconfirmedAdded"
	topicUnconfirmedRemoved subs.Topic = "unconfirmedRemoved"
	topicStatus             subs.Topic = "status"
	topicPartialAdded       subs.Topic = "partialAdded"
	topicPartialRemoved     subs.Topic = "partialRemoved"
	topicCosignature        subs.Topic = "cosignature"
	topicDriveState         subs.Topic = "driveState"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

type (
	// Subscribe path
	//Path string

	CatapultWebsocketClientImpl struct {
		UID    string
		config *sdk.Config

		conn *websocket.Conn

		blockSubs          subs.SubscribersPool[*sdk.BlockInfo]
		cosignatureSubs    subs.SubscribersPool[*sdk.SignerInfo]
		driveStateSubs     subs.SubscribersPool[*sdk.DriveStateInfo]
		confAddedSubs      subs.SubscribersPool[sdk.Transaction]
		partialAddedSubs   subs.SubscribersPool[*sdk.AggregateTransaction]
		partialRemovedSubs subs.SubscribersPool[*sdk.PartialRemovedInfo]
		statusSubs         subs.SubscribersPool[*sdk.StatusInfo]
		unconfAddedSubs    subs.SubscribersPool[sdk.Transaction]
		unconfRemovedSubs  subs.SubscribersPool[*sdk.UnconfirmedRemoved]

		publisher        *subs.Publisher
		messagePublisher MessagePublisher

		listening atomic.Bool
	}

	Client interface {
		io.Closer

		Listen(ctx context.Context)
	}

	CatapultClient interface {
		Client

		Config() *sdk.Config

		NewBlockSubscription() (sub <-chan *sdk.BlockInfo, subId int, err error)
		BlockUnsubscribe(subId int) error

		NewConfirmedAddedSubscription(address *sdk.Address) (sub <-chan sdk.Transaction, subId int, err error)
		ConfirmedAddedUnsubscribe(address *sdk.Address, subId int) error

		NewUnConfirmedAddedSubscription(address *sdk.Address) (sub <-chan sdk.Transaction, subId int, err error)
		UnConfirmedAddedUnsubscribe(address *sdk.Address, subId int) error

		NewUnConfirmedRemovedSubscription(address *sdk.Address) (sub <-chan *sdk.UnconfirmedRemoved, subId int, err error)
		UnConfirmedRemovedUnsubscribe(address *sdk.Address, subId int) error

		NewCosignatureSubscription(address *sdk.Address) (sub <-chan *sdk.SignerInfo, subId int, err error)
		CosignatureUnsubscribe(address *sdk.Address, subId int) error

		NewPartialAddedSubscription(address *sdk.Address) (sub <-chan *sdk.AggregateTransaction, subId int, err error)
		PartialAddedUnsubscribe(address *sdk.Address, subId int) error

		NewPartialRemovedSubscription(address *sdk.Address) (sub <-chan *sdk.PartialRemovedInfo, subId int, err error)
		PartialRemovedUnsubscribe(address *sdk.Address, subId int) error

		NewStatusSubscription(address *sdk.Address) (sub <-chan *sdk.StatusInfo, subId int, err error)
		DriveStateUnsubscribe(address *sdk.Address, subId int) error

		NewDriveStateSubscription(address *sdk.Address) (sub <-chan *sdk.DriveStateInfo, subId int, err error)
		StatusUnsubscribe(address *sdk.Address, subId int) error
	}
)

func NewClient(cfg *sdk.Config) (CatapultClient, error) {
	newSubs := make(map[subs.Topic]subs.Notifier)

	blockSubPools := subs.NewSubscribersPool[*sdk.BlockInfo](sdk.NewMapper[*sdk.BlockInfo](cfg.GenerationHash, sdk.BlockMapperFunc))
	newSubs[subs.TopicBlock] = blockSubPools

	cosignatureSubs := subs.NewSubscribersPool[*sdk.SignerInfo](sdk.NewMapper[*sdk.SignerInfo](cfg.GenerationHash, sdk.CosignatureMapperFunc))
	newSubs[subs.TopicCosignature] = cosignatureSubs

	driveStateSubs := subs.NewSubscribersPool[*sdk.DriveStateInfo](sdk.NewMapper[*sdk.DriveStateInfo](cfg.GenerationHash, sdk.DriveStateMapperFunc))
	newSubs[subs.TopicDriveState] = driveStateSubs

	confAddedSubs := subs.NewSubscribersPool[sdk.Transaction](sdk.NewMapper[sdk.Transaction](cfg.GenerationHash, sdk.TransactionMapperFunc))
	newSubs[subs.TopicConfirmedAdded] = confAddedSubs

	partialAddedSubs := subs.NewSubscribersPool[*sdk.AggregateTransaction](sdk.NewMapper[*sdk.AggregateTransaction](cfg.GenerationHash, sdk.AggregateTransactionMapperFunc))
	newSubs[subs.TopicPartialAdded] = partialAddedSubs

	partialRemovedSubs := subs.NewSubscribersPool[*sdk.PartialRemovedInfo](sdk.NewMapper[*sdk.PartialRemovedInfo](cfg.GenerationHash, sdk.PartialRemovedMapperFunc))
	newSubs[subs.TopicPartialRemoved] = partialRemovedSubs

	statusSubs := subs.NewSubscribersPool[*sdk.StatusInfo](sdk.NewMapper[*sdk.StatusInfo](cfg.GenerationHash, sdk.StatusMapperFunc))
	newSubs[subs.TopicStatus] = statusSubs

	unconfAddedSubs := subs.NewSubscribersPool[sdk.Transaction](sdk.NewMapper[sdk.Transaction](cfg.GenerationHash, sdk.TransactionMapperFunc))
	newSubs[subs.TopicUnconfirmedAdded] = unconfAddedSubs

	unconfRemovedSubs := subs.NewSubscribersPool[*sdk.UnconfirmedRemoved](sdk.NewMapper[*sdk.UnconfirmedRemoved](cfg.GenerationHash, sdk.UnconfirmedRemovedMapperFunc))
	newSubs[subs.TopicUnconfirmedRemoved] = unconfRemovedSubs

	var err error
	publisher := subs.NewPublisher()
	for topic, notifier := range newSubs {
		err = publisher.AddSubscriber(topic, notifier)
		if err != nil {
			return nil, err
		}
	}

	socketClient := &CatapultWebsocketClientImpl{
		config:             cfg,
		blockSubs:          blockSubPools,
		cosignatureSubs:    cosignatureSubs,
		driveStateSubs:     driveStateSubs,
		confAddedSubs:      confAddedSubs,
		partialAddedSubs:   partialAddedSubs,
		partialRemovedSubs: partialRemovedSubs,
		statusSubs:         statusSubs,
		unconfAddedSubs:    unconfAddedSubs,
		unconfRemovedSubs:  unconfRemovedSubs,
		publisher:          publisher,
	}

	if err := socketClient.initNewConnection(); err != nil {
		return nil, err
	}

	return socketClient, nil
}

func (c *CatapultWebsocketClientImpl) Listen(ctx context.Context) {
	if c.listening.Swap(true) {
		return
	}

	c.startMessageReading(ctx)
}

func (c *CatapultWebsocketClientImpl) Close() error {
	return c.closeConnection()
}

func (c *CatapultWebsocketClientImpl) Config() *sdk.Config {
	return c.config
}

func (c *CatapultWebsocketClientImpl) NewBlockSubscription() (sub <-chan *sdk.BlockInfo, subId int, err error) {
	return subscribe[*sdk.BlockInfo](topicBlock, nil, c.UID, c.messagePublisher, c.blockSubs)
}

func (c *CatapultWebsocketClientImpl) NewConfirmedAddedSubscription(address *sdk.Address) (sub <-chan sdk.Transaction, subId int, err error) {
	return subscribe[sdk.Transaction](topicConfirmedAdded, address, c.UID, c.messagePublisher, c.confAddedSubs)
}

func (c *CatapultWebsocketClientImpl) NewUnConfirmedAddedSubscription(address *sdk.Address) (sub <-chan sdk.Transaction, subId int, err error) {
	return subscribe[sdk.Transaction](topicUnconfirmedAdded, address, c.UID, c.messagePublisher, c.unconfAddedSubs)
}

func (c *CatapultWebsocketClientImpl) NewUnConfirmedRemovedSubscription(address *sdk.Address) (sub <-chan *sdk.UnconfirmedRemoved, subId int, err error) {
	return subscribe[*sdk.UnconfirmedRemoved](topicUnconfirmedRemoved, address, c.UID, c.messagePublisher, c.unconfRemovedSubs)
}

func (c *CatapultWebsocketClientImpl) NewCosignatureSubscription(address *sdk.Address) (sub <-chan *sdk.SignerInfo, subId int, err error) {
	return subscribe[*sdk.SignerInfo](topicCosignature, address, c.UID, c.messagePublisher, c.cosignatureSubs)
}

func (c *CatapultWebsocketClientImpl) NewPartialAddedSubscription(address *sdk.Address) (sub <-chan *sdk.AggregateTransaction, subId int, err error) {
	return subscribe[*sdk.AggregateTransaction](topicPartialAdded, address, c.UID, c.messagePublisher, c.partialAddedSubs)
}

func (c *CatapultWebsocketClientImpl) NewPartialRemovedSubscription(address *sdk.Address) (sub <-chan *sdk.PartialRemovedInfo, subId int, err error) {
	return subscribe[*sdk.PartialRemovedInfo](topicPartialRemoved, address, c.UID, c.messagePublisher, c.partialRemovedSubs)
}

func (c *CatapultWebsocketClientImpl) NewStatusSubscription(address *sdk.Address) (sub <-chan *sdk.StatusInfo, subId int, err error) {
	return subscribe[*sdk.StatusInfo](topicStatus, address, c.UID, c.messagePublisher, c.statusSubs)
}

func (c *CatapultWebsocketClientImpl) NewDriveStateSubscription(address *sdk.Address) (sub <-chan *sdk.DriveStateInfo, subId int, err error) {
	return subscribe[*sdk.DriveStateInfo](topicDriveState, address, c.UID, c.messagePublisher, c.driveStateSubs)
}

func (c *CatapultWebsocketClientImpl) BlockUnsubscribe(subId int) error {
	return unsubscribe[*sdk.BlockInfo](topicBlock, nil, c.UID, subId, c.blockSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) ConfirmedAddedUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[sdk.Transaction](topicConfirmedAdded, address, c.UID, subId, c.confAddedSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) UnConfirmedAddedUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[sdk.Transaction](topicUnconfirmedAdded, address, c.UID, subId, c.unconfAddedSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) UnConfirmedRemovedUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.UnconfirmedRemoved](topicUnconfirmedRemoved, address, c.UID, subId, c.unconfRemovedSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) CosignatureUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.SignerInfo](topicCosignature, address, c.UID, subId, c.cosignatureSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) PartialAddedUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.AggregateTransaction](topicPartialAdded, address, c.UID, subId, c.partialAddedSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) PartialRemovedUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.PartialRemovedInfo](topicPartialRemoved, address, c.UID, subId, c.partialRemovedSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) DriveStateUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.DriveStateInfo](topicDriveState, address, c.UID, subId, c.driveStateSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) StatusUnsubscribe(address *sdk.Address, subId int) error {
	return unsubscribe[*sdk.StatusInfo](topicStatus, address, c.UID, subId, c.statusSubs, c.messagePublisher)
}

func (c *CatapultWebsocketClientImpl) closeConnection() error {
	log.Println("closing connection...")
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			fmt.Println(fmt.Sprintf("websocket: disconnection error: %s", err))
			return err
		}
	}
	c.conn = nil

	return nil
}

func (c *CatapultWebsocketClientImpl) startMessageReading(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, resp, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Unexpected close error, attempting to reconnect...")
					c.reconnect(ctx)
					continue
				}
				return
			}

			err = c.publisher.Publish(ctx, resp)
			if err != nil {
				log.Printf("Cannot publish ws message:%s\n", err)
			}
		}
	}
}

func (c *CatapultWebsocketClientImpl) initNewConnection() error {
	var conn *websocket.Conn
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(newWSUrl(c.config.UsedBaseUrl).String(), nil)
	if err != nil {
		for _, u := range c.config.BaseURLs {

			if u == c.config.UsedBaseUrl {
				continue
			}

			conn, _, err = websocket.DefaultDialer.Dial(newWSUrl(u).String(), nil)
			if err != nil {
				continue
			}

			c.config.UsedBaseUrl = u
			break
		}
	}

	if conn == nil {
		return err
	}

	resp := new(wsConnectionResponse)
	if err = conn.ReadJSON(resp); err != nil {
		return err
	}

	c.UID = resp.Uid
	c.conn = conn

	c.messagePublisher = newMessagePublisher(c.conn)
	return nil
}

func (c *CatapultWebsocketClientImpl) reconnect(ctx context.Context) {
	c.closeConnection()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("websocket: connection is failed. Try again after wait period")
			err := c.initNewConnection()
			if err != nil {
				fmt.Println("websocket: connection is failed. Try again after wait period")
				select {
				case <-time.NewTicker(c.config.WsReconnectionTimeout).C:
					continue
				}
			}

			err = c.updateHandlers()
			if err != nil {
				fmt.Println("websocket: update handles is failed. Try again after timeout period")
				select {
				case <-time.NewTicker(c.config.WsReconnectionTimeout).C:
					continue
				}

			}

			fmt.Println(fmt.Sprintf("websocket: connection established: %s", c.config.UsedBaseUrl.String()))
			return
		}
	}
}

func (c *CatapultWebsocketClientImpl) updateHandlers() error {
	for _, path := range c.blockSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.confAddedSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.unconfAddedSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.unconfRemovedSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.cosignatureSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.partialAddedSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.partialRemovedSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.statusSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	for _, path := range c.driveStateSubs.GetPaths() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, path)
		if err != nil {
			return err
		}
	}

	return nil
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

func subscribe[T any](
	topic subs.Topic,
	address *sdk.Address,
	uid string,
	publisher MessagePublisher,
	subsPool subs.SubscribersPool[T]) (_ <-chan T, id int, err error) {

	path := subs.NewPath(topic, address)
	if !subsPool.HasSubscriptions(path) {
		err := publisher.PublishSubscribeMessage(uid, path.String())
		if err != nil {
			fmt.Printf("cannot subscribe on block topic: %s\n", err)
			return nil, 0, err
		}
	}

	sub, id := subsPool.NewSubscription(path)
	return sub, id, nil
}

func unsubscribe[T any](
	topic subs.Topic,
	address *sdk.Address,
	uid string,
	subId int,
	subsPool subs.SubscribersPool[T],
	publisher MessagePublisher) error {

	path := subs.NewPath(topic, address)
	subsPool.CloseSubscription(path, subId)
	if !subsPool.HasSubscriptions(path) {
		err := publisher.PublishUnsubscribeMessage(uid, path.String())
		if err != nil {
			fmt.Printf("cannot unsubscribe from confirmed added topic for %s: %s\n", path, err)
			return err
		}
	}

	return nil
}
