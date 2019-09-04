// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"testing"

	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestAddressService_GetAccountNames(t *testing.T) {

	networkType := sdk.MijinTest

	addresses := []*sdk.Address{
		{
			networkType,
			"SCWXLOABHP4FT2LWTT3Z6GDCHLLMUIKKFRBE2O3S",
		},
		{
			networkType,
			"SBKDKHFIRM72EAVDT6TI426CKUCP5DQIJV73XB5X",
		},
	}

	names, err := client.Account.GetAccountNames(
		ctx,
		addresses...)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(names), len(addresses))

	for i, accNames := range names {
		tests.ValidateStringers(t, addresses[i], accNames.Address)
	}
}
