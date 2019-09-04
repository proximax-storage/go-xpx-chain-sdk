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

	subscribers := make(map[string]map[*PartialAddedHandler]struct{})
	subscribers[address.Address] = make(map[*PartialAddedHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*PartialAddedHandler]struct{})

	tests := []struct {
		name    string
		e       *partialAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialAddedImpl{
				subscribers: nil,
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
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []PartialAddedHandler{
					partialAddedHandlerFunc1,
					partialAddedHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &partialAddedImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []PartialAddedHandler{
					partialAddedHandlerFunc1,
					partialAddedHandlerFunc2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.e.AddHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
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

	emptySubscribers := make(map[string]map[*PartialAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*PartialAddedHandler]struct{})

	partialAddedHandlerFunc1Ptr := PartialAddedHandler(partialAddedHandlerFunc1)
	partialdAddedHandlerFunc2Ptr := PartialAddedHandler(partialAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialdAddedHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*PartialAddedHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*PartialAddedHandler]struct{})
	oneSubsctiberStorage[address.Address][&partialAddedHandlerFunc1Ptr] = struct{}{}

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
				subscribers: nil,
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
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*PartialAddedHandler{
					&partialAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &partialAddedImpl{
				subscribers: hasSubscribersStorage,
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
			name: "success return true result",
			e: &partialAddedImpl{
				subscribers: oneSubsctiberStorage,
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
			got, err := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
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

	emptySubscribers := make(map[string]map[*PartialAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*PartialAddedHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *partialAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &partialAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &partialAddedImpl{
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.HasHandlers(tt.args.address)
			assert.Equal(t, got, tt.want)
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

	hasSubscribersStorage := make(map[string]map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *partialAddedImpl
		args args
		want map[*PartialAddedHandler]struct{}
	}{
		{
			name: "success",
			e: &partialAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &partialAddedImpl{
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
			if got := tt.e.GetHandlers(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("partialAddedImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
