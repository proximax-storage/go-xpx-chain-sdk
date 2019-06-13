// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	jsonLib "encoding/json"
	"errors"
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-catapult-sdk/transactions"
	"github.com/proximax-storage/go-xpx-utils"
	"github.com/proximax-storage/xpx-crypto-go"
	"strings"
	"sync"
)

type Transaction interface {
	GetAbstractTransaction() *AbstractTransaction
	String() string
	// number of bytes of serialized transaction
	Size() int
	generateBytes() ([]byte, error)
}

type transactionDto interface {
	toStruct() (Transaction, error)
}

type AbstractTransaction struct {
	*TransactionInfo
	NetworkType NetworkType
	Deadline    *Deadline
	Type        TransactionType
	Version     TransactionVersion
	MaxFee      Amount
	Signature   string
	Signer      *PublicAccount
}

func (tx *AbstractTransaction) IsUnconfirmed() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height == 0 && tx.TransactionInfo.Hash == tx.TransactionInfo.MerkleComponentHash
}

func (tx *AbstractTransaction) IsConfirmed() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height > 0
}

func (tx *AbstractTransaction) HasMissingSignatures() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height == 0 && tx.TransactionInfo.Hash != tx.TransactionInfo.MerkleComponentHash
}

func (tx *AbstractTransaction) IsUnannounced() bool {
	return tx.TransactionInfo == nil
}

func (tx *AbstractTransaction) ToAggregate(signer *PublicAccount) {
	tx.Signer = signer
}

func (tx *AbstractTransaction) String() string {
	return fmt.Sprintf(
		`
			"NetworkType": %s,
			"TransactionInfo": %s,
			"Type": %s,
			"Version": %d,
			"MaxFee": %s,
			"Deadline": %s,
			"Signature": %s,
			"Signer": %s
		`,
		tx.NetworkType,
		tx.TransactionInfo,
		tx.Type,
		tx.Version,
		tx.MaxFee,
		tx.Deadline,
		tx.Signature,
		tx.Signer,
	)
}

func (tx *AbstractTransaction) generateVectors(builder *flatbuffers.Builder) (v uint16, signatureV, signerV, dV, fV flatbuffers.UOffsetT, err error) {
	v = (uint16(tx.NetworkType) << 8) + uint16(tx.Version)
	signatureV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, SignatureSize))
	signerV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, SignerSize))
	dV = transactions.TransactionBufferCreateUint32Vector(builder, tx.Deadline.ToBlockchainTimestamp().toArray())
	fV = transactions.TransactionBufferCreateUint32Vector(builder, tx.MaxFee.toArray())
	return
}

func (tx *AbstractTransaction) buildVectors(builder *flatbuffers.Builder, v uint16, signatureV, signerV, dV, fV flatbuffers.UOffsetT) {
	transactions.TransactionBufferAddSignature(builder, signatureV)
	transactions.TransactionBufferAddSigner(builder, signerV)
	transactions.TransactionBufferAddVersion(builder, v)
	transactions.TransactionBufferAddType(builder, uint16(tx.Type))
	transactions.TransactionBufferAddMaxFee(builder, fV)
	transactions.TransactionBufferAddDeadline(builder, dV)
}

type abstractTransactionDTO struct {
	Type      TransactionType         `json:"type"`
	Version   uint64                  `json:"version"`
	MaxFee    *uint64DTO              `json:"maxFee"`
	Deadline  *blockchainTimestampDTO `json:"deadline"`
	Signature string                  `json:"signature"`
	Signer    string                  `json:"signer"`
}

func (dto *abstractTransactionDTO) toStruct(tInfo *TransactionInfo) (*AbstractTransaction, error) {
	nt := ExtractNetworkType(dto.Version)

	tv := TransactionVersion(ExtractVersion(dto.Version))

	pa, err := NewAccountFromPublicKey(dto.Signer, nt)
	if err != nil {
		return nil, err
	}

	var d *Deadline
	if dto.Deadline != nil {
		d = NewDeadlineFromBlockchainTimestamp(dto.Deadline.toStruct())
	}

	var f Amount
	if dto.MaxFee != nil {
		f = dto.MaxFee.toStruct()
	}

	return &AbstractTransaction{
		tInfo,
		nt,
		d,
		dto.Type,
		tv,
		f,
		dto.Signature,
		pa,
	}, nil
}

type TransactionInfo struct {
	Height              Height
	Index               uint32
	Id                  string
	Hash                Hash
	MerkleComponentHash Hash
	AggregateHash       Hash
	AggregateId         string
}

func (ti *TransactionInfo) String() string {
	return fmt.Sprintf(
		`
			"Height": %s,
			"Index": %d,
			"Id": %s,
			"Content": %s,
			"MerkleComponentHash:" %s,
			"AggregateHash": %s,
			"AggregateId": %s
		`,
		ti.Height,
		ti.Index,
		ti.Id,
		ti.Hash,
		ti.MerkleComponentHash,
		ti.AggregateHash,
		ti.AggregateId,
	)
}

type transactionInfoDTO struct {
	Height              uint64DTO `json:"height"`
	Index               uint32    `json:"index"`
	Id                  string    `json:"id"`
	Hash                Hash      `json:"hash"`
	MerkleComponentHash Hash      `json:"merkleComponentHash"`
	AggregateHash       Hash      `json:"aggregateHash,omitempty"`
	AggregateId         string    `json:"aggregateId,omitempty"`
}

func (dto *transactionInfoDTO) toStruct() *TransactionInfo {
	return &TransactionInfo{
		dto.Height.toStruct(),
		dto.Index,
		dto.Id,
		dto.Hash,
		dto.MerkleComponentHash,
		dto.AggregateHash,
		dto.AggregateId,
	}
}

type AccountPropertiesAddressModification struct {
	ModificationType PropertyModificationType
	Address          *Address
}

func (mod *AccountPropertiesAddressModification) String() string {
	return fmt.Sprintf(
		`
			"ModificationType": %d,
			"Address": %s,
		`,
		mod.ModificationType,
		mod.Address.Address,
	)
}

type AccountPropertiesAddressTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesAddressModification
}

// returns AccountPropertiesAddressTransaction from passed PropertyType and AccountPropertiesAddressModification's
func NewAccountPropertiesAddressTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesAddressModification, networkType NetworkType) (*AccountPropertiesAddressTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not be empty")
	}

	if propertyType&AllowAddress == 0 {
		return nil, errors.New("wrong propertyType for address account properties")
	}

	aptx := AccountPropertiesAddressTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountPropertyAddressVersion,
			Deadline:    deadline,
			Type:        AccountPropertyAddress,
			NetworkType: networkType,
		},
		PropertyType:  propertyType,
		Modifications: modifications,
	}

	return &aptx, nil
}

func (tx *AccountPropertiesAddressTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountPropertiesAddressTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"PropertyType": %d,
			"Modifications": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.PropertyType,
		tx.Modifications,
	)
}

func (tx *AccountPropertiesAddressTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	msb := make([]flatbuffers.UOffsetT, len(tx.Modifications))
	for i, m := range tx.Modifications {
		a, err := base32.StdEncoding.DecodeString(m.Address.Address)
		if err != nil {
			return nil, err
		}

		aV := transactions.TransactionBufferCreateByteVector(builder, a)

		transactions.PropertyModificationBufferStart(builder)
		transactions.PropertyModificationBufferAddModificationType(builder, uint8(m.ModificationType))
		transactions.PropertyModificationBufferAddValue(builder, aV)
		msb[i] = transactions.PropertyModificationBufferEnd(builder)
	}

	mV := transactions.TransactionBufferCreateUOffsetVector(builder, msb)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AccountPropertiesTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountPropertiesTransactionBufferAddPropertyType(builder, uint8(tx.PropertyType))
	transactions.AccountPropertiesTransactionBufferAddModificationCount(builder, uint8(len(tx.Modifications)))
	transactions.AccountPropertiesTransactionBufferAddModifications(builder, mV)
	t := transactions.AccountPropertiesTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountPropertyTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountPropertiesAddressTransaction) Size() int {
	return AccountPropertyAddressHeader + (AccountPropertiesAddressModificationSize * len(tx.Modifications))
}

type accountPropertiesAddressModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	Address          string                   `json:"value"`
}

func (dto *accountPropertiesAddressModificationDTO) toStruct() (*AccountPropertiesAddressModification, error) {
	a, err := NewAddressFromBase32(dto.Address)
	if err != nil {
		return nil, err
	}

	return &AccountPropertiesAddressModification{
		dto.ModificationType,
		a,
	}, nil
}

type accountPropertiesAddressTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                               `json:"propertyType"`
		Modifications []*accountPropertiesAddressModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesAddressTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesAddressModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesAddressTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type AccountPropertiesMosaicModification struct {
	ModificationType PropertyModificationType
	BlockchainId     BlockchainId
}

func (mod *AccountPropertiesMosaicModification) String() string {
	return fmt.Sprintf(
		`
			"ModificationType": %d,
			"BlockchainId": %s,
		`,
		mod.ModificationType,
		mod.BlockchainId,
	)
}

type AccountPropertiesMosaicTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesMosaicModification
}

// returns AccountPropertiesMosaicTransaction from passed PropertyType and AccountPropertiesMosaicModification's
func NewAccountPropertiesMosaicTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesMosaicModification, networkType NetworkType) (*AccountPropertiesMosaicTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not be empty")
	}

	if propertyType&AllowMosaic == 0 {
		return nil, errors.New("wrong propertyType for mosaic account properties")
	}

	aptx := AccountPropertiesMosaicTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountPropertyMosaicVersion,
			Deadline:    deadline,
			Type:        AccountPropertyMosaic,
			NetworkType: networkType,
		},
		PropertyType:  propertyType,
		Modifications: modifications,
	}

	return &aptx, nil
}

func (tx *AccountPropertiesMosaicTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountPropertiesMosaicTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"PropertyType": %d,
			"Modifications": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.PropertyType,
		tx.Modifications,
	)
}

func (tx *AccountPropertiesMosaicTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	msb := make([]flatbuffers.UOffsetT, len(tx.Modifications))
	for i, m := range tx.Modifications {
		mosaicB := make([]byte, MosaicSize)
		binary.LittleEndian.PutUint64(mosaicB, m.BlockchainId.Id())
		mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

		transactions.PropertyModificationBufferStart(builder)
		transactions.PropertyModificationBufferAddModificationType(builder, uint8(m.ModificationType))
		transactions.PropertyModificationBufferAddValue(builder, mV)
		msb[i] = transactions.PropertyModificationBufferEnd(builder)
	}

	mV := transactions.TransactionBufferCreateUOffsetVector(builder, msb)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AccountPropertiesTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountPropertiesTransactionBufferAddPropertyType(builder, uint8(tx.PropertyType))
	transactions.AccountPropertiesTransactionBufferAddModificationCount(builder, uint8(len(tx.Modifications)))
	transactions.AccountPropertiesTransactionBufferAddModifications(builder, mV)
	t := transactions.AccountPropertiesTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountPropertyTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountPropertiesMosaicTransaction) Size() int {
	return AccountPropertyMosaicHeader + (AccountPropertiesMosaicModificationSize * len(tx.Modifications))
}

type accountPropertiesMosaicModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	BlockchainId     blockchainIdDTO          `json:"value"`
}

func (dto *accountPropertiesMosaicModificationDTO) toStruct() (*AccountPropertiesMosaicModification, error) {
	blockchainId, err := dto.BlockchainId.toStruct()
	if err != nil {
		return nil, err
	}

	return &AccountPropertiesMosaicModification{
		dto.ModificationType,
		blockchainId,
	}, nil
}

type accountPropertiesMosaicTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                              `json:"propertyType"`
		Modifications []*accountPropertiesMosaicModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesMosaicTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesMosaicModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesMosaicTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type AccountPropertiesEntityTypeModification struct {
	ModificationType PropertyModificationType
	EntityType       TransactionType
}

func (mod *AccountPropertiesEntityTypeModification) String() string {
	return fmt.Sprintf(
		`
			"ModificationType": %d,
			"EntityType": %s,
		`,
		mod.ModificationType,
		mod.EntityType.String(),
	)
}

type AccountPropertiesEntityTypeTransaction struct {
	AbstractTransaction
	PropertyType  PropertyType
	Modifications []*AccountPropertiesEntityTypeModification
}

// returns AccountPropertiesEntityTypeTransaction from passed PropertyType and AccountPropertiesEntityTypeModification's
func NewAccountPropertiesEntityTypeTransaction(deadline *Deadline, propertyType PropertyType,
	modifications []*AccountPropertiesEntityTypeModification, networkType NetworkType) (*AccountPropertiesEntityTypeTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not be empty")
	}

	if propertyType&AllowTransaction == 0 {
		return nil, errors.New("wrong propertyType for entityType account properties")
	}

	aptx := AccountPropertiesEntityTypeTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountPropertyEntityTypeVersion,
			Deadline:    deadline,
			Type:        AccountPropertyEntityType,
			NetworkType: networkType,
		},
		PropertyType:  propertyType,
		Modifications: modifications,
	}

	return &aptx, nil
}

func (tx *AccountPropertiesEntityTypeTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountPropertiesEntityTypeTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"PropertyType": %d,
			"Modifications": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.PropertyType,
		tx.Modifications,
	)
}

func (tx *AccountPropertiesEntityTypeTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	msb := make([]flatbuffers.UOffsetT, len(tx.Modifications))
	for i, m := range tx.Modifications {
		typeB := make([]byte, 2)
		binary.LittleEndian.PutUint16(typeB, uint16(m.EntityType))
		mV := transactions.TransactionBufferCreateByteVector(builder, typeB)

		transactions.PropertyModificationBufferStart(builder)
		transactions.PropertyModificationBufferAddModificationType(builder, uint8(m.ModificationType))
		transactions.PropertyModificationBufferAddValue(builder, mV)
		msb[i] = transactions.PropertyModificationBufferEnd(builder)
	}

	mV := transactions.TransactionBufferCreateUOffsetVector(builder, msb)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AccountPropertiesTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountPropertiesTransactionBufferAddPropertyType(builder, uint8(tx.PropertyType))
	transactions.AccountPropertiesTransactionBufferAddModificationCount(builder, uint8(len(tx.Modifications)))
	transactions.AccountPropertiesTransactionBufferAddModifications(builder, mV)
	t := transactions.AccountPropertiesTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountPropertyTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountPropertiesEntityTypeTransaction) Size() int {
	return AccountPropertyEntityTypeHeader + (AccountPropertiesEntityModificationSize * len(tx.Modifications))
}

type accountPropertiesEntityTypeModificationDTO struct {
	ModificationType PropertyModificationType `json:"type"`
	EntityType       TransactionType          `json:"value"`
}

func (dto *accountPropertiesEntityTypeModificationDTO) toStruct() (*AccountPropertiesEntityTypeModification, error) {
	return &AccountPropertiesEntityTypeModification{
		dto.ModificationType,
		dto.EntityType,
	}, nil
}

type accountPropertiesEntityTypeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PropertyType  PropertyType                                  `json:"propertyType"`
		Modifications []*accountPropertiesEntityTypeModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountPropertiesEntityTypeTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	ms := make([]*AccountPropertiesEntityTypeModification, len(dto.Tx.Modifications))
	for i, m := range dto.Tx.Modifications {
		ms[i], err = m.toStruct()

		if err != nil {
			return nil, err
		}
	}

	return &AccountPropertiesEntityTypeTransaction{
		*atx,
		dto.Tx.PropertyType,
		ms,
	}, nil
}

type AliasTransaction struct {
	AbstractTransaction
	ActionType  AliasActionType
	NamespaceId *NamespaceId
}

func (tx *AliasTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"NamespaceId": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.NamespaceId.toHexString(),
	)
}

func (tx *AliasTransaction) generateBytes(builder *flatbuffers.Builder, aliasV flatbuffers.UOffsetT, sizeOfAlias int) ([]byte, error) {
	nV := transactions.TransactionBufferCreateUint32Vector(builder, tx.NamespaceId.toArray())

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AliasTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size()+sizeOfAlias)

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AliasTransactionBufferAddActionType(builder, uint8(tx.ActionType))
	transactions.AliasTransactionBufferAddNamespaceId(builder, nV)
	transactions.AliasTransactionBufferAddAliasId(builder, aliasV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return aliasTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AliasTransaction) Size() int {
	return AliasTransactionHeader
}

func (tx *AliasTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

type aliasTransactionDTO struct {
	abstractTransactionDTO
	NamespaceId namespaceIdDTO  `json:"namespaceId"`
	ActionType  AliasActionType `json:"action"`
}

func (dto *aliasTransactionDTO) toStruct(tInfo *TransactionInfo) (*AliasTransaction, error) {
	atx, err := dto.abstractTransactionDTO.toStruct(tInfo)
	if err != nil {
		return nil, err
	}

	namespaceId, err := dto.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &AliasTransaction{
		*atx,
		dto.ActionType,
		namespaceId,
	}, nil
}

type AddressAliasTransaction struct {
	AliasTransaction
	Address *Address
}

// returns AddressAliasTransaction from passed Address, NamespaceId and AliasActionType
func NewAddressAliasTransaction(deadline *Deadline, address *Address, namespaceId *NamespaceId, actionType AliasActionType, networkType NetworkType) (*AddressAliasTransaction, error) {
	if address == nil {
		return nil, errors.New("address must not be nil")
	}

	if namespaceId == nil {
		return nil, errors.New("namespaceId must not be nil")
	}

	aatx := AddressAliasTransaction{
		AliasTransaction: AliasTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     AddressAliasVersion,
				Deadline:    deadline,
				Type:        AddressAlias,
				NetworkType: networkType,
			},
			NamespaceId: namespaceId,
			ActionType:  actionType,
		},
		Address: address,
	}

	return &aatx, nil
}

func (tx *AddressAliasTransaction) String() string {
	return fmt.Sprintf(
		`
			"%s,
			"Address": %s,
		`,
		tx.AliasTransaction.String(),
		tx.Address,
	)
}

func (tx *AddressAliasTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	a, err := base32.StdEncoding.DecodeString(tx.Address.Address)
	if err != nil {
		return nil, err
	}

	aV := transactions.TransactionBufferCreateByteVector(builder, a)

	return tx.AliasTransaction.generateBytes(builder, aV, AddressSize)
}

func (tx *AddressAliasTransaction) Size() int {
	return tx.AliasTransaction.Size() + AddressSize
}

type addressAliasTransactionDTO struct {
	Tx struct {
		aliasTransactionDTO
		Address string `json:"address"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *addressAliasTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.aliasTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &AddressAliasTransaction{
		*atx,
		a,
	}, nil
}

type MosaicAliasTransaction struct {
	AliasTransaction
	MosaicId *MosaicId
}

// returns MosaicAliasTransaction from passed MosaicId, NamespaceId and AliasActionType
func NewMosaicAliasTransaction(deadline *Deadline, mosaicId *MosaicId, namespaceId *NamespaceId, actionType AliasActionType, networkType NetworkType) (*MosaicAliasTransaction, error) {
	if mosaicId == nil {
		return nil, errors.New("mosaicId must not bu nil")
	}

	if namespaceId == nil {
		return nil, errors.New("namespaceId must not bu nil")
	}

	matx := MosaicAliasTransaction{
		AliasTransaction: AliasTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     MosaicAliasVersion,
				Deadline:    deadline,
				Type:        MosaicAlias,
				NetworkType: networkType,
			},
			ActionType:  actionType,
			NamespaceId: namespaceId,
		},
		MosaicId: mosaicId,
	}

	return &matx, nil
}

func (tx *MosaicAliasTransaction) String() string {
	return fmt.Sprintf(
		`
			"%s,
			"MosaicId": %s,
		`,
		tx.AliasTransaction.String(),
		tx.MosaicId.toHexString(),
	)
}

func (tx *MosaicAliasTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, MosaicSize)
	binary.LittleEndian.PutUint64(mosaicB, tx.MosaicId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.AliasTransaction.generateBytes(builder, mV, MosaicSize)
}

func (tx *MosaicAliasTransaction) Size() int {
	return tx.AliasTransaction.Size() + MosaicSize
}

type mosaicAliasTransactionDTO struct {
	Tx struct {
		aliasTransactionDTO
		MosaicId *mosaicIdDTO `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicAliasTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.aliasTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicAliasTransaction{
		*atx,
		mosaicId,
	}, nil
}

type AccountLinkTransaction struct {
	AbstractTransaction
	RemoteAccount *PublicAccount
	LinkAction    AccountLinkAction
}

// returns AccountLinkTransaction from passed PublicAccount and AccountLinkAction
func NewAccountLinkTransaction(deadline *Deadline, remoteAccount *PublicAccount, linkAction AccountLinkAction, networkType NetworkType) (*AccountLinkTransaction, error) {
	if remoteAccount == nil {
		return nil, errors.New("remoteAccount must not be nil")
	}
	return &AccountLinkTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        LinkAccount,
			Version:     LinkAccountVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		RemoteAccount: remoteAccount,
		LinkAction:    linkAction,
	}, nil
}

