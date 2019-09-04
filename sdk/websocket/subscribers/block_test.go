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
	tests := []struct {
		name    string
		s       *blockSubscriberImpl
		args    args
		wantErr bool
	}{
		{
			name: "empty handlers",
			s: &blockSubscriberImpl{
				handlers: nil,
			},
			args: args{
				handlers: []BlockHandler{},
			},
			wantErr: false,
		},
		{
			name: "success",
			s: &blockSubscriberImpl{
				handlers: map[*BlockHandler]struct{}{},
			},
			args: args{
				handlers: []BlockHandler{
					blockHandlerFunc1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.AddHandlers(tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, len(tt.s.handlers), len(tt.args.handlers))
		})
	}
}

func Test_blockSubscriberImpl_RemoveHandlers(t *testing.T) {
	type args struct {
		handlers []*BlockHandler
	}

	emptyStorage := make(map[*BlockHandler]struct{})

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	oneHandlerStorage := make(map[*BlockHandler]struct{})
	oneHandlerStorage[&blockHandlerPtr1] = struct{}{}

	twoHandlersStorage := make(map[*BlockHandler]struct{})
	twoHandlersStorage[&blockHandlerPtr1] = struct{}{}
	twoHandlersStorage[&blockHandlerPtr2] = struct{}{}

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
				handlers: emptyStorage,
			},
			args: args{
				handlers: []*BlockHandler{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "after removal there are no handlers left",
			s: &blockSubscriberImpl{
				handlers: oneHandlerStorage,
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
				handlers: twoHandlersStorage,
			},
			args: args{
				handlers: []*BlockHandler{
					&blockHandlerPtr1,
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.RemoveHandlers(tt.args.handlers...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_blockSubscriberImpl_HasHandlers(t *testing.T) {
	emptyStorage := make(map[*BlockHandler]struct{})

	twoHandlersStorage := make(map[*BlockHandler]struct{})

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	twoHandlersStorage[&blockHandlerPtr1] = struct{}{}
	twoHandlersStorage[&blockHandlerPtr2] = struct{}{}

	tests := []struct {
		name string
		s    *blockSubscriberImpl
		want bool
	}{
		{
			name: "has handlers",
			s: &blockSubscriberImpl{
				handlers: twoHandlersStorage,
			},
			want: true,
		},
		{
			name: "empty handlers",
			s: &blockSubscriberImpl{
				handlers: emptyStorage,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.HasHandlers()
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_blockSubscriberImpl_GetHandlers(t *testing.T) {
	twoHandlersStorage := make(map[*BlockHandler]struct{})

	blockHandlerPtr1 := BlockHandler(blockHandlerFunc1)
	blockHandlerPtr2 := BlockHandler(blockHandlerFunc2)

	twoHandlersStorage[&blockHandlerPtr1] = struct{}{}
	twoHandlersStorage[&blockHandlerPtr2] = struct{}{}

	tests := []struct {
		name string
		s    *blockSubscriberImpl
		want map[*BlockHandler]struct{}
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
			got := tt.s.GetHandlers()
			assert.Equal(t, got, tt.want)
		})
	}
}
