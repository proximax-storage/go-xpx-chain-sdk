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
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subs"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
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
	pathDriveState         Path = "driveState"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

type (
	// Subscribe path
	Path string

	CatapultWebsocketClientImpl struct {
		UID    string
		config *sdk.Config

		conn *websocket.Conn

		blockSubPools      subs.SubscribersPool[*sdk.BlockInfo]
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
	}

	Client interface {
		io.Closer

		Listen(ctx context.Context)
	}

	CatapultClient interface {
		Client

		Config() *sdk.Config

		NewBlockSubscription() (sub <-chan *sdk.BlockInfo, subId int)
		BlockUnsubscribe(subId int)

		NewConfirmedAddedSubscription(address sdk.Address) (sub <-chan sdk.Transaction, subId int)
		ConfirmedAddedUnsubscribe(address sdk.Address, subId int)

		NewUnConfirmedAddedSubscription(address sdk.Address) (sub <-chan sdk.Transaction, subId int)
		UnConfirmedAddedUnsubscribe(address sdk.Address, subId int)

		NewUnConfirmedRemovedSubscription(address sdk.Address) (sub <-chan *sdk.UnconfirmedRemoved, subId int)
		UnConfirmedRemovedUnsubscribe(address sdk.Address, subId int)

		NewCosignatureSubscription(address sdk.Address) (sub <-chan *sdk.SignerInfo, subId int)
		CosignatureUnsubscribe(address sdk.Address, subId int)

		NewPartialAddedSubscription(address sdk.Address) (sub <-chan *sdk.AggregateTransaction, subId int)
		PartialAddedUnsubscribe(address sdk.Address, subId int)

		NewPartialRemovedSubscription(address sdk.Address) (sub <-chan *sdk.PartialRemovedInfo, subId int)
		PartialRemovedUnsubscribe(address sdk.Address, subId int)

		NewStatusSubscription(address sdk.Address) (sub <-chan *sdk.StatusInfo, subId int)
		DriveStateUnsubscribe(address sdk.Address, subId int)

		NewDriveStateSubscription(address sdk.Address) (sub <-chan *sdk.DriveStateInfo, subId int)
		StatusUnsubscribe(address sdk.Address, subId int)
	}
)

func NewClient(cfg *sdk.Config) (CatapultClient, error) {
	publisher := subs.NewPublisher()

	blockSubPools := subs.NewSubscribersPool[*sdk.BlockInfo](sdk.NewMapper[*sdk.BlockInfo](cfg.GenerationHash, sdk.BlockMapperFunc))
	err := publisher.AddSubscriber(subs.TopicBlock, blockSubPools)
	if err != nil {
		return nil, err
	}

	cosignatureSubs := subs.NewSubscribersPool[*sdk.SignerInfo](sdk.NewMapper[*sdk.SignerInfo](cfg.GenerationHash, sdk.CosignatureMapperFunc))
	err = publisher.AddSubscriber(subs.TopicCosignature, cosignatureSubs)
	if err != nil {
		return nil, err
	}

	driveStateSubs := subs.NewSubscribersPool[*sdk.DriveStateInfo](sdk.NewMapper[*sdk.DriveStateInfo](cfg.GenerationHash, sdk.DriveStateMapperFunc))
	err = publisher.AddSubscriber(subs.TopicDriveState, driveStateSubs)
	if err != nil {
		return nil, err
	}

	confAddedSubs := subs.NewSubscribersPool[sdk.Transaction](sdk.NewMapper[sdk.Transaction](cfg.GenerationHash, sdk.TransactionMapperFunc))
	err = publisher.AddSubscriber(subs.TopicConfirmedAdded, confAddedSubs)
	if err != nil {
		return nil, err
	}

	partialAddedSubs := subs.NewSubscribersPool[*sdk.AggregateTransaction](sdk.NewMapper[*sdk.AggregateTransaction](cfg.GenerationHash, sdk.AggregateTransactionMapperFunc))
	err = publisher.AddSubscriber(subs.TopicPartialAdded, partialAddedSubs)
	if err != nil {
		return nil, err
	}

	partialRemovedSubs := subs.NewSubscribersPool[*sdk.PartialRemovedInfo](sdk.NewMapper[*sdk.PartialRemovedInfo](cfg.GenerationHash, sdk.PartialRemovedMapperFunc))
	err = publisher.AddSubscriber(subs.TopicPartialRemoved, partialRemovedSubs)
	if err != nil {
		return nil, err
	}

	statusSubs := subs.NewSubscribersPool[*sdk.StatusInfo](sdk.NewMapper[*sdk.StatusInfo](cfg.GenerationHash, sdk.StatusMapperFunc))
	err = publisher.AddSubscriber(subs.TopicStatus, statusSubs)
	if err != nil {
		return nil, err
	}

	unconfAddedSubs := subs.NewSubscribersPool[sdk.Transaction](sdk.NewMapper[sdk.Transaction](cfg.GenerationHash, sdk.TransactionMapperFunc))
	err = publisher.AddSubscriber(subs.TopicUnconfirmedAdded, unconfAddedSubs)
	if err != nil {
		return nil, err
	}

	unconfRemovedSubs := subs.NewSubscribersPool[*sdk.UnconfirmedRemoved](sdk.NewMapper[*sdk.UnconfirmedRemoved](cfg.GenerationHash, sdk.UnconfirmedRemovedMapperFunc))
	err = publisher.AddSubscriber(subs.TopicUnconfirmedRemoved, unconfRemovedSubs)
	if err != nil {
		return nil, err
	}

	socketClient := &CatapultWebsocketClientImpl{
		config:             cfg,
		blockSubPools:      blockSubPools,
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
	go c.startMessageReading(ctx)
}

func (c *CatapultWebsocketClientImpl) Close() error {
	return c.closeConnection()
}

func (c *CatapultWebsocketClientImpl) Config() *sdk.Config {
	return c.config
}

func (c *CatapultWebsocketClientImpl) NewBlockSubscription() (sub <-chan *sdk.BlockInfo, subId int) {
	if !c.blockSubPools.HasSubscriptions("") {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s", pathBlock)))
		if err != nil {
			fmt.Printf("cannot subscribe on block topic: %s\n", err)
		}
	}

	return c.blockSubPools.NewSubscription("")
}

func (c *CatapultWebsocketClientImpl) NewConfirmedAddedSubscription(address sdk.Address) (sub <-chan sdk.Transaction, subId int) {
	if !c.confAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from confirmed added topic for %s: %s\n", address.Address, err)
		}
	}

	return c.confAddedSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewUnConfirmedAddedSubscription(address sdk.Address) (sub <-chan sdk.Transaction, subId int) {
	if !c.unconfAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from unconfirmed added topic for %s: %s\n", address.Address, err)
		}
	}

	return c.unconfAddedSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewUnConfirmedRemovedSubscription(address sdk.Address) (sub <-chan *sdk.UnconfirmedRemoved, subId int) {
	if !c.unconfAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from unconfirmed removed topic for %s: %s\n", address.Address, err)
		}
	}

	return c.unconfRemovedSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewCosignatureSubscription(address sdk.Address) (sub <-chan *sdk.SignerInfo, subId int) {
	if !c.cosignatureSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from cosignature topic for %s: %s\n", address.Address, err)
		}
	}

	return c.cosignatureSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewPartialAddedSubscription(address sdk.Address) (sub <-chan *sdk.AggregateTransaction, subId int) {
	if !c.partialAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from partial added topic for %s: %s\n", address.Address, err)
		}
	}

	return c.partialAddedSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewPartialRemovedSubscription(address sdk.Address) (sub <-chan *sdk.PartialRemovedInfo, subId int) {
	if !c.partialRemovedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from patrial removed topic for %s: %s\n", address.Address, err)
		}
	}

	return c.partialRemovedSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewStatusSubscription(address sdk.Address) (sub <-chan *sdk.StatusInfo, subId int) {
	c.statusSubs.CloseSubscription(address.Address, subId)
	if !c.statusSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from drive state topic for %s: %s\n", address.Address, err)
		}
	}

	return c.statusSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) NewDriveStateSubscription(address sdk.Address) (sub <-chan *sdk.DriveStateInfo, subId int) {
	if !c.statusSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathDriveState, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from status topic for %s: %s\n", address.Address, err)
		}
	}

	return c.driveStateSubs.NewSubscription(address.Address)
}

