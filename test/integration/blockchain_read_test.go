// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"context"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"testing"
	"time"
)

//
//const testUrl = "http://bcdev1.xpxsirius.io:3000"
//const privateKey = "451EA3199FE0520FB10B7F89D3A34BAF7E5C3B16FDFE2BC11A5CAC95CDB29ED6"

const testUrl = "http://127.0.0.1:3000"
const privateKey = "A31411BC4BA7267147DBBEDC034FA3D3C0B7294A0784507539C3BCE4EF70615A"
const nemesisPrivateKey = "C06B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716"

const timeout = 2 * time.Minute

var cfg *sdk.Config
var ctx context.Context
var client *sdk.Client
var wsc websocket.CatapultClient
var defaultAccount *sdk.Account
var nemesisAccount *sdk.Account

const iter = 1000

func init() {
	ctx = context.Background()

	cfg, err := sdk.NewConfig(ctx, []string{testUrl})
	if err != nil {
		panic(err)
	}
	cfg.FeeCalculationStrategy = 0

	client = sdk.NewClient(nil, cfg)

	wsc, err = websocket.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	defaultAccount, err = client.NewAccountFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	nemesisAccount, err = client.NewAccountFromPrivateKey(nemesisPrivateKey)
	if err != nil {
		panic(err)
	}
}

func TestMosaicService_GetMosaicsFromNamespaceExt(t *testing.T) {
	h, err := client.Blockchain.GetBlockchainHeight(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for i := sdk.Height(1); i < h && i <= iter; i++ {
		h := i
		trans, err := client.Blockchain.GetBlockTransactions(ctx, h)
		if err != nil {
			t.Fatal(err)
			continue
		}

		if len(trans) == 0 {
			t.Logf("%d block, empty transactions", h)
		}

		for j, val := range trans {
			if val == nil {
				t.Logf("empty trans #%d", j)
				continue
			}
			//t.Log(val.String())
			switch val.GetAbstractTransaction().Type {
			case sdk.MosaicDefinition:
				tran := val.(*sdk.MosaicDefinitionTransaction)

				if tran.MosaicId == nil {
					t.Logf("empty MosaicId")
					t.Log(tran)
					continue
				}
				mscInfo, err := client.Mosaic.GetMosaicInfo(ctx, tran.MosaicId)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%s", mscInfo)
			case sdk.MosaicSupplyChange:
				tran := val.(*sdk.MosaicSupplyChangeTransaction)

				if tran.AssetId == nil {
					t.Logf("empty MosaicId")
					t.Log(tran)
					continue
				}
				mscInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, tran.AssetId)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%s", mscInfo)
			case sdk.Transfer:
				tran := val.(*sdk.TransferTransaction)

				if tran.Mosaics == nil {
					t.Logf("empty Mosaics")
					t.Log(tran)
					continue
				}

				assetIds := make([]sdk.AssetId, len(tran.Mosaics))
				for i, val := range tran.Mosaics {
					assetIds[i] = val.AssetId
				}

				if len(assetIds) > 0 {
					mscInfoArr, err := client.Resolve.GetMosaicInfosByAssetIds(ctx, assetIds...)
					if err != nil {
						t.Fatal(err)
					}

					for _, mscInfo := range mscInfoArr {
						t.Logf("%s", mscInfo)
					}
				}
			case sdk.RegisterNamespace:
				tran := val.(*sdk.RegisterNamespaceTransaction)
				nsInfo, err := client.Namespace.GetNamespaceInfo(ctx, tran.NamespaceId)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%s", nsInfo)
			default:
				t.Log(val)
			}
		}

	}
}
