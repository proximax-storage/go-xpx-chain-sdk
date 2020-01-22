package subscribers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

var blockHandlerFunc1 = func(blockInfo *sdk.BlockInfo) bool {
	return false
}

var blockHandlerFunc2 = func(blockInfo *sdk.BlockInfo) bool {
	return false
}

func Test_blockSubscriberImpl_AddHandlers(t *testing.T) {
	type args struct {
		handlers []BlockHandler
	}
	bh := BlockHandler(blockHandlerFunc1)
	tests := []struct {
		name    string
		s       *blockSubscriberImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers",
			s: &blockSubscriberImpl{
				handlers:           nil,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			args: args{
				handlers: []BlockHandler{},
			},
			wantErr: false,
		},
		{
			name: "success",
			s: &blockSubscriberImpl{
				handlers:           []*BlockHandler{},
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			args: args{
				handlers: []BlockHandler{
					bh,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.s.handleNewSubscription()
			err := tt.s.AddHandlers(tt.args.handlers...)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, len(tt.args.handlers), len(tt.s.handlers))
		})
	}
}

func Test_blockSubscriberImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handlers []*BlockHandler
	}

	emptyStorage := make([]*BlockHandler, 0)

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	oneHandlerStorage := make([]*BlockHandler, 1)
	oneHandlerStorage[0] = &blockHandlerPtr1

	twoHandlersStorage := make([]*BlockHandler, 2)
	twoHandlersStorage[0] = &blockHandlerPtr1
	twoHandlersStorage[1] = &blockHandlerPtr2

	tests := []struct {
		name    string
		s       *blockSubscriberImpl
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "handler not found",
			s: &blockSubscriberImpl{
				handlers:           emptyStorage,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			args: args{
				handlers: []*BlockHandler{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "after removal there are no handlers left",
			s: &blockSubscriberImpl{
				handlers:           oneHandlerStorage,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			args: args{
				handlers: []*BlockHandler{
					&blockHandlerPtr1,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "after removal handlers left",
			s: &blockSubscriberImpl{
				handlers:           twoHandlersStorage,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			args: args{
				handlers: []*BlockHandler{
					&blockHandlerPtr1,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.s.handleNewSubscription()
			got := tt.s.RemoveHandlers(tt.args.handlers...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_blockSubscriberImpl_HasHandlers(t *testing.T) {
	emptyStorage := make([]*BlockHandler, 0)

	twoHandlersStorage := make([]*BlockHandler, 2)

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	twoHandlersStorage[0] = &blockHandlerPtr1
	twoHandlersStorage[1] = &blockHandlerPtr2

	tests := []struct {
		name string
		s    *blockSubscriberImpl
		want bool
	}{
		{
			name: "has handlers",
			s: &blockSubscriberImpl{
				handlers:           twoHandlersStorage,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			want: true,
		},
		{
			name: "empty handlers",
			s: &blockSubscriberImpl{
				handlers:           emptyStorage,
				newSubscriberCh:    make(chan *blockSubscription),
				removeSubscriberCh: make(chan *blockSubscription),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.s.handleNewSubscription()
			got := tt.s.HasHandlers()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_blockSubscriberImpl_GetHandlers(t *testing.T) {
	twoHandlersStorage := make([]*BlockHandler, 2)

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	twoHandlersStorage[0] = &blockHandlerPtr1
	twoHandlersStorage[1] = &blockHandlerPtr2

	tests := []struct {
		name string
		s    *blockSubscriberImpl
		want []*BlockHandler
	}{
		{
			name: "success",
			s: &blockSubscriberImpl{
				handlers: twoHandlersStorage,
			},
			want: twoHandlersStorage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.s.handleNewSubscription()
			got := tt.s.GetHandlers()
			assert.Equal(t, tt.want, got)
		})
	}
}
