package handlers

import (
	mappers "github.com/proximax-storage/go-xpx-chain-sdk/mocks"
	mocksSubscribers "github.com/proximax-storage/go-xpx-chain-sdk/mocks/websocket/subscribers"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_driveStateHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.DriveStateMapper
		handlers      subscribers.DriveState
	}
	type args struct {
		handle *sdk.TransactionChannelHandle
		resp   []byte
	}

	handle := sdk.NewTransactionChannelHandleFromAddress(new(sdk.Address))

	obj := new(sdk.DriveStateInfo)
	messageMapperMock := new(mappers.DriveStateMapper)
	messageMapperMock.On("MapDriveState", mock.Anything).Return(obj, nil)

	handlerFunc1 := func(*sdk.DriveStateInfo) bool {
		return false
	}

	handlerFunc2 := func(*sdk.DriveStateInfo) bool {
		return true
	}

	blockHandler1 := subscribers.DriveStateHandler(handlerFunc1)
	blockHandler2 := subscribers.DriveStateHandler(handlerFunc2)

	blockHandlers := []*subscribers.DriveStateHandler{
		&blockHandler1,
		&blockHandler2,
	}

	blockHandlersMock := new(mocksSubscribers.DriveState)
	blockHandlersMock.On("GetHandlers", mock.Anything).Return(nil).Once().
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
				handlers:      blockHandlersMock,
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
				handlers:      blockHandlersMock,
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
			h := &driveStateHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
			}
			got := h.Handle(tt.args.handle, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
