// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "reflect"

const (
	BalanceTransferReceiptEntityType         EntityType = 0x1143
	NamespaceRentalReceiptEntityType         EntityType = 0x124E
	MosaicRentalReceiptEntityType            EntityType = 0x124D
	BalanceChangeReceiptEntityType           EntityType = 0x2143
	BalanceDebitReceiptEntityType            EntityType = 0x3143
	LockHashCreatedReceiptEntityType         EntityType = 0x3148
	LockHashCompletedReceiptEntityType       EntityType = 0x2248
	LockHashExpiredReceiptEntityType         EntityType = 0x2348
	LockSecretCreatedReceiptEntityType       EntityType = 0x3152
	LockSecretCompletedReceiptEntityType     EntityType = 0x2252
	LockSecretExpiredReceiptEntityType       EntityType = 0x2352
	MosaicArtifactExpiryReceiptEntityType    EntityType = 0x414D
	NamespaceArtifactExpiryReceiptEntityType EntityType = 0x414E
	InflationReceiptEntityType               EntityType = 0x5143
	SignerImportanceReceiptEntityType        EntityType = 0x8143
	GlobalStateTrackingReceiptEntityType     EntityType = 0x8243
	TotalStakedReceiptEntityType             EntityType = 0x8162
	OperationStartedReceiptEntityType        EntityType = 0x715F
	OperationEndedReceiptEntityType          EntityType = 0x725F
	OperationExpiredReceiptEntityType        EntityType = 0x735F

	DriveStateReceiptEntityType                EntityType = 0x615B
	DriveDepositCreditReceiptEntityType        EntityType = 0x625B
	DriveDepositDebitReceiptEntityType         EntityType = 0x635B
	DriveRewardTransferCreditReceiptEntityType EntityType = 0x645B
	DriveRewardTransferDebitReceiptEntityType  EntityType = 0x655B
	DriveDownloadStartedReceiptEntityType      EntityType = 0x665B
	DriveDownloadCompletedReceiptEntityType    EntityType = 0x675B
	DriveDownloadExpiredReceiptEntityType      EntityType = 0x685B
)

type ReceiptDefinition struct {
	Type    reflect.Type
	DtoType reflect.Type
	Version EntityVersion
	Size    uint32
}

func MakeDefinition(eType reflect.Type, dtoType reflect.Type, version EntityVersion, size int) ReceiptDefinition {

	return ReceiptDefinition{
		Type:    eType,
		DtoType: dtoType,
		Version: version,
		Size:    uint32(size),
	}
}

const (
	ReceiptHeaderSize              int = 10
	BalanceTransferReceiptSize     int = ReceiptHeaderSize + KeySize + AddressSize + MosaicIdSize + AmountSize
	BalanceChangeReceiptSize       int = ReceiptHeaderSize + KeySize + MosaicIdSize + AmountSize
	BalanceDebitReceiptSize        int = ReceiptHeaderSize + KeySize + MosaicIdSize + AmountSize
	ArtifactExpiryReceiptSize      int = ReceiptHeaderSize + BaseInt64Size
	DriveReceiptSize               int = ReceiptHeaderSize + KeySize + ByteFlagsSize
	InflationReceiptSize           int = ReceiptHeaderSize + MosaicIdSize + AmountSize
	SignerImportanceReceiptSize    int = ReceiptHeaderSize + AmountSize + MosaicIdSize
	GlobalStateTrackingReceiptSize int = ReceiptHeaderSize + BaseInt64Size
	TotalStakedReceiptSize         int = ReceiptHeaderSize + AmountSize
)

