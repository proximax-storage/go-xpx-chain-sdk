package handlers

import (
	"github.com/proximax-storage/go-xpx-catapult-sdk/mocks/mappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	mocksSubscribers "github.com/proximax-storage/go-xpx-catapult-sdk/mocks/subscribers"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket/subscribers"
)

func Test_partialRemovedHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.PartialRemovedMapper
		handlers      subscribers.PartialRemoved
	}
	type args struct {
		address *sdk.Address
		resp    []byte
	}

	address := new(sdk.Address)

	obj := new(sdk.PartialRemovedInfo)
	messageMapperMock := new(mappers.PartialRemovedMapper)
	messageMapperMock.On("MapPartialRemoved", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(*sdk.PartialRemovedInfo) bool {
		return false
	}

	handlerFunc2 := func(*sdk.PartialRemovedInfo) bool {
		return true
	}

	handler1 := subscribers.PartialRemovedHandler(handlerFunc1)
	handler2 := subscribers.PartialRemovedHandler(handlerFunc2)

	blockHandlers := map[*subscribers.PartialRemovedHandler]struct{}{
		&handler1: {},
		&handler2: {},
	}

	HandlersMock := new(mocksSubscribers.PartialRemoved)
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
			h := &partialRemovedHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
			}
			got := h.Handle(tt.args.address, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
