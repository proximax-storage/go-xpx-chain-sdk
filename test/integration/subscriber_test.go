// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	baseUrl = "http://localhost:3000"
)

func TestAddConfirmedAddedHandlers(t *testing.T) {
	wg := sync.WaitGroup{}
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	cfg, err := sdk.NewConfig(ctx, []string{baseUrl})
	assert.Nil(t, err)

	webSock, err := websocket.NewClient(ctx, cfg)
	assert.Nil(t, err)

	wg.Add(1)

	go webSock.Listen()

	err = webSock.AddConfirmedAddedHandlers(testAccount.Address, func(transaction sdk.Transaction) bool {
		wg.Done()
		return true
	})
	assert.Nil(t, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewAccountPropertiesEntityTypeTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.BlockTransaction,
			[]*sdk.AccountPropertiesEntityTypeModification{
				{
					sdk.AddProperty,
					sdk.LinkAccount,
				},
			},
		)
	}, testAccount)
	assert.Nil(t, result.error)
	wg.Wait()
}

func TestAddPartialAddedHandlers(t *testing.T) {
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount()
	assert.Nil(t, err)
	acc2, err := client.NewAccount()
	assert.Nil(t, err)

	cfg, err := sdk.NewConfig(ctx, []string{baseUrl})
	assert.Nil(t, err)

	webSock, err := websocket.NewClient(ctx, cfg)
	assert.Nil(t, err)

	wg.Add(1)

	go webSock.Listen()

	err = webSock.AddPartialAddedHandlers(acc1.Address, func(transaction *sdk.AggregateTransaction) bool {
		wg.Done()
		return false
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(multisigAccount)

	multTxs, err := client.NewModifyMultisigAccountTransaction(
		sdk.NewDeadline(time.Hour),
		2,
		1,
		[]*sdk.MultisigCosignatoryModification{
			{
				sdk.Add,
				acc1.PublicAccount,
			},
			{
				sdk.Add,
				acc2.PublicAccount,
			},
		},
	)
	assert.Nil(t, err)
	multTxs.ToAggregate(multisigAccount.PublicAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount.PublicAccount)

	result := sendAggregateTransaction(t, func() (*sdk.AggregateTransaction, error) {
		return client.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}

func TestAddCosignatureHandlers(t *testing.T) {
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount()
	assert.Nil(t, err)
	acc2, err := client.NewAccount()
	assert.Nil(t, err)

	cfg, err := sdk.NewConfig(ctx, []string{baseUrl})
	assert.Nil(t, err)

	webSock, err := websocket.NewClient(ctx, cfg)
	assert.Nil(t, err)

	wg.Add(1)

	go webSock.Listen()

	err = webSock.AddCosignatureHandlers(acc2.Address, func(info *sdk.SignerInfo) bool {
		wg.Done()
		return true
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(multisigAccount)

	multTxs, err := client.NewModifyMultisigAccountTransaction(
		sdk.NewDeadline(time.Hour),
		2,
		1,
		[]*sdk.MultisigCosignatoryModification{
			{
				sdk.Add,
				acc1.PublicAccount,
			},
			{
				sdk.Add,
				acc2.PublicAccount,
			},
		},
	)
	assert.Nil(t, err)
	multTxs.ToAggregate(multisigAccount.PublicAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount.PublicAccount)

	result := sendAggregateTransaction(t, func() (*sdk.AggregateTransaction, error) {
		return client.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}
