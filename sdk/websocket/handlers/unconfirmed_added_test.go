package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mappers "github.com/proximax-storage/go-xpx-chain-sdk/mocks"

	mocksSubscribers "github.com/proximax-storage/go-xpx-chain-sdk/mocks/websocket/subscribers"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
)

func Test_unconfirmedAddedHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.UnconfirmedAddedMapper
		handlers      subscribers.UnconfirmedAdded
	}
	type args struct {
		handle *sdk.CompoundChannelHandle
		resp   []byte
	}

	handle := sdk.NewCompoundChannelHandleFromAddress(new(sdk.Address))

	obj := new(sdk.TransferTransaction)
	messageMapperMock := new(mappers.UnconfirmedAddedMapper)
	messageMapperMock.On("MapUnconfirmedAdded", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(sdk.Transaction) bool {
		return false
	}

	handlerFunc2 := func(sdk.Transaction) bool {
		return true
	}

	handler1 := subscribers.UnconfirmedAddedHandler(handlerFunc1)
	handler2 := subscribers.UnconfirmedAddedHandler(handlerFunc2)

	handlers := []*subscribers.UnconfirmedAddedHandler{
		&handler1,
		&handler2,
	}

	HandlersMock := new(mocksSubscribers.UnconfirmedAdded)
	HandlersMock.On("GetHandlers", mock.Anything).Return(nil).Once().
		On("GetHandlers", mock.Anything).Return(handlers).
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
				handle: handle,
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
				handle: handle,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &unconfirmedAddedHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
			}

			got := h.Handle(tt.args.handle, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
