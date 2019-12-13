// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	addExchangeOfferTransactionSerializationCorr = []byte{0x9c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x41, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	addExchangeOfferTransactionToAggregateCorr = []byte{0x4c, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x41, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	addExchangeOfferTransactionSigningCorr = "9C00000008097BEDBFED6DE077A631FFD1252E6128FA42CCB69255B2B0FD5C7B865005C0B2C63CB2FE4A608B49193E8F6394325804219B5AD763CA39F1D7823FB53414061026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF37755010000905D41000000000000000000BAFD560000000001513874B40E8D4FBC02000000000000000200000000000000000100000000000000"

	exchangeOfferTransactionSerializationCorr = []byte{0xb4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x42, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24}

	exchangeOfferTransactionToAggregateCorr = []byte{0x64, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x42, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xdd, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24}

	exchangeOfferTransactionSigningCorr = "B4000000AFC70E8B1085688F703EF656B108CB3E28431C12D73124FA043EA6800976EC2C26FA105946E07C69B56D245268875E573722C549695846262CB4D6F0525D14031026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF37755010000905D42000000000000000000BAFD560000000001513874B40E8D4FBC02000000000000000200000000000000009A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24"

	removeExchangeOfferTransactionSerializationCorr = []byte{0x84, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xba, 0xfd, 0x56, 0x0, 0x0, 0x0, 0x0, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x0}

	removeExchangeOfferTransactionToAggregateCorr = []byte{0x34, 0x0, 0x0, 0x0, 0x9a, 0x49, 0x36, 0x64, 0x6, 0xac, 0xa9, 0x52, 0xb8, 0x8b, 0xad, 0xf5, 0xf1, 0xe9, 0xbe, 0x6c, 0xe4, 0x96, 0x81, 0x41, 0x3, 0x5a, 0x60, 0xbe, 0x50, 0x32, 0x73, 0xea, 0x65, 0x45, 0x6b, 0x24, 0x1, 0x0, 0x0, 0x90, 0x5d, 0x43, 0x1, 0x51, 0x38, 0x74, 0xb4, 0xe, 0x8d, 0x4f, 0xbc, 0x0}

	removeExchangeOfferTransactionSigningCorr = "84000000C7DDA232E1E56500F734C115451A6E22A067A4D074F1719C09719094AE1A1A95A0EBA6EEFFA9E8DC8C9668C5E6EBC3C5F340447B80242EADCBDB3ADB44DF220B1026D70E1954775749C6811084D6450A3184D977383F0E4282CD47118AF37755010000905D43000000000000000000BAFD560000000001513874B40E8D4FBC00"
)

