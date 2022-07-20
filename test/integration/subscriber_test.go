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

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	baseUrl = "http://localhost:3000"
)

func TestAddConfirmedAddedHandlers(t *testing.T) {
	wg := sync.WaitGroup{}
	testAccount, err := client.NewAccount(ctx)
	assert.Nil(t, err)

	fmt.Println(testAccount)

	wg.Add(1)

	err = wsc.AddConfirmedAddedHandlers(testAccount.Address, func(transaction sdk.Transaction) bool {
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
	SkipIfEntityNotSupportedAtVersion(client, t, sdk.AggregateCompletedV1, sdk.AggregateCompletedV1Version)
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount(ctx)
	assert.Nil(t, err)
	acc2, err := client.NewAccount(ctx)
	assert.Nil(t, err)

	wg.Add(1)

	err = wsc.AddPartialAddedHandlers(acc1.Address, func(transaction sdk.Transaction) bool {
		wg.Done()
		return false
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount(ctx)
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
	multTxs.ToAggregate(multisigAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount)

	result := sendAggregateTransactionV1(t, func() (*sdk.AggregateTransactionV1, error) {
		return client.NewBondedAggregateV1Transaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}

func TestAddPartialAddedHandlersV2(t *testing.T) {
	SkipIfEntityNotSupportedAtVersion(client, t, sdk.AggregateCompletedV2, sdk.AggregateCompletedV2Version)
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount(ctx)
	assert.Nil(t, err)
	acc2, err := client.NewAccount(ctx)
	assert.Nil(t, err)

	wg.Add(1)

	err = wsc.AddPartialAddedHandlers(acc1.Address, func(transaction sdk.Transaction) bool {
		wg.Done()
		return false
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount(ctx)
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
	multTxs.ToAggregate(multisigAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount)

	result := sendAggregateTransactionV2(t, func() (*sdk.AggregateTransactionV2, error) {
		return client.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}

func TestAddCosignatureHandlersV1(t *testing.T) {
	SkipIfEntityNotSupportedAtVersion(client, t, sdk.AggregateCompletedV1, sdk.AggregateCompletedV1Version)
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount(ctx)
	assert.Nil(t, err)
	acc2, err := client.NewAccount(ctx)
	assert.Nil(t, err)

	wg.Add(1)

	err = wsc.AddCosignatureHandlers(acc2.Address, func(info *sdk.SignerInfo) bool {
		wg.Done()
		return true
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount(ctx)
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
	multTxs.ToAggregate(multisigAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount)

	result := sendAggregateTransactionV1(t, func() (*sdk.AggregateTransactionV1, error) {
		return client.NewBondedAggregateV1Transaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}

func TestAddCosignatureHandlersV2(t *testing.T) {
	SkipIfEntityNotSupportedAtVersion(client, t, sdk.AggregateCompletedV2, sdk.AggregateCompletedV2Version)
	wg := sync.WaitGroup{}

	acc1, err := client.NewAccount(ctx)
	assert.Nil(t, err)
	acc2, err := client.NewAccount(ctx)
	assert.Nil(t, err)

	wg.Add(1)

	err = wsc.AddCosignatureHandlers(acc2.Address, func(info *sdk.SignerInfo) bool {
		wg.Done()
		return true
	})
	assert.Nil(t, err)

	multisigAccount, err := client.NewAccount(ctx)
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
	multTxs.ToAggregate(multisigAccount)

	fackeTxs, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount)

	result := sendAggregateTransactionV2(t, func() (*sdk.AggregateTransactionV2, error) {
		return client.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
	wg.Wait()
}
