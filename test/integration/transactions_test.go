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
	"math/big"
	math "math/rand"
	"sync"
	"testing"
	"time"
)

const timeout = 2 * time.Minute

var cfg, _ = sdk.NewConfig(testUrl, networkType)
var ctx = context.Background()

var client = sdk.NewClient(nil, cfg)
var wsc, _ = websocket.NewClient(ctx, cfg)
var listening = false

type CreateTransaction func() (sdk.Transaction, error)

func initListeners(t *testing.T, account *sdk.Account) *sync.WaitGroup {
	if !listening {
		// Starting listening messages from websocket
		go wsc.Listen()
		listening = true
	}

	var wg sync.WaitGroup

	// Register handlers functions for needed topics
	wg.Add(1)
	if err := wsc.AddConfirmedAddedHandlers(account.Address, func(transaction sdk.Transaction) bool {
		fmt.Printf("ConfirmedAdded Tx Content: %v \n", transaction)
		fmt.Println("Successful!")
		wg.Done()
		return true
	}); err != nil {
		panic(err)
	}

	if err := wsc.AddStatusHandlers(account.Address, func(info *sdk.StatusInfo) bool {
		fmt.Printf("Got error: %v \n", info)
		t.Error()
		wg.Done()
		return true
	}); err != nil {
		panic(err)
	}

	return &wg
}

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
	tx, err := createTransaction()
	assert.Nil(t, err)

	signTx, err := account.Sign(tx)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	assert.Nil(t, err)
	wg := initListeners(t, account)
	_, err = client.Transaction.Announce(ctx, signTx)
	assert.Nil(t, err)

	waitTimeout(t, wg, timeout)
}

func sendAggregateTransaction(t *testing.T, createTransaction func() (*sdk.AggregateTransaction, error), account *sdk.Account, cosignatories ...*sdk.Account) {
	tx, err := createTransaction()
	assert.Nil(t, err)

	signTx, err := account.SignWithCosignatures(tx, cosignatories)
	assert.Nil(t, err)

	stx := &sdk.SignedTransaction{sdk.AggregateBonded, "payload", signTx.Hash}

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewLockFundsTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XpxRelative(10),
			big.NewInt(100),
			stx,
			networkType,
		)
	}, account)

	wg := initListeners(t, account)
	_, err = client.Transaction.AnnounceAggregateBonded(ctx, signTx)
	assert.Nil(t, err)

	waitTimeout(t, wg, timeout)
}

func TestAccountLinkTransaction(t *testing.T) {
	rootAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	fmt.Println(rootAccount)
	childAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	fmt.Println(childAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewAccountLinkTransaction(
			sdk.NewDeadline(time.Hour),
			childAccount.PublicAccount,
			sdk.AccountLink,
			networkType)
	}, rootAccount)
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

	multisigAccount, err := sdk.NewAccount(networkType)
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

	contractAccount, err := sdk.NewAccount(networkType)
	fmt.Println(contractAccount)

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
			sdk.XpxRelative(10),
			big.NewInt(100),
			stx,
			networkType,
		)
	}, defaultAccount)
}

func TestSecretTransaction(t *testing.T) {
	for _, hashType := range []sdk.HashType{sdk.SHA_256, sdk.KECCAK_256, sdk.SHA3_256, sdk.RIPEMD_160} {
		proofB := make([]byte, 8)
		_, err := rand.Read(proofB)
		assert.Nil(t, err)

		proof := sdk.NewProofFromBytes(proofB)
		secret, err := proof.Secret(hashType)
		assert.Nil(t, err)
		recipient := defaultAccount.PublicAccount.Address

		sendTransaction(t, func() (sdk.Transaction, error) {
			return sdk.NewSecretLockTransaction(
				sdk.NewDeadline(time.Hour),
				sdk.XpxRelative(10),
				big.NewInt(100),
				secret,
				recipient,
				networkType,
			)
		}, defaultAccount)

		sendTransaction(t, func() (sdk.Transaction, error) {
			return sdk.NewSecretProofTransaction(
				sdk.NewDeadline(time.Hour),
				hashType,
				proof,
				networkType,
			)
		}, defaultAccount)
	}
}

func TestCompleteAggregateTransaction(t *testing.T) {
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

func TestAggregateBoundedTransaction(t *testing.T) {
	receiverAccount, err := sdk.NewAccount(networkType)

	ttx1, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		receiverAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("test-message"),
		networkType,
	)
	assert.Nil(t, err)
	ttx1.ToAggregate(defaultAccount.PublicAccount)

	ttx2, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("test-message"),
		networkType,
	)
	assert.Nil(t, err)
	ttx2.ToAggregate(receiverAccount.PublicAccount)

	sendAggregateTransaction(t, func() (*sdk.AggregateTransaction, error) {
		return sdk.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{ttx1, ttx2},
			networkType,
		)
	}, defaultAccount, receiverAccount)
}

