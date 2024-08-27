// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package websocket

import (
	"context"
	"errors"
	"testing"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subs"
	"github.com/stretchr/testify/assert"
)

func TestSubscribersPool(t *testing.T) {
	t.Run("Notify_Success", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}
		payload := []byte("test payload")
		ctx := context.Background()

		mockPool.On("Notify", ctx, path, payload).Return(nil)

		err := mockPool.Notify(ctx, path, payload)
		assert.NoError(t, err)

		mockPool.AssertCalled(t, "Notify", ctx, path, payload)
	})

	t.Run("Notify_Failure", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}
		payload := []byte("test payload")
		ctx := context.Background()

		expectedErr := errors.New("notification error")

		mockPool.On("Notify", ctx, path, payload).Return(expectedErr)

		err := mockPool.Notify(ctx, path, payload)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockPool.AssertCalled(t, "Notify", ctx, path, payload)
	})

	t.Run("NewSubscription_Success", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}

		ch := make(<-chan string)

		mockPool.On("NewSubscription", path).Return(ch, 1)

		subChan, id := mockPool.NewSubscription(path)
		assert.Equal(t, ch, subChan, "Expected channel to match the mock's channel")
		assert.Equal(t, 1, id, "Expected subscription ID to be 1")

		mockPool.AssertCalled(t, "NewSubscription", path)
	})

	t.Run("NewSubscription_Failure", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}

		mockPool.On("NewSubscription", path).Return(nil, 0)

		subChan, id := mockPool.NewSubscription(path)
		assert.Nil(t, subChan)
		assert.Equal(t, 0, id)

		mockPool.AssertCalled(t, "NewSubscription", path)
	})

	t.Run("CloseSubscription", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}

		mockPool.On("CloseSubscription", path, 1).Return()

		mockPool.CloseSubscription(path, 1)

		mockPool.AssertCalled(t, "CloseSubscription", path, 1)
	})

	t.Run("GetPaths_Success", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])

		expectedPaths := []string{"path1", "path2"}
		mockPool.On("GetPaths").Return(expectedPaths)

		paths := mockPool.GetPaths()
		assert.Equal(t, expectedPaths, paths)

		mockPool.AssertCalled(t, "GetPaths")
	})

	t.Run("GetPaths_Failure", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])

		mockPool.On("GetPaths").Return(nil)

		paths := mockPool.GetPaths()
		assert.Nil(t, paths)

		mockPool.AssertCalled(t, "GetPaths")
	})

	t.Run("HasSubscriptions_Success", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}

		mockPool.On("HasSubscriptions", path).Return(true)

		hasSubs := mockPool.HasSubscriptions(path)
		assert.True(t, hasSubs)

		mockPool.AssertCalled(t, "HasSubscriptions", path)
	})

	t.Run("HasSubscriptions_Failure", func(t *testing.T) {
		mockPool := new(MockSubscribersPool[string])
		path := &subs.Path{}

		mockPool.On("HasSubscriptions", path).Return(false)

		hasSubs := mockPool.HasSubscriptions(path)
		assert.False(t, hasSubs)

		mockPool.AssertCalled(t, "HasSubscriptions", path)
	})
}