func (c *CatapultWebsocketClientImpl) BlockUnsubscribe(subId int) {
	c.blockSubPools.CloseSubscription("", subId)
	if !c.blockSubPools.HasSubscriptions("") {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s", pathBlock)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from block topic: %s\n", err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) ConfirmedAddedUnsubscribe(address sdk.Address, subId int) {
	c.confAddedSubs.CloseSubscription(address.Address, subId)
	if !c.confAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from confirmed added topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) UnConfirmedAddedUnsubscribe(address sdk.Address, subId int) {
	c.unconfAddedSubs.CloseSubscription(address.Address, subId)
	if !c.unconfAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from unconfirmed added topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) UnConfirmedRemovedUnsubscribe(address sdk.Address, subId int) {
	c.unconfRemovedSubs.CloseSubscription(address.Address, subId)
	if !c.unconfAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from unconfirmed removed topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) CosignatureUnsubscribe(address sdk.Address, subId int) {
	c.cosignatureSubs.CloseSubscription(address.Address, subId)
	if !c.cosignatureSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from cosignature topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) PartialAddedUnsubscribe(address sdk.Address, subId int) {
	c.partialAddedSubs.CloseSubscription(address.Address, subId)
	if !c.partialAddedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from partial added topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) PartialRemovedUnsubscribe(address sdk.Address, subId int) {
	c.partialRemovedSubs.CloseSubscription(address.Address, subId)
	if !c.partialRemovedSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from patrial removed topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) DriveStateUnsubscribe(address sdk.Address, subId int) {
	c.driveStateSubs.CloseSubscription(address.Address, subId)
	if !c.driveStateSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathDriveState, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from drive state topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) StatusUnsubscribe(address sdk.Address, subId int) {
	c.statusSubs.CloseSubscription(address.Address, subId)
	if !c.statusSubs.HasSubscriptions(address.Address) {
		err := c.messagePublisher.PublishUnsubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, address.Address)))
		if err != nil {
			fmt.Printf("cannot unsubscribe from status topic for %s: %s\n", address.Address, err)
		}
	}
}

func (c *CatapultWebsocketClientImpl) closeConnection() error {
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
				}
				return
			}

			err = c.publisher.Publish(resp)
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
	c.conn.Close()

	for {
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
		go c.startMessageReading(ctx)
		return
	}
}

func (c *CatapultWebsocketClientImpl) updateHandlers() error {
	if c.blockSubPools.HasSubscriptions("") {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s", pathBlock))); err != nil {
			return err
		}
	}

	for _, addr := range c.confAddedSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.unconfAddedSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.unconfRemovedSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.cosignatureSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathCosignature, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.partialAddedSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialAdded, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.partialRemovedSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.statusSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathStatus, addr)))
		if err != nil {
			return err
		}
	}

	for _, addr := range c.driveStateSubs.GetAddresses() {
		err := c.messagePublisher.PublishSubscribeMessage(c.UID, Path(fmt.Sprintf("%s/%s", pathDriveState, addr)))
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
