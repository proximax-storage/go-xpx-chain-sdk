// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
