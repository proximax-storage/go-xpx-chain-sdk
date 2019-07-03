// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"context"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"testing"
)

const iter = 1000
const testUrl = "http://127.0.0.1:3000"
const networkType = sdk.MijinTest
const privateKey = "FFCA7367B4EE8E14041091F438B7BC6C3206ADE2D312670B218C7516395C7D7F"
const nemesisPrivateKey = "C06B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716"

var defaultAccount, _ = sdk.NewAccountFromPrivateKey(privateKey, networkType)
var nemesisAccount, _ = sdk.NewAccountFromPrivateKey(nemesisPrivateKey, networkType)

func TestMosaicService_GetMosaicsFromNamespaceExt(t *testing.T) {
	cfg, _ := sdk.NewConfig([]string{testUrl}, networkType, sdk.WebsocketReconnectionDefaultTimeout)
	ctx := context.TODO()

	serv := sdk.NewClient(nil, cfg)
	h, err := serv.Blockchain.GetBlockchainHeight(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for i := sdk.Height(1); i < h && i <= iter; i++ {
		h := i
		trans, err := serv.Blockchain.GetBlockTransactions(ctx, h)
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
				mscInfo, err := serv.Mosaic.GetMosaicInfo(ctx, tran.MosaicId)
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
				mscInfo, err := serv.Resolve.GetMosaicInfoByAssetId(ctx, tran.AssetId)
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
					mscInfoArr, err := serv.Resolve.GetMosaicInfosByAssetIds(ctx, assetIds...)
					if err != nil {
						t.Fatal(err)
					}

					for _, mscInfo := range mscInfoArr {
						t.Logf("%s", mscInfo)
					}
				}
			case sdk.RegisterNamespace:
				tran := val.(*sdk.RegisterNamespaceTransaction)
				nsInfo, err := serv.Namespace.GetNamespaceInfo(ctx, tran.NamespaceId)
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
