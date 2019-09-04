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

	subscribers := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	subscribers[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*UnconfirmedRemovedHandler]struct{})

	tests := []struct {
		name    string
		e       *unconfirmedRemovedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedRemovedImpl{
				subscribers: nil,
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
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1,
					unconfirmedRemovedHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1,
					unconfirmedRemovedHandlerFunc1,
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

func Test_unconfirmedRemovedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*UnconfirmedRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})
	oneSubsctiberStorage[address.Address][&unconfirmedRemovedHandlerFunc1Ptr] = struct{}{}

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
				subscribers: nil,
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
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedRemovedHandler{
					&unconfirmedRemovedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &unconfirmedRemovedImpl{
				subscribers: hasSubscribersStorage,
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
			name: "success return true result",
			e: &unconfirmedRemovedImpl{
				subscribers: oneSubsctiberStorage,
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

	emptySubscribers := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &unconfirmedRemovedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &unconfirmedRemovedImpl{
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

func Test_unconfirmedRemovedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	nilSubscribers := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedRemovedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedRemovedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want map[*UnconfirmedRemovedHandler]struct{}
	}{
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &unconfirmedRemovedImpl{
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
			got := tt.e.GetHandlers(tt.args.address)
			assert.Equal(t, got, tt.want)
		})
	}
}
