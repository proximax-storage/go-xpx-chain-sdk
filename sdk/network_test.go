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
			"blockChainConfig": "[network]\n\nidentifier = mijin-test\npublicKey = B4F12E7C9F6946091E2CB8B6D3A12B50D17CCBBF646386EA27CE2946A7423DCF\ngenerationHash = 86258172F90639811F2ABD055747D1E11B55A64B68AED2CEA9A34FBD6C0BE790\n\n",
			"supportedEntityVersions": "{\n    \"entities\": [\n\t\t{\n\t\t\t\"name\": \"Block\",\n\t\t\t\"type\": \"33091\",\n\t\t\t\"supportedVersions\": [3]\n\t\t}]\n}"
		}
	}`
)

var (
	networkVersion = &NetworkVersion{
		StartedHeight:   Height(206),
		CatapultVersion: NewCatapultVersion(0, 4, 0, 0),
	}

	networkConfig = &NetworkConfig{
		StartedHeight: Height(144),
		BlockChainConfig: &BlockChainConfig{
			Sections: map[string]*ConfigBag{
				"network": &ConfigBag{
					Name:    "network",
					Comment: "",
					Fields: map[string]*Field{
						"identifier": &Field{
							Comment: "\n",
							Key:     "identifier",
							Value:   "mijin-test",
						},
						"publicKey": &Field{
							Comment: "",
							Key:     "publicKey",
							Value:   "B4F12E7C9F6946091E2CB8B6D3A12B50D17CCBBF646386EA27CE2946A7423DCF",
						},
						"generationHash": &Field{
							Comment: "",
							Key:     "generationHash",
							Value:   "86258172F90639811F2ABD055747D1E11B55A64B68AED2CEA9A34FBD6C0BE790",
						},
					},
				},
			},
		},
		SupportedEntityVersions: &SupportedEntities{
			Entities: map[EntityType]*Entity{
				Block: &Entity{
					Name:              "Block",
					Type:              Block,
					SupportedVersions: []TransactionVersion{3},
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
	i := uint64(36888)

	nt := ExtractNetworkType(i)

	assert.Equal(t, MijinTest, nt)
}

func TestNetworkService_GetNetworkVersionAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(upgradeRoute, Height(210)),
		RespBody: catapultUpgrade,
	})

	nVersion, err := mockServer.getPublicTestClientUnsafe().Network.GetNetworkVersionAtHeight(ctx, Height(210))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkVersion, nVersion)
}

func TestNetworkService_GetNetworkConfigAtHeight(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(configRoute, Height(150)),
		RespBody: catapultConfig,
	})

	nConfig, err := mockServer.getPublicTestClientUnsafe().Network.GetNetworkConfigAtHeight(ctx, Height(150))

	assert.Nil(t, err)
	tests.ValidateStringers(t, networkConfig, nConfig)
}
