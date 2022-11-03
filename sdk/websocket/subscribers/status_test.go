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
		handle   *sdk.CompoundChannelHandle
		handlers []StatusHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)
	subscribers := make(map[string][]*StatusHandler)
	subscribers[handle.String()] = make([]*StatusHandler, 0)

	subscribersNilHandlers := make(map[string][]*StatusHandler)

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
				handle:   handle,
				handlers: []StatusHandler{},
			},
			wantErr: false,
		},
		{
			name: "nil handlers",
			e: &statusImpl{
				subscribers:        subscribersNilHandlers,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
				handlers: []StatusHandler{
					statusHandlerFunc1Ptr,
					statusHandlerFunc2Ptr,
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			e: &statusImpl{
				subscribers:        subscribers,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
				handlers: []StatusHandler{
					statusHandlerFunc1Ptr,
					statusHandlerFunc2Ptr,
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

func Test_statusImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handle   *sdk.CompoundChannelHandle
		handlers []*StatusHandler
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	emptySubscribers := make(map[string][]*StatusHandler)
	emptySubscribers[handle.String()] = make([]*StatusHandler, 0)

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	hasSubscribersStorage := make(map[string][]*StatusHandler)
	hasSubscribersStorage[handle.String()] = make([]*StatusHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc2Ptr

	oneSubsctiberStorage := make(map[string][]*StatusHandler)
	oneSubsctiberStorage[handle.String()] = make([]*StatusHandler, 1)
	oneSubsctiberStorage[handle.String()][0] = &statusHandlerFunc1Ptr

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
				subscribers:        nil,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*StatusHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty handlers storage for handle",
			e: &statusImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
				handlers: []*StatusHandler{
					&statusHandlerFunc1Ptr,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return false result",
			e: &statusImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle:   handle,
				handlers: []*StatusHandler{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "success return true result",
			e: &statusImpl{
				subscribers:        oneSubsctiberStorage,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
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
			go tt.e.handleNewSubscription()
			got := tt.e.RemoveHandlers(tt.args.handle, tt.args.handlers...)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_statusImpl_HasHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	emptySubscribers := make(map[string][]*StatusHandler)
	emptySubscribers[handle.String()] = make([]*StatusHandler, 0)

	hasSubscribersStorage := make(map[string][]*StatusHandler)
	hasSubscribersStorage[handle.String()] = make([]*StatusHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *statusImpl
		args args
		want bool
	}{
		{
			name: "true result",
			e: &statusImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
			},
			want: true,
		},
		{
			name: "false result",
			e: &statusImpl{
				subscribers:        emptySubscribers,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
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

func Test_statusImpl_GetHandlers(t *testing.T) {
	type args struct {
		handle *sdk.CompoundChannelHandle
	}

	address := &sdk.Address{}
	address.Address = "test-address"
	handle := sdk.NewCompoundChannelHandleFromAddress(address)

	statusHandlerFunc1Ptr := StatusHandler(statusHandlerFunc1)
	statusHandlerFunc2Ptr := StatusHandler(statusHandlerFunc2)

	nilSubscribers := make(map[string][]*StatusHandler)
	nilSubscribers[handle.String()] = nil

	hasSubscribersStorage := make(map[string][]*StatusHandler)
	hasSubscribersStorage[handle.String()] = make([]*StatusHandler, 2)
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc1Ptr
	hasSubscribersStorage[handle.String()][0] = &statusHandlerFunc2Ptr

	tests := []struct {
		name string
		e    *statusImpl
		args args
		want []*StatusHandler
	}{
		{
			name: "success",
			e: &statusImpl{
				subscribers:        hasSubscribersStorage,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
			},
			args: args{
				handle: handle,
			},
			want: hasSubscribersStorage[handle.String()],
		},
		{
			name: "nil result",
			e: &statusImpl{
				subscribers:        nil,
				newSubscriberCh:    make(chan *statusSubscription),
				removeSubscriberCh: make(chan *statusSubscription),
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
			if got := tt.e.GetHandlers(tt.args.handle); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statusImpl.GetHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}