func TestAddressAliasTransaction(t *testing.T) {
	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.Nil(t, err)
	nameHex := hex.EncodeToString(name)

	nsId, err := sdk.NewNamespaceIdFromName(nameHex)
	assert.Nil(t, err)

	registerTx, err := sdk.NewRegisterRootNamespaceTransaction(
		sdk.NewDeadline(time.Hour),
		nameHex,
		big.NewInt(10),
		networkType,
	)
	assert.Nil(t, err)
	registerTx.ToAggregate(defaultAccount.PublicAccount)

	aliasTx, err := sdk.NewAddressAliasTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.PublicAccount.Address,
		nsId,
		sdk.AliasLink,
		networkType,
	)
	assert.Nil(t, err)
	aliasTx.ToAggregate(defaultAccount.PublicAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registerTx, aliasTx},
			networkType,
		)
	}, defaultAccount)

	senderAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewTransferTransactionWithNamespace(
			sdk.NewDeadline(time.Hour),
			nsId,
			[]*sdk.Mosaic{},
			sdk.NewPlainMessage("Test"),
			networkType,
		)
	}, senderAccount)
}

func TestMosaicAliasTransaction(t *testing.T) {
	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.Nil(t, err)
	nameHex := hex.EncodeToString(name)

	nsId, err := sdk.NewNamespaceIdFromName(nameHex)
	assert.Nil(t, err)

	registerTx, err := sdk.NewRegisterRootNamespaceTransaction(
		sdk.NewDeadline(time.Hour),
		nameHex,
		big.NewInt(10),
		networkType,
	)
	assert.Nil(t, err)
	registerTx.ToAggregate(defaultAccount.PublicAccount)

	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	mosaicId, err := sdk.NewMosaicIdFromNonceAndOwner(nonce, defaultAccount.PublicAccount.PublicKey)
	assert.Nil(t, err)
	mosaicDefinitionTx, err := sdk.NewMosaicDefinitionTransaction(
		sdk.NewDeadline(time.Hour),
		nonce,
		defaultAccount.PublicAccount.PublicKey,
		sdk.NewMosaicProperties(true, true, true, 4, big.NewInt(1)),
		networkType,
	)
	assert.Nil(t, err)
	mosaicDefinitionTx.ToAggregate(defaultAccount.PublicAccount)

	aliasTx, err := sdk.NewMosaicAliasTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		nsId,
		sdk.AliasLink,
		networkType,
	)
	assert.Nil(t, err)
	aliasTx.ToAggregate(defaultAccount.PublicAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registerTx, mosaicDefinitionTx, aliasTx},
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

	time.Sleep(2 * time.Second)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMetadataMosaicTransaction(
			sdk.NewDeadline(time.Hour),
			mosaicId,
			[]*sdk.MetadataModification{
				{
					sdk.RemoveMetadata,
					"hello",
					"",
				},
			},
			networkType)
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

	time.Sleep(2 * time.Second)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMetadataNamespaceTransaction(
			sdk.NewDeadline(time.Hour),
			namespaceId,
			[]*sdk.MetadataModification{
				{
					sdk.RemoveMetadata,
					"hello",
					"",
				},
			},
			networkType)
	}, defaultAccount)
}

func TestAccountPropertiesAddressTransaction(t *testing.T) {
	blockAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(blockAccount, testAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewAccountPropertiesAddressTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.BlockAddress,
			[]*sdk.AccountPropertiesAddressModification{
				{
					sdk.AddProperty,
					blockAccount.Address,
				},
			},
			networkType,
		)
	}, testAccount)
}

func TestAccountPropertiesMosaicTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	mosaicId, err := sdk.NewMosaicIdFromNonceAndOwner(nonce, defaultAccount.PublicAccount.PublicKey)
	assert.Nil(t, err)

	fmt.Println(mosaicId.String())

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewMosaicDefinitionTransaction(
			sdk.NewDeadline(time.Hour),
			nonce,
			defaultAccount.PublicAccount.PublicKey,
			sdk.NewMosaicProperties(true, true, true, 4, big.NewInt(1)),
			networkType,
		)
	}, defaultAccount)

	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(testAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewAccountPropertiesMosaicTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.BlockMosaic,
			[]*sdk.AccountPropertiesMosaicModification{
				{
					sdk.AddProperty,
					mosaicId,
				},
			},
			networkType,
		)
	}, testAccount)
}

func TestAccountPropertiesEntityTypeTransaction(t *testing.T) {
	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(testAccount)

	sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewAccountPropertiesEntityTypeTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.BlockTransaction,
			[]*sdk.AccountPropertiesEntityTypeModification{
				{
					sdk.AddProperty,
					sdk.LinkAccount,
				},
			},
			networkType,
		)
	}, testAccount)
}
