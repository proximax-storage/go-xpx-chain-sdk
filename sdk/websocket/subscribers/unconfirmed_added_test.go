package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var unconfirmedAddedHandlerFunc1 = func(tx sdk.Transaction) bool {
	return false
}

var unconfirmedAddedHandlerFunc2 = func(tx sdk.Transaction) bool {
	return false
}

func Test_unconfirmedAddedImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []UnconfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	subscribers[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*UnconfirmedAddedHandler]struct{})

	tests := []struct {
		name    string
		e       *unconfirmedAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedAddedImpl{
				subscribers: nil,
			},
			args: args{
				address:  address,
				handlers: []UnconfirmedAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &unconfirmedAddedImpl{
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []UnconfirmedAddedHandler{
					unconfirmedAddedHandlerFunc1,
					unconfirmedAddedHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &unconfirmedAddedImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []UnconfirmedAddedHandler{
					unconfirmedAddedHandlerFunc1,
					unconfirmedAddedHandlerFunc1,
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

func Test_unconfirmedAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*UnconfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)

	hasSubscribersStorage := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})
	oneSubsctiberStorage[address.Address][&unconfirmedAddedHandlerFunc1Ptr] = struct{}{}

	tests := []struct {
		name    string
		e       *unconfirmedAddedImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedAddedImpl{
				subscribers: nil,
			},
			args: args{
				address:  address,
				handlers: []*UnconfirmedAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &unconfirmedAddedImpl{
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedAddedHandler{
					&unconfirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &unconfirmedAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedAddedHandler{
					&unconfirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &unconfirmedAddedImpl{
				subscribers: oneSubsctiberStorage,
			},
			args: args{
				address: address,
				handlers: []*UnconfirmedAddedHandler{
					&unconfirmedAddedHandlerFunc1Ptr,
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

func Test_unconfirmedAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc2)

	emptySubscribers := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *unconfirmedAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &unconfirmedAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &unconfirmedAddedImpl{
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

func Test_unconfirmedAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc2)

	nilSubscribers := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*UnconfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&unconfirmedAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *unconfirmedAddedImpl
		args args
		want map[*UnconfirmedAddedHandler]struct{}
	}{
		{
			name: "success",
			e: &unconfirmedAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &unconfirmedAddedImpl{
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
