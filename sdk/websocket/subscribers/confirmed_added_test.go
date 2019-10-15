package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var confirmedAddedHandlerFunc1 = func(tx sdk.Transaction) bool {
	return false
}

var confirmedAddedHandlerFunc2 = func(tx sdk.Transaction) bool {
	return false
}

func Test_confirmedAddedImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*ConfirmedAddedHandler)
	subscribers[address.Address] = make([]*ConfirmedAddedHandler, 0)

	subscribersNilHandlers := make(map[string][]*ConfirmedAddedHandler)

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	tests := []struct {
		name    string
		e       *confirmedAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &confirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*ConfirmedAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &confirmedAddedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
					&confirmedAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
					&confirmedAddedHandlerFunc2Ptr,
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

func Test_confirmedAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string][]*ConfirmedAddedHandler)
	emptySubscribers[address.Address] = make([]*ConfirmedAddedHandler, 0)

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)
	confirmedAddedHandlerFunc3Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &confirmedAddedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*ConfirmedAddedHandler)
	oneSubsctiberStorage[address.Address] = make([]*ConfirmedAddedHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &confirmedAddedHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *confirmedAddedImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &confirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*ConfirmedAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &confirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc3Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &confirmedAddedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got, err := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_confirmedAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	emptySubscribers := make(map[string][]*ConfirmedAddedHandler)
	emptySubscribers[address.Address] = make([]*ConfirmedAddedHandler, 0)

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &confirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &confirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
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

func Test_confirmedAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	nilSubscribers := make(map[string][]*ConfirmedAddedHandler)
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[address.Address] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &confirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want []*ConfirmedAddedHandler
	}{
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &confirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
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
			got := tt.e.GetHandlers(tt.args.address)
			assert.Equal(t, tt.want, got)
		})
	}
}
