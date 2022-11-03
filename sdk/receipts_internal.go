// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type AddressResolutionStatementDto struct {
	Height            uint64DTO
	UnresolvedAddress string
	ResolutionEntries []string
}

func (a *AddressResolutionStatementDto) toStruct(networkType NetworkType) (*AddressResolutionStatement, error) {
	address := NewAddress(a.UnresolvedAddress, networkType)
	entries := make([]*Address, len(a.ResolutionEntries))
	for i := 0; i < len(a.ResolutionEntries); i++ {
		entries[i] = NewAddress(a.ResolutionEntries[i], networkType)
	}
	return &AddressResolutionStatement{
		Height:            a.Height.toStruct(),
		UnresolvedAddress: address,
		ResolutionEntries: entries,
	}, nil
}

type MosaicResolutionStatementDto struct {
	Height            uint64DTO
	UnresolvedMosaic  uint64DTO
	ResolutionEntries []uint64DTO
}

func (a *MosaicResolutionStatementDto) toStruct(networkType NetworkType) (*MosaicResolutionStatement, error) {
	mosaicId, err := NewMosaicId(a.UnresolvedMosaic.toUint64())
	if err != nil {
		return nil, err
	}
	entries := make([]*MosaicId, len(a.ResolutionEntries))
	for i := 0; i < len(a.ResolutionEntries); i++ {
		entries[i], err = NewMosaicId(a.ResolutionEntries[i].toUint64())
		if err != nil {
			return nil, err
		}
	}
	return &MosaicResolutionStatement{
		Height:            a.Height.toStruct(),
		UnresolvedMosaic:  mosaicId,
		ResolutionEntries: entries,
	}, nil
}

type BalanceChangeReceiptDto struct {
	Account  string
	Amount   uint64DTO
	MosaicId uint64DTO
}

func (a *BalanceChangeReceiptDto) toStruct(networkType NetworkType) (*BalanceChangeReceipt, error) {
	account, err := NewAccountFromPublicKey(a.Account, networkType)
	if err != nil {
		return nil, err
	}
	mosaicId, err := NewMosaicId(a.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}
	return &BalanceChangeReceipt{
		Account:  account,
		MosaicId: mosaicId,
		Amount:   a.Amount.toStruct(),
	}, nil
}

type BalanceDebitReceiptDto struct {
	Account  string
	Amount   uint64DTO
	MosaicId uint64DTO
}

func (a *BalanceDebitReceiptDto) toStruct(networkType NetworkType) (*BalanceDebitReceipt, error) {
	account, err := NewAccountFromPublicKey(a.Account, networkType)
	if err != nil {
		return nil, err
	}
	mosaicId, err := NewMosaicId(a.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}
	return &BalanceDebitReceipt{
		Account:  account,
		MosaicId: mosaicId,
		Amount:   a.Amount.toStruct(),
	}, nil
}

type BalanceTransferReceiptDto struct {
	Sender    string
	Recipient string
	Amount    uint64DTO
	MosaicId  uint64DTO
}

func (a *BalanceTransferReceiptDto) toStruct(networkType NetworkType) (*BalanceTransferReceipt, error) {
	sender, err := NewAccountFromPublicKey(a.Sender, networkType)
	if err != nil {
		return nil, err
	}
	recipient := NewAddress(a.Recipient, networkType)
	mosaicId, err := NewMosaicId(a.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}
	return &BalanceTransferReceipt{
		Sender:    sender,
		Recipient: recipient,
		MosaicId:  mosaicId,
		Amount:    a.Amount.toStruct(),
	}, nil
}

type InflationReceiptDto struct {
	Amount   uint64DTO
	MosaicId uint64DTO
}

func (a *InflationReceiptDto) toStruct(networkType NetworkType) (*InflationReceipt, error) {
	mosaicId, err := NewMosaicId(a.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}
	return &InflationReceipt{
		MosaicId: mosaicId,
		Amount:   a.Amount.toStruct(),
	}, nil
}

type ArtifactExpiryReceiptDto struct {
	ArtifactId uint64DTO
}

func (a *ArtifactExpiryReceiptDto) toStruct(networkType NetworkType) (*ArtifactExpiryReceipt, error) {
	return &ArtifactExpiryReceipt{
		ArtifactId: a.ArtifactId.toUint64(),
	}, nil
}

type DriveStateReceiptDto struct {
	DriveKey   string
	DriveState uint8
}

func (a *DriveStateReceiptDto) toStruct(networkType NetworkType) (*DriveStateReceipt, error) {
	key, err := NewAccountFromPublicKey(a.DriveKey, networkType)
	if err != nil {
		return nil, err
	}
	return &DriveStateReceipt{
		DriveKey:   key,
		DriveState: a.DriveState,
	}, nil
}

type SignerBalanceReceiptDto struct {
	Amount       uint64DTO
	LockedAmount uint64DTO
}

func (a *SignerBalanceReceiptDto) toStruct(networkType NetworkType) (*SignerBalanceReceipt, error) {
	return &SignerBalanceReceipt{
		Amount:       a.Amount.toStruct(),
		LockedAmount: a.LockedAmount.toStruct(),
	}, nil
}

type GlobalStateChangeReceiptDto struct {
	Flags uint64DTO
}

func (a *GlobalStateChangeReceiptDto) toStruct(networkType NetworkType) (*GlobalStateChangeReceipt, error) {
	return &GlobalStateChangeReceipt{
		Flags: a.Flags.toUint64(),
	}, nil
}

