// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
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

	catapultUpgrade = `{
  		"catapultUpgrade": {
    		"height": [
      			206,
      			0
    		],
			"catapultVersion": [
      			0,
      			4
    		]
		}
	}`

	catapultConfig = `{
  		"catapultConfig": {
    		"height": [
      			144,
      			0
    		],
			"blockChainConfig": "Hello",
			"supportedEntityVersions": "world"
		}
	}`
)

var (
	networkVersion = &NetworkVersion{
		StartedHeight:   Height(206),
		CatapultVersion: NewCatapultVersion(0, 4, 0, 0),
	}

	networkConfig = &NetworkConfig{
		StartedHeight:           Height(144),
		BlockChainConfig:        "Hello",
		SupportedEntityVersions: "world",
	}
)

var networkClient = mockServer.getPublicTestClientUnsafe().Network

func TestNetworkService_GetNetworkType(t *testing.T) {
	t.Run("mijin", func(t *testing.T) {
		mockServ := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: mijinRoute,
		})

		defer mockServ.Close()

		netType, err := networkClient.GetNetworkType(ctx)

		assert.Nilf(t, err, "NetworkService.GetNetworkType returned error=%s", err)
		assert.Equal(t, netType, Mijin)
	})

	t.Run("mijinTest", func(t *testing.T) {
		mock := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: mijinTestRoute,
		})

		defer mock.Close()

		netType, err := networkClient.GetNetworkType(ctx)

		assert.Nilf(t, err, "NetworkService.GetNetworkType should return error")
		assert.Equal(t, netType, MijinTest)
	})

	t.Run("notSupported", func(t *testing.T) {
		mock := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: notSupportedRoute,
		})

		defer mock.Close()

		netType, err := networkClient.GetNetworkType(ctx)

		assert.NotNil(t, err, "NetworkService.GetNetworkType should return error")
		assert.Equal(t, netType, NotSupportedNet)
	})
}

func TestExtractNetworkType(t *testing.T) {
	i := uint64(36888)

	nt := ExtractNetworkType(i)

	assert.Equal(t, MijinTest, nt)
}

func TestNetworkService_GetNetworkVersionAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(upgradeRoute, Height(210)),
		RespBody: catapultUpgrade,
	})

	nVersion, err := networkClient.GetNetworkVersionAtHeight(ctx, Height(210))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkVersion, nVersion)
}

func TestNetworkService_GetNetworkConfigAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(configRoute, Height(150)),
		RespBody: catapultConfig,
	})

	nConfig, err := networkClient.GetNetworkConfigAtHeight(ctx, Height(150))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkConfig, nConfig)
}
