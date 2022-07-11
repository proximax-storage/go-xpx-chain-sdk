package subscribers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var partialAddedHandlerFunc1 = func(atx sdk.Transaction) bool {
	return false
}

var partialAddedHandlerFunc2 = func(atx sdk.Transaction) bool {
	return false
}

func Test_partialAddedImpl_AddHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.TransactionChannelHandle
		handlers []PartialAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewTransactionChannelHandleFromAddress(address)
	subscribers := make(map[string][]*PartialAddedHandler)
	subscribers[handle.String()] = make([]*PartialAddedHandler, 0)

	subscribersNilHandlers := make(map[string][]*PartialAddedHandler)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*PartialAddedHandler, 0)

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialdAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	tests := []struct {
		name    string
		e       *partialAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialAddedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []PartialAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &partialAddedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []PartialAddedHandler{
					nil,
					nil,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &partialAddedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []PartialAddedHandler{
					partialAddedHandlerFunc1Ptr,
					partialdAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			err := tt.e.AddHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_partialAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.TransactionChannelHandle
		handlers []*PartialAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewTransactionChannelHandleFromAddress(address)
	subscribers := make(map[string][]*PartialAddedHandler)
	subscribers[handle.String()] = make([]*PartialAddedHandler, 0)

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialdAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*PartialAddedHandler, 2)

	hasSubscribersStorage[handle.String()][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &partialdAddedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*PartialAddedHandler)
	oneSubsctiberStorage[handle.String()] = make([]*PartialAddedHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &partialAddedHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *partialAddedImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*PartialAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &partialAddedImpl{
				subscribers:        make(map[string][]*PartialAddedHandler),
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*PartialAddedHandler{
					&partialAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &partialAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*PartialAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &partialAddedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*PartialAddedHandler{
					&partialAddedHandlerFunc1Ptr,
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
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_partialAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.TransactionChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewTransactionChannelHandleFromAddress(address)
	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	emptySubscribers := make(map[string][]*PartialAddedHandler)
	emptySubscribers[handle.String()] = make([]*PartialAddedHandler, 0)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*PartialAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &partialAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *partialAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &partialAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &partialAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
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
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_partialAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.TransactionChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewTransactionChannelHandleFromAddress(address)
	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	nilSubscribers := make(map[string]map[*PartialAddedHandler]struct{})
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*PartialAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &partialAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *partialAddedImpl
		args args
		want []*PartialAddedHandler
	}{
		{
			name: "success",
			e: &partialAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &partialAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
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
			if got := tt.e.GetHandlers(tt.args.handle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("partialAddedImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
