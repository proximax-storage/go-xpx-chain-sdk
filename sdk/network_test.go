// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/proximax-storage/proximax-utils-go/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	mijinRoute = `{
  			"name": "mijin",
  			"description": "catapult development network"
  	}`
	mijinTestRoute = `{
  			"name": "mijinTest",
  			"description": "catapult development network"
  	}`
	notSupportedRoute = `{
			"name": "",
			"description": "catapult development network"
	}`
)

func TestNetworkService_GetNetworkType(t *testing.T) {
	t.Run("mijin", func(t *testing.T) {
		mockServ := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: mijinRoute,
		})

		defer mockServ.Close()

		netType, err := mockServ.getPublicTestClientUnsafe().Network.GetNetworkType(ctx)

		assert.Nilf(t, err, "NetworkService.GetNetworkType returned error=%s", err)

		assert.Equal(t, netType, MijinTest)
	})

	t.Run("mijinTest", func(t *testing.T) {
		mockServ := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: mijinTestRoute,
		})

		defer mockServ.Close()

		netType, err := mockServ.getPublicTestClientUnsafe().Network.GetNetworkType(ctx)

		assert.Nilf(t, err, "NetworkService.GetNetworkType returned error=%s", err)

		assert.Equal(t, netType, MijinTest)
	})

	t.Run("NotSupportedNet", func(t *testing.T) {
		mock := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: notSupportedRoute,
		})

		defer mock.Close()

		netType, err := mock.getPublicTestClientUnsafe().Network.GetNetworkType(ctx)

		assert.NotNil(t, err, "NetworkService.GetNetworkType should return error")
		assert.Equal(t, netType, NotSupportedNet)
	})
}

func TestExtractNetworkType(t *testing.T) {
	i := uint64(36888)

	nt := ExtractNetworkType(i)

	assert.Equal(t, MijinTest, nt)
}
