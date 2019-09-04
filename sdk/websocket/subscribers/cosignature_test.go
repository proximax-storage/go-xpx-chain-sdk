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

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string]map[*CosignatureHandler]struct{})
	subscribers[address.Address] = make(map[*CosignatureHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*CosignatureHandler]struct{})

	tests := []struct {
		name    string
		e       *cosignatureImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &cosignatureImpl{
				subscribers: nil,
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
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []CosignatureHandler{
					cosignatureHandlerFunc1,
					cosignatureHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &cosignatureImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []CosignatureHandler{
					cosignatureHandlerFunc1,
					cosignatureHandlerFunc2,
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

func Test_cosignatureImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*CosignatureHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*CosignatureHandler]struct{})
	emptySubscribers[address.Address] = make(map[*CosignatureHandler]struct{})

	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*CosignatureHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*CosignatureHandler]struct{})
	oneSubsctiberStorage[address.Address][&cosignatureHandlerFunc1Ptr] = struct{}{}

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
				subscribers: nil,
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
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*CosignatureHandler{
					&cosignatureHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &cosignatureImpl{
				subscribers: hasSubscribersStorage,
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
			name: "success return true result",
			e: &cosignatureImpl{
				subscribers: oneSubsctiberStorage,
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
			got, err := tt.e.RemoveHandlers(tt.args.address, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
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

	emptySubscribers := make(map[string]map[*CosignatureHandler]struct{})
	emptySubscribers[address.Address] = make(map[*CosignatureHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *cosignatureImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &cosignatureImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &cosignatureImpl{
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

func Test_cosignatureImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	cosignatureHandlerFunc1Ptr := CosignatureHandler(cosignatureHandlerFunc1)
	cosignatureHandlerFunc2Ptr := CosignatureHandler(cosignatureHandlerFunc2)

	nilSubscribers := make(map[string]map[*CosignatureHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*CosignatureHandler]struct{})
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&cosignatureHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *cosignatureImpl
		args args
		want map[*CosignatureHandler]struct{}
	}{
		{
			name: "success",
			e: &cosignatureImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &cosignatureImpl{
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
