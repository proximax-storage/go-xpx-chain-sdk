package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var unconfirmedRemovedHandlerFunc1 = func(*sdk.UnconfirmedRemoved) bool {
	return false
}

var unconfirmedRemovedHandlerFunc2 = func(*sdk.UnconfirmedRemoved) bool {
	return false
}

func Test_unconfirmedRemovedImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []UnconfirmedRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*UnconfirmedRemovedHandler)
	subscribers[address.Address] = make([]*UnconfirmedRemovedHandler, 0)

	subscribersNilHandlers := make(map[string][]*UnconfirmedRemovedHandler)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	tests := []struct {
		name    string
		e       *unconfirmedRemovedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []UnconfirmedRemovedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &unconfirmedRemovedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1Ptr,
					unconfirmedRemovedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1Ptr,
					unconfirmedRemovedHandlerFunc2Ptr,
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

func Test_unconfirmedRemovedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*UnconfirmedRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	emptySubscribers[address.Address] = make([]*UnconfirmedRemovedHandler, 0)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &unconfirmedRemovedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*UnconfirmedRemovedHandler)
	oneSubsctiberStorage[address.Address] = make([]*UnconfirmedRemovedHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &unconfirmedRemovedHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *unconfirmedRemovedImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*UnconfirmedRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &unconfirmedRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedRemovedHandler{
					&unconfirmedRemovedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address:  address,
				handlers: []*UnconfirmedRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &unconfirmedRemovedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedRemovedHandler{
					&unconfirmedRemovedHandlerFunc1Ptr,
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
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_unconfirmedRemovedImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)

	emptySubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	emptySubscribers[address.Address] = make([]*UnconfirmedRemovedHandler, 0)

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &unconfirmedRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &unconfirmedRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
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

func Test_unconfirmedRemovedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	nilSubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[address.Address] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[address.Address][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &unconfirmedRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want []*UnconfirmedRemovedHandler
	}{
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &unconfirmedRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
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
			assert.Equal(t, got, tt.want)
		})
	}
}