func (tx *AccountLinkTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountLinkTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RemoteAccount": %s,
			"LinkAction": %d
		`,
		tx.AbstractTransaction.String(),
		tx.RemoteAccount.String(),
		tx.LinkAction,
	)
}

func (tx *AccountLinkTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	b, err := utils.HexDecodeStringOdd(tx.RemoteAccount.PublicKey)
	if err != nil {
		return nil, err
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, b)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AccountLinkTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)
	transactions.AccountLinkTransactionBufferAddRemoteAccountKey(builder, pV)
	transactions.AccountLinkTransactionBufferAddLinkAction(builder, uint8(tx.LinkAction))
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return accountLinkTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountLinkTransaction) Size() int {
	return AccountLinkTransactionSize
}

type accountLinkTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		RemoteAccountKey string            `json:"remoteAccountKey"`
		Action           AccountLinkAction `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *accountLinkTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.RemoteAccountKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &AccountLinkTransaction{
		*atx,
		acc,
		dto.Tx.Action,
	}, nil
}

type AggregateTransaction struct {
	AbstractTransaction
	InnerTransactions []Transaction
	Cosignatures      []*AggregateTransactionCosignature
}

// returns complete AggregateTransaction from passed array of own Transaction's to be included in
func NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransaction, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateCompleted,
			Version:     AggregateCompletedVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

// returns bounded AggregateTransaction from passed array of transactions to be included in
func NewBondedAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransaction, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateBonded,
			Version:     AggregateBondedVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

func (tx *AggregateTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AggregateTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"InnerTransactions": %s,
			"Cosignatures": %s
		`,
		tx.AbstractTransaction.String(),
		tx.InnerTransactions,
		tx.Cosignatures,
	)
}

func (tx *AggregateTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	var txsb []byte
	for _, itx := range tx.InnerTransactions {
		txb, err := toAggregateTransactionBytes(itx)
		if err != nil {
			return nil, err
		}
		txsb = append(txsb, txb...)
	}
	tV := transactions.TransactionBufferCreateByteVector(builder, txsb)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AggregateTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)
	transactions.AggregateTransactionBufferAddTransactionsSize(builder, uint32(len(txsb)))
	transactions.AggregateTransactionBufferAddTransactions(builder, tV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return aggregateTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AggregateTransaction) Size() int {
	sizeOfInnerTransactions := 0
	for _, itx := range tx.InnerTransactions {
		sizeOfInnerTransactions += itx.Size() - SignatureSize - MaxFeeSize - DeadLineSize
	}
	return AggregateBondedHeader + sizeOfInnerTransactions
}

type aggregateTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Cosignatures      []*aggregateTransactionCosignatureDTO `json:"cosignatures"`
		InnerTransactions []map[string]interface{}              `json:"transactions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *aggregateTransactionDTO) toStruct() (Transaction, error) {
	txsr, err := json.Marshal(dto.Tx.InnerTransactions)
	if err != nil {
		return nil, err
	}

	txs, err := MapTransactions(bytes.NewBuffer(txsr))
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	as := make([]*AggregateTransactionCosignature, len(dto.Tx.Cosignatures))
	for i, a := range dto.Tx.Cosignatures {
		as[i], err = a.toStruct(atx.NetworkType)
	}
	if err != nil {
		return nil, err
	}

	for _, tx := range txs {
		iatx := tx.GetAbstractTransaction()
		iatx.Deadline = atx.Deadline
		iatx.Signature = atx.Signature
		iatx.MaxFee = atx.MaxFee
		if iatx.TransactionInfo == nil {
			iatx.TransactionInfo = atx.TransactionInfo
		}
	}

	return &AggregateTransaction{
		*atx,
		txs,
		as,
	}, nil
}

type ModifyMetadataTransaction struct {
	AbstractTransaction
	MetadataType  MetadataType
	Modifications []*MetadataModification
}

func (tx *ModifyMetadataTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MetadataType": %s,
			"Modifications": %s
		`,
		tx.AbstractTransaction.String(),
		tx.MetadataType,
		tx.Modifications,
	)
}

func (tx *ModifyMetadataTransaction) generateBytes(builder *flatbuffers.Builder, metadataV flatbuffers.UOffsetT, sizeOfMetadata int) ([]byte, error) {

	mV, err := metadataModificationArrayToBuffer(builder, tx.Modifications)
	if err != nil {
		return nil, err
	}

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.ModifyMetadataTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size()+sizeOfMetadata)

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ModifyMetadataTransactionBufferAddMetadataType(builder, uint8(tx.MetadataType))
	transactions.ModifyMetadataTransactionBufferAddMetadataId(builder, metadataV)
	transactions.ModifyMetadataTransactionBufferAddModifications(builder, mV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyMetadataTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ModifyMetadataTransaction) Size() int {
	sizeOfModifications := 0
	for _, m := range tx.Modifications {
		sizeOfModifications += m.Size()
	}
	return MetadataHeaderSize + sizeOfModifications
}

func (tx *ModifyMetadataTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

type modifyMetadataTransactionDTO struct {
	abstractTransactionDTO
	MetadataType  MetadataType               `json:"metadataType"`
	Modifications []*metadataModificationDTO `json:"modifications"`
}

func (dto *modifyMetadataTransactionDTO) toStruct(tInfo *TransactionInfo) (*ModifyMetadataTransaction, error) {
	atx, err := dto.abstractTransactionDTO.toStruct(tInfo)
	if err != nil {
		return nil, err
	}

	ms, err := metadataDTOArrayToStruct(dto.Modifications, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataTransaction{
		*atx,
		dto.MetadataType,
		ms,
	}, nil
}

type ModifyMetadataAddressTransaction struct {
	ModifyMetadataTransaction
	Address *Address
}

// returns ModifyMetadataAddressTransaction from passed Address to be modified, and an array of MetadataModification's
func NewModifyMetadataAddressTransaction(deadline *Deadline, address *Address, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataAddressTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not empty")
	}

	mmatx := ModifyMetadataAddressTransaction{
		ModifyMetadataTransaction: ModifyMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     MetadataAddressVersion,
				Deadline:    deadline,
				Type:        MetadataAddress,
				NetworkType: networkType,
			},
			MetadataType:  MetadataAddressType,
			Modifications: modifications,
		},
		Address: address,
	}

	return &mmatx, nil
}

func (tx *ModifyMetadataAddressTransaction) String() string {
	return fmt.Sprintf(
		`
			"%s,
			"Address": %s,
		`,
		tx.ModifyMetadataTransaction.String(),
		tx.Address,
	)
}

func (tx *ModifyMetadataAddressTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	a, err := base32.StdEncoding.DecodeString(tx.Address.Address)
	if err != nil {
		return nil, err
	}

	aV := transactions.TransactionBufferCreateByteVector(builder, a)

	return tx.ModifyMetadataTransaction.generateBytes(builder, aV, AddressSize)
}

func (tx *ModifyMetadataAddressTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + AddressSize
}

type modifyMetadataAddressTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		Address string `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataAddressTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataAddressTransaction{
		*atx,
		a,
	}, nil
}

type ModifyMetadataMosaicTransaction struct {
	ModifyMetadataTransaction
	MosaicId *MosaicId
}

// returns ModifyMetadataMosaicTransaction from passed MosaicId to be modified, and an array of MetadataModification's
func NewModifyMetadataMosaicTransaction(deadline *Deadline, mosaicId *MosaicId, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataMosaicTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not empty")
	}

	mmatx := ModifyMetadataMosaicTransaction{
		ModifyMetadataTransaction: ModifyMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     MetadataMosaicVersion,
				Deadline:    deadline,
				Type:        MetadataMosaic,
				NetworkType: networkType,
			},
			MetadataType:  MetadataMosaicType,
			Modifications: modifications,
		},
		MosaicId: mosaicId,
	}

	return &mmatx, nil
}

func (tx *ModifyMetadataMosaicTransaction) String() string {
	return fmt.Sprintf(
		`
			"%s,
			"MosaicId": %s,
		`,
		tx.ModifyMetadataTransaction.String(),
		tx.MosaicId,
	)
}

func (tx *ModifyMetadataMosaicTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, MosaicSize)
	binary.LittleEndian.PutUint64(mosaicB, tx.MosaicId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.ModifyMetadataTransaction.generateBytes(builder, mV, MosaicSize)
}

func (tx *ModifyMetadataMosaicTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + MosaicSize
}

type modifyMetadataMosaicTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		MosaicId *mosaicIdDTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataMosaicTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataMosaicTransaction{
		*atx,
		mosaicId,
	}, nil
}

type ModifyMetadataNamespaceTransaction struct {
	ModifyMetadataTransaction
	NamespaceId *NamespaceId
}

// returns ModifyMetadataNamespaceTransaction from passed NamespaceId to be modified, and an array of MetadataModification's
func NewModifyMetadataNamespaceTransaction(deadline *Deadline, namespaceId *NamespaceId, modifications []*MetadataModification, networkType NetworkType) (*ModifyMetadataNamespaceTransaction, error) {
	if len(modifications) == 0 {
		return nil, errors.New("modifications must not empty")
	}

	mmatx := ModifyMetadataNamespaceTransaction{
		ModifyMetadataTransaction: ModifyMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     MetadataNamespaceVersion,
				Deadline:    deadline,
				Type:        MetadataNamespace,
				NetworkType: networkType,
			},
			MetadataType:  MetadataNamespaceType,
			Modifications: modifications,
		},
		NamespaceId: namespaceId,
	}

	return &mmatx, nil
}

func (tx *ModifyMetadataNamespaceTransaction) String() string {
	return fmt.Sprintf(
		`
			"%s,
			"NamespaceId": %s,
		`,
		tx.ModifyMetadataTransaction.String(),
		tx.NamespaceId,
	)
}

func (tx *ModifyMetadataNamespaceTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	namespaceB := make([]byte, NamespaceSize)
	binary.LittleEndian.PutUint64(namespaceB, tx.NamespaceId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, namespaceB)

	return tx.ModifyMetadataTransaction.generateBytes(builder, mV, NamespaceSize)
}

func (tx *ModifyMetadataNamespaceTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + NamespaceSize
}

type modifyMetadataNamespaceTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		NamespaceId *namespaceIdDTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataNamespaceTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	namespaceId, err := dto.Tx.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataNamespaceTransaction{
		*atx,
		namespaceId,
	}, nil
}

type MosaicDefinitionTransaction struct {
	AbstractTransaction
	*MosaicProperties
	MosaicNonce uint32
	*MosaicId
}

// returns MosaicDefinitionTransaction from passed nonce, public key of announcer and MosaicProperties
func NewMosaicDefinitionTransaction(deadline *Deadline, nonce uint32, ownerPublicKey string, mosaicProps *MosaicProperties, networkType NetworkType) (*MosaicDefinitionTransaction, error) {
	if len(ownerPublicKey) != 64 {
		return nil, ErrInvalidOwnerPublicKey
	}

	if mosaicProps == nil {
		return nil, ErrNilMosaicProperties
	}

	// Signer of transaction must be the same with ownerPublicKey
	mosaicId, err := NewMosaicIdFromNonceAndOwner(nonce, ownerPublicKey)
	if err != nil {
		return nil, err
	}

	return &MosaicDefinitionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     MosaicDefinitionVersion,
			Deadline:    deadline,
			Type:        MosaicDefinition,
			NetworkType: networkType,
		},
		MosaicProperties: mosaicProps,
		MosaicNonce:      nonce,
		MosaicId:         mosaicId,
	}, nil
}

func (tx *MosaicDefinitionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *MosaicDefinitionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicProperties": %s,
			"MosaicNonce": %d,
			"MosaicId": %s
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicProperties,
		tx.MosaicNonce,
		tx.MosaicId,
	)
}

func (tx *MosaicDefinitionTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	var f uint8 = 0
	if tx.MosaicProperties.SupplyMutable {
		f += 1
	}
	if tx.MosaicProperties.Transferable {
		f += 2
	}
	if tx.MosaicProperties.LevyMutable {
		f += 4
	}

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicProperties.Duration.toArray())

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MosaicDefinitionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicDefinitionTransactionBufferAddMosaicNonce(builder, tx.MosaicNonce)
	transactions.MosaicDefinitionTransactionBufferAddMosaicId(builder, mV)
	transactions.MosaicDefinitionTransactionBufferAddNumOptionalProperties(builder, 1)
	transactions.MosaicDefinitionTransactionBufferAddFlags(builder, f)
	transactions.MosaicDefinitionTransactionBufferAddDivisibility(builder, tx.MosaicProperties.Divisibility)
	transactions.MosaicDefinitionTransactionBufferAddIndicateDuration(builder, 2)
	transactions.MosaicDefinitionTransactionBufferAddDuration(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)
	return mosaicDefinitionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicDefinitionTransaction) Size() int {
	return MosaicDefinitionTransactionSize
}

type mosaicDefinitionTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Properties  mosaicDefinitonTransactionPropertiesDTO `json:"properties"`
		MosaicNonce int32                                   `json:"mosaicNonce"`
		MosaicId    *mosaicIdDTO                            `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicDefinitionTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicDefinitionTransaction{
		*atx,
		dto.Tx.Properties.toStruct(),
		uint32(dto.Tx.MosaicNonce),
		mosaicId,
	}, nil
}

type MosaicSupplyChangeTransaction struct {
	AbstractTransaction
	MosaicSupplyType
	BlockchainId
	Delta Amount
}

// returns MosaicSupplyChangeTransaction from passed BlockchainId, MosaicSupplyTypeand supply delta
func NewMosaicSupplyChangeTransaction(deadline *Deadline, blockchainId BlockchainId, supplyType MosaicSupplyType, delta Duration, networkType NetworkType) (*MosaicSupplyChangeTransaction, error) {
	if blockchainId == nil || blockchainId.Id() == 0 {
		return nil, ErrNilBlockchainId
	}

	if !(supplyType == Increase || supplyType == Decrease) {
		return nil, errors.New("supplyType must not be nil")
	}

	return &MosaicSupplyChangeTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     MosaicSupplyChangeVersion,
			Deadline:    deadline,
			Type:        MosaicSupplyChange,
			NetworkType: networkType,
		},
		BlockchainId:     blockchainId,
		MosaicSupplyType: supplyType,
		Delta:            delta,
	}, nil
}

func (tx *MosaicSupplyChangeTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *MosaicSupplyChangeTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicSupplyType": %s,
			"BlockchainId": %s,
			"Delta": %d
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicSupplyType,
		tx.BlockchainId,
		tx.Delta,
	)
}

func (tx *MosaicSupplyChangeTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.BlockchainId.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Delta.toArray())

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MosaicSupplyChangeTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicSupplyChangeTransactionBufferAddMosaicId(builder, mV)
	transactions.MosaicSupplyChangeTransactionBufferAddDirection(builder, uint8(tx.MosaicSupplyType))
	transactions.MosaicSupplyChangeTransactionBufferAddDelta(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return mosaicSupplyChangeTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicSupplyChangeTransaction) Size() int {
	return MosaicSupplyChangeTransactionSize
}

type mosaicSupplyChangeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicSupplyType `json:"direction"`
		BlockchainId     *blockchainIdDTO `json:"mosaicId"`
		Delta            uint64DTO        `json:"delta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicSupplyChangeTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	blockchainId, err := dto.Tx.BlockchainId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicSupplyChangeTransaction{
		*atx,
		dto.Tx.MosaicSupplyType,
		blockchainId,
		dto.Tx.Delta.toStruct(),
	}, nil
}

type TransferTransaction struct {
	AbstractTransaction
	Message   Message
	Mosaics   []*Mosaic
	Recipient *Address
}

// returns a TransferTransaction from passed transfer recipient Adderess, array of Mosaic's to transfer and transfer Message
func NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message Message, networkType NetworkType) (*TransferTransaction, error) {
	if recipient == nil {
		return nil, errors.New("recipient must not be nil")
	}
	if mosaics == nil {
		return nil, errors.New("mosaics must not be nil")
	}
	if message == nil {
		return nil, errors.New("message must not be nil, but could be with empty payload")
	}

	return &TransferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     TransferVersion,
			Deadline:    deadline,
			Type:        Transfer,
			NetworkType: networkType,
		},
		Recipient: recipient,
		Mosaics:   mosaics,
		Message:   message,
	}, nil
}

// returns TransferTransaction from passed recipient NamespaceId, Mosaic's and transfer Message
func NewTransferTransactionWithNamespace(deadline *Deadline, recipient *NamespaceId, mosaics []*Mosaic, message Message, networkType NetworkType) (*TransferTransaction, error) {
	if recipient == nil {
		return nil, errors.New("recipient namespace must not be nil")
	}
	if mosaics == nil {
		return nil, errors.New("mosaics must not be nil")
	}
	if message == nil {
		return nil, errors.New("message must not be nil, but could be with empty payload")
	}

	address, err := NewAddressFromNamespace(recipient)
	if err != nil {
		return nil, err
	}

	return &TransferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     TransferVersion,
			Deadline:    deadline,
			Type:        Transfer,
			NetworkType: networkType,
		},
		Recipient: address,
		Mosaics:   mosaics,
		Message:   message,
	}, nil
}

func (tx *TransferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *TransferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Mosaics": %s,
			"Address": %s,
			"Message": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaics,
		tx.Recipient,
		tx.Message,
	)
}

func (tx *TransferTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	ml := len(tx.Mosaics)
	mb := make([]flatbuffers.UOffsetT, ml)
	for i, mos := range tx.Mosaics {
		id := transactions.TransactionBufferCreateUint32Vector(builder, mos.BlockchainId.toArray())
		am := transactions.TransactionBufferCreateUint32Vector(builder, mos.Amount.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}

	mp := transactions.TransactionBufferCreateByteVector(builder, tx.Message.Payload())
	transactions.MessageBufferStart(builder)
	transactions.MessageBufferAddType(builder, uint8(tx.Message.Type()))
	transactions.MessageBufferAddPayload(builder, mp)
	m := transactions.TransactionBufferEnd(builder)

	r, err := base32.StdEncoding.DecodeString(tx.Recipient.Address)
	if err != nil {
		return nil, err
	}

	rV := transactions.TransactionBufferCreateByteVector(builder, r)
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.TransferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.TransferTransactionBufferAddRecipient(builder, rV)
	transactions.TransferTransactionBufferAddNumMosaics(builder, uint8(ml))
	transactions.TransferTransactionBufferAddMessageSize(builder, uint16(tx.MessageSize()))
	transactions.TransferTransactionBufferAddMessage(builder, m)
	transactions.TransferTransactionBufferAddMosaics(builder, mV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return transferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *TransferTransaction) Size() int {
	return TransferHeaderSize + ((MosaicSize + AmountSize) * len(tx.Mosaics)) + tx.MessageSize()
}

func (tx *TransferTransaction) MessageSize() int {
	// Message + MessageType
	return len(tx.Message.Payload()) + 1
}

type transferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Message messageDTO   `json:"message"`
		Mosaics []*mosaicDTO `json:"mosaics"`
		Address string       `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *transferTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.Mosaics))

	for i, mosaic := range dto.Tx.Mosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	a, err := NewAddressFromBase32(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	m, err := dto.Tx.Message.toStruct()
	if err != nil {
		return nil, err
	}

	return &TransferTransaction{
		*atx,
		m,
		mosaics,
		a,
	}, nil
}

type ModifyMultisigAccountTransaction struct {
	AbstractTransaction
	MinApprovalDelta int8
	MinRemovalDelta  int8
	Modifications    []*MultisigCosignatoryModification
}

// returns a ModifyMultisigAccountTransaction from passed min approval and removal deltas and array of MultisigCosignatoryModification's
func NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta int8, minRemovalDelta int8, modifications []*MultisigCosignatoryModification, networkType NetworkType) (*ModifyMultisigAccountTransaction, error) {
	if len(modifications) == 0 && minApprovalDelta == 0 && minRemovalDelta == 0 {
		return nil, errors.New("modifications must not empty")
	}

	mmatx := ModifyMultisigAccountTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ModifyMultisigVersion,
			Deadline:    deadline,
			Type:        ModifyMultisig,
			NetworkType: networkType,
		},
		MinRemovalDelta:  minRemovalDelta,
		MinApprovalDelta: minApprovalDelta,
		Modifications:    modifications,
	}

	return &mmatx, nil
}

func (tx *ModifyMultisigAccountTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ModifyMultisigAccountTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MinApprovalDelta": %d,
			"MinRemovalDelta": %d,
			"Modifications": %s
		`,
		tx.AbstractTransaction.String(),
		tx.MinApprovalDelta,
		tx.MinRemovalDelta,
		tx.Modifications,
	)
}

