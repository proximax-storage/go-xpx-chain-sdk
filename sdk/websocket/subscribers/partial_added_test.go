package subscribers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var partialAddedHandlerFunc1 = func(atx *sdk.AggregateTransaction) bool {
	return false
}

var partialAddedHandlerFunc2 = func(atx *sdk.AggregateTransaction) bool {
	return false
}

func Test_partialAddedImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []PartialAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*PartialAddedHandler)
	subscribers[address.Address] = make([]*PartialAddedHandler, 0)

	subscribersNilHandlers := make(map[string][]*PartialAddedHandler)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialAddedHandler, 0)

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
				address:  address,
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
				address: address,
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
				address: address,
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
			err := tt.e.AddHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_partialAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*PartialAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*PartialAddedHandler)
	subscribers[address.Address] = make([]*PartialAddedHandler, 0)

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialdAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialAddedHandler, 2)

	hasSubscribersStorage[address.Address][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialdAddedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*PartialAddedHandler)
	oneSubsctiberStorage[address.Address] = make([]*PartialAddedHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &partialAddedHandlerFunc1Ptr

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
				address:  address,
				handlers: []*PartialAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &partialAddedImpl{
				subscribers:        make(map[string][]*PartialAddedHandler),
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
			},
			args: args{
				address: address,
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
				address:  address,
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
				address: address,
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
			got := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_partialAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	emptySubscribers := make(map[string][]*PartialAddedHandler)
	emptySubscribers[address.Address] = make([]*PartialAddedHandler, 0)

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialAddedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialAddedHandlerFunc2Ptr

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
				address: address,
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
				address: address,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.HasHandlers(tt.args.address)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_partialAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	nilSubscribers := make(map[string]map[*PartialAddedHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*PartialAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*PartialAddedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &partialAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &partialAddedHandlerFunc2Ptr

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
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &partialAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *partialAddedSubscription),
				removeSubscriberCh: make(chan *partialAddedSubscription),
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
				t.Errorf("partialAddedImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
