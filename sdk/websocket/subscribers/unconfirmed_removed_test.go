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
		handle   *sdk.CompoundChannelHandle
		handlers []UnconfirmedRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	subscribers := make(map[string][]*UnconfirmedRemovedHandler)
	subscribers[handle.String()] = make([]*UnconfirmedRemovedHandler, 0)

	subscribersNilHandlers := make(map[string][]*UnconfirmedRemovedHandler)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	tests := []struct {
		name    string
		e       *unconfirmedRemovedImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers arg",
			e: &unconfirmedRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []UnconfirmedRemovedHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &unconfirmedRemovedImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1Ptr,
					unconfirmedRemovedHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []UnconfirmedRemovedHandler{
					unconfirmedRemovedHandlerFunc1Ptr,
					unconfirmedRemovedHandlerFunc2Ptr,
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

func Test_unconfirmedRemovedImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*UnconfirmedRemovedHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	emptySubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	emptySubscribers[handle.String()] = make([]*UnconfirmedRemovedHandler, 0)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedRemovedHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*UnconfirmedRemovedHandler)
	oneSubsctiberStorage[handle.String()] = make([]*UnconfirmedRemovedHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &unconfirmedRemovedHandlerFunc1Ptr

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
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*UnconfirmedRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &unconfirmedRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*UnconfirmedRemovedHandler{
					&unconfirmedRemovedHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*UnconfirmedRemovedHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &unconfirmedRemovedImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
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
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_unconfirmedRemovedImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)

	emptySubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	emptySubscribers[handle.String()] = make([]*UnconfirmedRemovedHandler, 0)

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &unconfirmedRemovedImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
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

func Test_unconfirmedRemovedImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	unconfirmedRemovedHandlerFunc1Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc1)
	unconfirmedRemovedHandlerFunc2Ptr := UnconfirmedRemovedHandler(unconfirmedRemovedHandlerFunc2)

	nilSubscribers := make(map[string][]*UnconfirmedRemovedHandler)
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*UnconfirmedRemovedHandler)
	hasSubscribersStorage[handle.String()] = make([]*UnconfirmedRemovedHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &unconfirmedRemovedHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][1] = &unconfirmedRemovedHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *unconfirmedRemovedImpl
		args args
		want []*UnconfirmedRemovedHandler
	}{
		{
			name: "success",
			e: &unconfirmedRemovedImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &unconfirmedRemovedImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *unconfirmedRemovedSubscription),
				removeSubscriberCh: make(chan *unconfirmedRemovedSubscription),
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
