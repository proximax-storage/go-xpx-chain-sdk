package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var confirmedAddedHandlerFunc1 = func(tx sdk.Transaction) bool {
	return false
}

var confirmedAddedHandlerFunc2 = func(tx sdk.Transaction) bool {
	return false
}

func Test_confirmedAddedImpl_AddHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	subscribers := make(map[string][]*ConfirmedAddedHandler)
	subscribers[handle.String()] = make([]*ConfirmedAddedHandler, 0)

	subscribersNilHandlers := make(map[string][]*ConfirmedAddedHandler)

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	tests := []struct {
		name    string
		e       *confirmedAddedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &confirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []ConfirmedAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &confirmedAddedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []ConfirmedAddedHandler{
					confirmedAddedHandlerFunc1Ptr,
					confirmedAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []ConfirmedAddedHandler{
					confirmedAddedHandlerFunc1Ptr,
					confirmedAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			err := tt.e.AddHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_confirmedAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*ConfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	emptySubscribers := make(map[string][]*ConfirmedAddedHandler)
	emptySubscribers[handle.String()] = make([]*ConfirmedAddedHandler, 0)

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)
	confirmedAddedHandlerFunc3Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &confirmedAddedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*ConfirmedAddedHandler)
	oneSubsctiberStorage[handle.String()] = make([]*ConfirmedAddedHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &confirmedAddedHandlerFunc1Ptr

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
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*ConfirmedAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &confirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*ConfirmedAddedHandler{
					&confirmedAddedHandlerFunc3Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &confirmedAddedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
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
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_confirmedAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	handle := sdk.NewCompoundChannelHandleFromAddress(address)
	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	emptySubscribers := make(map[string][]*ConfirmedAddedHandler)
	emptySubscribers[handle.String()] = make([]*ConfirmedAddedHandler, 0)

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &confirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &confirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.HasHandlers(tt.args.handle)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_confirmedAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"

	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	confirmedAddedHandlerFunc1Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc1)
	confirmedAddedHandlerFunc2Ptr := ConfirmedAddedHandler(confirmedAddedHandlerFunc2)

	nilSubscribers := make(map[string][]*ConfirmedAddedHandler)
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*ConfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*ConfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &confirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &confirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *confirmedAddedImpl
		args args
		want []*ConfirmedAddedHandler
	}{
		{
			name: "success",
			e: &confirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &confirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *confirmedAddedSubscription),
				removeSubscriberCh: make(chan *confirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			got := tt.e.GetHandlers(tt.args.handle)
			assert.Equal(t, tt.want, got)
		})
	}
}
