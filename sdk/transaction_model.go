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
	"github.com/proximax-storage/go-xpx-catapult-sdk/utils"
	"github.com/proximax-storage/go-xpx-utils/str"
	"github.com/proximax-storage/nem2-crypto-go"
	"math/big"
	"strings"
	"sync"
	"time"
)

// Models
// Transaction
type Transaction interface {
	GetAbstractTransaction() *AbstractTransaction
	String() string
	generateBytes() ([]byte, error)
}

// AbstractTransaction
type AbstractTransaction struct {
	*TransactionInfo
	NetworkType NetworkType
	Deadline    *Deadline
	Type        TransactionType
	Version     TransactionVersion
	Fee         *big.Int
	Signature   string
	Signer      *PublicAccount
}

func (tx *AbstractTransaction) IsUnconfirmed() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height.Int64() == 0 && tx.TransactionInfo.Hash == tx.TransactionInfo.MerkleComponentHash
}

func (tx *AbstractTransaction) IsConfirmed() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height.Int64() > 0
}

func (tx *AbstractTransaction) HasMissingSignatures() bool {
	return tx.TransactionInfo != nil && tx.TransactionInfo.Height.Int64() == 0 && tx.TransactionInfo.Hash != tx.TransactionInfo.MerkleComponentHash
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
			"Fee": %d,
			"Deadline": %s,
			"Signature": %s,
			"Signer": %s
		`,
		tx.NetworkType,
		tx.TransactionInfo.String(),
		tx.Type,
		tx.Version,
		tx.Fee,
		tx.Deadline,
		tx.Signature,
		tx.Signer,
	)
}

func (tx *AbstractTransaction) generateVectors(builder *flatbuffers.Builder) (v uint16, signatureV, signerV, dV, fV flatbuffers.UOffsetT, err error) {
	v = (uint16(tx.NetworkType) << 8) + uint16(tx.Version)
	signatureV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, 64))
	signerV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, 32))
	dV = transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(big.NewInt(tx.Deadline.GetInstant())))
	fV = transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Fee))
	return
}

func (tx *AbstractTransaction) buildVectors(builder *flatbuffers.Builder, v uint16, signatureV, signerV, dV, fV flatbuffers.UOffsetT) {
	transactions.TransactionBufferAddSignature(builder, signatureV)
	transactions.TransactionBufferAddSigner(builder, signerV)
	transactions.TransactionBufferAddVersion(builder, v)
	transactions.TransactionBufferAddType(builder, tx.Type.Hex())
	transactions.TransactionBufferAddFee(builder, fV)
	transactions.TransactionBufferAddDeadline(builder, dV)
}

type abstractTransactionDTO struct {
	Type      uint32     `json:"type"`
	Version   uint64     `json:"version"`
	Fee       *uint64DTO `json:"fee"`
	Deadline  *uint64DTO `json:"deadline"`
	Signature string     `json:"signature"`
	Signer    string     `json:"signer"`
}

func (dto *abstractTransactionDTO) toStruct(tInfo *TransactionInfo) (*AbstractTransaction, error) {
	t, err := TransactionTypeFromRaw(dto.Type)
	if err != nil {
		return nil, err
	}

	nt := ExtractNetworkType(dto.Version)

	tv := TransactionVersion(ExtractVersion(dto.Version))

	pa, err := NewAccountFromPublicKey(dto.Signer, nt)
	if err != nil {
		return nil, err
	}

	var d *Deadline
	if dto.Deadline != nil {
		d = &Deadline{time.Unix(0, dto.Deadline.toBigInt().Int64()*int64(time.Millisecond))}
	}

	var f *big.Int
	if dto.Fee != nil {
		f = dto.Fee.toBigInt()
	}

	return &AbstractTransaction{
		tInfo,
		nt,
		d,
		t,
		tv,
		f,
		dto.Signature,
		pa,
	}, nil
}

// Transaction Info
type TransactionInfo struct {
	Height              *big.Int
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
			"Height": %d,
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
	Height              *uint64DTO `json:"height"`
	Index               uint32     `json:"index"`
	Id                  string     `json:"id"`
	Hash                Hash       `json:"hash"`
	MerkleComponentHash Hash       `json:"merkleComponentHash"`
	AggregateHash       Hash       `json:"aggregateHash,omitempty"`
	AggregateId         string     `json:"aggregateId,omitempty"`
}

func (dto *transactionInfoDTO) toStruct() *TransactionInfo {
	height := big.NewInt(0)
	if dto.Height != nil {
		height = dto.Height.toBigInt()
	}
	return &TransactionInfo{
		height,
		dto.Index,
		dto.Id,
		dto.Hash,
		dto.MerkleComponentHash,
		dto.AggregateHash,
		dto.AggregateId,
	}
}

// AggregateTransaction
type AggregateTransaction struct {
	AbstractTransaction
	InnerTransactions []Transaction
	Cosignatures      []*AggregateTransactionCosignature
}

// Create an aggregate complete transaction
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
	transactions.TransactionBufferAddSize(builder, 120+4+len(txsb))
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)
	transactions.AggregateTransactionBufferAddTransactionsSize(builder, uint32(len(txsb)))
	transactions.AggregateTransactionBufferAddTransactions(builder, tV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return aggregateTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type aggregateTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Cosignatures      []*aggregateTransactionCosignatureDTO `json:"cosignatures"`
		InnerTransactions []map[string]interface{}              `json:"transactions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *aggregateTransactionDTO) toStruct() (*AggregateTransaction, error) {
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
		iatx.Fee = atx.Fee
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

// ModifyMetadataTransaction
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
		tx.MetadataType.String(),
		tx.Modifications,
	)
}

