package subscribers

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
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
		handlers []ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string]map[*ConfirmedAddedHandler]struct{})
	subscribers[address.Address] = make(map[*ConfirmedAddedHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*ConfirmedAddedHandler]struct{})

	tests := []struct {
		name    string
		e       *confirmedAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &confirmedAddedImpl{
				subscribers: nil,
			},
			args: args{
				address:  address,
				handlers: []ConfirmedAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &confirmedAddedImpl{
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []ConfirmedAddedHandler{
					confirmedAddedHandlerFunc1,
					confirmedAddedHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []ConfirmedAddedHandler{
					confirmedAddedHandlerFunc1,
					confirmedAddedHandlerFunc2,
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

func Test_confirmedAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*ConfirmedAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*ConfirmedAddedHandler]struct{})

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*ConfirmedAddedHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*ConfirmedAddedHandler]struct{})
	oneSubsctiberStorage[address.Address][&confirmedAddedHandlerFunc1Ptr] = struct{}{}

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
				subscribers: nil,
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
				subscribers: emptySubscribers,
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
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &confirmedAddedImpl{
				subscribers: oneSubsctiberStorage,
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
			got, err := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
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

	emptySubscribers := make(map[string]map[*ConfirmedAddedHandler]struct{})
	emptySubscribers[address.Address] = make(map[*ConfirmedAddedHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &confirmedAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &confirmedAddedImpl{
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

func Test_confirmedAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	nilSubscribers := make(map[string]map[*ConfirmedAddedHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*ConfirmedAddedHandler]struct{})
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&confirmedAddedHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want map[*ConfirmedAddedHandler]struct{}
	}{
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &confirmedAddedImpl{
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
