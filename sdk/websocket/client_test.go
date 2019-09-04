// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks2 "github.com/proximax-storage/go-xpx-chain-sdk/mocks"
	mocks "github.com/proximax-storage/go-xpx-chain-sdk/mocks/subscribers"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func TestCatapultWebsocketClientImpl_AddBlockHandlers(t *testing.T) {
	type fields struct {
		UID              string
		blockSubscriber  subscribers.Block
		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		handlers []subscribers.BlockHandler
	}

	uid := "123456"

	handler1 := func(_ *sdk.BlockInfo) bool {
		return false
	}

	handler2 := func(_ *sdk.BlockInfo) bool {
		return false
	}

	userHandlers := []subscribers.BlockHandler{
		handler1,
		handler2,
	}

	emptyTopicHandler := make(topicHandlers)

	publishSubscribeMessageError := errors.New("PublishSubscribeMessage error")
	messagePublisherErrorObj := new(MockMessagePublisher)
	messagePublisherErrorObj.On("PublishSubscribeMessage", uid, Path("block")).Return(publishSubscribeMessageError)

	messagePublisherSuccessObj := new(MockMessagePublisher)
	messagePublisherSuccessObj.On("PublishSubscribeMessage", uid, Path("block")).Return(nil)

	blockSubscriberEmptyHandlersObj := new(mocks.Block)
	blockSubscriberEmptyHandlersObj.On("HasHandlers").Return(false)

	blockSubscriberError := errors.New("block subscription error")
	blockSubscriberErrorObj := new(mocks.Block)
	blockSubscriberErrorObj.On("HasHandlers").Return(true).Once().
		On("AddHandlers", mock.Anything, mock.Anything).Return(blockSubscriberError)

	blockSubscriberSuccessObj := new(mocks.Block)
	blockSubscriberSuccessObj.On("HasHandlers").Return(true).Once().
		On("AddHandlers", mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "empty handlers arg",
			fields: fields{},
			args: args{
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:              uid,
				blockSubscriber:  blockSubscriberEmptyHandlersObj,
				topicHandlers:    emptyTopicHandler,
				messagePublisher: messagePublisherErrorObj,
			},
			args: args{
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "block subscriber add handler error error",
			fields: fields{
				UID:              uid,
				blockSubscriber:  blockSubscriberErrorObj,
				topicHandlers:    emptyTopicHandler,
				messagePublisher: nil,
			},
			args: args{
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:              uid,
				blockSubscriber:  blockSubscriberSuccessObj,
				topicHandlers:    emptyTopicHandler,
				messagePublisher: nil,
			},
			args: args{
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:              tt.fields.UID,
				blockSubscriber:  tt.fields.blockSubscriber,
				topicHandlers:    tt.fields.topicHandlers,
				messagePublisher: tt.fields.messagePublisher,
			}

			err := c.AddBlockHandlers(tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}

}

func TestCatapultWebsocketClientImpl_AddConfirmedAddedHandlers(t *testing.T) {
	type fields struct {
		UID                       string
		confirmedAddedSubscribers subscribers.ConfirmedAdded
		topicHandlers             TopicHandlersStorage
		messagePublisher          MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.ConfirmedAddedHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(sdk.Transaction) bool {
		return false
	}

	handler2 := func(sdk.Transaction) bool {
		return false
	}

	userHandlers := []subscribers.ConfirmedAddedHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.ConfirmedAdded)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				confirmedAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				confirmedAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				confirmedAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:                       tt.fields.UID,
				confirmedAddedSubscribers: tt.fields.confirmedAddedSubscribers,
				topicHandlers:             tt.fields.topicHandlers,
				messagePublisher:          tt.fields.messagePublisher,
			}
			err := c.AddConfirmedAddedHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddUnconfirmedAddedHandlers(t *testing.T) {
	type fields struct {
		UID                         string
		unconfirmedAddedSubscribers subscribers.UnconfirmedAdded
		topicHandlers               TopicHandlersStorage
		messagePublisher            MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.UnconfirmedAddedHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(sdk.Transaction) bool {
		return false
	}

	handler2 := func(sdk.Transaction) bool {
		return false
	}

	userHandlers := []subscribers.UnconfirmedAddedHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.UnconfirmedAdded)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                         uid,
				topicHandlers:               mockTopicHandler,
				messagePublisher:            mockMessagePublisher,
				unconfirmedAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                         uid,
				topicHandlers:               mockTopicHandler,
				messagePublisher:            mockMessagePublisher,
				unconfirmedAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{

				UID:                         tt.fields.UID,
				unconfirmedAddedSubscribers: tt.fields.unconfirmedAddedSubscribers,
				topicHandlers:               tt.fields.topicHandlers,
				messagePublisher:            tt.fields.messagePublisher,
			}
			err := c.AddUnconfirmedAddedHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddUnconfirmedRemovedHandlers(t *testing.T) {
	type fields struct {
		UID                           string
		unconfirmedRemovedSubscribers subscribers.UnconfirmedRemoved

		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.UnconfirmedRemovedHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(*sdk.UnconfirmedRemoved) bool {
		return false
	}

	handler2 := func(*sdk.UnconfirmedRemoved) bool {
		return false
	}

	userHandlers := []subscribers.UnconfirmedRemovedHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.UnconfirmedRemoved)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                           uid,
				topicHandlers:                 mockTopicHandler,
				messagePublisher:              mockMessagePublisher,
				unconfirmedRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                           uid,
				topicHandlers:                 mockTopicHandler,
				messagePublisher:              mockMessagePublisher,
				unconfirmedRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:                           uid,
				topicHandlers:                 mockTopicHandler,
				messagePublisher:              mockMessagePublisher,
				unconfirmedRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:                           tt.fields.UID,
				topicHandlers:                 tt.fields.topicHandlers,
				messagePublisher:              tt.fields.messagePublisher,
				unconfirmedRemovedSubscribers: tt.fields.unconfirmedRemovedSubscribers,
			}
			err := c.AddUnconfirmedRemovedHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddPartialAddedHandlers(t *testing.T) {
	type fields struct {
		UID                     string
		partialAddedSubscribers subscribers.PartialAdded

		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.PartialAddedHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(*sdk.AggregateTransaction) bool {
		return false
	}

	handler2 := func(*sdk.AggregateTransaction) bool {
		return false
	}

	userHandlers := []subscribers.PartialAddedHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.PartialAdded)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                     uid,
				topicHandlers:           mockTopicHandler,
				messagePublisher:        mockMessagePublisher,
				partialAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                     uid,
				topicHandlers:           mockTopicHandler,
				messagePublisher:        mockMessagePublisher,
				partialAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:                     uid,
				topicHandlers:           mockTopicHandler,
				messagePublisher:        mockMessagePublisher,
				partialAddedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{

				UID:                     tt.fields.UID,
				partialAddedSubscribers: tt.fields.partialAddedSubscribers,
				topicHandlers:           tt.fields.topicHandlers,
				messagePublisher:        tt.fields.messagePublisher,
			}
			err := c.AddPartialAddedHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddPartialRemovedHandlers(t *testing.T) {
	type fields struct {
		UID                       string
		partialRemovedSubscribers subscribers.PartialRemoved

		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.PartialRemovedHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(*sdk.PartialRemovedInfo) bool {
		return false
	}

	handler2 := func(*sdk.PartialRemovedInfo) bool {
		return false
	}

	userHandlers := []subscribers.PartialRemovedHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.PartialRemoved)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				partialRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				partialRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:                       uid,
				topicHandlers:             mockTopicHandler,
				messagePublisher:          mockMessagePublisher,
				partialRemovedSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:                       tt.fields.UID,
				partialRemovedSubscribers: tt.fields.partialRemovedSubscribers,
				topicHandlers:             tt.fields.topicHandlers,
				messagePublisher:          tt.fields.messagePublisher,
			}
			err := c.AddPartialRemovedHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddStatusHandlers(t *testing.T) {
	type fields struct {
		UID               string
		statusSubscribers subscribers.Status

		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.StatusHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(*sdk.StatusInfo) bool {
		return false
	}

	handler2 := func(*sdk.StatusInfo) bool {
		return false
	}

	userHandlers := []subscribers.StatusHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.Status)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:               uid,
				topicHandlers:     mockTopicHandler,
				messagePublisher:  mockMessagePublisher,
				statusSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:               uid,
				topicHandlers:     mockTopicHandler,
				messagePublisher:  mockMessagePublisher,
				statusSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:               uid,
				topicHandlers:     mockTopicHandler,
				messagePublisher:  mockMessagePublisher,
				statusSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:               tt.fields.UID,
				statusSubscribers: tt.fields.statusSubscribers,
				topicHandlers:     tt.fields.topicHandlers,
				messagePublisher:  tt.fields.messagePublisher,
			}
			err := c.AddStatusHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_AddCosignatureHandlers(t *testing.T) {
	type fields struct {
		UID                    string
		cosignatureSubscribers subscribers.Cosignature

		topicHandlers    TopicHandlersStorage
		messagePublisher MessagePublisher
	}
	type args struct {
		address  *sdk.Address
		handlers []subscribers.CosignatureHandler
	}

	uid := "123456"
	address := new(sdk.Address)

	handler1 := func(*sdk.SignerInfo) bool {
		return false
	}

	handler2 := func(*sdk.SignerInfo) bool {
		return false
	}

	userHandlers := []subscribers.CosignatureHandler{
		handler1,
		handler2,
	}

	messagePublisherError := errors.New("message publisher error")
	mockAddress := new(sdk.Address)
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.On("PublishSubscribeMessage", uid, mock.Anything).Return(messagePublisherError).Once().
		On("PublishSubscribeMessage", uid, mock.Anything).Return(nil)

	mockTopicHandler := new(MockTopicHandlersStorage)
	mockTopicHandler.On("HasHandler", mock.Anything).Return(false).Once().
		On("HasHandler", mock.Anything).Return(true).
		On("SetTopicHandler", mock.Anything, mock.Anything).Return(nil)

	subscribersAddHandlersError := errors.New("error adding handlers")
	mockSubscribers := new(mocks.Cosignature)
	mockSubscribers.On("HasHandlers", address).Return(false).Once().
		On("HasHandlers", address).Return(true).
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(subscribersAddHandlersError).Once().
		On("AddHandlers", address, mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "nil handlers",
			fields: fields{
				UID: uid,
			},
			args: args{
				address:  mockAddress,
				handlers: nil,
			},
			wantErr: false,
		},
		{
			name: "message publisher error",
			fields: fields{
				UID:                    uid,
				topicHandlers:          mockTopicHandler,
				messagePublisher:       mockMessagePublisher,
				cosignatureSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "add handlers error",
			fields: fields{
				UID:                    uid,
				topicHandlers:          mockTopicHandler,
				messagePublisher:       mockMessagePublisher,
				cosignatureSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				UID:                    uid,
				topicHandlers:          mockTopicHandler,
				messagePublisher:       mockMessagePublisher,
				cosignatureSubscribers: mockSubscribers,
			},
			args: args{
				address:  mockAddress,
				handlers: userHandlers,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				UID:                    tt.fields.UID,
				cosignatureSubscribers: tt.fields.cosignatureSubscribers,
				topicHandlers:          tt.fields.topicHandlers,
				messagePublisher:       tt.fields.messagePublisher,
			}
			err := c.AddCosignatureHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestCatapultWebsocketClientImpl_reconnect(t *testing.T) {
	type fields struct {
		config                        *sdk.Config
		conn                          *websocket.Conn
		ctx                           context.Context
		cancelFunc                    context.CancelFunc
		UID                           string
		blockSubscriber               subscribers.Block
		confirmedAddedSubscribers     subscribers.ConfirmedAdded
		unconfirmedAddedSubscribers   subscribers.UnconfirmedAdded
		unconfirmedRemovedSubscribers subscribers.UnconfirmedRemoved
		partialAddedSubscribers       subscribers.PartialAdded
		partialRemovedSubscribers     subscribers.PartialRemoved
		statusSubscribers             subscribers.Status
		cosignatureSubscribers        subscribers.Cosignature
		topicHandlers                 TopicHandlersStorage
		messageRouter                 Router
		messagePublisher              MessagePublisher
		connectFn                     func(cfg *sdk.Config) (*websocket.Conn, string, error)
		alreadyListening              bool
	}

	errorConnectFn := func(cfg *sdk.Config) (*websocket.Conn, string, error) {
		return nil, "", errors.New("test error")
	}

	successConnectFn := func(cfg *sdk.Config) (*websocket.Conn, string, error) {
		return nil, "test-uid", nil
	}

	mockRouter := new(mocks2.Router)
	mockRouter.On("SetUid", mock.Anything).Return(nil)

	publisherError := errors.New("test publisher error")
	mockMessagePublisher := new(MockMessagePublisher)
	mockMessagePublisher.
		On("SetConn", mock.Anything).Return().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s", pathBlock))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s", pathBlock))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, "test confirmed added address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathConfirmedAdded, "test confirmed added address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathCosignature, "test cosignature address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathCosignature, "test cosignature address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathPartialAdded, "test partial added address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathPartialAdded, "test partial added address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, "test partial removed address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathPartialRemoved, "test partial removed address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathStatus, "test status address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathStatus, "test status address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, "test unconfirmed added address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathUnconfirmedAdded, "test unconfirmed added address"))).Return(nil).
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, "test unconfirmed removed address"))).Return(publisherError).Once().
		On("PublishSubscribeMessage", mock.Anything, Path(fmt.Sprintf("%s/%s", pathUnconfirmedRemoved, "test unconfirmed removed address"))).Return(nil)

	mockBlockSubscriber := new(mocks.Block)
	mockBlockSubscriber.On("HasHandlers").Return(true)

	mockConfirmedAddedSubscriber := new(mocks.ConfirmedAdded)
	mockConfirmedAddedSubscriber.On("GetAddresses").Return([]string{"test confirmed added address"})

	mockCosignatureSubscriber := new(mocks.Cosignature)
	mockCosignatureSubscriber.On("GetAddresses").Return([]string{"test cosignature address"})

	mockPartialAddedSubscriber := new(mocks.PartialAdded)
	mockPartialAddedSubscriber.On("GetAddresses").Return([]string{"test partial added address"})

	mockPartialRemovedSubscriber := new(mocks.PartialRemoved)
	mockPartialRemovedSubscriber.On("GetAddresses").Return([]string{"test partial removed address"})

	mockStatusSubscriber := new(mocks.Status)
	mockStatusSubscriber.On("GetAddresses").Return([]string{"test status address"})

	mockUnconfirmedAddedSubscriber := new(mocks.UnconfirmedAdded)
	mockUnconfirmedAddedSubscriber.On("GetAddresses").Return([]string{"test unconfirmed added address"})

	mockUnconfirmedRemovedSubscriber := new(mocks.UnconfirmedRemoved)
	mockUnconfirmedRemovedSubscriber.On("GetAddresses").Return([]string{"test unconfirmed removed address"})

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "connection error",
			fields: fields{
				config:    &sdk.Config{},
				connectFn: errorConnectFn,
			},
			wantErr: true,
		},
		{
			name: "block message publisher error",
			fields: fields{
				config:           &sdk.Config{},
				connectFn:        successConnectFn,
				blockSubscriber:  mockBlockSubscriber,
				messagePublisher: mockMessagePublisher,
				messageRouter:    mockRouter,
			},
			wantErr: true,
		},
		{
			name: "confirmed added message publisher error",
			fields: fields{
				config:                    &sdk.Config{},
				connectFn:                 successConnectFn,
				blockSubscriber:           mockBlockSubscriber,
				confirmedAddedSubscribers: mockConfirmedAddedSubscriber,
				messagePublisher:          mockMessagePublisher,
				messageRouter:             mockRouter,
			},
			wantErr: true,
		},
		{
			name: "confirmed added message publisher error",
			fields: fields{
				config:                    &sdk.Config{},
				connectFn:                 successConnectFn,
				blockSubscriber:           mockBlockSubscriber,
				confirmedAddedSubscribers: mockConfirmedAddedSubscriber,
				cosignatureSubscribers:    mockCosignatureSubscriber,
				messagePublisher:          mockMessagePublisher,
				messageRouter:             mockRouter,
			},
			wantErr: true,
		},
		{
			name: "partial added message publisher error",
			fields: fields{
				config:                    &sdk.Config{},
				connectFn:                 successConnectFn,
				blockSubscriber:           mockBlockSubscriber,
				confirmedAddedSubscribers: mockConfirmedAddedSubscriber,
				cosignatureSubscribers:    mockCosignatureSubscriber,
				partialAddedSubscribers:   mockPartialAddedSubscriber,
				messagePublisher:          mockMessagePublisher,
				messageRouter:             mockRouter,
			},
			wantErr: true,
		},
		{
			name: "partial removed message publisher error",
			fields: fields{
				config:                    &sdk.Config{},
				connectFn:                 successConnectFn,
				blockSubscriber:           mockBlockSubscriber,
				confirmedAddedSubscribers: mockConfirmedAddedSubscriber,
				cosignatureSubscribers:    mockCosignatureSubscriber,
				partialAddedSubscribers:   mockPartialAddedSubscriber,
				partialRemovedSubscribers: mockPartialRemovedSubscriber,
				messagePublisher:          mockMessagePublisher,
				messageRouter:             mockRouter,
			},
			wantErr: true,
		},
		{
			name: "status message publisher error",
			fields: fields{
				config:                      &sdk.Config{},
				connectFn:                   successConnectFn,
				blockSubscriber:             mockBlockSubscriber,
				confirmedAddedSubscribers:   mockConfirmedAddedSubscriber,
				cosignatureSubscribers:      mockCosignatureSubscriber,
				partialAddedSubscribers:     mockPartialAddedSubscriber,
				partialRemovedSubscribers:   mockPartialRemovedSubscriber,
				statusSubscribers:           mockStatusSubscriber,
				unconfirmedAddedSubscribers: mockUnconfirmedAddedSubscriber,
				messagePublisher:            mockMessagePublisher,
				messageRouter:               mockRouter,
			},
			wantErr: true,
		},
		{
			name: "unconfirmed added message publisher error",
			fields: fields{
				config:                      &sdk.Config{},
				connectFn:                   successConnectFn,
				blockSubscriber:             mockBlockSubscriber,
				confirmedAddedSubscribers:   mockConfirmedAddedSubscriber,
				cosignatureSubscribers:      mockCosignatureSubscriber,
				partialAddedSubscribers:     mockPartialAddedSubscriber,
				partialRemovedSubscribers:   mockPartialRemovedSubscriber,
				statusSubscribers:           mockStatusSubscriber,
				unconfirmedAddedSubscribers: mockUnconfirmedAddedSubscriber,
				messagePublisher:            mockMessagePublisher,
				messageRouter:               mockRouter,
			},
			wantErr: true,
		},
		{
			name: "unconfirmed removed message publisher error",
			fields: fields{
				config:                        &sdk.Config{},
				connectFn:                     successConnectFn,
				blockSubscriber:               mockBlockSubscriber,
				confirmedAddedSubscribers:     mockConfirmedAddedSubscriber,
				cosignatureSubscribers:        mockCosignatureSubscriber,
				partialAddedSubscribers:       mockPartialAddedSubscriber,
				partialRemovedSubscribers:     mockPartialRemovedSubscriber,
				statusSubscribers:             mockStatusSubscriber,
				unconfirmedAddedSubscribers:   mockUnconfirmedAddedSubscriber,
				unconfirmedRemovedSubscribers: mockUnconfirmedRemovedSubscriber,
				messagePublisher:              mockMessagePublisher,
				messageRouter:                 mockRouter,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				config:                        &sdk.Config{},
				connectFn:                     successConnectFn,
				blockSubscriber:               mockBlockSubscriber,
				confirmedAddedSubscribers:     mockConfirmedAddedSubscriber,
				cosignatureSubscribers:        mockCosignatureSubscriber,
				partialAddedSubscribers:       mockPartialAddedSubscriber,
				partialRemovedSubscribers:     mockPartialRemovedSubscriber,
				statusSubscribers:             mockStatusSubscriber,
				unconfirmedAddedSubscribers:   mockUnconfirmedAddedSubscriber,
				unconfirmedRemovedSubscribers: mockUnconfirmedRemovedSubscriber,
				messagePublisher:              mockMessagePublisher,
				messageRouter:                 mockRouter,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CatapultWebsocketClientImpl{
				config:                        tt.fields.config,
				conn:                          tt.fields.conn,
				ctx:                           tt.fields.ctx,
				cancelFunc:                    tt.fields.cancelFunc,
				UID:                           tt.fields.UID,
				blockSubscriber:               tt.fields.blockSubscriber,
				confirmedAddedSubscribers:     tt.fields.confirmedAddedSubscribers,
				unconfirmedAddedSubscribers:   tt.fields.unconfirmedAddedSubscribers,
				unconfirmedRemovedSubscribers: tt.fields.unconfirmedRemovedSubscribers,
				partialAddedSubscribers:       tt.fields.partialAddedSubscribers,
				partialRemovedSubscribers:     tt.fields.partialRemovedSubscribers,
				statusSubscribers:             tt.fields.statusSubscribers,
				cosignatureSubscribers:        tt.fields.cosignatureSubscribers,
				topicHandlers:                 tt.fields.topicHandlers,
				messageRouter:                 tt.fields.messageRouter,
				messagePublisher:              tt.fields.messagePublisher,
				alreadyListening:              tt.fields.alreadyListening,
				connectFn:                     tt.fields.connectFn,
			}

			err := c.reconnect()
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