func (tx *ModifyMetadataTransaction) generateBytes(builder *flatbuffers.Builder, metadataV flatbuffers.UOffsetT, sizeOfMetadata uint32) ([]byte, error) {

	mV, sizeOfModifications, err := metadataModificationArrayToBuffer(builder, tx.Modifications)
	if err != nil {
		return nil, err
	}

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.ModifyMetadataTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, int(120+1+sizeOfMetadata+sizeOfModifications))

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ModifyMetadataTransactionBufferAddMetadataType(builder, uint8(tx.MetadataType))
	transactions.ModifyMetadataTransactionBufferAddMetadataId(builder, metadataV)
	transactions.ModifyMetadataTransactionBufferAddModifications(builder, mV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyMetadataTransactionSchema().serialize(builder.FinishedBytes()), nil
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

// ModifyMetadataAddressTransaction
type ModifyMetadataAddressTransaction struct {
	ModifyMetadataTransaction
	Address *Address
}

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

	return tx.ModifyMetadataTransaction.generateBytes(builder, aV, 25)
}

type modifyMetadataAddressTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		Address string `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataAddressTransactionDTO) toStruct() (*ModifyMetadataAddressTransaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromEncoded(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataAddressTransaction{
		*atx,
		a,
	}, nil
}

// ModifyMetadataMosaicTransaction
type ModifyMetadataMosaicTransaction struct {
	ModifyMetadataTransaction
	MosaicId *MosaicId
}

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
		tx.MosaicId.String(),
	)
}

func (tx *ModifyMetadataMosaicTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, 8)
	binary.LittleEndian.PutUint64(mosaicB, mosaicIdToBigInt(tx.MosaicId).Uint64())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.ModifyMetadataTransaction.generateBytes(builder, mV, 8)
}

type modifyMetadataMosaicTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		MosaicId *uint64DTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataMosaicTransactionDTO) toStruct() (*ModifyMetadataMosaicTransaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataMosaicTransaction{
		*atx,
		mosaicId,
	}, nil
}

// ModifyMetadataNamespaceTransaction
type ModifyMetadataNamespaceTransaction struct {
	ModifyMetadataTransaction
	NamespaceId *NamespaceId
}

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
		tx.NamespaceId.String(),
	)
}

func (tx *ModifyMetadataNamespaceTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, 8)
	binary.LittleEndian.PutUint64(mosaicB, namespaceIdToBigInt(tx.NamespaceId).Uint64())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.ModifyMetadataTransaction.generateBytes(builder, mV, 8)
}

type modifyMetadataNamespaceTransactionDTO struct {
	Tx struct {
		modifyMetadataTransactionDTO
		NamespaceId *uint64DTO `json:"metadataId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMetadataNamespaceTransactionDTO) toStruct() (*ModifyMetadataNamespaceTransaction, error) {
	atx, err := dto.Tx.modifyMetadataTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	namespaceId, err := NewNamespaceId(dto.Tx.NamespaceId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &ModifyMetadataNamespaceTransaction{
		*atx,
		namespaceId,
	}, nil
}

// MosaicDefinitionTransaction
type MosaicDefinitionTransaction struct {
	AbstractTransaction
	*MosaicProperties
	MosaicNonce uint32
	*MosaicId
}

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
			"MosaicId": [ %s ]
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicProperties.String(),
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

	mV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(mosaicIdToBigInt(tx.MosaicId)))
	dV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.MosaicProperties.Duration))

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MosaicDefinitionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, 120+24)
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

type mosaicDefinitionTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Properties  mosaicDefinitonTransactionPropertiesDTO `json:"properties"`
		MosaicNonce uint32                                  `json:"mosaicNonce"`
		MosaicId    *uint64DTO                              `json:"mosaicId"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicDefinitionTransactionDTO) toStruct() (*MosaicDefinitionTransaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &MosaicDefinitionTransaction{
		*atx,
		dto.Tx.Properties.toStruct(),
		dto.Tx.MosaicNonce,
		mosaicId,
	}, nil
}

// MosaicSupplyChangeTransaction
type MosaicSupplyChangeTransaction struct {
	AbstractTransaction
	MosaicSupplyType
	*MosaicId
	Delta *big.Int
}

func NewMosaicSupplyChangeTransaction(deadline *Deadline, mosaicId *MosaicId, supplyType MosaicSupplyType, delta *big.Int, networkType NetworkType) (*MosaicSupplyChangeTransaction, error) {
	if mosaicId == nil || mosaicIdToBigInt(mosaicId).Int64() == 0 {
		return nil, ErrNilMosaicId
	}

	if !(supplyType == Increase || supplyType == Decrease) {
		return nil, errors.New("supplyType must not be nil")
	}
	if delta == nil {
		return nil, errors.New("delta must not be nil")
	}

	return &MosaicSupplyChangeTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     MosaicSupplyChangeVersion,
			Deadline:    deadline,
			Type:        MosaicSupplyChange,
			NetworkType: networkType,
		},
		MosaicId:         mosaicId,
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
			"MosaicId": [ %v ],
			"Delta": %d
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicSupplyType.String(),
		tx.MosaicId,
		tx.Delta,
	)
}

func (tx *MosaicSupplyChangeTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(mosaicIdToBigInt(tx.MosaicId)))
	dV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Delta))

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MosaicSupplyChangeTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, 137)
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicSupplyChangeTransactionBufferAddMosaicId(builder, mV)
	transactions.MosaicSupplyChangeTransactionBufferAddDirection(builder, uint8(tx.MosaicSupplyType))
	transactions.MosaicSupplyChangeTransactionBufferAddDelta(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return mosaicSupplyChangeTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type mosaicSupplyChangeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicSupplyType `json:"direction"`
		MosaicId         *uint64DTO `json:"mosaicId"`
		Delta            *uint64DTO `json:"delta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *mosaicSupplyChangeTransactionDTO) toStruct() (*MosaicSupplyChangeTransaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &MosaicSupplyChangeTransaction{
		*atx,
		dto.Tx.MosaicSupplyType,
		mosaicId,
		dto.Tx.Delta.toBigInt(),
	}, nil
}

// TransferTransaction
type TransferTransaction struct {
	AbstractTransaction
	*Message
	Mosaics   []*Mosaic
	Recipient *Address
}

// Create a transfer transaction
func NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message *Message, networkType NetworkType) (*TransferTransaction, error) {
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
		tx.Message.String(),
	)
}