var ReceiptDefinitionMap = map[EntityType]ReceiptDefinition{
	BalanceTransferReceiptEntityType:         MakeDefinition(reflect.TypeOf((*BalanceTransferReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceTransferReceiptDto)(nil)).Elem(), BalanceTransferReceiptVersion, BalanceTransferReceiptSize),
	NamespaceRentalReceiptEntityType:         MakeDefinition(reflect.TypeOf((*BalanceTransferReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceTransferReceiptDto)(nil)).Elem(), NamespaceRentalReceiptVersion, BalanceTransferReceiptSize),
	MosaicRentalReceiptEntityType:            MakeDefinition(reflect.TypeOf((*BalanceTransferReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceTransferReceiptDto)(nil)).Elem(), MosaicRentalReceiptVersion, BalanceTransferReceiptSize),
	BalanceChangeReceiptEntityType:           MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), BalanceChangeReceiptVersion, BalanceChangeReceiptSize),
	BalanceDebitReceiptEntityType:            MakeDefinition(reflect.TypeOf((*BalanceDebitReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceDebitReceiptDto)(nil)).Elem(), BalanceDebitReceiptVersion, BalanceDebitReceiptSize),
	LockHashCreatedReceiptEntityType:         MakeDefinition(reflect.TypeOf((*BalanceDebitReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceDebitReceiptDto)(nil)).Elem(), LockHashCreatedReceiptVersion, BalanceDebitReceiptSize),
	LockHashCompletedReceiptEntityType:       MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), LockHashCompletedReceiptVersion, BalanceChangeReceiptSize),
	LockHashExpiredReceiptEntityType:         MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), LockHashExpiredReceiptVersion, BalanceChangeReceiptSize),
	LockSecretCreatedReceiptEntityType:       MakeDefinition(reflect.TypeOf((*BalanceDebitReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceDebitReceiptDto)(nil)).Elem(), LockSecretCreatedReceiptVersion, BalanceDebitReceiptSize),
	LockSecretCompletedReceiptEntityType:     MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), LockSecretCompletedReceiptVersion, BalanceChangeReceiptSize),
	LockSecretExpiredReceiptEntityType:       MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), LockSecretExpiredReceiptVersion, BalanceChangeReceiptSize),
	MosaicArtifactExpiryReceiptEntityType:    MakeDefinition(reflect.TypeOf((*ArtifactExpiryReceipt)(nil)).Elem(), reflect.TypeOf((*ArtifactExpiryReceiptDto)(nil)).Elem(), MosaicArtifactExpiryReceiptVersion, ArtifactExpiryReceiptSize),
	NamespaceArtifactExpiryReceiptEntityType: MakeDefinition(reflect.TypeOf((*ArtifactExpiryReceipt)(nil)).Elem(), reflect.TypeOf((*ArtifactExpiryReceiptDto)(nil)).Elem(), NamespaceArtifactExpiryReceiptVersion, ArtifactExpiryReceiptSize),
	InflationReceiptEntityType:               MakeDefinition(reflect.TypeOf((*InflationReceipt)(nil)).Elem(), reflect.TypeOf((*InflationReceiptDto)(nil)).Elem(), InflationReceiptVersion, InflationReceiptSize),
	SignerImportanceReceiptEntityType:        MakeDefinition(reflect.TypeOf((*SignerBalanceReceipt)(nil)).Elem(), reflect.TypeOf((*SignerBalanceReceiptDto)(nil)).Elem(), SignerImportanceReceiptVersion, SignerImportanceReceiptSize),
	GlobalStateTrackingReceiptEntityType:     MakeDefinition(reflect.TypeOf((*GlobalStateChangeReceipt)(nil)).Elem(), reflect.TypeOf((*GlobalStateChangeReceiptDto)(nil)).Elem(), GlobalStateTrackingReceiptVersion, GlobalStateTrackingReceiptSize),
	TotalStakedReceiptEntityType:             MakeDefinition(reflect.TypeOf((*TotalStakedReceipt)(nil)).Elem(), reflect.TypeOf((*TotalStakedReceiptDto)(nil)).Elem(), TotalStakedReceiptVersion, TotalStakedReceiptSize),
	OperationStartedReceiptEntityType:        MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), OperationStartedReceiptVersion, BalanceChangeReceiptSize),
	OperationEndedReceiptEntityType:          MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), OperationEndedReceiptVersion, BalanceChangeReceiptSize),
	OperationExpiredReceiptEntityType:        MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), OperationExpiredReceiptVersion, BalanceChangeReceiptSize),

	DriveStateReceiptEntityType:                MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveStateReceiptVersion, DriveReceiptSize),
	DriveDepositCreditReceiptEntityType:        MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveDepositCreditReceiptVersion, DriveReceiptSize),
	DriveDepositDebitReceiptEntityType:         MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveDepositDebitReceiptVersion, DriveReceiptSize),
	DriveRewardTransferCreditReceiptEntityType: MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveRewardTransferCreditReceiptVersion, DriveReceiptSize),
	DriveRewardTransferDebitReceiptEntityType:  MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveRewardTransferDebitReceiptVersion, DriveReceiptSize),
	DriveDownloadStartedReceiptEntityType:      MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveDownloadStartedReceiptVersion, DriveReceiptSize),
	DriveDownloadCompletedReceiptEntityType:    MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveDownloadCompletedReceiptVersion, DriveReceiptSize),
	DriveDownloadExpiredReceiptEntityType:      MakeDefinition(reflect.TypeOf((*BalanceChangeReceipt)(nil)).Elem(), reflect.TypeOf((*BalanceChangeReceiptDto)(nil)).Elem(), DriveDownloadExpiredReceiptVersion, DriveReceiptSize),
}

