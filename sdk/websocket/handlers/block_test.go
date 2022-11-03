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

func Test_blockHandler_Handle(t *testing.T) {
	type fields struct {
		messageMapper sdk.BlockMapper
		handlers      subscribers.Block
	}
	type args struct {
		handle *sdk.CompoundChannelHandle
		resp   []byte
	}

	blockInfo := new(sdk.BlockInfo)
	messageMapperMock := new(mappers.BlockMapper)
	messageMapperMock.On("MapBlock", mock.Anything).Return(blockInfo, nil)

	handlerFunc1 := func(info *sdk.BlockInfo) bool {
		return false
	}

	handlerFunc2 := func(info *sdk.BlockInfo) bool {
		return true
	}

	blockHandler1 := subscribers.BlockHandler(handlerFunc1)
	blockHandler2 := subscribers.BlockHandler(handlerFunc2)

	blockHandlers := []*subscribers.BlockHandler{
		&blockHandler1,
		&blockHandler2,
	}

	blockHandlersMock := new(mocksSubscribers.Block)
	blockHandlersMock.On("GetHandlers").Return(nil).Once().
		On("GetHandlers").Return(blockHandlers).
		On("RemoveHandlers", mock.Anything).Return(true, nil).
		On("HasHandlers").Return(true, nil)

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
			args: args{},
			want: true,
		},
		{
			name: "remove handlers without error",
			fields: fields{
				handlers:      blockHandlersMock,
				messageMapper: messageMapperMock,
			},
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &blockHandler{
				messageMapper: tt.fields.messageMapper,
				handlers:      tt.fields.handlers,
			}
			got := h.Handle(tt.args.handle, tt.args.resp)
			assert.Equal(t, got, tt.want)
		})
	}
}
