// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk"
	"github.com/proximax-storage/go-xpx-catapult-sdk/sdk/websocket"
	"github.com/stretchr/testify/assert"
	math "math/rand"
	"testing"
	"time"
)

const timeout = 2 * time.Minute

var cfg, _ = sdk.NewConfig([]string{testUrl}, networkType, sdk.WebsocketReconnectionDefaultTimeout)
var ctx = context.Background()

var client = sdk.NewClient(nil, cfg)
var wsc, _ = websocket.NewClient(ctx, cfg)
var listening = false

type CreateTransaction func() (sdk.Transaction, error)

type Result struct {
	sdk.Transaction
	error
}

func initListeners(t *testing.T, account *sdk.Account) <-chan Result {
	if !listening {
		// Starting listening messages from websocket
		go wsc.Listen()
		listening = true
	}

	out := make(chan Result)

	// Register handlers functions for needed topics
	if err := wsc.AddConfirmedAddedHandlers(account.Address, func(transaction sdk.Transaction) bool {
		fmt.Printf("ConfirmedAdded Tx Content: %v \n", transaction)
		fmt.Println("Successful!")
		out <- Result{transaction, nil}
		return true
	}); err != nil {
		panic(err)
	}

	if err := wsc.AddStatusHandlers(account.Address, func(info *sdk.StatusInfo) bool {
		fmt.Printf("Got error: %v \n", info)
		t.Error()
		out <- Result{nil, errors.New(info.Status)}
		return true
	}); err != nil {
		panic(err)
	}

	return out
}

func waitTimeout(t *testing.T, wg <-chan Result, timeout time.Duration) Result {
	select {
	case result := <-wg:
		return result
	case <-time.After(timeout):
		t.Error("Timeout request")
		return Result{nil, errors.New("Timeout error")}
	}
}

func sendTransaction(t *testing.T, createTransaction CreateTransaction, account *sdk.Account) Result {
	tx, err := createTransaction()
	assert.Nil(t, err)

	signTx, err := account.Sign(tx, GenerationHash)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	assert.Nil(t, err)
	wg := initListeners(t, account)
	_, err = client.Transaction.Announce(ctx, signTx)
	assert.Nil(t, err)

	return waitTimeout(t, wg, timeout)
}

func sendAggregateTransaction(t *testing.T, createTransaction func() (*sdk.AggregateTransaction, error), account *sdk.Account, cosignatories ...*sdk.Account) Result {
	tx, err := createTransaction()
	assert.Nil(t, err)

	signTx, err := account.SignWithCosignatures(tx, cosignatories, GenerationHash)
	assert.Nil(t, err)

	stx := &sdk.SignedTransaction{sdk.AggregateBonded, "payload", signTx.Hash}

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewLockFundsTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XpxRelative(10),
			sdk.Duration(100),
			stx,
			networkType,
		)
	}, account)

	if result.error != nil {
		return result
	}

	time.Sleep(2 * time.Second)

	wg := initListeners(t, account)
	_, err = client.Transaction.AnnounceAggregateBonded(ctx, signTx)
	assert.Nil(t, err)

	return waitTimeout(t, wg, timeout)
}

func TestAccountLinkTransaction(t *testing.T) {
	rootAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	fmt.Println(rootAccount)
	childAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	fmt.Println(childAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewAccountLinkTransaction(
			sdk.NewDeadline(time.Hour),
			childAccount.PublicAccount,
			sdk.AccountLink,
			networkType)
	}, rootAccount)
	assert.Nil(t, result.error)
}

func TestMosaicDefinitionTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewMosaicDefinitionTransaction(
			sdk.NewDeadline(time.Hour),
			nonce,
			defaultAccount.PublicAccount.PublicKey,
			sdk.NewMosaicProperties(true, true, 4, sdk.Duration(1)),
			networkType)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestTransferTransaction(t *testing.T) {
	recipientAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			recipientAccount.Address,
			[]*sdk.Mosaic{},
			sdk.NewPlainMessage("Test"),
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestTransferTransaction_SecureMessage(t *testing.T) {
	const message = "I love you forever"
	recipientAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	secureMessage, err := sdk.NewSecureMessageFromPlaintText(message, defaultAccount.PrivateKey, recipientAccount.KeyPair.PublicKey)
	assert.Nil(t, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			recipientAccount.PublicAccount.Address,
			[]*sdk.Mosaic{},
			secureMessage,
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	transfer := result.Transaction.(*sdk.TransferTransaction)
	plainMessage, err := recipientAccount.DecryptMessage(
		transfer.Message.(*sdk.SecureMessage),
		defaultAccount.PublicAccount,
	)

	assert.Equal(t, message, plainMessage.Message())
}

func TestModifyMultisigTransaction(t *testing.T) {
	acc1, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	acc2, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	multisigAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	fmt.Println(multisigAccount)

	multTxs, err := sdk.NewModifyMultisigAccountTransaction(
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
		networkType,
	)
	assert.Nil(t, err)
	multTxs.ToAggregate(multisigAccount.PublicAccount)

	fackeTxs, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		multisigAccount.PublicAccount.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("I wan't to create multisig"),
		networkType,
	)
	assert.Nil(t, err)
	fackeTxs.ToAggregate(defaultAccount.PublicAccount)

	result := sendAggregateTransaction(t, func() (*sdk.AggregateTransaction, error) {
		return sdk.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{multTxs, fackeTxs},
			networkType,
		)
	}, defaultAccount, multisigAccount, acc1, acc2)
	assert.Nil(t, result.error)
}

