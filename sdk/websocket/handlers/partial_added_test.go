package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/proximax-storage/go-xpx-chain-sdk/mocks/mappers"

	mocksSubscribers "github.com/proximax-storage/go-xpx-chain-sdk/mocks/websocket/subscribers"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func Test_partialAddedHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.PartialAddedMapper
		handlers      subscribers.PartialAdded
	}
	type args struct {
		address *sdk.Address
		resp    []byte
	}

	address := new(sdk.Address)

	obj := new(sdk.AggregateTransaction)
	messageMapperMock := new(mappers.PartialAddedMapper)
	messageMapperMock.On("MapPartialAdded", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(*sdk.AggregateTransaction) bool {
		return false
	}

	handlerFunc2 := func(*sdk.AggregateTransaction) bool {
		return true
	}

	blockHandler1 := subscribers.PartialAddedHandler(handlerFunc1)
	blockHandler2 := subscribers.PartialAddedHandler(handlerFunc2)

	blockHandlers := []*subscribers.PartialAddedHandler{
		&blockHandler1,
		&blockHandler2,
	}

	HandlersMock := new(mocksSubscribers.PartialAdded)
	HandlersMock.On("GetHandlers", mock.Anything).Return(nil).Once().
		On("GetHandlers", mock.Anything).Return(blockHandlers).
		On("RemoveHandlers", mock.Anything, mock.Anything).Return(true, nil).
		On("HasHandlers", mock.Anything).Return(true, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty handlers",
			fields: fields{
				handlers:      HandlersMock,
				messageMapper: messageMapperMock,
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
			},
			args: args{
				address: address,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &partialAddedHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
			}
			got := h.Handle(tt.args.address, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
