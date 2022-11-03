package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var receiptHandlerFunc1 = func(*sdk.AnonymousReceipt) bool {
	return false
}

var receiptHandlerFunc2 = func(*sdk.AnonymousReceipt) bool {
	return false
}

func Test_receiptImpl_AddHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []ReceiptHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	subscribers := make(map[string][]*ReceiptHandler)
	subscribers[handle.String()] = make([]*ReceiptHandler, 0)

	subscribersNilHandlers := make(map[string][]*ReceiptHandler)

	receiptHandlerFunc1Ptr := ReceiptHandler(receiptHandlerFunc1)
	receiptHandlerFunc2Ptr := ReceiptHandler(receiptHandlerFunc2)

	tests := []struct {
		name    string
		e       *receiptImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &receiptImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []ReceiptHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &receiptImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
				handlers: []ReceiptHandler{
					receiptHandlerFunc1Ptr,
					receiptHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &receiptImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
				handlers: []ReceiptHandler{
					receiptHandlerFunc1Ptr,
					receiptHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			err := tt.e.AddHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func Test_receiptImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*ReceiptHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	emptySubscribers := make(map[string][]*ReceiptHandler)
	emptySubscribers[handle.String()] = make([]*ReceiptHandler, 0)

	receiptHandlerFunc1Ptr := ReceiptHandler(receiptHandlerFunc1)
	receiptHandlerFunc2Ptr := ReceiptHandler(receiptHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*ReceiptHandler)
	hasSubscribersStorage[handle.String()] = make([]*ReceiptHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &receiptHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &receiptHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*ReceiptHandler)
	oneSubsctiberStorage[handle.String()] = make([]*ReceiptHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &receiptHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *receiptImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &receiptImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*ReceiptHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &receiptImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*ReceiptHandler{
					&receiptHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &receiptImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*ReceiptHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &receiptImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*ReceiptHandler{
					&receiptHandlerFunc1Ptr,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_receiptImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	receiptHandlerFunc1Ptr := ReceiptHandler(receiptHandlerFunc1)
	receiptHandlerFunc2Ptr := ReceiptHandler(receiptHandlerFunc1)

	emptySubscribers := make(map[string][]*ReceiptHandler)
	emptySubscribers[handle.String()] = make([]*ReceiptHandler, 0)

	hasSubscribersStorage := make(map[string][]*ReceiptHandler)
	hasSubscribersStorage[handle.String()] = make([]*ReceiptHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &receiptHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &receiptHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *receiptImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &receiptImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &receiptImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.HasHandlers(tt.args.handle)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_receiptImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	receiptHandlerFunc1Ptr := ReceiptHandler(receiptHandlerFunc1)
	receiptHandlerFunc2Ptr := ReceiptHandler(receiptHandlerFunc2)

	nilSubscribers := make(map[string][]*ReceiptHandler)
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*ReceiptHandler)
	hasSubscribersStorage[handle.String()] = make([]*ReceiptHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &receiptHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &receiptHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *receiptImpl
		args args
		want []*ReceiptHandler
	}{
		{
			name: "success",
			e: &receiptImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &receiptImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *receiptSubscription),
				removeSubscriberCh: make(chan *receiptSubscription),
			},
			args: args{
				handle: handle,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.GetHandlers(tt.args.handle)
			assert.Equal(t, got, tt.want)
		})
	}
}
