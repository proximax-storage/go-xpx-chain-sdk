// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
)

var Block *SubscribeBlock

type SubscribeService serviceWs

// const routers path for methods SubscribeService
const (
	pathBlock              = "block"
	pathConfirmedAdded     = "confirmedAdded"
	pathUnconfirmedAdded   = "unconfirmedAdded"
	pathUnconfirmedRemoved = "unconfirmedRemoved"
	pathStatus             = "status"
	pathPartialAdded       = "partialAdded"
	pathPartialRemoved     = "partialRemoved"
	pathCosignature        = "cosignature"
)

func (s *subscribe) closeChannel() error {
	switch s.Ch.(type) {
	case chan *BlockInfo:
		chType := s.Ch.(chan *BlockInfo)
		close(chType)

	case chan *StatusInfo:
		chType := s.Ch.(chan *StatusInfo)
		delete(statusInfoChannels, s.getAdd())
		close(chType)

	case chan *HashInfo:
		chType := s.Ch.(chan *HashInfo)
		delete(unconfirmedRemovedChannels, s.getAdd())
		close(chType)

	case chan *PartialRemovedInfo:
		chType := s.Ch.(chan *PartialRemovedInfo)
		delete(partialRemovedInfoChannels, s.getAdd())
		close(chType)

	case chan *SignerInfo:
		chType := s.Ch.(chan *SignerInfo)
		delete(signerInfoChannels, s.getAdd())
		close(chType)

	case chan *ErrorInfo:
		chType := s.Ch.(chan *ErrorInfo)
		delete(errChannels, s.getAdd())
		close(chType)

	case chan Transaction:
		chType := s.Ch.(chan Transaction)
		if s.getSubscribe() == "partialAdded" {
			delete(partialAddedChannels, s.getAdd())
		} else if s.getSubscribe() == "unconfirmedAdded" {
			delete(unconfirmedAddedChannels, s.getAdd())
		} else {
			delete(confirmedAddedChannels, s.getAdd())
		}
		close(chType)

	default:
		return errors.New("WRONG TYPE CHANNEL")
	}
	return nil
}

// Unsubscribe terminates the specified subscription.
// It does not have any specific param.
func (c *subscribe) unsubscribe() error {
	c.conn = connectsWs[c.getAdd()]
	if err := websocket.JSON.Send(c.conn, sendJson{
		Uid:       c.Uid,
		Subscribe: c.Subscribe,
	}); err != nil {
		return err
	}

	if err := c.closeChannel(); err != nil {
		return err
	}

	return nil
}

// Generate a new channel and subscribe to the websocket.
// param route A subscription channel route.
// return A pointer Subscribe struct
func (c *SubscribeService) newSubscribe(route string) (*subscribe, error) {
	subMsg := c.client.buildSubscribe(route)

	err := c.client.subsChannel(subMsg)
	if err != nil {
		return nil, err
	}
	return subMsg, nil
}

func (c *SubscribeService) getClient(add *Address) *ClientWebsocket {
	if len(connectsWs) == 0 {
		connectsWs[add.Address] = c.client.client
		return c.client
	} else if _, exist := connectsWs[add.Address]; exist {
		return c.client
	} else {
		client, err := NewConnectWs(c.client.config.BaseURL.String(), *c.client.duration)
		if err != nil {
			fmt.Println(err)
		}
		connectsWs[add.Address] = client.client
		return client
	}
}

// returns entity from which you access channel with block infos
// block info gets into channel when new block is harvested
func (c *SubscribeService) Block() (*SubscribeBlock, error) {
	subBlock := new(SubscribeBlock)
	Block = subBlock
	subBlock.Ch = make(chan *BlockInfo)
	subscribe, err := c.newSubscribe(pathBlock)
	subBlock.subscribe = subscribe
	subscribe.Ch = subBlock.Ch
	return subBlock, err
}

// returns an entity from which you can access channel with Transaction infos for passed address
// Transaction info gets into channel when it is included in a block
func (c *SubscribeService) ConfirmedAdded(add *Address) (*SubscribeTransaction, error) {
	c.client = c.getClient(add)
	subTransaction := new(SubscribeTransaction)
	subTransaction.Ch = make(chan Transaction)
	confirmedAddedChannels[add.Address] = subTransaction.Ch
	subscribe, err := c.newSubscribe(pathConfirmedAdded + "/" + add.Address)
	subTransaction.subscribe = subscribe
	subscribe.Ch = subTransaction.Ch
	return subTransaction, err
}

