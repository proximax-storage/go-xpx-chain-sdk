// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io"
	"sync"
)

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

	var resp *wsConnectionResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	return &CatapultWebsocketClientImpl{
		conn:             conn,
		UID:              resp.Uid,
		eventSubscribers: newEventSubscribers(),
		messageProcessor: newMessageProcessor(MapTransaction),
		messagePublisher: newMessagePublisher(conn),
		errorsChan:       make(chan error, 100),
	}, nil
}

type CatapultWebsocketClient interface {
	Listen(wg *sync.WaitGroup)

	AddBlockHandlers(handlers ...blockHandler) error
	AddConfirmedAddedHandlers(address *Address, handlers ...confirmedAddedHandler) error
	AddUnconfirmedAddedHandlers(address *Address, handlers ...unconfirmedAddedHandler) error
	AddUnconfirmedRemovedHandlers(address *Address, handlers ...unconfirmedRemovedHandler) error
	AddStatusHandlers(address *Address, handlers ...statusHandler) error
	AddPartialAddedHandlers(address *Address, handlers ...partialAddedHandler) error
	AddPartialRemovedHandlers(address *Address, handlers ...partialRemovedHandler) error
	AddCosignatureHandlers(address *Address, handlers ...cosignatureHandler) error

	GetErrorsChan() (chan error, error)
}

type CatapultWebsocketClientImpl struct {
	conn             *websocket.Conn
	UID              string
	eventSubscribers eventsSubscribers
	messageProcessor messageProcessor
	messagePublisher messagePublisher

	errorsChan chan error
}