func (tx *TransferTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	ml := len(tx.Mosaics)
	mb := make([]flatbuffers.UOffsetT, ml)
	for i, mos := range tx.Mosaics {
		id := transactions.MosaicBufferCreateIdVector(builder, FromBigInt(mosaicIdToBigInt(mos.MosaicId)))
		am := transactions.MosaicBufferCreateAmountVector(builder, FromBigInt(mos.Amount))
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}

	p := []byte(tx.Payload)
	pl := len(p)
	mp := transactions.TransactionBufferCreateByteVector(builder, p)
	transactions.MessageBufferStart(builder)
	transactions.MessageBufferAddType(builder, tx.Message.Type)
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
	transactions.TransactionBufferAddSize(builder, 148+1+(16*ml)+pl)
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.TransferTransactionBufferAddRecipient(builder, rV)
	transactions.TransferTransactionBufferAddNumMosaics(builder, uint8(ml))
	transactions.TransferTransactionBufferAddMessageSize(builder, uint16(pl+1))
	transactions.TransferTransactionBufferAddMessage(builder, m)
	transactions.TransferTransactionBufferAddMosaics(builder, mV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return transferTransactionSchema().serialize(builder.FinishedBytes()), nil
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

func (dto *transferTransactionDTO) toStruct() (*TransferTransaction, error) {
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

	a, err := NewAddressFromEncoded(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &TransferTransaction{
		*atx,
		dto.Tx.Message.toStruct(),
		mosaics,
		a,
	}, nil
}

// ModifyMultisigAccountTransaction
type ModifyMultisigAccountTransaction struct {
	AbstractTransaction
	MinApprovalDelta uint8
	MinRemovalDelta  uint8
	Modifications    []*MultisigCosignatoryModification
}

func NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta uint8, minRemovalDelta uint8, modifications []*MultisigCosignatoryModification, networkType NetworkType) (*ModifyMultisigAccountTransaction, error) {
	if len(modifications) == 0 {
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
	transactions.TransactionBufferAddSize(builder, 123+(33*len(tx.Modifications)))
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ModifyMultisigAccountTransactionBufferAddMinRemovalDelta(builder, tx.MinRemovalDelta)
	transactions.ModifyMultisigAccountTransactionBufferAddMinApprovalDelta(builder, tx.MinApprovalDelta)
	transactions.ModifyMultisigAccountTransactionBufferAddNumModifications(builder, uint8(len(tx.Modifications)))
	transactions.ModifyMultisigAccountTransactionBufferAddModifications(builder, mV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyMultisigAccountTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type modifyMultisigAccountTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MinApprovalDelta uint8                                 `json:"minApprovalDelta"`
		MinRemovalDelta  uint8                                 `json:"minRemovalDelta"`
		Modifications    []*multisigCosignatoryModificationDTO `json:"modifications"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyMultisigAccountTransactionDTO) toStruct() (*ModifyMultisigAccountTransaction, error) {
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

// ModifyContractTransaction
type ModifyContractTransaction struct {
	AbstractTransaction
	DurationDelta int64
	Hash          string
	Customers     []*MultisigCosignatoryModification
	Executors     []*MultisigCosignatoryModification
	Verifiers     []*MultisigCosignatoryModification
}

func NewModifyContractTransaction(
	deadline *Deadline, durationDelta int64, hash string,
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

	durationV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(big.NewInt(tx.DurationDelta)))
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
	transactions.TransactionBufferAddSize(builder, 120+ // AbstractTransaction
		8+32+1+1+1+ // Fields of current transaction
		((32+1)*(len(tx.Customers)+len(tx.Executors)+len(tx.Verifiers))))
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

type modifyContractTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DurationDelta *uint64DTO                            `json:"duration"`
		Hash          string                                `json:"hash"`
		Customers     []*multisigCosignatoryModificationDTO `json:"customers"`
		Executors     []*multisigCosignatoryModificationDTO `json:"executors"`
		Verifiers     []*multisigCosignatoryModificationDTO `json:"verifiers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyContractTransactionDTO) toStruct() (*ModifyContractTransaction, error) {
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
		dto.Tx.DurationDelta.toBigInt().Int64(),
		dto.Tx.Hash,
		customers,
		executors,
		verifiers,
	}, nil
}

// RegisterNamespaceTransaction
type RegisterNamespaceTransaction struct {
	AbstractTransaction
	*NamespaceId
	NamespaceType
	NamspaceName string
	Duration     *big.Int
	ParentId     *NamespaceId
}

func NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration *big.Int, networkType NetworkType) (*RegisterNamespaceTransaction, error) {
	if len(namespaceName) == 0 {
		return nil, ErrInvalidNamespaceName
	}

	nsId, err := NewNamespaceIdFromName(namespaceName)
	if err != nil {
		return nil, err
	}

	if duration == nil {
		return nil, errors.New("duration must not be nil")
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

func NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId, networkType NetworkType) (*RegisterNamespaceTransaction, error) {
	if len(namespaceName) == 0 {
		return nil, ErrInvalidNamespaceName
	}

	if parentId == nil || namespaceIdToBigInt(parentId).Int64() == 0 {
		return nil, ErrNilNamespaceId
	}

	nsId, err := generateNamespaceId(namespaceName, namespaceIdToBigInt(parentId))
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
		NamespaceId:   bigIntToNamespaceId(nsId),
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
			"Duration": %d
		`,
		tx.AbstractTransaction.String(),
		tx.NamspaceName,
		tx.Duration,
	)
}

func (tx *RegisterNamespaceTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	nV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(namespaceIdToBigInt(tx.NamespaceId)))
	var dV flatbuffers.UOffsetT
	if tx.NamespaceType == Root {
		dV = transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Duration))
	} else {
		dV = transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(namespaceIdToBigInt(tx.ParentId)))
	}
	n := builder.CreateString(tx.NamspaceName)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.RegisterNamespaceTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, 138+len(tx.NamspaceName))
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

func (dto *registerNamespaceTransactionDTO) toStruct() (*RegisterNamespaceTransaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	d := big.NewInt(0)
	n := &NamespaceId{}

	if dto.Tx.NamespaceType == Root {
		d = dto.Tx.Duration.toBigInt()
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

// LockFundsTransaction
type LockFundsTransaction struct {
	AbstractTransaction
	*Mosaic
	Duration *big.Int
	*SignedTransaction
}

func NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration *big.Int, signedTx *SignedTransaction, networkType NetworkType) (*LockFundsTransaction, error) {
	if mosaic == nil {
		return nil, errors.New("mosaic must not be nil")
	}
	if duration == nil {
		return nil, errors.New("duration must not be nil")
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
			"MosaicId": %s,
			"Duration": %d,
			"SignedTxHash": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaic.String(),
		tx.Duration,
		tx.SignedTransaction.Hash,
	)
}

func (tx *LockFundsTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mv := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(mosaicIdToBigInt(tx.Mosaic.MosaicId)))
	maV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Mosaic.Amount))
	dV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Duration))

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
	transactions.TransactionBufferAddSize(builder, 176)
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.LockFundsTransactionBufferAddMosaicId(builder, mv)
	transactions.LockFundsTransactionBufferAddMosaicAmount(builder, maV)
	transactions.LockFundsTransactionBufferAddDuration(builder, dV)
	transactions.LockFundsTransactionBufferAddHash(builder, hV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return lockFundsTransactionSchema().serialize(builder.FinishedBytes()), nil
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

func (dto *lockFundsTransactionDTO) toStruct() (*LockFundsTransaction, error) {
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
		dto.Tx.Duration.toBigInt(),
		&SignedTransaction{Lock, "", dto.Tx.Hash},
	}, nil
}

// SecretLockTransaction
type SecretLockTransaction struct {
	AbstractTransaction
	*Mosaic
	HashType
	Duration  *big.Int
	Secret    string
	Recipient *Address
}

func NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration *big.Int, hashType HashType, secret string, recipient *Address, networkType NetworkType) (*SecretLockTransaction, error) {
	if mosaic == nil {
		return nil, errors.New("mosaic must not be nil")
	}
	if duration == nil {
		return nil, errors.New("duration must not be nil")
	}
	if secret == "" {
		return nil, errors.New("secret must not be empty")
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
		HashType:  hashType,
		Secret:    secret, // TODO Add secret validation
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
			"MosaicId": %s,
			"Duration": %d,
			"HashType": %s,
			"Secret": %s,
			"Recipient": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaic.String(),
		tx.Duration,
		tx.HashType.String(),
		tx.Secret,
		tx.Recipient,
	)
}

func (tx *SecretLockTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(mosaicIdToBigInt(tx.Mosaic.MosaicId)))
	maV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Mosaic.Amount))
	dV := transactions.TransactionBufferCreateUint32Vector(builder, FromBigInt(tx.Duration))

	s, err := hex.DecodeString(tx.Secret)
	if err != nil {
		return nil, err
	}
	sV := transactions.TransactionBufferCreateByteVector(builder, s)

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
	transactions.TransactionBufferAddSize(builder, 202)
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.SecretLockTransactionBufferAddMosaicId(builder, mV)
	transactions.SecretLockTransactionBufferAddMosaicAmount(builder, maV)
	transactions.SecretLockTransactionBufferAddDuration(builder, dV)
	transactions.SecretLockTransactionBufferAddHashAlgorithm(builder, byte(tx.HashType))
	transactions.SecretLockTransactionBufferAddSecret(builder, sV)
	transactions.SecretLockTransactionBufferAddRecipient(builder, rV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return secretLockTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type secretLockTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		MosaicId  *uint64DTO `json:"mosaicId"`
		Amount    *uint64DTO `json:"amount"`
		HashType  `json:"hashAlgorithm"`
		Duration  uint64DTO `json:"duration"`
		Secret    string    `json:"secret"`
		Recipient string    `json:"recipient"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretLockTransactionDTO) toStruct() (*SecretLockTransaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromEncoded(dto.Tx.Recipient)
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	mosaic, err := NewMosaic(mosaicId, dto.Tx.Amount.toBigInt())
	if err != nil {
		return nil, err
	}

	return &SecretLockTransaction{
		*atx,
		mosaic,
		dto.Tx.HashType,
		dto.Tx.Duration.toBigInt(),
		dto.Tx.Secret,
		a,
	}, nil
}

// SecretProofTransaction
type SecretProofTransaction struct {
	AbstractTransaction
	HashType
	Secret string
	Proof  string
}

func NewSecretProofTransaction(deadline *Deadline, hashType HashType, secret string, proof string, networkType NetworkType) (*SecretProofTransaction, error) {
	if proof == "" {
		return nil, errors.New("proof must not be empty")
	}
	if secret == "" {
		return nil, errors.New("secret must not be empty")
	}

	return &SecretProofTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     SecretProofVersion,
			Deadline:    deadline,
			Type:        SecretProof,
			NetworkType: networkType,
		},
		HashType: hashType,
		Secret:   secret, // TODO Add secret validation
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
			"Secret": %s,
			"Proof": %s
		`,
		tx.AbstractTransaction.String(),
		tx.HashType.String(),
		tx.Secret,
		tx.Proof,
	)
}

func (tx *SecretProofTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	s, err := hex.DecodeString(tx.Secret)
	if err != nil {
		return nil, err
	}
	sV := transactions.TransactionBufferCreateByteVector(builder, s)

	p, err := hex.DecodeString(tx.Proof)
	if err != nil {
		return nil, err
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, p)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.SecretProofTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, 155+len(p))
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.SecretProofTransactionBufferAddHashAlgorithm(builder, byte(tx.HashType))
	transactions.SecretProofTransactionBufferAddSecret(builder, sV)
	transactions.SecretProofTransactionBufferAddProofSize(builder, uint16(len(p)))
	transactions.SecretProofTransactionBufferAddProof(builder, pV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return secretProofTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type secretProofTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		HashType `json:"hashAlgorithm"`
		Secret   string `json:"secret"`
		Proof    string `json:"proof"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *secretProofTransactionDTO) toStruct() (*SecretProofTransaction, error) {
	atx, err := dto.Tx.abstractTransactionDTO.toStruct(dto.TDto.toStruct())
	if err != nil {
		return nil, err
	}

	return &SecretProofTransaction{
		*atx,
		dto.Tx.HashType,
		dto.Tx.Secret,
		dto.Tx.Proof,
	}, nil
}

type CosignatureTransaction struct {
	TransactionToCosign *AggregateTransaction
}

func NewCosignatureTransaction(txToCosign *AggregateTransaction) (*CosignatureTransaction, error) {
	if txToCosign == nil {
		return nil, errors.New("txToCosign must not be nil")
	}
	return &CosignatureTransaction{txToCosign}, nil
}

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
	return fmt.Sprintf(`"TransactionToCosign": %s`, tx.TransactionToCosign.String())
}

// SignedTransaction
type SignedTransaction struct {
	TransactionType `json:"transactionType"`
	Payload         string `json:"payload"`
	Hash            Hash   `json:"hash"`
}

// CosignatureSignedTransaction
type CosignatureSignedTransaction struct {
	ParentHash Hash   `json:"parentHash"`
	Signature  string `json:"signature"`
	Signer     string `json:"signer"`
}

// AggregateTransactionCosignature
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

// MultisigCosignatoryModification
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
		m.Type.String(),
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

// MetadataModification
type MetadataModification struct {
	Type  MetadataModificationType
	Key   string
	Value string
}

func (m *MetadataModification) String() string {
	return fmt.Sprintf(
		`
			"Type"	: %s,
			"Key" 	: %s,
			"Value" : %s
		`,
		m.Type.String(),
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
	flags := "00" + dto[0].Value.toBigInt().Text(2)
	bitMapFlags := flags[len(flags)-3:]

	duration := big.NewInt(0)
	if len(dto) == 3 {
		duration = dto[2].Value.toBigInt()
	}

	return NewMosaicProperties(bitMapFlags[2] == '1',
		bitMapFlags[1] == '1',
		bitMapFlags[0] == '1',
		byte(dto[1].Value.toBigInt().Int64()),
		duration,
	)
}

// TransactionStatus
type TransactionStatus struct {
	Deadline *Deadline
	Group    string
	Status   string
	Hash     Hash
	Height   *big.Int
}

func (ts *TransactionStatus) String() string {
	return fmt.Sprintf(
		`
			"Group:" %s,
			"Status:" %s,
			"Content": %s,
			"Deadline": %s,
			"Height": %d
		`,
		ts.Group,
		ts.Status,
		ts.Hash,
		ts.Deadline,
		ts.Height,
	)
}

type transactionStatusDTO struct {
	Group    string    `json:"group"`
	Status   string    `json:"status"`
	Hash     Hash      `json:"hash"`
	Deadline uint64DTO `json:"deadline"`
	Height   uint64DTO `json:"height"`
}

func (dto *transactionStatusDTO) toStruct() (*TransactionStatus, error) {
	return &TransactionStatus{
		&Deadline{time.Unix(dto.Deadline.toBigInt().Int64(), int64(time.Millisecond))},
		dto.Group,
		dto.Status,
		dto.Hash,
		dto.Height.toBigInt(),
	}, nil
}

// TransactionIds
type TransactionIdsDTO struct {
	Ids []string `json:"transactionIds"`
}

// TransactionHashes
type TransactionHashesDTO struct {
	Hashes []string `json:"hashes"`
}

var TimestampNemesisBlock = time.Unix(1459468800, 0)

// Deadline
type Deadline struct {
	time.Time
}

func (d *Deadline) GetInstant() int64 {
	return (d.Time.UnixNano() / 1e6) - (TimestampNemesisBlock.UnixNano() / 1e6)
}

// Create deadline model
func NewDeadline(d time.Duration) *Deadline {
	return &Deadline{time.Now().Add(d)}
}

// Message
type Message struct {
	Type    uint8
	Payload string
}

// The transaction message of 1024 characters.
func NewPlainMessage(payload string) *Message {
	return &Message{0, payload}
}

func (m *Message) String() string {
	return str.StructToString(
		"Message",
		str.NewField("Type", str.IntPattern, m.Type),
		str.NewField("Payload", str.StringPattern, m.Payload),
	)
}

type messageDTO struct {
	Type    uint8  `json:"type"`
	Payload string `json:"payload"`
}

func (m *messageDTO) toStruct() *Message {
	b, err := hex.DecodeString(m.Payload)

	if err != nil {
		return &Message{0, ""}
	}

	return &Message{m.Type, string(b)}
}

type transactionTypeStruct struct {
	transactionType TransactionType
	raw             uint32
	hex             uint16
}

var transactionTypes = []transactionTypeStruct{
	{AggregateCompleted, 16705, 0x4141},
	{AggregateBonded, 16961, 0x4241},
	{MetadataAddress, 16701, 0x413d},
	{MetadataMosaic, 16957, 0x423d},
	{MetadataNamespace, 17213, 0x433d},
	{MosaicDefinition, 16717, 0x414d},
	{MosaicSupplyChange, 16973, 0x424d},
	{ModifyMultisig, 16725, 0x4155},
	{ModifyContract, 16727, 0x4157},
	{RegisterNamespace, 16718, 0x414e},
	{Transfer, 16724, 0x4154},
	{Lock, 16712, 0x4148},
	{SecretLock, 16722, 0x4152},
	{SecretProof, 16978, 0x4252},
}

type TransactionType uint16

// TransactionType enums
const (
	AggregateCompleted TransactionType = iota
	AggregateBonded
	MetadataAddress
	MetadataMosaic
	MetadataNamespace
	MosaicDefinition
	MosaicSupplyChange
	ModifyMultisig
	ModifyContract
	RegisterNamespace
	Transfer
	Lock
	SecretLock
	SecretProof
)

type TransactionVersion uint8

// TransactionVersion enums
const (
	AggregateCompletedVersion TransactionVersion = 2
	AggregateBondedVersion    TransactionVersion = 2
	MetadataAddressVersion    TransactionVersion = 1
	MetadataMosaicVersion     TransactionVersion = 1
	MetadataNamespaceVersion  TransactionVersion = 1
	MosaicDefinitionVersion   TransactionVersion = 3
	MosaicSupplyChangeVersion TransactionVersion = 2
	ModifyMultisigVersion     TransactionVersion = 3
	ModifyContractVersion     TransactionVersion = 3
	RegisterNamespaceVersion  TransactionVersion = 2
	TransferVersion           TransactionVersion = 3
	LockVersion               TransactionVersion = 1
	SecretLockVersion         TransactionVersion = 1
	SecretProofVersion        TransactionVersion = 1
)

func (t TransactionType) Hex() uint16 {
	return transactionTypes[t].hex
}

func (t TransactionType) Raw() uint32 {
	return transactionTypes[t].raw
}

func (t TransactionType) String() string {
	return fmt.Sprintf("%d", t.Raw())
}

// TransactionType error
var transactionTypeError = errors.New("wrong raw TransactionType int")

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

type HashType uint8

func (ht HashType) String() string {
	return fmt.Sprintf("%d", ht)
}

const SHA3_256 HashType = 0

func ExtractVersion(version uint64) uint8 {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, version)

	return uint8(b[0])
}

func TransactionTypeFromRaw(value uint32) (TransactionType, error) {
	for _, t := range transactionTypes {
		if t.raw == value {
			return t.transactionType, nil
		}
	}
	return 0, transactionTypeError
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

func MapTransaction(b *bytes.Buffer) (Transaction, error) {
	rawT := struct {
		Transaction struct {
			Type uint32
		}
	}{}

	err := json.Unmarshal(b.Bytes(), &rawT)
	if err != nil {
		return nil, err
	}

	t, err := TransactionTypeFromRaw(rawT.Transaction.Type)
	if err != nil {
		return nil, err
	}

	switch t {
	case AggregateBonded:
		return mapAggregateTransaction(b)
	case AggregateCompleted:
		return mapAggregateTransaction(b)
	case MetadataAddress:
		dto := modifyMetadataAddressTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case MetadataMosaic:
		dto := modifyMetadataMosaicTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case MetadataNamespace:
		dto := modifyMetadataNamespaceTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case MosaicDefinition:
		dto := mosaicDefinitionTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case MosaicSupplyChange:
		dto := mosaicSupplyChangeTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case ModifyMultisig:
		dto := modifyMultisigAccountTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case ModifyContract:
		dto := modifyContractTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case RegisterNamespace:
		dto := registerNamespaceTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case Transfer:
		dto := transferTransactionDTO{}
		err := json.Unmarshal(b.Bytes(), &dto)

		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case Lock:
		dto := lockFundsTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case SecretLock:
		dto := secretLockTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	case SecretProof:
		dto := secretProofTransactionDTO{}

		err := json.Unmarshal(b.Bytes(), &dto)
		if err != nil {
			return nil, err
		}

		tx, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		return tx, nil
	}

	return nil, nil
}

func mapAggregateTransaction(b *bytes.Buffer) (*AggregateTransaction, error) {
	dto := aggregateTransactionDTO{}

	err := json.Unmarshal(b.Bytes(), &dto)
	if err != nil {
		return nil, err
	}

	tx, err := dto.toStruct()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func createTransactionHash(p string) (string, error) {
	b, err := hex.DecodeString(p)
	if err != nil {
		return "", err
	}
	sb := make([]byte, len(b)-36)
	copy(sb[:32], b[4:32+4])
	copy(sb[32:], b[68:])

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

	rB := make([]byte, len(b)-64-16)
	copy(rB[4:32+4], sb[:32])
	copy(rB[32+4:32+4+4], b[100:104])
	copy(rB[32+4+4:32+4+4+len(b)-120], b[100+2+2+16:100+2+2+16+len(b)-120])

	s := big.NewInt(int64(len(b) - 64 - 16)).Bytes()
	utils.ReverseByteArray(s)

	copy(rB[:len(s)], s)

	return rB, nil
}

func signTransactionWith(tx Transaction, a *Account) (*SignedTransaction, error) {
	s := crypto.NewSignerFromKeyPair(a.KeyPair, nil)
	b, err := tx.generateBytes()
	if err != nil {
		return nil, err
	}
	sb := make([]byte, len(b)-100)
	copy(sb, b[100:])
	signature, err := s.Sign(sb)
	if err != nil {
		return nil, err
	}

	p := make([]byte, len(b))
	copy(p[:4], b[:4])
	copy(p[4:64+4], signature.Bytes())
	copy(p[64+4:64+4+32], a.KeyPair.PublicKey.Raw)
	copy(p[100:], b[100:])

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

	s := big.NewInt(int64(len(pb))).Bytes()
	utils.ReverseByteArray(s)

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

func metadataModificationArrayToBuffer(builder *flatbuffers.Builder, modifications []*MetadataModification) (flatbuffers.UOffsetT, uint32, error) {
	msb := make([]flatbuffers.UOffsetT, len(modifications))
	allSize := uint32(0)
	for i, m := range modifications {
		keySize := len(m.Key)

		if keySize == 0 {
			return 0, 0, errors.New("key must not empty")
		}

		pKey := transactions.TransactionBufferCreateByteVector(builder, []byte(m.Key))
		valueSize := len(m.Value)

		// it is hack, because we can have case when size of the value is zero(in RemoveData modification),
		// but flattbuffer doesn't store int(0) like 4 bytes, it stores like one byte
		valueB := make([]byte, 2)
		binary.LittleEndian.PutUint16(valueB, uint16(valueSize))
		pValueSize := transactions.TransactionBufferCreateByteVector(builder, valueB)

		pValue := transactions.TransactionBufferCreateByteVector(builder, []byte(m.Value))

		size := uint32(4 + 1 + 1 + 2 + keySize + valueSize)

		transactions.MetadataModificationBufferStart(builder)
		transactions.MetadataModificationBufferAddSize(builder, size)
		transactions.MetadataModificationBufferAddModificationType(builder, uint8(m.Type))
		transactions.MetadataModificationBufferAddKeySize(builder, uint8(keySize))
		transactions.MetadataModificationBufferAddValueSize(builder, pValueSize)
		transactions.MetadataModificationBufferAddKey(builder, pKey)
		transactions.MetadataModificationBufferAddValue(builder, pValue)

		msb[i] = transactions.MetadataModificationBufferEnd(builder)

		allSize = allSize + size
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), allSize, nil
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