func (tx *ModifyMultisigAccountTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV, err := cosignatoryModificationArrayToBuffer(builder, tx.Modifications)
	if err != nil {
		return nil, err
	}

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.ModifyMultisigAccountTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ModifyMultisigAccountTransactionBufferAddMinRemovalDelta(builder, tx.MinRemovalDelta)
	transactions.ModifyMultisigAccountTransactionBufferAddMinApprovalDelta(builder, tx.MinApprovalDelta)
	transactions.ModifyMultisigAccountTransactionBufferAddNumModifications(builder, uint8(len(tx.Modifications)))
	transactions.ModifyMultisigAccountTransactionBufferAddModifications(builder, mV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyMultisigAccountTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ModifyMultisigAccountTransaction) Size() int {
	return ModifyMultisigHeaderSize + ((KeySize + 1 /* MultisigModificationType size */) * len(tx.Modifications))
}

type modifyMultisigAccountTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MinApprovalDelta int8                                  `json:"minApprovalDelta"`
		MinRemovalDelta  int8                                  `json:"minRemovalDelta"`
		Modifications    []*multisigCosignatoryModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMultisigAccountTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	ms, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Modifications, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyMultisigAccountTransaction{
		*atx,
		dto.Tx.MinApprovalDelta,
		dto.Tx.MinRemovalDelta,
		ms,
	}, nil
}

type ModifyContractTransaction struct {
	AbstractTransaction
	DurationDelta Duration
	Hash          string
	Customers     []*MultisigCosignatoryModification
	Executors     []*MultisigCosignatoryModification
	Verifiers     []*MultisigCosignatoryModification
}

// returns ModifyContractTransaction from passed duration delta in blocks, file hash, arrays of customers, replicators and verificators
func NewModifyContractTransaction(
	deadline *Deadline, durationDelta Duration, hash string,
	customers []*MultisigCosignatoryModification,
	executors []*MultisigCosignatoryModification,
	verifiers []*MultisigCosignatoryModification,
	networkType NetworkType) (*ModifyContractTransaction, error) {

	if len(customers) == 0 {
		return nil, errors.New("customers must not empty")
	}
	if len(executors) == 0 {
		return nil, errors.New("executors must not empty")
	}
	if len(verifiers) == 0 {
		return nil, errors.New("verifiers must not empty")
	}

	mctx := ModifyContractTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ModifyContractVersion,
			Deadline:    deadline,
			Type:        ModifyContract,
			NetworkType: networkType,
		},
		DurationDelta: durationDelta,
		Hash:          hash,
		Customers:     customers,
		Executors:     executors,
		Verifiers:     verifiers,
	}

	return &mctx, nil
}

func (tx *ModifyContractTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ModifyContractTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DurationDelta": %d,
			"Content": %s,
			"Customers": %s,
			"Executors": %s,
			"Verifiers": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DurationDelta,
		tx.Hash,
		tx.Customers,
		tx.Executors,
		tx.Verifiers,
	)
}

func (tx *ModifyContractTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	durationV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DurationDelta.toArray())
	hashV := stringToBuffer(builder, tx.Hash)

	customersV, err := cosignatoryModificationArrayToBuffer(builder, tx.Customers)
	if err != nil {
		return nil, err
	}

	executorsV, err := cosignatoryModificationArrayToBuffer(builder, tx.Executors)
	if err != nil {
		return nil, err
	}

	verifiersV, err := cosignatoryModificationArrayToBuffer(builder, tx.Verifiers)
	if err != nil {
		return nil, err
	}

	transactions.ModifyContractTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ModifyContractTransactionBufferAddDurationDelta(builder, durationV)
	transactions.ModifyContractTransactionBufferAddHash(builder, hashV)

	transactions.ModifyContractTransactionBufferAddNumCustomers(builder, uint8(len(tx.Customers)))
	transactions.ModifyContractTransactionBufferAddNumExecutors(builder, uint8(len(tx.Executors)))
	transactions.ModifyContractTransactionBufferAddNumVerifiers(builder, uint8(len(tx.Verifiers)))
	transactions.ModifyContractTransactionBufferAddCustomers(builder, customersV)
	transactions.ModifyContractTransactionBufferAddExecutors(builder, executorsV)
	transactions.ModifyContractTransactionBufferAddVerifiers(builder, verifiersV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyContractTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ModifyContractTransaction) Size() int {
	return ModifyContractHeaderSize + ((KeySize + 1 /* ContractModificationType size */) * (len(tx.Customers) + len(tx.Executors) + len(tx.Verifiers)))
}

type modifyContractTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DurationDelta uint64DTO                             `json:"durationDelta"`
		Hash          string                                `json:"hash"`
		Customers     []*multisigCosignatoryModificationDTO `json:"customers"`
		Executors     []*multisigCosignatoryModificationDTO `json:"executors"`
		Verifiers     []*multisigCosignatoryModificationDTO `json:"verifiers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyContractTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	customers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Customers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	executors, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Executors, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	verifiers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Verifiers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyContractTransaction{
		*atx,
		dto.Tx.DurationDelta.toStruct(),
		dto.Tx.Hash,
		customers,
		executors,
		verifiers,
	}, nil
}

type RegisterNamespaceTransaction struct {
	AbstractTransaction
	*NamespaceId
	NamespaceType
	NamspaceName string
	Duration     Duration
	ParentId     *NamespaceId
}

// returns a RegisterNamespaceTransaction from passed namespace name and duration in blocks
func NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration Duration, networkType NetworkType) (*RegisterNamespaceTransaction, error) {
	if len(namespaceName) == 0 {
		return nil, ErrInvalidNamespaceName
	}

	nsId, err := NewNamespaceIdFromName(namespaceName)
	if err != nil {
		return nil, err
	}

	return &RegisterNamespaceTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     RegisterNamespaceVersion,
			Deadline:    deadline,
			Type:        RegisterNamespace,
			NetworkType: networkType,
		},
		NamspaceName:  namespaceName,
		NamespaceId:   nsId,
		NamespaceType: Root,
		Duration:      duration,
	}, nil
}

// returns a RegisterNamespaceTransaction from passed namespace name and parent NamespaceId
func NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId, networkType NetworkType) (*RegisterNamespaceTransaction, error) {
	if len(namespaceName) == 0 {
		return nil, ErrInvalidNamespaceName
	}

	if parentId == nil || parentId.Id() == 0 {
		return nil, ErrNilNamespaceId
	}

	nsId, err := generateNamespaceId(namespaceName, parentId)
	if err != nil {
		return nil, err
	}

	return &RegisterNamespaceTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     RegisterNamespaceVersion,
			Deadline:    deadline,
			Type:        RegisterNamespace,
			NetworkType: networkType,
		},
		NamspaceName:  namespaceName,
		NamespaceId:   nsId,
		NamespaceType: Sub,
		ParentId:      parentId,
	}, nil
}

func (tx *RegisterNamespaceTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *RegisterNamespaceTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"NamespaceName": %s,
			"NamespaceId": %s,
			"Duration": %d
		`,
		tx.AbstractTransaction.String(),
		tx.NamspaceName,
		tx.NamespaceId,
		tx.Duration,
	)
}

func (tx *RegisterNamespaceTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	nV := transactions.TransactionBufferCreateUint32Vector(builder, tx.NamespaceId.toArray())
	var dV flatbuffers.UOffsetT
	if tx.NamespaceType == Root {
		dV = transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())
	} else {
		dV = transactions.TransactionBufferCreateUint32Vector(builder, tx.ParentId.toArray())
	}
	n := builder.CreateString(tx.NamspaceName)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.RegisterNamespaceTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.RegisterNamespaceTransactionBufferAddNamespaceType(builder, uint8(tx.NamespaceType))
	transactions.RegisterNamespaceTransactionBufferAddDurationParentId(builder, dV)
	transactions.RegisterNamespaceTransactionBufferAddNamespaceId(builder, nV)
	transactions.RegisterNamespaceTransactionBufferAddNamespaceNameSize(builder, byte(len(tx.NamspaceName)))
	transactions.RegisterNamespaceTransactionBufferAddNamespaceName(builder, n)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return registerNamespaceTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *RegisterNamespaceTransaction) Size() int {
	return RegisterNamespaceHeaderSize + len(tx.NamspaceName)
}

type registerNamespaceTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Id            namespaceIdDTO `json:"namespaceId"`
		NamespaceType `json:"namespaceType"`
		NamspaceName  string    `json:"name"`
		Duration      uint64DTO `json:"duration"`
		ParentId      namespaceIdDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *registerNamespaceTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	d := Duration(0)
	n := NewNamespaceIdNoCheck(0)

	if dto.Tx.NamespaceType == Root {
		d = dto.Tx.Duration.toStruct()
	} else {
		n, err = dto.Tx.ParentId.toStruct()
		if err != nil {
			return nil, err
		}
	}

	nsId, err := dto.Tx.Id.toStruct()
	if err != nil {
		return nil, err
	}

	return &RegisterNamespaceTransaction{
		*atx,
		nsId,
		dto.Tx.NamespaceType,
		dto.Tx.NamspaceName,
		d,
		n,
	}, nil
}

type LockFundsTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration Duration
	*SignedTransaction
}