func TestModifyContracTransaction(t *testing.T) {
	acc1, err := sdk.NewAccountFromPublicKey("68b3fbb18729c1fde225c57f8ce080fa828f0067e451a3fd81fa628842b0b763", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)
	acc2, err := sdk.NewAccountFromPublicKey("cf893ffcc47c33e7f68ab1db56365c156b0736824a0c1e273f9e00b8df8f01eb", networkType)
	assert.Nilf(t, err, "NewAccountFromPublicKey returned error: %s", err)

	contractAccount, err := sdk.NewAccount(networkType)
	fmt.Println(contractAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyContractTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(2),
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
	assert.Nil(t, result.error)
}

func TestRegisterRootNamespaceTransaction(t *testing.T) {
	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.Nil(t, err)
	nameHex := hex.EncodeToString(name)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewRegisterRootNamespaceTransaction(
			sdk.NewDeadline(time.Hour),
			nameHex,
			sdk.Duration(1),
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestLockFundsTransactionTransaction(t *testing.T) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	assert.Nil(t, err)
	hash := sdk.Hash(hex.EncodeToString(key))

	stx := &sdk.SignedTransaction{sdk.AggregateBonded, "payload", hash}

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewLockFundsTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.XpxRelative(10),
			sdk.Duration(100),
			stx,
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestSecretTransaction(t *testing.T) {
	for _, hashType := range []sdk.HashType{sdk.SHA_256, sdk.KECCAK_256, sdk.SHA3_256, sdk.HASH_160} {
		proofB := make([]byte, 8)
		_, err := rand.Read(proofB)
		assert.Nil(t, err)

		proof := sdk.NewProofFromBytes(proofB)
		secret, err := proof.Secret(hashType)
		assert.Nil(t, err)
		recipient := defaultAccount.PublicAccount.Address

		result := sendTransaction(t, func() (sdk.Transaction, error) {
			return sdk.NewSecretLockTransaction(
				sdk.NewDeadline(time.Hour),
				sdk.XpxRelative(10),
				sdk.Duration(100),
				secret,
				recipient,
				networkType,
			)
		}, defaultAccount)
		assert.Nil(t, result.error)

		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return sdk.NewSecretProofTransaction(
				sdk.NewDeadline(time.Hour),
				hashType,
				proof,
				recipient,
				networkType,
			)
		}, defaultAccount)
		assert.Nil(t, result.error)
	}
}

func TestCompleteAggregateTransaction(t *testing.T) {
	acc, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	ttx, err := sdk.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		acc.Address,
		[]*sdk.Mosaic{},
		sdk.NewPlainMessage("test-message"),
		networkType,
	)
	assert.Nil(t, err)
	ttx.ToAggregate(defaultAccount.PublicAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{ttx},
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
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

	result := sendAggregateTransaction(t, func() (*sdk.AggregateTransaction, error) {
		return sdk.NewBondedAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{ttx1, ttx2},
			networkType,
		)
	}, defaultAccount, receiverAccount)
	assert.Nil(t, result.error)
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
		sdk.Duration(10),
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

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registerTx, aliasTx},
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	senderAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewTransferTransactionWithNamespace(
			sdk.NewDeadline(time.Hour),
			nsId,
			[]*sdk.Mosaic{},
			sdk.NewPlainMessage("Test"),
			networkType,
		)
	}, senderAccount)
	assert.Nil(t, result.error)
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
		sdk.Duration(10),
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
		sdk.NewMosaicProperties(true, true, 4, sdk.Duration(1)),
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

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registerTx, mosaicDefinitionTx, aliasTx},
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestModifyAddressMetadataTransaction(t *testing.T) {
	fmt.Println(defaultAccount.PublicAccount.Address)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewModifyMetadataAddressTransaction(
			sdk.NewDeadline(time.Hour),
			defaultAccount.PublicAccount.Address,
			[]*sdk.MetadataModification{
				{
					sdk.AddMetadata,
					"jora229",
					"I Love you!",
				},
			},
			networkType)
	}, defaultAccount)
	assert.Nil(t, result.error)

	time.Sleep(2 * time.Second)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
}

func TestModifyMosaicMetadataTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	mosaicDefinitionTx, err := sdk.NewMosaicDefinitionTransaction(
		sdk.NewDeadline(time.Hour),
		nonce,
		defaultAccount.PublicAccount.PublicKey,
		sdk.NewMosaicProperties(true, true, 4, sdk.Duration(1)),
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

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{mosaicDefinitionTx, metadataTx},
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	time.Sleep(2 * time.Second)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
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
		sdk.Duration(10),
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

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registrNamespaceTx, modifyMetadataTx},
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	time.Sleep(2 * time.Second)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
}

func TestAccountPropertiesAddressTransaction(t *testing.T) {
	blockAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)
	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(blockAccount, testAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
}

func TestAccountPropertiesMosaicTransaction(t *testing.T) {
	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()

	mosaicId, err := sdk.NewMosaicIdFromNonceAndOwner(nonce, defaultAccount.PublicAccount.PublicKey)
	assert.Nil(t, err)

	fmt.Println(mosaicId.String())

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return sdk.NewMosaicDefinitionTransaction(
			sdk.NewDeadline(time.Hour),
			nonce,
			defaultAccount.PublicAccount.PublicKey,
			sdk.NewMosaicProperties(true, true, 4, sdk.Duration(1)),
			networkType,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(testAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
}

func TestAccountPropertiesEntityTypeTransaction(t *testing.T) {
	testAccount, err := sdk.NewAccount(networkType)
	assert.Nil(t, err)

	fmt.Println(testAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
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
	assert.Nil(t, result.error)
}
