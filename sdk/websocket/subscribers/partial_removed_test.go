package subscribers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var partialRemovedHandlerFunc1 = func(info *sdk.PartialRemovedInfo) bool {
	return false
}

var partialRemovedHandlerFunc2 = func(info *sdk.PartialRemovedInfo) bool {
	return false
}

func Test_partialRemovedImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []PartialRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc2)
	subscribers := make(map[string][]*PartialRemovedHandler)
	subscribers[address.Address] = make([]*PartialRemovedHandler, 0)

	subscribersNilHandlers := make(map[string][]*PartialRemovedHandler)

	tests := []struct {
		name    string
		e       *partialRemovedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []PartialRemovedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &partialRemovedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []PartialRemovedHandler{
					partialRemovedHandlerFunc1Ptr,
					partialRemovedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &partialRemovedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []PartialRemovedHandler{
					partialRemovedHandlerFunc1Ptr,
					partialRemovedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			err := tt.e.AddHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func Test_partialRemovedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*PartialRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string][]*PartialRemovedHandler)
	emptySubscribers[address.Address] = make([]*PartialRemovedHandler, 0)

	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*PartialRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &partialRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialRemovedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*PartialRemovedHandler)
	oneSubsctiberStorage[address.Address] = make([]*PartialRemovedHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &partialRemovedHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *partialRemovedImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*PartialRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &partialRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []*PartialRemovedHandler{
					&partialRemovedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &partialRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*PartialRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &partialRemovedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []*PartialRemovedHandler{
					&partialRemovedHandlerFunc1Ptr,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_partialRemovedImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc2)

	emptySubscribers := make(map[string][]*PartialRemovedHandler)
	emptySubscribers[address.Address] = make([]*PartialRemovedHandler, 0)

	hasSubscribersStorage := make(map[string][]*PartialRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &partialRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *partialRemovedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &partialRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &partialRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *partialRemovedSubscription),
				removeSubscriberCh: make(chan *partialRemovedSubscription),
			},
			args: args{
				address: address,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.HasHandlers(tt.args.address)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_partialRemovedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)

	nilSubscribers := make(map[string][]*PartialRemovedHandler)
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*PartialRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &partialRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *partialRemovedImpl
		args args
		want []*PartialRemovedHandler
	}{
		{
			name: "success",
			e: &partialRemovedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &partialRemovedImpl{
				subscribers: nil,
			},
			args: args{
				address: address,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			if got := tt.e.GetHandlers(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("partialRemovedImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