// returns a LockFundsTransaction from passed Mosaic, duration in blocks and SignedTransaction
func NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, signedTx *SignedTransaction, networkType NetworkType) (*LockFundsTransaction, error) {
	if mosaic == nil {
		return nil, errors.New("mosaic must not be nil")
	}

	if signedTx == nil {
		return nil, errors.New("signedTx must not be nil")
	}

	if signedTx.TransactionType != AggregateBonded {
		return nil, errors.New("signedTx must be of type AggregateBonded")
	}

	return &LockFundsTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     LockVersion,
			Deadline:    deadline,
			Type:        Lock,
			NetworkType: networkType,
		},
		Mosaic:            mosaic,
		Duration:          duration,
		SignedTransaction: signedTx,
	}, nil
}

func (tx *LockFundsTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *LockFundsTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Mosaic": %s,
			"Duration": %d,
			"SignedTxHash": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaic,
		tx.Duration,
		tx.SignedTransaction.Hash,
	)
}

func (tx *LockFundsTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mv := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.BlockchainId.toArray())
	maV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.Amount.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())

	h, err := hex.DecodeString((string)(tx.SignedTransaction.Hash))
	if err != nil {
		return nil, err
	}
	hV := transactions.TransactionBufferCreateByteVector(builder, h)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.LockFundsTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.LockFundsTransactionBufferAddMosaicId(builder, mv)
	transactions.LockFundsTransactionBufferAddMosaicAmount(builder, maV)
	transactions.LockFundsTransactionBufferAddDuration(builder, dV)
	transactions.LockFundsTransactionBufferAddHash(builder, hV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return lockFundsTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *LockFundsTransaction) Size() int {
	return LockSize
}

type lockFundsTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Mosaic   mosaicDTO `json:"mosaic"`
		Duration uint64DTO `json:"duration"`
		Hash     Hash      `json:"hash"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *lockFundsTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaic, err := dto.Tx.Mosaic.toStruct()
	if err != nil {
		return nil, err
	}

	return &LockFundsTransaction{
		*atx,
		mosaic,
		dto.Tx.Duration.toStruct(),
		&SignedTransaction{Lock, "", dto.Tx.Hash},
	}, nil
}

type SecretLockTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration  Duration
	Secret    *Secret
	Recipient *Address
}

// returns a SecretLockTransaction from passed Mosaic, duration in blocks, Secret and mosaic recipient Address
func NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, secret *Secret, recipient *Address, networkType NetworkType) (*SecretLockTransaction, error) {
	if mosaic == nil {
		return nil, errors.New("mosaic must not be nil")
	}

	if secret == nil {
		return nil, errors.New("secret must not be nil")
	}

	if recipient == nil {
		return nil, errors.New("recipient must not be nil")
	}

	return &SecretLockTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     SecretLockVersion,
			Deadline:    deadline,
			Type:        SecretLock,
			NetworkType: networkType,
		},
		Mosaic:    mosaic,
		Duration:  duration,
		Secret:    secret,
		Recipient: recipient,
	}, nil
}

func (tx *SecretLockTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *SecretLockTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Mosaic": %s,
			"Duration": %d,
			"Secret": %s,
			"Recipient": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaic,
		tx.Duration,
		tx.Secret,
		tx.Recipient,
	)
}

func (tx *SecretLockTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.BlockchainId.toArray())
	maV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.Amount.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())

	sV := transactions.TransactionBufferCreateByteVector(builder, tx.Secret.Hash)

	addr, err := base32.StdEncoding.DecodeString(tx.Recipient.Address)
	if err != nil {
		return nil, err
	}
	rV := transactions.TransactionBufferCreateByteVector(builder, addr)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.SecretLockTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.SecretLockTransactionBufferAddMosaicId(builder, mV)
	transactions.SecretLockTransactionBufferAddMosaicAmount(builder, maV)
	transactions.SecretLockTransactionBufferAddDuration(builder, dV)
	transactions.SecretLockTransactionBufferAddHashAlgorithm(builder, byte(tx.Secret.Type))
	transactions.SecretLockTransactionBufferAddSecret(builder, sV)
	transactions.SecretLockTransactionBufferAddRecipient(builder, rV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return secretLockTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *SecretLockTransaction) Size() int {
	return SecretLockSize
}

type secretLockTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		BlockchainId *blockchainIdDTO `json:"mosaicId"`
		Amount       *uint64DTO       `json:"amount"`
		HashType     HashType         `json:"hashAlgorithm"`
		Duration     uint64DTO        `json:"duration"`
		Secret       string           `json:"secret"`
		Recipient    string           `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretLockTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromBase32(dto.Tx.Recipient)
	if err != nil {
		return nil, err
	}

	blockchainId, err := dto.Tx.BlockchainId.toStruct()
	if err != nil {
		return nil, err
	}

	mosaic, err := NewMosaic(blockchainId, dto.Tx.Amount.toStruct())
	if err != nil {
		return nil, err
	}

	secret, err := NewSecretFromHexString(dto.Tx.Secret, dto.Tx.HashType)
	if err != nil {
		return nil, err
	}

	return &SecretLockTransaction{
		*atx,
		mosaic,
		dto.Tx.Duration.toStruct(),
		secret,
		a,
	}, nil
}

type SecretProofTransaction struct {
	AbstractTransaction
	HashType
	Proof *Proof
}

// returns a SecretProofTransaction from passed HashType and Proof
func NewSecretProofTransaction(deadline *Deadline, hashType HashType, proof *Proof, networkType NetworkType) (*SecretProofTransaction, error) {
	if proof == nil {
		return nil, errors.New("proof must not be nil")
	}

	return &SecretProofTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     SecretProofVersion,
			Deadline:    deadline,
			Type:        SecretProof,
			NetworkType: networkType,
		},
		HashType: hashType,
		Proof:    proof,
	}, nil
}

func (tx *SecretProofTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *SecretProofTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"HashType": %s,
			"Proof": %s
		`,
		tx.AbstractTransaction.String(),
		tx.HashType,
		tx.Proof,
	)
}

func (tx *SecretProofTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	secret, err := tx.Proof.Secret(tx.HashType)
	if err != nil {
		return nil, err
	}
	sV := transactions.TransactionBufferCreateByteVector(builder, secret.Hash)

	pV := transactions.TransactionBufferCreateByteVector(builder, tx.Proof.Data)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.SecretProofTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.SecretProofTransactionBufferAddHashAlgorithm(builder, byte(tx.HashType))
	transactions.SecretProofTransactionBufferAddSecret(builder, sV)
	transactions.SecretProofTransactionBufferAddProofSize(builder, uint16(tx.Proof.Size()))
	transactions.SecretProofTransactionBufferAddProof(builder, pV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return secretProofTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *SecretProofTransaction) Size() int {
	return SecretProofHeaderSize + tx.Proof.Size()
}

type secretProofTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		HashType `json:"hashAlgorithm"`
		Proof    string `json:"proof"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretProofTransactionDTO) toStruct() (Transaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	proof, err := NewProofFromHexString(dto.Tx.Proof)
	if err != nil {
		return nil, err
	}

	return &SecretProofTransaction{
		*atx,
		dto.Tx.HashType,
		proof,
	}, nil
}

type CosignatureTransaction struct {
	TransactionToCosign *AggregateTransaction
}

// returns a CosignatureTransaction from passed AggregateTransaction
func NewCosignatureTransaction(txToCosign *AggregateTransaction) (*CosignatureTransaction, error) {
	if txToCosign == nil {
		return nil, errors.New("txToCosign must not be nil")
	}
	return &CosignatureTransaction{txToCosign}, nil
}

// returns a CosignatureTransaction from passed hash of bounded AggregateTransaction
func NewCosignatureTransactionFromHash(hash Hash) *CosignatureTransaction {
	return &CosignatureTransaction{
		TransactionToCosign: &AggregateTransaction{
			AbstractTransaction: AbstractTransaction{
				TransactionInfo: &TransactionInfo{
					Hash: hash,
				},
			},
		},
	}
}

func (tx *CosignatureTransaction) String() string {
	return fmt.Sprintf(`"TransactionToCosign": %s`, tx.TransactionToCosign)
}

type SignedTransaction struct {
	TransactionType `json:"transactionType"`
	Payload         string `json:"payload"`
	Hash            Hash   `json:"hash"`
}

type CosignatureSignedTransaction struct {
	ParentHash Hash   `json:"parentHash"`
	Signature  string `json:"signature"`
	Signer     string `json:"signer"`
}

type AggregateTransactionCosignature struct {
	Signature string
	Signer    *PublicAccount
}

type aggregateTransactionCosignatureDTO struct {
	Signature string `json:"signature"`
	Signer    string
}

func (dto *aggregateTransactionCosignatureDTO) toStruct(networkType NetworkType) (*AggregateTransactionCosignature, error) {
	acc, err := NewAccountFromPublicKey(dto.Signer, networkType)
	if err != nil {
		return nil, err
	}
	return &AggregateTransactionCosignature{
		dto.Signature,
		acc,
	}, nil
}

func (agt *AggregateTransactionCosignature) String() string {
	return fmt.Sprintf(
		`
			"Signature": %s,
			"Signer": %s
		`,
		agt.Signature,
		agt.Signer,
	)
}

type MultisigCosignatoryModification struct {
	Type MultisigCosignatoryModificationType
	*PublicAccount
}

func (m *MultisigCosignatoryModification) String() string {
	return fmt.Sprintf(
		`
			"Type": %s,
			"PublicAccount": %s
		`,
		m.Type,
		m.PublicAccount,
	)
}

type multisigCosignatoryModificationDTO struct {
	Type          MultisigCosignatoryModificationType `json:"type"`
	PublicAccount string                              `json:"cosignatoryPublicKey"`
}

func (dto *multisigCosignatoryModificationDTO) toStruct(networkType NetworkType) (*MultisigCosignatoryModification, error) {
	acc, err := NewAccountFromPublicKey(dto.PublicAccount, networkType)
	if err != nil {
		return nil, err
	}

	return &MultisigCosignatoryModification{
		dto.Type,
		acc,
	}, nil
}

type MetadataModification struct {
	Type  MetadataModificationType
	Key   string
	Value string
}

func (m *MetadataModification) Size() int {
	return SizeSize + 1 /* MetadataModificationType size */ + 1 /* KeySize size */ + 2 /* ValueSize size */ + len(m.Key) + len(m.Value)
}

func (m *MetadataModification) String() string {
	return fmt.Sprintf(
		`
			"Type"	: %s,
			"Key" 	: %s,
			"Value" : %s
		`,
		m.Type,
		m.Key,
		m.Value,
	)
}