func (c *CatapultWebsocketClientImpl) Listen(wg *sync.WaitGroup) {
	wg.Add(1)

	var resp []byte
	for {
		err := websocket.Message.Receive(c.conn, &resp)

		if err == io.EOF {
			wg.Done()
			return
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		messageInfo, err := c.getMessageInfo(resp)
		if err != nil {
			fmt.Println(fmt.Errorf("error getting websocket message info: %s", err))
			continue
		}

		if err := c.routeMessage(messageInfo, resp); err != nil {
			fmt.Println(fmt.Errorf("error routing message websocket message info: %s", err))
		}
	}
}

func (c *CatapultWebsocketClientImpl) AddBlockHandlers(handlers ...blockHandler) error {
	if !c.eventSubscribers.IsBlockSubscribed() {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, pathBlock); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddBlockHandlers(handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddConfirmedAddedHandlers(address *Address, handlers ...confirmedAddedHandler) error {
	if !c.eventSubscribers.IsAddConfirmedAddedSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathConfirmedAdded, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddConfirmedAddedHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedAddedHandlers(address *Address, handlers ...unconfirmedAddedHandler) error {
	if !c.eventSubscribers.IsAddUnconfirmedAddedSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddUnconfirmedAddedHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddUnconfirmedRemovedHandlers(address *Address, handlers ...unconfirmedRemovedHandler) error {
	if !c.eventSubscribers.IsAddUnconfirmedRemovedSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddUnconfirmedRemovedHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddStatusHandlers(address *Address, handlers ...statusHandler) error {
	if !c.eventSubscribers.IsAddStatusInfoSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathStatus, address.Address)); err != nil {
			return err
		}

	}

	err := c.eventSubscribers.AddStatusInfoHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialAddedHandlers(address *Address, handlers ...partialAddedHandler) error {
	if !c.eventSubscribers.IsAddPartialAddedSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathPartialAdded, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddPartialAddedHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddPartialRemovedHandlers(address *Address, handlers ...partialRemovedHandler) error {
	if !c.eventSubscribers.IsAddPartialRemovedSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathPartialRemoved, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddPartialRemovedHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) AddCosignatureHandlers(address *Address, handlers ...cosignatureHandler) error {
	if c.eventSubscribers.IsAddCosignatureSubscribed(address) {
		if err := c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathCosignature, address.Address)); err != nil {
			return err
		}
	}

	err := c.eventSubscribers.AddCosignatureHandlers(address, handlers...)
	if err != nil {
		return err
	}

	return nil
}

func (c *CatapultWebsocketClientImpl) GetErrorsChan() (chan error, error) {
	if c.errorsChan == nil {
		return nil, errors.New("channel is not initialized")
	}

	return c.errorsChan, nil
}

func (c *CatapultWebsocketClientImpl) getMessageInfo(m []byte) (*wsMessageInfo, error) {
	var messageInfoDTO wsMessageInfoDTO
	if err := json.Unmarshal(m, &messageInfoDTO); err != nil {
		return nil, err
	}

	return messageInfoDTO.toStruct()
}

func (c *CatapultWebsocketClientImpl) routeMessage(messageInfo *wsMessageInfo, resp []byte /* err error*/) error {
	switch messageInfo.ChannelName {
	case pathBlock:
		res, err := c.messageProcessor.ProcessBlock(resp)
		if err != nil {
			return err
		}

		handlers := c.eventSubscribers.GetBlockHandlers()

		for f := range handlers {
			go func(callFuncPtr *blockHandler, err error) {
				callFunc := *callFuncPtr

				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveBlockHandlers(callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishUnsubscribeMessage(c.UID, pathBlock)
				return

			}(f, err)
		}

		return err
	case pathConfirmedAdded:
		res, err := c.messageProcessor.ProcessConfirmedAdded(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetConfirmedAddedHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *confirmedAddedHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveConfirmedAddedHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathConfirmedAdded, address))
				return

			}(messageInfo.Address, f, err)
		}

		return err

	case pathUnconfirmedAdded:
		res, err := c.messageProcessor.ProcessUnconfirmedAdded(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetUnconfirmedAddedHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *unconfirmedAddedHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveUnconfirmedAddedHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathUnconfirmedAdded, address))
				return
			}(messageInfo.Address, f, err)
		}

		return err

	case pathUnconfirmedRemoved:
		res, err := c.messageProcessor.ProcessUnconfirmedRemoved(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetUnconfirmedRemovedHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *unconfirmedRemovedHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveUnconfirmedRemovedHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, address))
				return
			}(messageInfo.Address, f, err)
		}
		return err

	case pathStatus:
		res, err := c.messageProcessor.ProcessStatus(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetStatusInfoHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *statusHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveStatusInfoHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathStatus, address))
				return
			}(messageInfo.Address, f, err)
		}

		return err
	case pathPartialAdded:
		res, err := c.messageProcessor.ProcessPartialAdded(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetPartialAddedHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *partialAddedHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemovePartialAddedHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathPartialAdded, address))
				return
			}(messageInfo.Address, f, err)
		}

		return err
	case pathPartialRemoved:
		res, err := c.messageProcessor.ProcessPartialRemoved(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetPartialRemovedHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *partialRemovedHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemovePartialRemovedHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathPartialRemoved, address))
				return
			}(messageInfo.Address, f, err)
		}

		return err
	case pathCosignature:
		res, err := c.messageProcessor.ProcessCosignature(resp)
		if err != nil {
			return err
		}

		handlers, err := c.eventSubscribers.GetCosignatureHandlers(messageInfo.Address)
		if err != nil {
			return err
		}

		for f := range handlers {
			go func(address *Address, callFuncPtr *cosignatureHandler, err error) {
				callFunc := *callFuncPtr
				if rm := callFunc(res); !rm {
					return
				}

				shouldUnsubscribe, err := c.eventSubscribers.RemoveCosignatureHandlers(address, callFuncPtr)
				if err != nil {
					return
				}

				if !shouldUnsubscribe {
					return
				}

				err = c.messagePublisher.PublishSubscribeMessage(c.UID, fmt.Sprintf("%s/%s", pathCosignature, address))
				return
			}(messageInfo.Address, f, err)
		}
		return err
	default:
		return unsupportedMessageTypeError
	}
}
