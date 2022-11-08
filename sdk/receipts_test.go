package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testBlockStatementJson = `{
		"addressResolutionStatements": [
			{
				"height": [123, 0],
				"unresolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE0",
				"resolutionEntries": [
					{
						"resolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE1"
					},
					{
						"resolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE2"
					}
				]
			},
			{
				"height": [123, 0],
				"unresolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE1",
				"resolutionEntries": [
					{
						"resolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE2"
					},
					{
						"resolved": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE3"
					}
				]
			}
		],
		"mosaicResolutionStatements": [
			{
				"height": [123, 0],
				"unresolved": [123, 0],
				"resolutionEntries": [
					{
						"resolved": [124, 0]
					},
					{
						"resolved": [125, 0]
					}
				]
			},
			{
				"height": [123, 0],
				"unresolved": [124, 0],
				"resolutionEntries": [
					{
						"resolved": [125, 0]
					},
					{
						"resolved": [126, 0]
					}
				]
			}
		],
		"transactionStatements": [
			{
				"height": [123, 0],
				"receipts": [
					{
						"type": 8515,
						"version": 1,
						"account": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 4419,
						"version": 1,
						"sender": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"recipient": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 16718,
						"version": 1,
						"artifactId": [123, 0]
					},
					{
						"type": 20803,
						"version": 1,
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 33091,
						"version": 1,
						"amount": [123, 0],
						"lockedAmount": [123, 0]
					},
					{
						"type": 33122,
						"version": 1,
						"amount": [123, 0]
					},
					{
						"type": 33347,
						"version": 1,
						"flags": [123, 0]
					}
				]
			}
		],
		"publicKeyStatements": [
			{
				"height": [123, 0],
				"receipts": [
					{
						"type": 8515,
						"version": 1,
						"account": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 4419,
						"version": 1,
						"sender": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"recipient": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 16718,
						"version": 1,
						"artifactId": [123, 0]
					},
					{
						"type": 20803,
						"version": 1,
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 33091,
						"version": 1,
						"amount": [123, 0],
						"lockedAmount": [123, 0]
					},
					{
						"type": 33122,
						"version": 1,
						"amount": [123, 0]
					},
					{
						"type": 33347,
						"version": 1,
						"flags": [123, 0]
					}
				]
			}
		],
		"blockchainStateStatements": [
			{
				"height": [123, 0],
				"receipts": [
					{
						"type": 8515,
						"version": 1,
						"account": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 4419,
						"version": 1,
						"sender": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716",
						"recipient": "906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4",
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 16718,
						"version": 1,
						"artifactId": [123, 0]
					},
					{
						"type": 20803,
						"version": 1,
						"mosaicId": [123, 0],
						"amount": [123, 0]
					},
					{
						"type": 33091,
						"version": 1,
						"amount": [123, 0],
						"lockedAmount": [123, 0]
					},
					{
						"type": 33122,
						"version": 1,
						"amount": [123, 0]
					},
					{
						"type": 33347,
						"version": 1,
						"flags": [123, 0]
					}
				]
			}
		]
	}`

	testBlockStatementJsonArr = "[" + testKeyRecordGroupEntryJson + "]"
)