func GetReceiptHeader(entityType EntityType) ReceiptHeader {
	return MakeHeader(entityType, ReceiptDefinitionMap[entityType].Version, ReceiptDefinitionMap[entityType].Size)
}

const (
	BalanceTransferReceiptVersion         EntityVersion = 1
	NamespaceRentalReceiptVersion         EntityVersion = 1
	MosaicRentalReceiptVersion            EntityVersion = 1
	BalanceChangeReceiptVersion           EntityVersion = 1
	BalanceDebitReceiptVersion            EntityVersion = 1
	LockHashCreatedReceiptVersion         EntityVersion = 1
	LockHashCompletedReceiptVersion       EntityVersion = 1
	LockHashExpiredReceiptVersion         EntityVersion = 1
	LockSecretCreatedReceiptVersion       EntityVersion = 1
	LockSecretCompletedReceiptVersion     EntityVersion = 1
	LockSecretExpiredReceiptVersion       EntityVersion = 1
	MosaicArtifactExpiryReceiptVersion    EntityVersion = 1
	NamespaceArtifactExpiryReceiptVersion EntityVersion = 1
	InflationReceiptVersion               EntityVersion = 1
	SignerImportanceReceiptVersion        EntityVersion = 1
	GlobalStateTrackingReceiptVersion     EntityVersion = 1
	TotalStakedReceiptVersion             EntityVersion = 1
	OperationStartedReceiptVersion        EntityVersion = 1
	OperationEndedReceiptVersion          EntityVersion = 1
	OperationExpiredReceiptVersion        EntityVersion = 1

	DriveStateReceiptVersion                EntityVersion = 1
	DriveDepositCreditReceiptVersion        EntityVersion = 1
	DriveDepositDebitReceiptVersion         EntityVersion = 1
	DriveRewardTransferCreditReceiptVersion EntityVersion = 1
	DriveRewardTransferDebitReceiptVersion  EntityVersion = 1
	DriveDownloadStartedReceiptVersion      EntityVersion = 1
	DriveDownloadCompletedReceiptVersion    EntityVersion = 1
	DriveDownloadExpiredReceiptVersion      EntityVersion = 1
)

type IReceipt interface {
	ReceiptHeader() *ReceiptHeader
	ReceiptBody() interface{}
}

type Receipt struct {
	Header ReceiptHeader
	Body   interface{}
}

func (r *Receipt) ReceiptHeader() *ReceiptHeader {
	return &r.Header
}

func (r *Receipt) ReceiptBody() interface{} {
	return r.Body
}

func MakeHeader(entityType EntityType, version EntityVersion, size uint32) ReceiptHeader {
	return ReceiptHeader{
		Type:    entityType,
		Version: version,
		Size:    size,
	}
}

type AddressResolutionStatement struct {
	Height            Height
	UnresolvedAddress *Address
	ResolutionEntries []*Address
}

type MosaicResolutionStatement struct {
	Height            Height
	UnresolvedMosaic  *MosaicId
	ResolutionEntries []*MosaicId
}

type ReceiptStatement struct {
	Height   Height
	Receipts []*Receipt
}

type TransactionStatement ReceiptStatement
type PublicKeyStatement ReceiptStatement
type BlockchainStateStatement ReceiptStatement

type ReceiptHeader struct {
	Size    uint32
	Version EntityVersion
	Type    EntityType
}
type AnonymousReceipt struct {
	Header  ReceiptHeader
	Receipt []uint8
}

func (r *AnonymousReceipt) ReceiptHeader() *ReceiptHeader {
	return &r.Header
}

type BalanceChangeReceipt struct {
	Account  *PublicAccount
	MosaicId *MosaicId
	Amount   Amount
}

type BalanceDebitReceipt struct {
	Account  *PublicAccount
	MosaicId *MosaicId
	Amount   Amount
}

type BalanceTransferReceipt struct {
	Sender    *PublicAccount
	Recipient *Address
	MosaicId  *MosaicId
	Amount    Amount
}

type InflationReceipt struct {
	MosaicId *MosaicId
	Amount   Amount
}

type ArtifactExpiryReceipt struct {
	ArtifactId uint64
}

type DriveStateReceipt struct {
	DriveKey   *PublicAccount
	DriveState uint8
}

type SignerBalanceReceipt struct {
	Amount       Amount
	LockedAmount Amount
}

type GlobalStateChangeReceipt struct {
	Flags uint64
}

type TotalStakedReceipt struct {
	Amount Amount
}

type BlockStatement struct {
	TransactionStatements       []*TransactionStatement
	AddressResolutionStatements []*AddressResolutionStatement
	MosaicResolutionStatements  []*MosaicResolutionStatement
	PublicKeyStatements         []*PublicKeyStatement
	BlockchainStateStatements   []*BlockchainStateStatement
}
