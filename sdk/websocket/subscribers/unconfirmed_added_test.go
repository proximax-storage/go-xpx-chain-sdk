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
		handle   *sdk.CompoundChannelHandle
		handlers []UnconfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)

	subscribers := make(map[string][]*UnconfirmedAddedHandler)
	subscribers[handle.String()] = make([]*UnconfirmedAddedHandler, 0)

	subscribersNilHandlers := make(map[string][]*UnconfirmedAddedHandler)

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
				handle:   handle,
				handlers: []UnconfirmedAddedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &unconfirmedAddedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []UnconfirmedAddedHandler{
					unconfirmedAddedHandlerFunc1Ptr,
					unconfirmedAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &unconfirmedAddedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []UnconfirmedAddedHandler{
					unconfirmedAddedHandlerFunc1Ptr,
					unconfirmedAddedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.e.handleNewSubscription()
			err := tt.e.AddHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func Test_unconfirmedAddedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*UnconfirmedAddedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	emptySubscribers := make(map[string][]*UnconfirmedAddedHandler)
	emptySubscribers[handle.String()] = make([]*UnconfirmedAddedHandler, 0)

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)

	hasSubscribersStorage := make(map[string][]*UnconfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedAddedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*UnconfirmedAddedHandler)
	oneSubsctiberStorage[handle.String()] = make([]*UnconfirmedAddedHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &unconfirmedAddedHandlerFunc1Ptr

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
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*UnconfirmedAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &unconfirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*UnconfirmedAddedHandler{
					&unconfirmedAddedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &unconfirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*UnconfirmedAddedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &unconfirmedAddedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
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
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_unconfirmedAddedImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc2)

	emptySubscribers := make(map[string][]*UnconfirmedAddedHandler)
	emptySubscribers[handle.String()] = make([]*UnconfirmedAddedHandler, 0)

	hasSubscribersStorage := make(map[string][]*UnconfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedAddedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &unconfirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &unconfirmedAddedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
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
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_unconfirmedAddedImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	unconfirmedAddedHandlerFunc1Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc1)
	unconfirmedAddedHandlerFunc2Ptr := UnconfirmedAddedHandler(unconfirmedAddedHandlerFunc2)

	nilSubscribers := make(map[string][]*UnconfirmedAddedHandler)
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*UnconfirmedAddedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedAddedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedAddedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedAddedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedAddedImpl
		args args
		want []*UnconfirmedAddedHandler
	}{
		{
			name: "success",
			e: &unconfirmedAddedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &unconfirmedAddedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedAddedSubscription),
				removeSubscriberCh: make(chan *unconfirmedAddedSubscription),
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
			assert.Equal(t, got, tt.want)
		})
	}
}
