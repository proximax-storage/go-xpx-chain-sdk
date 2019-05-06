// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"math/big"
	math "math/rand"
	"sync"
	"testing"
	"time"
)

const timeout = 2 * time.Minute
const networkType = sdk.MijinTest

var cfg, _ = sdk.NewConfig([]string{testUrl}, networkType, sdk.WebsocketReconnectionDefaultTimeout)
var ctx = context.Background()

var client = sdk.NewClient(nil, cfg)
var wsc, _ = websocket.NewClient(ctx, cfg)

type CreateTransaction func() (sdk.Transaction, error)

func waitTimeout(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return
	case <-time.After(timeout):
		t.Error("Timeout request")
		return
	}
}

func sendTransaction(t *testing.T, createTransaction CreateTransaction, account *sdk.Account) {

	// Starting listening messages from websocket
	go wsc.Listen()

	var wg sync.WaitGroup

	// Register handlers functions for needed topics
	wg.Add(1)
	if err := wsc.AddConfirmedAddedHandlers(account.Address, func(transaction sdk.Transaction) bool {
		fmt.Printf("ConfirmedAdded Tx Content: %v \n", transaction.GetAbstractTransaction().Hash)
		fmt.Println("Successful!")
		wg.Done()
		return true
	}); err != nil {
		panic(err)
	}

	if err := wsc.AddStatusHandlers(account.Address, func(info *sdk.StatusInfo) bool {
		fmt.Printf("Got error: %v \n", info)
		wg.Done()
		t.Error()
		return true
	}); err != nil {
		panic(err)
	}

	tx, err := createTransaction()

	signTx, err := account.Sign(tx)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	_, err = client.Transaction.Announce(ctx, signTx)
	assert.Nil(t, err)

	waitTimeout(t, &wg, timeout)
}

func TestMosaicDefinitionTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewMosaicDefinitionTransaction(
			sdk.NewDeadline(time.Hour),
			nonce,
			defaultAccount.PublicAccount.PublicKey,
			sdk.NewMosaicProperties(true, true, true, 4, big.NewInt(1)),
			networkType)
	}, defaultAccount)
}

func TestTransferTransaction(t *testing.T) {
	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.NewAddress("SDUP5PLHDXKBX3UU5Q52LAY4WYEKGEWC6IB3VBFM", networkType),
			[]*sdk.Mosaic{},
			sdk.NewPlainMessage("Test"),
			networkType,
		)

	}, defaultAccount)
}

func TestModifyMultisigTransaction(t *testing.T) {
	acc1, err := sdk.NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)
	acc2, err := sdk.NewAccountFromPublicKey("cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	multisigAccount, err := sdk.NewAccount(sdk.MijinTest)
	fmt.Println(multisigAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMultisigAccountTransaction(
			sdk.NewDeadline(time.Hour),
			2,
			1,
			[]*sdk.MultisigCosignatoryModification{
				{
					sdk.Add,
					acc1,
				},
				{
					sdk.Add,
					acc2,
				},
			},
			networkType,
		)
	}, multisigAccount)
}

func TestModifyContracTransaction(t *testing.T) {
	acc1, err := sdk.NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)
	acc2, err := sdk.NewAccountFromPublicKey("cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	contractAccount, err := sdk.NewAccount(sdk.MijinTest)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyContractTransaction(
			sdk.NewDeadline(time.Hour),
			2,
			"cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb",
			[]*sdk.MultisigCosignatoryModification{
				{
					sdk.Add,
					acc1,
				},
				{
					sdk.Add,
					acc2,
				},
			},
			[]*sdk.MultisigCosignatoryModification{
				{
					sdk.Add,
					acc1,
				},
				{
					sdk.Add,
					acc2,
				},
			},
			[]*sdk.MultisigCosignatoryModification{
				{
					sdk.Add,
					acc1,
				},
				{
					sdk.Add,
					acc2,
				},
			},
			networkType,
		)
	}, contractAccount)
}

func TestRegisterRootNamespaceTransaction(t *testing.T) {
	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.Nil(t, err)
	nameHex := hex.EncodeToString(name)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewRegisterRootNamespaceTransaction(
			sdk.NewDeadline(time.Hour),
			nameHex,
			big.NewInt(1),
			networkType,
		)
	}, defaultAccount)
}