type metadataModificationDTO struct {
	Type  MetadataModificationType `json:"modificationType"`
	Key   string                   `json:"key"`
	Value string                   `json:"value"`
}

func (dto *metadataModificationDTO) toStruct(networkType NetworkType) (*MetadataModification, error) {
	return &MetadataModification{
		dto.Type,
		dto.Key,
		dto.Value,
	}, nil
}

type mosaicDefinitonTransactionPropertiesDTO []struct {
	Key   int
	Value uint64DTO
}

func (dto mosaicDefinitonTransactionPropertiesDTO) toStruct() *MosaicProperties {
	flags := dto[0].Value.toUint64()
	duration := Duration(0)
	if len(dto) == 3 {
		duration = dto[2].Value.toStruct()
	}
	return NewMosaicProperties(
		hasBits(flags, Supply_Mutable),
		hasBits(flags, Transferable),
		hasBits(flags, LevyMutable),
		byte(dto[1].Value.toUint64()),
		duration,
	)
}

type TransactionStatus struct {
	Deadline *Deadline
	Group    string
	Status   string
	Hash     Hash
	Height   Height
}

func (ts *TransactionStatus) String() string {
	return fmt.Sprintf(
		`
			"Group:" %s,
			"Status:" %s,
			"Content": %s,
			"Deadline": %s,
			"Height": %s
		`,
		ts.Group,
		ts.Status,
		ts.Hash,
		ts.Deadline,
		ts.Height,
	)
}

type transactionStatusDTO struct {
	Group    string                 `json:"group"`
	Status   string                 `json:"status"`
	Hash     Hash                   `json:"hash"`
	Deadline blockchainTimestampDTO `json:"deadline"`
	Height   uint64DTO              `json:"height"`
}

func (dto *transactionStatusDTO) toStruct() (*TransactionStatus, error) {
	return &TransactionStatus{
		NewDeadlineFromBlockchainTimestamp(dto.Deadline.toStruct()),
		dto.Group,
		dto.Status,
		dto.Hash,
		dto.Height.toStruct(),
	}, nil
}

type TransactionIdsDTO struct {
	Ids []string `json:"transactionIds"`
}

type TransactionHashesDTO struct {
	Hashes []string `json:"hashes"`
}

const (
	AddressSize                              int = 25
	AmountSize                               int = 8
	KeySize                                  int = 32
	Hash256                                  int = 32
	MosaicSize                               int = 8
	NamespaceSize                            int = 8
	SizeSize                                 int = 4
	SignerSize                               int = KeySize
	SignatureSize                            int = 64
	VersionSize                              int = 2
	TypeSize                                 int = 2
	MaxFeeSize                               int = 8
	DeadLineSize                             int = 8
	DurationSize                             int = 8
	TransactionHeaderSize                    int = SizeSize + SignerSize + SignatureSize + VersionSize + TypeSize + MaxFeeSize + DeadLineSize
	PropertyTypeSize                         int = 2
	PropertyModificationTypeSize             int = 1
	AccountPropertiesAddressModificationSize int = PropertyModificationTypeSize + AddressSize
	AccountPropertiesMosaicModificationSize  int = PropertyModificationTypeSize + MosaicSize
	AccountPropertiesEntityModificationSize  int = PropertyModificationTypeSize + TypeSize
	AccountPropertyAddressHeader             int = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyMosaicHeader              int = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyEntityTypeHeader          int = TransactionHeaderSize + PropertyTypeSize
	LinkActionSize                           int = 1
	AccountLinkTransactionSize               int = TransactionHeaderSize + KeySize + LinkActionSize
	AliasActionSize                          int = 1
	AliasTransactionHeader                   int = TransactionHeaderSize + NamespaceSize + AliasActionSize
	AggregateBondedHeader                    int = TransactionHeaderSize + SizeSize
	HashTypeSize                             int = 1
	LockSize                                 int = TransactionHeaderSize + MosaicSize + AmountSize + DurationSize + Hash256
	MetadataTypeSize                         int = 1
	MetadataHeaderSize                       int = TransactionHeaderSize + MetadataTypeSize
	ModificationsSizeSize                    int = 1
	ModifyContractHeaderSize                 int = TransactionHeaderSize + DurationSize + Hash256 + 3*ModificationsSizeSize
	MinApprovalSize                          int = 1
	MinRemovalSize                           int = 1
	ModifyMultisigHeaderSize                 int = TransactionHeaderSize + MinApprovalSize + MinRemovalSize + ModificationsSizeSize
	MosaicNonceSize                          int = 4
	MosaicPropertySize                       int = 4
	MosaicDefinitionTransactionSize          int = TransactionHeaderSize + MosaicNonceSize + MosaicSize + DurationSize + MosaicPropertySize
	MosaicSupplyDirectionSize                int = 1
	MosaicSupplyChangeTransactionSize        int = TransactionHeaderSize + MosaicSize + AmountSize + MosaicSupplyDirectionSize
	NamespaceTypeSize                        int = 1
	NamespaceNameSizeSize                    int = 1
	RegisterNamespaceHeaderSize              int = TransactionHeaderSize + NamespaceTypeSize + DurationSize + NamespaceSize + NamespaceNameSizeSize
	SecretLockSize                           int = TransactionHeaderSize + MosaicSize + AmountSize + DurationSize + HashTypeSize + Hash256 + AddressSize
	ProofSizeSize                            int = 2
	SecretProofHeaderSize                    int = TransactionHeaderSize + HashTypeSize + Hash256 + ProofSizeSize
	MosaicsSizeSize                          int = 1
	MessageSizeSize                          int = 2
	TransferHeaderSize                       int = TransactionHeaderSize + AddressSize + MosaicsSizeSize + MessageSizeSize
)

type TransactionType uint16

const (
	AccountPropertyAddress    TransactionType = 0x4150
	AccountPropertyMosaic     TransactionType = 0x4250
	AccountPropertyEntityType TransactionType = 0x4350
	AddressAlias              TransactionType = 0x424e
	AggregateBonded           TransactionType = 0x4241
	AggregateCompleted        TransactionType = 0x4141
	LinkAccount               TransactionType = 0x414c
	Lock                      TransactionType = 0x4148
	MetadataAddress           TransactionType = 0x413d
	MetadataMosaic            TransactionType = 0x423d
	MetadataNamespace         TransactionType = 0x433d
	ModifyContract            TransactionType = 0x4157
	ModifyMultisig            TransactionType = 0x4155
	MosaicAlias               TransactionType = 0x434e
	MosaicDefinition          TransactionType = 0x414d
	MosaicSupplyChange        TransactionType = 0x424d
	RegisterNamespace         TransactionType = 0x414e
	SecretLock                TransactionType = 0x4152
	SecretProof               TransactionType = 0x4252
	Transfer                  TransactionType = 0x4154
)

func (t TransactionType) String() string {
	return fmt.Sprintf("%x", uint16(t))
}

var transactionTypeError = errors.New("wrong raw TransactionType int")

type TransactionVersion uint8

const (
	AccountPropertyAddressVersion    TransactionVersion = 1
	AccountPropertyMosaicVersion     TransactionVersion = 1
	AccountPropertyEntityTypeVersion TransactionVersion = 1
	AddressAliasVersion              TransactionVersion = 1
	AggregateBondedVersion           TransactionVersion = 2
	AggregateCompletedVersion        TransactionVersion = 2
	LinkAccountVersion               TransactionVersion = 2
	LockVersion                      TransactionVersion = 1
	MetadataAddressVersion           TransactionVersion = 1
	MetadataMosaicVersion            TransactionVersion = 1
	MetadataNamespaceVersion         TransactionVersion = 1
	ModifyContractVersion            TransactionVersion = 3
	ModifyMultisigVersion            TransactionVersion = 3
	MosaicAliasVersion               TransactionVersion = 1
	MosaicDefinitionVersion          TransactionVersion = 3
	MosaicSupplyChangeVersion        TransactionVersion = 2
	RegisterNamespaceVersion         TransactionVersion = 2
	SecretLockVersion                TransactionVersion = 1
	SecretProofVersion               TransactionVersion = 1
	TransferVersion                  TransactionVersion = 3
)

type AccountLinkAction uint8

// AccountLinkAction enums
const (
	AccountLink AccountLinkAction = iota
	AccountUnlink
)

type AliasActionType uint8

// AliasActionType enums
const (
	AliasLink AliasActionType = iota
	AliasUnlink
)

type AliasType uint8

// AliasType enums
const (
	NoneAliasType AliasType = iota
	MosaicAliasType
	AddressAliasType
)

type PropertyModificationType uint8

// PropertyModificationType enums
const (
	AddProperty PropertyModificationType = iota
	RemoveProperty
)

type PropertyType uint8

// Account property type
// 0x01	The property type is an address.
// 0x02	The property type is mosaic id.
// 0x04	The property type is a transaction type.
// 0x05	Property type sentinel.
// 0x80 + type	The property is interpreted as a blocking operation.
const (
	AllowAddress     PropertyType = 0x01
	AllowMosaic      PropertyType = 0x02
	AllowTransaction PropertyType = 0x04
	Sentinel         PropertyType = 0x05
	BlockAddress     PropertyType = 0x80 + 0x01
	BlockMosaic      PropertyType = 0x80 + 0x02
	BlockTransaction PropertyType = 0x80 + 0x04
)

type MultisigCosignatoryModificationType uint8

func (t MultisigCosignatoryModificationType) String() string {
	return fmt.Sprintf("%d", t)
}

const (
	Add MultisigCosignatoryModificationType = iota
	Remove
)

type MetadataModificationType uint8

func (t MetadataModificationType) String() string {
	return fmt.Sprintf("%d", t)
}

const (
	AddMetadata MetadataModificationType = iota
	RemoveMetadata
)

type MetadataType uint8

func (t MetadataType) String() string {
	return fmt.Sprintf("%d", t)
}

const (
	MetadataNone MetadataType = iota
	MetadataAddressType
	MetadataMosaicType
	MetadataNamespaceType
)

type Hash string

func (h Hash) String() string {
	return (string)(h)
}

func ExtractVersion(version uint64) uint8 {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, version)

	return uint8(b[0])
}

