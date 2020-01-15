package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var driveStateHandlerFunc1 = func(tx *sdk.DriveStateInfo) bool {
	return false
}

var driveStateHandlerFunc2 = func(tx *sdk.DriveStateInfo) bool {
	return false
}

func Test_driveStateImpl_AddHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []DriveStateHandler
	}
	driveStateHandlerFunc1Ptr := DriveStateHandler(driveStateHandlerFunc1)
	driveStateHandlerFunc2Ptr := DriveStateHandler(driveStateHandlerFunc2)
	address := &sdk.Address{}
	address.Address = "test-address"

	subscribers := make(map[string][]*DriveStateHandler)
	subscribers[address.Address] = make([]*DriveStateHandler, 0)

	driveStateNilHandlers := make(map[string][]*DriveStateHandler)

	tests := []struct {
		name    string
		e       *driveStateImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &driveStateImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address:  address,
				handlers: []DriveStateHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &driveStateImpl{
				subscribers:        driveStateNilHandlers,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
				handlers: []DriveStateHandler{
					driveStateHandlerFunc1Ptr,
					driveStateHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &driveStateImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
				handlers: []DriveStateHandler{
					driveStateHandlerFunc1Ptr,
					driveStateHandlerFunc2Ptr,
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

func Test_driveStateImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		address  *sdk.Address
		handlers []*DriveStateHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	emptySubscribers := make(map[string][]*DriveStateHandler)
	emptySubscribers[address.Address] = make([]*DriveStateHandler, 0)

	driveStateHandlerFunc1Ptr := DriveStateHandler(driveStateHandlerFunc1)
	driveStateHandlerFunc2Ptr := DriveStateHandler(driveStateHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*DriveStateHandler)
	hasSubscribersStorage[address.Address] = make([]*DriveStateHandler, 2)
	hasSubscribersStorage[address.Address][0] = &driveStateHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &driveStateHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*DriveStateHandler)
	oneSubsctiberStorage[address.Address] = make([]*DriveStateHandler, 1)
	oneSubsctiberStorage[address.Address][0] = &driveStateHandlerFunc1Ptr

	tests := []struct {
		name    string
		e       *driveStateImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &driveStateImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address:  address,
				handlers: []*DriveStateHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for address",
			e: &driveStateImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
				handlers: []*DriveStateHandler{
					&driveStateHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &driveStateImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address:  address,
				handlers: []*DriveStateHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &driveStateImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
				handlers: []*DriveStateHandler{
					&driveStateHandlerFunc1Ptr,
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

func Test_driveStateImpl_HasHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	driveStateHandlerFunc1Ptr := DriveStateHandler(driveStateHandlerFunc1)
	driveStateHandlerFunc2Ptr := DriveStateHandler(driveStateHandlerFunc2)

	emptySubscribers := make(map[string][]*DriveStateHandler)
	emptySubscribers[address.Address] = make([]*DriveStateHandler, 0)

	hasSubscribersStorage := make(map[string][]*DriveStateHandler)
	hasSubscribersStorage[address.Address] = make([]*DriveStateHandler, 2)
	hasSubscribersStorage[address.Address][0] = &driveStateHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &driveStateHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *driveStateImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &driveStateImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
			},
			want: true,
		},
		{
			name: "false result",
			e: &driveStateImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
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

func Test_driveStateImpl_GetHandlers(t *testing.T) {
	type args struct {
		address *sdk.Address
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	driveStateHandlerFunc1Ptr := DriveStateHandler(driveStateHandlerFunc1)
	driveStateHandlerFunc2Ptr := DriveStateHandler(driveStateHandlerFunc2)

	nilSubscribers := make(map[string][]*DriveStateHandler)
	nilSubscribers[address.Address] = nil

	hasSubscribersStorage := make(map[string][]*DriveStateHandler)
	hasSubscribersStorage[address.Address] = make([]*DriveStateHandler, 2)
	hasSubscribersStorage[address.Address][0] = &driveStateHandlerFunc1Ptr
	hasSubscribersStorage[address.Address][1] = &driveStateHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *driveStateImpl
		args args
		want []*DriveStateHandler
	}{
		{
			name: "success",
			e: &driveStateImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
			},
			args: args{
				address: address,
			},
			want: hasSubscribersStorage[address.Address],
		},
		{
			name: "nil result",
			e: &driveStateImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *driveStateSubscription),
				removeSubscriberCh: make(chan *driveStateSubscription),
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