var (
	testBlockStatement = &BlockStatement{
		TransactionStatements: []*TransactionStatement{
			&TransactionStatement{
				Height: Height(123),
				Receipts: []*Receipt{
					&Receipt{
						Header: GetReceiptHeader(BalanceChangeReceiptEntityType),
						Body: &BalanceChangeReceipt{
							Account:  newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(BalanceTransferReceiptEntityType),
						Body: &BalanceTransferReceipt{
							Sender:    newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							Recipient: newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4"),
							MosaicId:  newMosaicIdPanic(123),
							Amount:    Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(NamespaceArtifactExpiryReceiptEntityType),
						Body: &ArtifactExpiryReceipt{
							ArtifactId: 123,
						},
					},
					&Receipt{
						Header: GetReceiptHeader(InflationReceiptEntityType),
						Body: &InflationReceipt{
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(SignerImportanceReceiptEntityType),
						Body: &SignerBalanceReceipt{
							Amount:       Amount(123),
							LockedAmount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(TotalStakedReceiptEntityType),
						Body: &TotalStakedReceipt{
							Amount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(GlobalStateTrackingReceiptEntityType),
						Body: &GlobalStateChangeReceipt{
							Flags: 123,
						},
					},
				},
			},
		},
		AddressResolutionStatements: []*AddressResolutionStatement{
			&AddressResolutionStatement{
				Height:            Height(123),
				UnresolvedAddress: newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE0"),
				ResolutionEntries: []*Address{
					newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE1"),
					newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE2"),
				},
			},
			&AddressResolutionStatement{
				Height:            Height(123),
				UnresolvedAddress: newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE1"),
				ResolutionEntries: []*Address{
					newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE2"),
					newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE3"),
				},
			},
		},
		MosaicResolutionStatements: []*MosaicResolutionStatement{
			&MosaicResolutionStatement{
				Height:           Height(123),
				UnresolvedMosaic: newMosaicIdPanic(123),
				ResolutionEntries: []*MosaicId{
					newMosaicIdPanic(124),
					newMosaicIdPanic(125),
				},
			},
			&MosaicResolutionStatement{
				Height:           Height(123),
				UnresolvedMosaic: newMosaicIdPanic(124),
				ResolutionEntries: []*MosaicId{
					newMosaicIdPanic(125),
					newMosaicIdPanic(126),
				},
			},
		},
		PublicKeyStatements: []*PublicKeyStatement{
			&PublicKeyStatement{
				Height: Height(123),
				Receipts: []*Receipt{
					&Receipt{
						Header: GetReceiptHeader(BalanceChangeReceiptEntityType),
						Body: &BalanceChangeReceipt{
							Account:  newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(BalanceTransferReceiptEntityType),
						Body: &BalanceTransferReceipt{
							Sender:    newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							Recipient: newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4"),
							MosaicId:  newMosaicIdPanic(123),
							Amount:    Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(NamespaceArtifactExpiryReceiptEntityType),
						Body: &ArtifactExpiryReceipt{
							ArtifactId: 123,
						},
					},
					&Receipt{
						Header: GetReceiptHeader(InflationReceiptEntityType),
						Body: &InflationReceipt{
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(SignerImportanceReceiptEntityType),
						Body: &SignerBalanceReceipt{
							Amount:       Amount(123),
							LockedAmount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(TotalStakedReceiptEntityType),
						Body: &TotalStakedReceipt{
							Amount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(GlobalStateTrackingReceiptEntityType),
						Body: &GlobalStateChangeReceipt{
							Flags: 123,
						},
					},
				},
			},
		},
		BlockchainStateStatements: []*BlockchainStateStatement{
			&BlockchainStateStatement{
				Height: Height(123),
				Receipts: []*Receipt{
					&Receipt{
						Header: GetReceiptHeader(BalanceChangeReceiptEntityType),
						Body: &BalanceChangeReceipt{
							Account:  newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(BalanceTransferReceiptEntityType),
						Body: &BalanceTransferReceipt{
							Sender:    newAccountFromPublicKey("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", NetworkType(168)),
							Recipient: newAddressFromHexString("906B2CC5D7B66900B2493CF68BE10B7AA8690D973B7FD0DAE4"),
							MosaicId:  newMosaicIdPanic(123),
							Amount:    Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(NamespaceArtifactExpiryReceiptEntityType),
						Body: &ArtifactExpiryReceipt{
							ArtifactId: 123,
						},
					},
					&Receipt{
						Header: GetReceiptHeader(InflationReceiptEntityType),
						Body: &InflationReceipt{
							MosaicId: newMosaicIdPanic(123),
							Amount:   Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(SignerImportanceReceiptEntityType),
						Body: &SignerBalanceReceipt{
							Amount:       Amount(123),
							LockedAmount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(TotalStakedReceiptEntityType),
						Body: &TotalStakedReceipt{
							Amount: Amount(123),
						},
					},
					&Receipt{
						Header: GetReceiptHeader(GlobalStateTrackingReceiptEntityType),
						Body: &GlobalStateChangeReceipt{
							Flags: 123,
						},
					},
				},
			},
		},
	}
)

func Test_ReceiptService_GetBlockStatement(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(blockStatementsByHeight, Height(120)),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testBlockStatementJson,
	})
	receiptClient := mock.getPublicTestClientUnsafe().Receipt

	defer mock.Close()

	record, err := receiptClient.GetBlockStatementAtHeight(ctx, Height(120))
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testBlockStatement, record)
}
