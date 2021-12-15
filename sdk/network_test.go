// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
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

	blockchainUpgrade = `{
  		"blockchainUpgrade": {
    		"height": [
      			206,
      			0
    		],
			"blockChainVersion": [
      			0,
      			4
    		]
		}
	}`

	networkConfigJson = `{
  		"networkConfig": {
    		"height": [
      			144,
      			0
    		],
			"networkConfig": "[network]\n\nidentifier = mijin-test\npublicKey = B4F12E7C9F6946091E2CB8B6D3A12B50D17CCBBF646386EA27CE2946A7423DCF\ngenerationHash = 86258172F90639811F2ABD055747D1E11B55A64B68AED2CEA9A34FBD6C0BE790\n\n",
			"supportedEntityVersions": "{\n    \"entities\": [\n\t\t{\n\t\t\t\"name\": \"Block\",\n\t\t\t\"type\": \"33091\",\n\t\t\t\"supportedVersions\": [3]\n\t\t}]\n}"
		}
	}`
)

var (
	networkVersion = &NetworkVersion{
		StartedHeight:     Height(206),
		BlockChainVersion: NewBlockChainVersion(0, 4, 0, 0),
	}

	networkConfig = &BlockchainConfig{
		StartedHeight: Height(144),
		NetworkConfig: &NetworkConfig{
			Sections: map[string]*ConfigBag{
				"network": {
					Name:    "network",
					Comment: "",
					Fields: map[string]*Field{
						"identifier": {
							Comment: "\n",
							Key:     "identifier",
							Value:   "mijin-test",
							Index:   0,
						},
						"publicKey": {
							Comment: "",
							Key:     "publicKey",
							Value:   "B4F12E7C9F6946091E2CB8B6D3A12B50D17CCBBF646386EA27CE2946A7423DCF",
							Index:   1,
						},
						"generationHash": {
							Comment: "",
							Key:     "generationHash",
							Value:   "86258172F90639811F2ABD055747D1E11B55A64B68AED2CEA9A34FBD6C0BE790",
							Index:   2,
						},
					},
				},
			},
		},
		SupportedEntityVersions: &SupportedEntities{
			Entities: map[EntityType]*Entity{
				Block: {
					Name:              "Block",
					Type:              Block,
					SupportedVersions: []EntityVersion{3},
				},
			},
		},
	}
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
		assert.Equal(t, netType, Mijin)
	})

	t.Run("mijinTest", func(t *testing.T) {
		mockServ := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: mijinTestRoute,
		})

		defer mockServ.Close()

		netType, err := mockServ.getPublicTestClientUnsafe().Network.GetNetworkType(ctx)

		assert.Nilf(t, err, "NetworkService.GetNetworkType should return error")
		assert.Equal(t, netType, MijinTest)
	})

	t.Run("notSupported", func(t *testing.T) {
		mockServ := newSdkMockWithRouter(&mock.Router{
			Path:     networkRoute,
			RespBody: notSupportedRoute,
		})

		defer mockServ.Close()

		netType, err := mockServ.getPublicTestClientUnsafe().Network.GetNetworkType(ctx)

		assert.NotNil(t, err, "NetworkService.GetNetworkType should return error")
		assert.Equal(t, netType, NotSupportedNet)
	})
}

func TestExtractNetworkType(t *testing.T) {
	i := int64(-1879048189)

	nt := ExtractNetworkType(i)

	assert.Equal(t, MijinTest, nt)
	i = int64(2415919106)

	nt = ExtractNetworkType(i)

	assert.Equal(t, MijinTest, nt)
}

func TestNetworkService_GetNetworkVersionAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(upgradeRoute, Height(210)),
		RespBody: blockchainUpgrade,
	})

	nVersion, err := mockServer.getPublicTestClientUnsafe().Network.GetNetworkVersionAtHeight(ctx, Height(210))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkVersion, nVersion)
}

func TestNetworkService_GetNetworkConfigAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(configRoute, Height(150)),
		RespBody: networkConfigJson,
	})

	nConfig, err := mockServer.getPublicTestClientUnsafe().Network.GetNetworkConfigAtHeight(ctx, Height(150))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkConfig, nConfig)
}