func MapTransactions(b *bytes.Buffer) ([]Transaction, error) {
	var wg sync.WaitGroup
	var err error

	var m []jsonLib.RawMessage

	err = json.Unmarshal(b.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	txs := make([]Transaction, len(m))
	errs := make([]error, len(m))
	for i, t := range m {
		wg.Add(1)
		go func(i int, t jsonLib.RawMessage) {
			defer wg.Done()
			txs[i], errs[i] = MapTransaction(bytes.NewBuffer([]byte(t)))
		}(i, t)
	}
	wg.Wait()

	for _, err = range errs {
		if err != nil {
			return txs, err
		}
	}

	return txs, nil
}

func dtoToTransaction(b *bytes.Buffer, dto transactionDto) (Transaction, error) {
	if dto == nil {
		return nil, errors.New("dto can't be nil")
	}

	err := json.Unmarshal(b.Bytes(), dto)
	if err != nil {
		return nil, err
	}

	tx, err := dto.toStruct()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func MapTransaction(b *bytes.Buffer) (Transaction, error) {
	rawT := struct {
		Transaction struct {
			Type TransactionType
		}
	}{}

	err := json.Unmarshal(b.Bytes(), &rawT)
	if err != nil {
		return nil, err
	}

	var dto transactionDto = nil

	switch rawT.Transaction.Type {
	case AccountPropertyAddress:
		dto = &accountPropertiesAddressTransactionDTO{}
	case AccountPropertyMosaic:
		dto = &accountPropertiesMosaicTransactionDTO{}
	case AccountPropertyEntityType:
		dto = &accountPropertiesEntityTypeTransactionDTO{}
	case AddressAlias:
		dto = &addressAliasTransactionDTO{}
	case AggregateBonded, AggregateCompleted:
		dto = &aggregateTransactionDTO{}
	case LinkAccount:
		dto = &accountLinkTransactionDTO{}
	case Lock:
		dto = &lockFundsTransactionDTO{}
	case MetadataAddress:
		dto = &modifyMetadataAddressTransactionDTO{}
	case MetadataMosaic:
		dto = &modifyMetadataMosaicTransactionDTO{}
	case MetadataNamespace:
		dto = &modifyMetadataNamespaceTransactionDTO{}
	case ModifyContract:
		dto = &modifyContractTransactionDTO{}
	case ModifyMultisig:
		dto = &modifyMultisigAccountTransactionDTO{}
	case MosaicAlias:
		dto = &mosaicAliasTransactionDTO{}
	case MosaicDefinition:
		dto = &mosaicDefinitionTransactionDTO{}
	case MosaicSupplyChange:
		dto = &mosaicSupplyChangeTransactionDTO{}
	case RegisterNamespace:
		dto = &registerNamespaceTransactionDTO{}
	case SecretLock:
		dto = &secretLockTransactionDTO{}
	case SecretProof:
		dto = &secretProofTransactionDTO{}
	case Transfer:
		dto = &transferTransactionDTO{}
	}

	return dtoToTransaction(b, dto)
}

func createTransactionHash(p string) (string, error) {
	b, err := hex.DecodeString(p)
	if err != nil {
		return "", err
	}

	const HalfOfSignature = SignatureSize / 2

	sb := make([]byte, len(b)-SizeSize-HalfOfSignature)
	copy(sb[:HalfOfSignature], b[SizeSize:SizeSize+HalfOfSignature])
	copy(sb[HalfOfSignature:], b[SizeSize+SignatureSize:])

	r, err := crypto.HashesSha3_256(sb)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(r)), nil
}

func toAggregateTransactionBytes(tx Transaction) ([]byte, error) {
	if tx.GetAbstractTransaction().Signer == nil {
		return nil, fmt.Errorf("some of the transaction does not have a signer")
	}
	sb, err := hex.DecodeString(tx.GetAbstractTransaction().Signer.PublicKey)
	if err != nil {
		return nil, err
	}
	b, err := tx.generateBytes()
	if err != nil {
		return nil, err
	}

	rB := make([]byte, len(b)-SignatureSize-MaxFeeSize-DeadLineSize)
	copy(rB[SizeSize:SignerSize+SizeSize], sb[:SignerSize])
	copy(
		rB[SignerSize+SizeSize:SignerSize+SizeSize+VersionSize+TypeSize],
		b[SizeSize+SignerSize+SignatureSize:SizeSize+SignerSize+SignatureSize+VersionSize+TypeSize],
	)
	copy(rB[SignerSize+SizeSize+VersionSize+TypeSize:], b[TransactionHeaderSize:])

	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, uint32(len(rB)))

	copy(rB[:len(s)], s)

	return rB, nil
}

func signTransactionWith(tx Transaction, a *Account) (*SignedTransaction, error) {
	s := crypto.NewSignerFromKeyPair(a.KeyPair, nil)
	b, err := tx.generateBytes()
	if err != nil {
		return nil, err
	}
	sb := make([]byte, len(b)-SizeSize-SignerSize-SignatureSize)
	copy(sb, b[SizeSize+SignerSize+SignatureSize:])
	signature, err := s.Sign(sb)
	if err != nil {
		return nil, err
	}

	p := make([]byte, len(b))
	copy(p[:SizeSize], b[:SizeSize])
	copy(p[SizeSize:SizeSize+SignatureSize], signature.Bytes())
	copy(p[SizeSize+SignatureSize:SizeSize+SignatureSize+SignerSize], a.KeyPair.PublicKey.Raw)
	copy(p[SizeSize+SignatureSize+SignerSize:], b[SizeSize+SignatureSize+SignerSize:])

	ph := hex.EncodeToString(p)
	h, err := createTransactionHash(ph)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{tx.GetAbstractTransaction().Type, strings.ToUpper(ph), (Hash)(h)}, nil
}

func signTransactionWithCosignatures(tx *AggregateTransaction, a *Account, cosignatories []*Account) (*SignedTransaction, error) {
	stx, err := signTransactionWith(tx, a)
	if err != nil {
		return nil, err
	}

	p := stx.Payload

	b, err := hex.DecodeString((string)(stx.Hash))
	if err != nil {
		return nil, err
	}

	for _, cos := range cosignatories {
		s := crypto.NewSignerFromKeyPair(cos.KeyPair, nil)
		sb, err := s.Sign(b)
		if err != nil {
			return nil, err
		}
		p += cos.KeyPair.PublicKey.String() + hex.EncodeToString(sb.Bytes())
	}

	pb, err := hex.DecodeString(p)
	if err != nil {
		return nil, err
	}

	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, uint32(len(pb)))

	copy(pb[:len(s)], s)

	return &SignedTransaction{tx.Type, hex.EncodeToString(pb), stx.Hash}, nil
}

func signCosignatureTransaction(a *Account, tx *CosignatureTransaction) (*CosignatureSignedTransaction, error) {
	if tx.TransactionToCosign.TransactionInfo == nil || tx.TransactionToCosign.TransactionInfo.Hash == "" {
		return nil, errors.New("cosignature transaction hash is nil")
	}

	s := crypto.NewSignerFromKeyPair(a.KeyPair, nil)
	b, err := hex.DecodeString((string)(tx.TransactionToCosign.TransactionInfo.Hash))
	if err != nil {
		return nil, err
	}

	sb, err := s.Sign(b)
	if err != nil {
		return nil, err
	}

	return &CosignatureSignedTransaction{tx.TransactionToCosign.TransactionInfo.Hash, hex.EncodeToString(sb.Bytes()), a.PublicAccount.PublicKey}, nil
}

func cosignatoryModificationArrayToBuffer(builder *flatbuffers.Builder, modifications []*MultisigCosignatoryModification) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(modifications))
	for i, m := range modifications {
		b, err := utils.HexDecodeStringOdd(m.PublicAccount.PublicKey)
		if err != nil {
			return 0, err
		}
		pV := transactions.TransactionBufferCreateByteVector(builder, b)
		transactions.CosignatoryModificationBufferStart(builder)
		transactions.CosignatoryModificationBufferAddType(builder, uint8(m.Type))
		transactions.CosignatoryModificationBufferAddCosignatoryPublicKey(builder, pV)
		msb[i] = transactions.TransactionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func metadataModificationArrayToBuffer(builder *flatbuffers.Builder, modifications []*MetadataModification) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(modifications))
	for i, m := range modifications {
		keySize := len(m.Key)

		if keySize == 0 {
			return 0, errors.New("key must not empty")
		}

		pKey := transactions.TransactionBufferCreateByteVector(builder, []byte(m.Key))
		valueSize := len(m.Value)

		// it is hack, because we can have case when size of the value is zero(in RemoveData modification),
		// but flattbuffer doesn't store int(0) like 4 bytes, it stores like one byte
		valueB := make([]byte, 2)
		binary.LittleEndian.PutUint16(valueB, uint16(valueSize))
		pValueSize := transactions.TransactionBufferCreateByteVector(builder, valueB)

		pValue := transactions.TransactionBufferCreateByteVector(builder, []byte(m.Value))

		transactions.MetadataModificationBufferStart(builder)
		transactions.MetadataModificationBufferAddSize(builder, uint32(m.Size()))
		transactions.MetadataModificationBufferAddModificationType(builder, uint8(m.Type))
		transactions.MetadataModificationBufferAddKeySize(builder, uint8(keySize))
		transactions.MetadataModificationBufferAddValueSize(builder, pValueSize)
		transactions.MetadataModificationBufferAddKey(builder, pKey)
		transactions.MetadataModificationBufferAddValue(builder, pValue)

		msb[i] = transactions.MetadataModificationBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func stringToBuffer(builder *flatbuffers.Builder, hash string) flatbuffers.UOffsetT {
	b := utils.MustHexDecodeString(hash)
	pV := transactions.TransactionBufferCreateByteVector(builder, b)

	return pV
}

func metadataDTOArrayToStruct(Modifications []*metadataModificationDTO, NetworkType NetworkType) ([]*MetadataModification, error) {
	ms := make([]*MetadataModification, len(Modifications))
	var err error = nil
	for i, m := range Modifications {
		ms[i], err = m.toStruct(NetworkType)

		if err != nil {
			return nil, err
		}
	}

	return ms, err
}

func multisigCosignatoryDTOArrayToStruct(Modifications []*multisigCosignatoryModificationDTO, NetworkType NetworkType) ([]*MultisigCosignatoryModification, error) {
	ms := make([]*MultisigCosignatoryModification, len(Modifications))
	var err error = nil
	for i, m := range Modifications {
		ms[i], err = m.toStruct(NetworkType)

		if err != nil {
			return nil, err
		}
	}

	return ms, err
}