// returns an entity from which you can access channel with Transaction infos for passed address
// Transaction info gets into channel when it is in unconfirmed state and waiting to be included into a block
func (c *SubscribeService) UnconfirmedAdded(add *Address) (*SubscribeTransaction, error) {
	c.client = c.getClient(add)
	subTransaction := new(SubscribeTransaction)
	subTransaction.Ch = make(chan Transaction)
	unconfirmedAddedChannels[add.Address] = subTransaction.Ch
	subscribe, err := c.newSubscribe(pathUnconfirmedAdded + "/" + add.Address)
	subTransaction.subscribe = subscribe
	subscribe.Ch = unconfirmedAddedChannels[add.Address]
	return subTransaction, err
}

// returns an entity from which you can access channel with Transaction infos for passed address
// Transaction info gets into channel when it was in unconfirmed state but not anymore
func (c *SubscribeService) UnconfirmedRemoved(add *Address) (*SubscribeHash, error) {
	c.client = c.getClient(add)
	subHash := new(SubscribeHash)
	subHash.Ch = make(chan *HashInfo)
	unconfirmedRemovedChannels[add.Address] = subHash.Ch
	subscribe, err := c.newSubscribe(pathUnconfirmedRemoved + "/" + add.Address)
	subHash.subscribe = subscribe
	subscribe.Ch = unconfirmedRemovedChannels[add.Address]
	return subHash, err
}

// returns an entity from which you can access channel with Transaction status infos for passed address
// Transaction info gets into channel when it rises an error
func (c *SubscribeService) Status(add *Address) (*SubscribeStatus, error) {
	c.client = c.getClient(add)
	subStatus := new(SubscribeStatus)
	subStatus.Ch = make(chan *StatusInfo)
	statusInfoChannels[add.Address] = subStatus.Ch
	subscribe, err := c.newSubscribe(pathStatus + "/" + add.Address)
	subStatus.subscribe = subscribe
	subscribe.Ch = statusInfoChannels[add.Address]
	return subStatus, err
}

// returns an entity from which you can access channel with Aggregate Bonded Transaction info
// Aggregate Bonded Transaction info gets into channel when it is in partial state
// and waiting for actors to send all required cosignature transactions
func (c *SubscribeService) PartialAdded(add *Address) (*SubscribeBonded, error) {
	c.client = c.getClient(add)
	subTransaction := new(SubscribeBonded)
	subTransaction.Ch = make(chan *AggregateTransaction)
	partialAddedChannels[add.Address] = subTransaction.Ch
	subscribe, err := c.newSubscribe(pathPartialAdded + "/" + add.Address)
	subTransaction.subscribe = subscribe
	subscribe.Ch = partialAddedChannels[add.Address]
	return subTransaction, err
}

// returns an entity from which you can access channel with Aggregate Bonded Transaction hash related to passed address
// Aggregate Bonded Transaction hash gets into channel when it was in partial state but not anymore
func (c *SubscribeService) PartialRemoved(add *Address) (*SubscribePartialRemoved, error) {
	c.client = c.getClient(add)
	subPartialRemoved := new(SubscribePartialRemoved)
	subPartialRemoved.Ch = make(chan *PartialRemovedInfo)
	partialRemovedInfoChannels[add.Address] = subPartialRemoved.Ch
	subscribe, err := c.newSubscribe(pathPartialRemoved + "/" + add.Address)
	subPartialRemoved.subscribe = subscribe
	subscribe.Ch = partialRemovedInfoChannels[add.Address]
	return subPartialRemoved, err
}

// returns an entity from which you can access channel with cosignature transaction is added to an
// aggregate bounded transaction with partial state related to passed address
func (c *SubscribeService) Cosignature(add *Address) (*SubscribeSigner, error) {
	c.client = c.getClient(add)
	subCosignature := new(SubscribeSigner)
	subCosignature.Ch = make(chan *SignerInfo)
	signerInfoChannels[add.Address] = subCosignature.Ch
	subscribe, err := c.newSubscribe(pathCosignature + "/" + add.Address)
	subCosignature.subscribe = subscribe
	subscribe.Ch = signerInfoChannels[add.Address]
	return subCosignature, err
}

// returns an entity from which you can access channel with errors related to passed address
func (c *SubscribeService) Error(add *Address) *SubscribeError {
	c.client = c.getClient(add)
	subError := new(SubscribeError)
	subError.Ch = make(chan *ErrorInfo)
	errChannels[add.Address] = subError.Ch
	subscribe := new(subscribe)
	subscribe.Subscribe = "error/" + add.Address
	subError.subscribe = subscribe
	subscribe.Ch = errChannels[add.Address]
	return subError
}
