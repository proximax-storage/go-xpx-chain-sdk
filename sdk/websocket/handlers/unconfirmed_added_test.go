package handlers

import (
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-catapult-sdk/mocks/mappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	mocksSubscribers "github.com/proximax-storage/go-xpx-catapult-sdk/mocks/subscribers"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
)

func Test_unconfirmedAddedHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.UnconfirmedAddedMapper
		handlers      subscribers.UnconfirmedAdded
		errCh         chan<- error
	}
	type args struct {
		address *sdk.Address
		resp    []byte
	}

	errCh := make(chan error, 10)

	address := new(sdk.Address)

	mappingError := errors.New("block mapping error")
	obj := new(sdk.TransferTransaction)
	messageMapperMock := new(mappers.UnconfirmedAddedMapper)
	messageMapperMock.On("MapUnconfirmedAdded", mock.Anything).Return(nil, mappingError).Once().
		On("MapUnconfirmedAdded", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(sdk.Transaction) bool {
		return false
	}

	handlerFunc2 := func(sdk.Transaction) bool {
		return true
	}

	handler1 := subscribers.UnconfirmedAddedHandler(handlerFunc1)
	handler2 := subscribers.UnconfirmedAddedHandler(handlerFunc2)

	handlers := map[*subscribers.UnconfirmedAddedHandler]struct{}{
		&handler1: {},
		&handler2: {},
	}

	removingHandlerError := errors.New("removing handler error")
	HandlersMock := new(mocksSubscribers.UnconfirmedAdded)
	HandlersMock.On("GetHandlers", mock.Anything).Return(nil).Once().
		On("GetHandlers", mock.Anything).Return(handlers).
		On("RemoveHandlers", mock.Anything, mock.Anything).Return(true, removingHandlerError).Once().
		On("RemoveHandlers", mock.Anything, mock.Anything).Return(true, nil).
		On("HasHandlers", mock.Anything).Return(true, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "message mapper error",
			fields: fields{
				messageMapper: messageMapperMock,
				errCh:         errCh,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "empty handlers",
			fields: fields{
				handlers:      HandlersMock,
				messageMapper: messageMapperMock,
				errCh:         errCh,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "remove handlers with error",
			fields: fields{
				handlers:      HandlersMock,
				messageMapper: messageMapperMock,
				errCh:         errCh,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "remove handlers without error",
			fields: fields{
				handlers:      HandlersMock,
				messageMapper: messageMapperMock,
				errCh:         errCh,
			},
			args: args{
				address: address,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &unconfirmedAddedHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
				errCh:         tt.fields.errCh,
			}

			got := h.Handle(tt.args.address, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
