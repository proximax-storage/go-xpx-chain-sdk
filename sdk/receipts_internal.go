// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type AddressResolutionEntryDto struct {
	Resolved AddressDTO `json:"resolved"`
}

func (a *AddressResolutionEntryDto) toStruct() (*Address, error) {
	address, err := NewAddressFromRaw(a.Resolved.ToString())
	if err != nil {
		return nil, err
	}
	return address, nil
}

type MosaicResolutionEntryDto struct {
	Resolved uint64DTO `json:"resolved"`
}

func (a *MosaicResolutionEntryDto) toStruct() uint64 {
	return a.Resolved.toUint64()
}

type AddressResolutionStatementDto struct {
	Height            uint64DTO                   `json:"height"`
	UnresolvedAddress AddressDTO                  `json:"unresolved"`
	ResolutionEntries []AddressResolutionEntryDto `json:"resolutionEntries"`
}

func (a *AddressResolutionStatementDto) toStruct(networkType NetworkType) (*AddressResolutionStatement, error) {
	address, err := NewAddressFromRaw(a.UnresolvedAddress.ToString())
	if err != nil {
		return nil, err
	}
	entries := make([]*Address, len(a.ResolutionEntries))
	for i := 0; i < len(a.ResolutionEntries); i++ {
		entries[i], err = a.ResolutionEntries[i].toStruct()
		if err != nil {
			return nil, err
		}
	}
	return &AddressResolutionStatement{
		Height:            a.Height.toStruct(),
		UnresolvedAddress: address,
		ResolutionEntries: entries,
	}, nil
}

type MosaicResolutionStatementDto struct {
	Height            uint64DTO                  `json:"height"`
	UnresolvedMosaic  uint64DTO                  `json:"unresolved"`
	ResolutionEntries []MosaicResolutionEntryDto `json:"resolutionEntries"`
}

func (a *MosaicResolutionStatementDto) toStruct(networkType NetworkType) (*MosaicResolutionStatement, error) {
	mosaicId, err := NewMosaicId(a.UnresolvedMosaic.toUint64())
	if err != nil {
		return nil, err
	}
	entries := make([]*MosaicId, len(a.ResolutionEntries))
	for i := 0; i < len(a.ResolutionEntries); i++ {
		entries[i], err = NewMosaicId(a.ResolutionEntries[i].toStruct())
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
	Account  PublicKeyDTO `json:"account"`
	Amount   uint64DTO    `json:"amount"`
	MosaicId uint64DTO    `json:"mosaicId"`
}

func (a *BalanceChangeReceiptDto) ToStruct(networkType NetworkType) (*BalanceChangeReceipt, error) {
	account, err := NewAccountFromPublicKey(a.Account.ToString(), networkType)
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
	Account  PublicKeyDTO `json:"account"`
	Amount   uint64DTO    `json:"amount"`
	MosaicId uint64DTO    `json:"mosaicId"`
}

func (a *BalanceDebitReceiptDto) ToStruct(networkType NetworkType) (*BalanceDebitReceipt, error) {
	account, err := NewAccountFromPublicKey(a.Account.ToString(), networkType)
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
	Sender    PublicKeyDTO `json:"sender"`
	Recipient string       `json:"recipient"`
	Amount    uint64DTO    `json:"amount"`
	MosaicId  uint64DTO    `json:"mosaicId"`
}

func (a *BalanceTransferReceiptDto) ToStruct(networkType NetworkType) (*BalanceTransferReceipt, error) {
	sender, err := NewAccountFromPublicKey(a.Sender.ToString(), networkType)
	if err != nil {
		return nil, err
	}
	address, err := HexToBase32(a.Recipient)
	if err != nil {
		return nil, err
	}
	recipient, err := NewAddressFromRaw(*address)
	if err != nil {
		return nil, err
	}
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
	Amount   uint64DTO `json:"amount"`
	MosaicId uint64DTO `json:"mosaicId"`
}

func (a *InflationReceiptDto) ToStruct(networkType NetworkType) (*InflationReceipt, error) {
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
	ArtifactId uint64DTO `json:"artifactId"`
}

func (a *ArtifactExpiryReceiptDto) ToStruct(networkType NetworkType) (*ArtifactExpiryReceipt, error) {
	return &ArtifactExpiryReceipt{
		ArtifactId: a.ArtifactId.toUint64(),
	}, nil
}

type DriveStateReceiptDto struct {
	DriveKey   string `json:"driveKey"`
	DriveState uint8  `json:"driveState"`
}

func (a *DriveStateReceiptDto) ToStruct(networkType NetworkType) (*DriveStateReceipt, error) {
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
	Amount       uint64DTO `json:"amount"`
	LockedAmount uint64DTO `json:"lockedAmount"`
}

func (a *SignerBalanceReceiptDto) ToStruct(networkType NetworkType) (*SignerBalanceReceipt, error) {
	return &SignerBalanceReceipt{
		Amount:       a.Amount.toStruct(),
		LockedAmount: a.LockedAmount.toStruct(),
	}, nil
}

type GlobalStateChangeReceiptDto struct {
	Flags uint64DTO `json:"flags"`
}

func (a *GlobalStateChangeReceiptDto) ToStruct(networkType NetworkType) (*GlobalStateChangeReceipt, error) {
	return &GlobalStateChangeReceipt{
		Flags: a.Flags.toUint64(),
	}, nil
}

type TotalStakedReceiptDto struct {
	Amount uint64DTO `json:"amount"`
}

func (a *TotalStakedReceiptDto) ToStruct(networkType NetworkType) (*TotalStakedReceipt, error) {
	return &TotalStakedReceipt{
		Amount: a.Amount.toStruct(),
	}, nil
}

type ReceiptBody interface {
	toStruct() interface{}
}
type ReceiptDto struct {
	Type    uint16      `json:"type"`
	Version uint32      `json:"version"`
	Body    interface{} `json:"body"`
}

func (r *ReceiptDto) toStruct(networkType NetworkType) (*Receipt, error) {
	value := reflect.ValueOf(r.Body)
	results := value.MethodByName("ToStruct").Call([]reflect.Value{reflect.ValueOf(networkType)})
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
	Height   uint64DTO    `json:"height"`
	Receipts []ReceiptDto `json:"receipts"`
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
	b.Body = reflect.New(ReceiptDefinitionMap[EntityType(b.Type)].DtoType).Interface()
	mapstructure.Decode(v, b.Body)
	return nil
}
