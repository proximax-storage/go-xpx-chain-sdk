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

func Test_confirmedAddedHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.ConfirmedAddedMapper
		handlers      subscribers.ConfirmedAdded
		errCh         chan<- error
	}
	type args struct {
		address *sdk.Address
		resp    []byte
	}

	errCh := make(chan error, 10)

	address := new(sdk.Address)

	blockMappingError := errors.New("block mapping error")
	obj := new(sdk.TransferTransaction)
	messageMapperMock := new(mappers.ConfirmedAddedMapper)
	messageMapperMock.On("MapConfirmedAdded", mock.Anything).Return(nil, blockMappingError).Once().
		On("MapConfirmedAdded", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(tx sdk.Transaction) bool {
		return false
	}

	handlerFunc2 := func(tx sdk.Transaction) bool {
		return true
	}

	blockHandler1 := subscribers.ConfirmedAddedHandler(handlerFunc1)
	blockHandler2 := subscribers.ConfirmedAddedHandler(handlerFunc2)

	blockHandlers := map[*subscribers.ConfirmedAddedHandler]struct{}{
		&blockHandler1: {},
		&blockHandler2: {},
	}

	removingHandlerError := errors.New("removing handler error")
	blockHandlersMock := new(mocksSubscribers.ConfirmedAdded)
	blockHandlersMock.On("GetHandlers", mock.Anything).Return(nil).Once().
		On("GetHandlers", mock.Anything).Return(blockHandlers).
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
				handlers:      blockHandlersMock,
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
				handlers:      blockHandlersMock,
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
				handlers:      blockHandlersMock,
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
			h := &confirmedAddedHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
				errCh:         tt.fields.errCh,
			}
			got := h.Handle(tt.args.address, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}