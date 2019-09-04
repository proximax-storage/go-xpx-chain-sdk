package subscribers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var statusHandlerFunc1 = func(tx *sdk.StatusInfo) bool {
	return false
}

var statusHandlerFunc2 = func(tx *sdk.StatusInfo) bool {
	return false
}

func Test_statusImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []StatusHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string]map[*StatusHandler]struct{})
	subscribers[address.Address] = make(map[*StatusHandler]struct{})

	subscribersNilHandlers := make(map[string]map[*StatusHandler]struct{})

	tests := []struct {
		name    string
		e       *statusImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &statusImpl{
				subscribers: nil,
			},
			args: args{
				address:  address,
				handlers: []StatusHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &statusImpl{
				subscribers: subscribersNilHandlers,
			},
			args: args{
				address: address,
				handlers: []StatusHandler{
					statusHandlerFunc1,
					statusHandlerFunc2,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &statusImpl{
				subscribers: subscribers,
			},
			args: args{
				address: address,
				handlers: []StatusHandler{
					statusHandlerFunc1,
					statusHandlerFunc2,
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

func Test_statusImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*StatusHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string]map[*StatusHandler]struct{})
	emptySubscribers[address.Address] = make(map[*StatusHandler]struct{})

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	hasSubscribersStorage := make(map[string]map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address][&statusHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&statusHandlerFunc2Ptr] = struct{}{}

	oneSubsctiberStorage := make(map[string]map[*StatusHandler]struct{})
	oneSubsctiberStorage[address.Address] = make(map[*StatusHandler]struct{})
	oneSubsctiberStorage[address.Address][&statusHandlerFunc1Ptr] = struct{}{}

	tests := []struct {
		name    string
		e       *statusImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &statusImpl{
				subscribers: nil,
			},
			args: args{
				address:  address,
				handlers: []*StatusHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &statusImpl{
				subscribers: emptySubscribers,
			},
			args: args{
				address: address,
				handlers: []*StatusHandler{
					&statusHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success return false result",
			e: &statusImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
				handlers: []*StatusHandler{
					&statusHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &statusImpl{
				subscribers: oneSubsctiberStorage,
			},
			args: args{
				address: address,
				handlers: []*StatusHandler{
					&statusHandlerFunc1Ptr,
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

func Test_statusImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	emptySubscribers := make(map[string]map[*StatusHandler]struct{})
	emptySubscribers[address.Address] = make(map[*StatusHandler]struct{})

	hasSubscribersStorage := make(map[string]map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address][&statusHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&statusHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *statusImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &statusImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &statusImpl{
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

func Test_statusImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	nilSubscribers := make(map[string]map[*StatusHandler]struct{})
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string]map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address] = make(map[*StatusHandler]struct{})
	hasSubscribersStorage[address.Address][&statusHandlerFunc1Ptr] = struct{}{}
	hasSubscribersStorage[address.Address][&statusHandlerFunc2Ptr] = struct{}{}

	tests := []struct {
		name string
		e    *statusImpl
		args args
		want map[*StatusHandler]struct{}
	}{
		{
			name: "success",
			e: &statusImpl{
				subscribers: hasSubscribersStorage,
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &statusImpl{
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
				t.Errorf("statusImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
