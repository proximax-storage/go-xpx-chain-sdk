// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"context"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	cfg, err := sdk.NewConfigFromRemote([]string{testUrl})
	if err != nil {
		panic(err)
	}

	ctx = context.Background()
	client = sdk.NewClient(nil, cfg)

	wsc, err = websocket.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	defaultAccount, err = client.NewAccountFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
}

func TestAddressService_GetAccountNames(t *testing.T) {

	names, err := client.Account.GetAccountNames(
		ctx,
		&sdk.Address{sdk.Mijin, "SDRDGFTDLLCB67D4HPGIMIHPNSRYRJRT7DOBGWZY"},
		&sdk.Address{sdk.Mijin, "SBCPGZ3S2SCC3YHBBTYDCUZV4ZZEPHM2KGCP4QXX"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(names), 2)

}
