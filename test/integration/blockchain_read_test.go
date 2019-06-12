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
const testUrl = "http://bcstage1.xpxsirius.io:3000"
const networkType = sdk.PublicTest
const privateKey = "D54AC0CB0FF50FB44233782B3A6B5FDE2F1C83B9AE2F1352119F93713F3AB923"

var defaultAccount, _ = sdk.NewAccountFromPrivateKey(privateKey, networkType)

func TestMosaicService_GetMosaicsFromNamespaceExt(t *testing.T) {
	cfg, _ := sdk.NewConfig([]string{testUrl}, networkType, sdk.WebsocketReconnectionDefaultTimeout)
	ctx := context.TODO()

	serv := sdk.NewClient(nil, cfg)
	h, err := serv.Blockchain.GetBlockchainHeight(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for i := uint64(1); i < h.Uint64() && i <= iter; i++ {

		h := sdk.NewHeight(i)
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

				t.Logf("%+v", mscInfo)
			case sdk.MosaicSupplyChange:
				tran := val.(*sdk.MosaicSupplyChangeTransaction)

				if tran.MosaicId == nil {
					t.Logf("empty MosaicId")
					t.Log(tran)
					continue
				}
				mscInfo, err := serv.Mosaic.GetMosaicInfo(ctx, tran.MosaicId)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%+v", mscInfo)
			case sdk.Transfer:
				tran := val.(*sdk.TransferTransaction)
				if tran.Mosaics == nil {
					t.Logf("empty Mosaics")
					t.Log(tran)
					continue
				}
				mosaicIDs := make([]*sdk.MosaicId, len(tran.Mosaics))
				for _, val := range tran.Mosaics {
					mosaicIDs = append(mosaicIDs, val.MosaicId)
				}
				mscInfoArr, err := serv.Mosaic.GetMosaicInfos(ctx, mosaicIDs)
				if err != nil {
					t.Fatal(err)
				}

				for _, mscInfo := range mscInfoArr {
					t.Logf("%+v", mscInfo)
				}
			case sdk.RegisterNamespace:
				tran := val.(*sdk.RegisterNamespaceTransaction)
				nsInfo, err := serv.Namespace.GetNamespaceInfo(ctx, tran.NamespaceId)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("%#v", nsInfo)
			default:
				t.Log(val)
			}
		}

	}
}