func TestAddExchangeOfferTransactionSerialization(t *testing.T) {
	tx, err := NewAddExchangeOfferTransaction(
		fakeDeadline,
		[]*AddOffer{
			&AddOffer{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Duration: Duration(1),
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewAddExchangeOfferTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "AddExchangeOfferTransaction.Bytes returned error: %s", err)
	assert.Equal(t, addExchangeOfferTransactionSerializationCorr, b)
}

func TestAddExchangeOfferTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	tx, err := NewAddExchangeOfferTransaction(
		fakeDeadline,
		[]*AddOffer{
			&AddOffer{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Duration: Duration(1),
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewAddExchangeOfferTransaction returned error: %s", err)
	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, addExchangeOfferTransactionToAggregateCorr, b)
}

func TestAddExchangeOfferTransactionSigning(t *testing.T) {
	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest, GenerationHash)
	assert.Nil(t, err)

	tx, err := NewAddExchangeOfferTransaction(
		fakeDeadline,
		[]*AddOffer{
			&AddOffer{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Duration: Duration(1),
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewAddExchangeOfferTransaction returned error: %s", err)

	b, err := acc.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, addExchangeOfferTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("aefce3842d5219d7d9bf573c345498645b82028719de34e687d92760db30d772"), b.Hash)
}

func TestExchangeOfferTransactionSerialization(t *testing.T) {
	owner, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	assert.Nil(t, err)
	tx, err := NewExchangeOfferTransaction(
		fakeDeadline,
		[]*ExchangeConfirmation{
			&ExchangeConfirmation{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Owner: owner,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewExchangeOfferTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "ExchangeOfferTransaction.Bytes returned error: %s", err)
	assert.Equal(t, exchangeOfferTransactionSerializationCorr, b)
}

func TestExchangeOfferTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	assert.Nil(t, err)

	owner, err := NewAccountFromPublicKey("DD49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	assert.Nil(t, err)
	tx, err := NewExchangeOfferTransaction(
		fakeDeadline,
		[]*ExchangeConfirmation{
			&ExchangeConfirmation{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Owner: owner,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewExchangeOfferTransaction returned error: %s", err)
	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, exchangeOfferTransactionToAggregateCorr, b)
}

func TestExchangeOfferTransactionSigning(t *testing.T) {
	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest, GenerationHash)
	assert.Nil(t, err)

	owner, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	assert.Nil(t, err)
	tx, err := NewExchangeOfferTransaction(
		fakeDeadline,
		[]*ExchangeConfirmation{
			&ExchangeConfirmation{
				Offer: Offer{
					Type: SellOffer,
					Cost: Amount(2),
					Mosaic: newMosaicPanic(StorageNamespaceId, Amount(2)),
				},
				Owner: owner,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewExchangeOfferTransaction returned error: %s", err)

	b, err := acc.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, exchangeOfferTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("0b8621c77e1e36b10fbcca489a82661cfc2330cdf3771a8c4fc53b1c318d745b"), b.Hash)
}

func TestRemoveExchangeOfferTransactionSerialization(t *testing.T) {
	tx, err := NewRemoveExchangeOfferTransaction(
		fakeDeadline,
		[]*RemoveOffer{
			&RemoveOffer{
				Type: SellOffer,
				AssetId: StorageNamespaceId,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewRemoveExchangeOfferTransaction returned error: %s", err)

	b, err := tx.Bytes()

	assert.Nilf(t, err, "RemoveExchangeOfferTransaction.Bytes returned error: %s", err)
	assert.Equal(t, removeExchangeOfferTransactionSerializationCorr, b)
}

func TestRemoveExchangeOfferTransactionToAggregate(t *testing.T) {
	p, err := NewAccountFromPublicKey("9A49366406ACA952B88BADF5F1E9BE6CE4968141035A60BE503273EA65456B24", MijinTest)
	assert.Nil(t, err)

	tx, err := NewRemoveExchangeOfferTransaction(
		fakeDeadline,
		[]*RemoveOffer{
			&RemoveOffer{
				Type: SellOffer,
				AssetId: StorageNamespaceId,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewRemoveExchangeOfferTransaction returned error: %s", err)
	tx.Signer = p

	b, err := toAggregateTransactionBytes(tx)

	assert.Nilf(t, err, "toAggregateTransactionBytes returned error: %s", err)
	assert.Equal(t, removeExchangeOfferTransactionToAggregateCorr, b)
}

func TestRemoveExchangeOfferTransactionSigning(t *testing.T) {
	acc, err := NewAccountFromPrivateKey("787225aaff3d2c71f4ffa32d4f19ec4922f3cd869747f267378f81f8e3fcb12d", MijinTest, GenerationHash)
	assert.Nil(t, err)

	tx, err := NewRemoveExchangeOfferTransaction(
		fakeDeadline,
		[]*RemoveOffer{
			&RemoveOffer{
				Type: SellOffer,
				AssetId: StorageNamespaceId,
			},
		},
		MijinTest,
	)
	assert.Nilf(t, err, "NewRemoveExchangeOfferTransaction returned error: %s", err)

	b, err := acc.Sign(tx)

	assert.Nilf(t, err, "Sign returned error: %s", err)
	assert.Equal(t, removeExchangeOfferTransactionSigningCorr, b.Payload)
	assert.Equal(t, stringToHashPanic("6cae5bc3b7f2f7793d4e82a70c00535cbd675a02793937e11fd5e5b35debf1fe"), b.Hash)
}
