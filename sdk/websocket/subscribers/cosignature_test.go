package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var cosignatureHandlerFunc1 = func(tx *sdk.SignerInfo) bool {
	return false
}

var cosignatureHandlerFunc2 = func(tx *sdk.SignerInfo) bool {
	return false
}

func Test_cosignatureImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []CosignatureHandler
	}
	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)
	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*CosignatureHandler)
	subscribers[address.Address] = make([]*CosignatureHandler, 0)

	subscribersNilHandlers := make(map[string][]*CosignatureHandler)

	tests := []struct {
		name    string
		e       *cosignatureImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &cosignatureImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address:  address,
				handlers: []CosignatureHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &cosignatureImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
				handlers: []CosignatureHandler{
					cosignatureHandlerFunc1Ptr,
					cosignatureHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &cosignatureImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
				handlers: []CosignatureHandler{
					cosignatureHandlerFunc1Ptr,
					cosignatureHandlerFunc2Ptr,
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

func Test_cosignatureImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*CosignatureHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string][]*CosignatureHandler)
	emptySubscribers[address.Address] = make([]*CosignatureHandler, 0)

	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*CosignatureHandler)
	hasSubscribersStorage[address.Address] = make([]*CosignatureHandler, 2)
	hasSubscribersStorage[address.Address][0] = &cosignatureHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &cosignatureHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*CosignatureHandler)
	oneSubsctiberStorage[address.Address] = make([]*CosignatureHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &cosignatureHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *cosignatureImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &cosignatureImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address:  address,
				handlers: []*CosignatureHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &cosignatureImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
				handlers: []*CosignatureHandler{
					&cosignatureHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &cosignatureImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address:  address,
				handlers: []*CosignatureHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &cosignatureImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
				handlers: []*CosignatureHandler{
					&cosignatureHandlerFunc1Ptr,
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

func Test_cosignatureImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)

	emptySubscribers := make(map[string][]*CosignatureHandler)
	emptySubscribers[address.Address] = make([]*CosignatureHandler, 0)

	hasSubscribersStorage := make(map[string][]*CosignatureHandler)
	hasSubscribersStorage[address.Address] = make([]*CosignatureHandler, 2)
	hasSubscribersStorage[address.Address][0] = &cosignatureHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &cosignatureHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *cosignatureImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &cosignatureImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &cosignatureImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
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

func Test_cosignatureImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)

	nilSubscribers := make(map[string][]*CosignatureHandler)
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*CosignatureHandler)
	hasSubscribersStorage[address.Address] = make([]*CosignatureHandler, 2)
	hasSubscribersStorage[address.Address][0] = &cosignatureHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &cosignatureHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *cosignatureImpl
		args args
		want []*CosignatureHandler
	}{
		{
			name: "success",
			e: &cosignatureImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &cosignatureImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *cosignatureSubscription),
				removeSubscriberCh: make(chan *cosignatureSubscription),
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
