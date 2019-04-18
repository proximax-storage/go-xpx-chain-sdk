package subscribers

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
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

	subscribers := make(map[string]map[*PartialRemovedHandler]struct{})
	subscribers[address.Address] = make(map[*PartialRemovedHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*PartialRemovedHandler]struct{})

	tests := []struct {
		name    string
		e       *partialRemovedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &partialRemovedImpl{
				subscribers: nil,
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
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []PartialRemovedHandler{
					partialRemovedHandlerFunc1,
					partialRemovedHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &partialRemovedImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []PartialRemovedHandler{
					partialRemovedHandlerFunc1,
					partialRemovedHandlerFunc2,
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

func Test_partialRemovedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*PartialRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*PartialRemovedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*PartialRemovedHandler]struct{})

	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*PartialRemovedHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*PartialRemovedHandler]struct{})
	oneSubsctiberStorage[address.Address][&partialRemovedHandlerFunc1Ptr] = struct{}{}

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
				subscribers: nil,
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
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*PartialRemovedHandler{
					&partialRemovedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &partialRemovedImpl{
				subscribers: hasSubscribersStorage,
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
			name: "success return true result",
			e: &partialRemovedImpl{
				subscribers: oneSubsctiberStorage,
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
			got, err := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
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

	emptySubscribers := make(map[string]map[*PartialRemovedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*PartialRemovedHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *partialRemovedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &partialRemovedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &partialRemovedImpl{
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

func Test_partialRemovedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	partialRemovedHandlerFunc1Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)
	partialRemovedHandlerFunc2Ptr := PartialRemovedHandler(partialRemovedHandlerFunc1)

	nilSubscribers := make(map[string]map[*PartialRemovedHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*PartialRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&partialRemovedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *partialRemovedImpl
		args args
		want map[*PartialRemovedHandler]struct{}
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
			if got := tt.e.GetHandlers(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("partialRemovedImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