func TestLockFundsTransactionTransaction(t *testing.T) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	assert.Nil(t, err)
	hash := sdk.Hash(hex.EncodeToString(key))

	stx := &sdk.SignedTransaction{sdk.AggregateBonded, "payload", hash}
	//id, err := sdk.NewMosaicId(big.NewInt(0x20B5A75C59C18264))
	//assert.Nil(t, err)
	//mosaic, err := sdk.NewMosaic(id, big.NewInt(10000000))
	//assert.Nil(t, err)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewLockFundsTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XemRelative(10),
			big.NewInt(100),
			stx,
			networkType,
		)
	}, defaultAccount)
}

func TestSecretTransactionTransaction(t *testing.T) {
	proof := make([]byte, 8)

	_, err := rand.Read(proof)
	assert.Nil(t, err)

	result := sha3.New256()
	_, err = result.Write(proof)
	assert.Nil(t, err)

	secret := hex.EncodeToString(result.Sum(nil))

	recipient := defaultAccount.PublicAccount.Address

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewSecretLockTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XemRelative(10),
			big.NewInt(100),
			sdk.SHA3_256,
			secret,
			recipient,
			networkType,
		)
	}, defaultAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewSecretProofTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.SHA3_256,
			secret,
			hex.EncodeToString(proof),
			networkType,
		)
	}, defaultAccount)
}

func TestCompleteAggregateTransactionTransaction(t *testing.T) {
	ttx, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		sdk.NewAddress("SBILTA367K2LX2FEXG5TFWAS7GEFYAGY7QLFBYKC", networkType),
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("test-message"),
		networkType,
	)
	assert.Nil(t, err)
	ttx.ToAggregate(defaultAccount.PublicAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{ttx},
			networkType,
		)
	}, defaultAccount)
}

func TestModifyAddressMetadataTransaction(t *testing.T) {
	fmt.Println(defaultAccount.PublicAccount.Address)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMetadataAddressTransaction(
			sdk.NewDeadline(time.Hour),
			defaultAccount.PublicAccount.Address,
			[]*sdk.MetadataModification{
				{
					sdk.AddMetadata,
					"jora229",
					"I Love you",
				},
			},
			networkType)
	}, defaultAccount)

	time.Sleep(2 * time.Second)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMetadataAddressTransaction(
			sdk.NewDeadline(time.Hour),
			defaultAccount.PublicAccount.Address,
			[]*sdk.MetadataModification{
				{
					sdk.RemoveMetadata,
					"jora229",
					"",
				},
			},
			networkType)
	}, defaultAccount)
}

func TestModifyMosaicMetadataTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	mosaicDefinitionTx, err := sdk.NewMosaicDefinitionTransaction(
		sdk.NewDeadline(time.Hour),
		nonce,
		defaultAccount.PublicAccount.PublicKey,
		sdk.NewMosaicProperties(true, true, true, 4, big.NewInt(1)),
		networkType)
	assert.Nil(t, err)
	mosaicDefinitionTx.ToAggregate(defaultAccount.PublicAccount)

	mosaicId, err := sdk.NewMosaicIdFromNonceAndOwner(nonce, defaultAccount.PublicAccount.PublicKey)
	assert.Nil(t, err)

	fmt.Println(mosaicId.String())

	metadataTx, err := sdk.NewModifyMetadataMosaicTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		[]*sdk.MetadataModification{
			{
				sdk.AddMetadata,
				"hello",
				"hell",
			},
		},
		networkType)
	assert.Nil(t, err)
	metadataTx.ToAggregate(defaultAccount.PublicAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{mosaicDefinitionTx, metadataTx},
			networkType,
		)
	}, defaultAccount)
}

func TestModifyNamespaceMetadataTransaction(t *testing.T) {
	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.Nil(t, err)
	nameHex := hex.EncodeToString(name)

	namespaceId, err := sdk.NewNamespaceIdFromName(nameHex)
	assert.Nil(t, err)
	fmt.Println(namespaceId)

	registrNamespaceTx, err := sdk.NewRegisterRootNamespaceTransaction(
		sdk.NewDeadline(time.Hour),
		nameHex,
		big.NewInt(10),
		networkType,
	)
	assert.Nil(t, err)
	registrNamespaceTx.ToAggregate(defaultAccount.PublicAccount)

	modifyMetadataTx, err := sdk.NewModifyMetadataNamespaceTransaction(
		sdk.NewDeadline(time.Hour),
		namespaceId,
		[]*sdk.MetadataModification{
			{
				sdk.AddMetadata,
				"hello",
				"world",
			},
		},
		networkType,
	)
	assert.Nil(t, err)
	modifyMetadataTx.ToAggregate(defaultAccount.PublicAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registrNamespaceTx, modifyMetadataTx},
			networkType,
		)
	}, defaultAccount)
}