type TotalStakedReceiptDto struct {
	TotalStaked uint64DTO
}

func (a *TotalStakedReceiptDto) toStruct(networkType NetworkType) (*TotalStakedReceipt, error) {
	return &TotalStakedReceipt{
		TotalStaked: a.TotalStaked.toStruct(),
	}, nil
}

type ReceiptBody interface {
	toStruct() interface{}
}
type ReceiptDto struct {
	Type    uint16
	Version uint32
	Body    interface{}
}

func (r *ReceiptDto) toStruct(networkType NetworkType) (*Receipt, error) {
	results := reflect.ValueOf(r.Body).MethodByName("toStruct").Call([]reflect.Value{reflect.ValueOf(networkType)})
	err := results[1].Interface()
	if err != nil {
		return nil, err.(error)
	}
	return &Receipt{
		Header: MakeHeader(EntityType(r.Type), EntityVersion(r.Version), ReceiptDefinitionMap[EntityType(r.Type)].Size),
		Body:   results[0].Interface(),
	}, nil
}

type ReceiptStatementDto struct {
	Height   uint64DTO
	Receipts []ReceiptDto
}

func (r *ReceiptStatementDto) toStruct(networkType NetworkType) (*ReceiptStatement, error) {
	entries := make([]*Receipt, len(r.Receipts))
	var err error
	for i := 0; i < len(r.Receipts); i++ {
		entries[i], err = r.Receipts[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	return &ReceiptStatement{
		Height:   r.Height.toStruct(),
		Receipts: entries,
	}, nil
}

type TransactionStatementDto ReceiptStatementDto

func (r *TransactionStatementDto) toStruct(networkType NetworkType) (*TransactionStatement, error) {
	res, err := (*ReceiptStatementDto)(r).toStruct(networkType)
	return (*TransactionStatement)(res), err
}

type PublicKeyStatementDto ReceiptStatementDto

func (r *PublicKeyStatementDto) toStruct(networkType NetworkType) (*PublicKeyStatement, error) {
	res, err := (*ReceiptStatementDto)(r).toStruct(networkType)
	return (*PublicKeyStatement)(res), err
}

type BlockchainStateStatementDto ReceiptStatementDto

func (r *BlockchainStateStatementDto) toStruct(networkType NetworkType) (*BlockchainStateStatement, error) {
	res, err := (*ReceiptStatementDto)(r).toStruct(networkType)
	return (*BlockchainStateStatement)(res), err
}

type BlockStatementDto struct {
	TransactionStatements       []TransactionStatementDto
	AddressResolutionStatements []AddressResolutionStatementDto
	MosaicResolutionStatements  []MosaicResolutionStatementDto
	PublicKeyStatements         []PublicKeyStatementDto
	BlockchainStateStatements   []BlockchainStateStatementDto
}

type anonymousReceiptDto struct {
	MetaDto struct {
		Size    uint32 `json:"size"`
		Version uint32 `json:"version"`
		Type    uint16 `json:"type"`
	} `json:"meta"`
	Receipt string `json:"receipt"`
}

func (dto *anonymousReceiptDto) toStruct() (*AnonymousReceipt, error) {
	data, err := hex.DecodeString(dto.Receipt)
	if err != nil {
		return nil, err
	}
	return &AnonymousReceipt{
		Header:  MakeHeader(EntityType(dto.MetaDto.Type), EntityVersion(dto.MetaDto.Version), dto.MetaDto.Size),
		Receipt: data,
	}, nil
}

func (r *BlockStatementDto) toStruct(networkType NetworkType) (*BlockStatement, error) {
	transactionStatements := make([]*TransactionStatement, len(r.TransactionStatements))
	var err error
	for i := 0; i < len(r.TransactionStatements); i++ {
		transactionStatements[i], err = r.TransactionStatements[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	addressResolutionStatements := make([]*AddressResolutionStatement, len(r.AddressResolutionStatements))
	for i := 0; i < len(r.AddressResolutionStatements); i++ {
		addressResolutionStatements[i], err = r.AddressResolutionStatements[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	mosaicResolutionStatements := make([]*MosaicResolutionStatement, len(r.MosaicResolutionStatements))
	for i := 0; i < len(r.MosaicResolutionStatements); i++ {
		mosaicResolutionStatements[i], err = r.MosaicResolutionStatements[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	publicKeyStatements := make([]*PublicKeyStatement, len(r.PublicKeyStatements))
	for i := 0; i < len(r.PublicKeyStatements); i++ {
		publicKeyStatements[i], err = r.PublicKeyStatements[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	blockchainStateStatements := make([]*BlockchainStateStatement, len(r.BlockchainStateStatements))
	for i := 0; i < len(r.BlockchainStateStatements); i++ {
		blockchainStateStatements[i], err = r.BlockchainStateStatements[i].toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}
	return &BlockStatement{
		TransactionStatements:       transactionStatements,
		AddressResolutionStatements: addressResolutionStatements,
		MosaicResolutionStatements:  mosaicResolutionStatements,
		PublicKeyStatements:         publicKeyStatements,
		BlockchainStateStatements:   blockchainStateStatements,
	}, nil
}

func (b *ReceiptDto) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	b.Type = uint16(v["type"].(float64))
	b.Version = uint32(v["version"].(float64))
	b.Body = reflect.New(ReceiptDefinitionMap[EntityType(b.Type)].DtoType)
	mapstructure.Decode(v, b.Body)
	return nil
}
