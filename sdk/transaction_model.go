// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	jsonLib "encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
	crypto "github.com/proximax-storage/go-xpx-crypto"
	utils "github.com/proximax-storage/go-xpx-utils"

	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

type Transaction interface {
	GetAbstractTransaction() *AbstractTransaction
	String() string
	// number of bytes of serialized transaction
	Size() int
	Bytes() ([]byte, error)
}

type transactionDto interface {
	toStruct(*Hash) (Transaction, error)
}

type AbstractTransaction struct {
	TransactionInfo
	NetworkType NetworkType    `json:"network_type"`
	Deadline    *Deadline      `json:"deadline"`
	Type        EntityType     `json:"entity_type"`
	Version     EntityVersion  `json:"version"`
	MaxFee      Amount         `json:"max_fee"`
	Signature   string         `json:"signature"`
	Signer      *PublicAccount `json:"signer"`
}

func (tx *AbstractTransaction) IsUnconfirmed() bool {
	return tx.TransactionInfo.Height == 0 && tx.TransactionInfo.TransactionHash.Equal(tx.TransactionInfo.MerkleComponentHash)
}

func (tx *AbstractTransaction) IsConfirmed() bool {
	return tx.TransactionInfo.Height > 0
}

func (tx *AbstractTransaction) HasMissingSignatures() bool {
	return tx.TransactionInfo.Height == 0 && !tx.TransactionInfo.TransactionHash.Equal(tx.TransactionInfo.MerkleComponentHash)
}

func (tx *AbstractTransaction) IsAnnounced() bool {
	return tx.TransactionInfo.TransactionHash != nil || tx.TransactionInfo.AggregateHash != nil
}

func (tx *AbstractTransaction) IsUnannounced() bool {
	return !tx.IsAnnounced()
}

func (tx *AbstractTransaction) ToAggregate(signer *Account) {
	tx.Signer = signer.PublicAccount
	derivationScheme := GetDerivationSchemeForAccountVersion(signer.Version).EngineDerivationScheme()
	tx.Version = EntityVersion(uint32(tx.Version) | (uint32(derivationScheme) << 16))
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
		tx.TransactionInfo.String(),
		tx.Type,
		tx.Version,
		tx.MaxFee,
		tx.Deadline,
		tx.Signature,
		tx.Signer,
	)
}

func (tx *AbstractTransaction) generateVectors(builder *flatbuffers.Builder) (v uint32, signatureV, signerV, dV, fV flatbuffers.UOffsetT, err error) {
	v = (uint32(tx.NetworkType) << 24) + uint32(tx.Version)
	signatureV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, SignatureSize))
	signerV = transactions.TransactionBufferCreateByteVector(builder, make([]byte, SignerSize))
	dV = transactions.TransactionBufferCreateUint32Vector(builder, tx.Deadline.ToBlockchainTimestamp().toArray())
	fV = transactions.TransactionBufferCreateUint32Vector(builder, tx.MaxFee.toArray())
	return
}

func (tx *AbstractTransaction) buildVectors(builder *flatbuffers.Builder, v uint32, signatureV, signerV, dV, fV flatbuffers.UOffsetT) {
	transactions.TransactionBufferAddSignature(builder, signatureV)
	transactions.TransactionBufferAddSigner(builder, signerV)
	transactions.TransactionBufferAddVersion(builder, v)
	transactions.TransactionBufferAddType(builder, uint16(tx.Type))
	transactions.TransactionBufferAddMaxFee(builder, fV)
	transactions.TransactionBufferAddDeadline(builder, dV)
}

type Pagination struct {
	TotalEntries uint64
	PageNumber   uint64
	PageSize     uint64
	TotalPages   uint64
}

type TransactionsPage struct {
	Transactions []Transaction
	Pagination   Pagination
}

type TransactionsPageOptions struct {
	Height           uint   `url:"height,omitempty"`
	FromHeight       uint64 `url:"fromHeight,omitempty"`
	ToHeight         uint64 `url:"toHeight,omitempty"`
	Address          string `url:"address,omitempty"`
	SignerPublicKey  string `url:"signerPublicKey,omitempty"`
	RecipientAddress string `url:"recipientAddress,omitempty"`
	Type             []uint `url:"type[],omitempty"`
	Embedded         bool   `url:"embedded,omitempty"`
	PublicKey        bool   `url:"publicKey,omitempty"`
	PaginationOrderingOptions
}

type TransactionInfo struct {
	Height              Height
	Index               uint32
	Id                  string
	TransactionHash     *Hash
	MerkleComponentHash *Hash
	AggregateHash       *Hash
	UniqueAggregateHash *Hash
	AggregateId         string
}

func (ti *TransactionInfo) String() string {
	return fmt.Sprintf(
		`
			"Height": %s,
			"Index": %d,
			"Id": %s,
			"TransactionHash": %s,
			"MerkleComponentHash:" %s,
			"AggregateHash": %s,
			"UniqueAggregateHash": %s,
			"AggregateId": %s
		`,
		ti.Height,
		ti.Index,
		ti.Id,
		ti.TransactionHash,
		ti.MerkleComponentHash,
		ti.AggregateHash,
		ti.UniqueAggregateHash,
		ti.AggregateId,
	)
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

func (tx *AccountPropertiesAddressTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	msb := make([]flatbuffers.UOffsetT, len(tx.Modifications))
	for i, m := range tx.Modifications {
		a, err := m.Address.Decode()
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

type AccountPropertiesMosaicModification struct {
	ModificationType PropertyModificationType
	AssetId          AssetId
}

func (mod *AccountPropertiesMosaicModification) String() string {
	return fmt.Sprintf(
		`
			"ModificationType": %d,
			"AssetId": %s,
		`,
		mod.ModificationType,
		mod.AssetId,
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

func (tx *AccountPropertiesMosaicTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	msb := make([]flatbuffers.UOffsetT, len(tx.Modifications))
	for i, m := range tx.Modifications {
		mosaicB := make([]byte, MosaicIdSize)
		binary.LittleEndian.PutUint64(mosaicB, m.AssetId.Id())
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

type AccountPropertiesEntityTypeModification struct {
	ModificationType PropertyModificationType
	EntityType       EntityType
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

func (tx *AccountPropertiesEntityTypeTransaction) Bytes() ([]byte, error) {
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

func (tx *AliasTransaction) Bytes(builder *flatbuffers.Builder, aliasV flatbuffers.UOffsetT, sizeOfAlias int) ([]byte, error) {
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
	return AliasTransactionHeaderSize
}

func (tx *AliasTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
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

func (tx *AddressAliasTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	a, err := tx.Address.Decode()
	if err != nil {
		return nil, err
	}

	aV := transactions.TransactionBufferCreateByteVector(builder, a)

	return tx.AliasTransaction.Bytes(builder, aV, AddressSize)
}

func (tx *AddressAliasTransaction) Size() int {
	return tx.AliasTransaction.Size() + AddressSize
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

func (tx *MosaicAliasTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, MosaicIdSize)
	binary.LittleEndian.PutUint64(mosaicB, tx.MosaicId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.AliasTransaction.Bytes(builder, mV, MosaicIdSize)
}

func (tx *MosaicAliasTransaction) Size() int {
	return tx.AliasTransaction.Size() + MosaicIdSize
}

type NodeKeyLinkTransaction struct {
	AbstractTransaction
	NodeKey    string
	LinkAction AccountLinkAction
}

// returns NodeKeyLinkTransaction from passed PublicAccount and NodeLinkAction
func NewNodeKeyLinkTransaction(deadline *Deadline, remoteAccount string, linkAction AccountLinkAction, networkType NetworkType) (*NodeKeyLinkTransaction, error) {
	if len(remoteAccount) == 0 {
		return nil, errors.New("remoteAccount must not be empty")
	}
	return &NodeKeyLinkTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        NodeKeyLink,
			Version:     NodeKeyLinkVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		NodeKey:    remoteAccount,
		LinkAction: linkAction,
	}, nil
}

func (tx *NodeKeyLinkTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *NodeKeyLinkTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RemoteAccount": %s,
			"LinkAction": %d
		`,
		tx.AbstractTransaction.String(),
		tx.NodeKey,
		tx.LinkAction,
	)
}

func (tx *NodeKeyLinkTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	b, err := utils.HexDecodeStringOdd(tx.NodeKey)
	if err != nil {
		return nil, err
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, b)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.NodeLinkTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)
	transactions.NodeLinkTransactionBufferAddRemoteAccountKey(builder, pV)
	transactions.NodeLinkTransactionBufferAddLinkAction(builder, uint8(tx.LinkAction))
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return nodeKeyLinkTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *NodeKeyLinkTransaction) Size() int {
	return NodeKeyLinkTransactionSize
}

type VrfKeyLinkTransaction struct {
	AbstractTransaction
	VrfAccount *PublicAccount
	LinkAction AccountLinkAction
}

// returns VrfKeyLinkTransaction from passed PublicAccount and LinkAction
func NewVrfKeyLinkTransaction(deadline *Deadline, vrfAccount *PublicAccount, linkAction AccountLinkAction, networkType NetworkType) (*VrfKeyLinkTransaction, error) {
	if vrfAccount == nil {
		return nil, errors.New("vrfAccount must not be empty")
	}
	return &VrfKeyLinkTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        VrfKeyLink,
			Version:     VrfKeyLinkVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		VrfAccount: vrfAccount,
		LinkAction: linkAction,
	}, nil
}

func (tx *VrfKeyLinkTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *VrfKeyLinkTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"VrfAccount": %s,
			"LinkAction": %d
		`,
		tx.AbstractTransaction.String(),
		tx.VrfAccount.String(),
		tx.LinkAction,
	)
}

func (tx *VrfKeyLinkTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	b, err := utils.HexDecodeStringOdd(tx.VrfAccount.PublicKey)
	if err != nil {
		return nil, err
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, b)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.VrfLinkTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)
	transactions.VrfLinkTransactionBufferAddRemoteAccountKey(builder, pV)
	transactions.VrfLinkTransactionBufferAddLinkAction(builder, uint8(tx.LinkAction))
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return vrfKeyLinkTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *VrfKeyLinkTransaction) Size() int {
	return VrfKeyLinkTransactionSize
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

func (tx *AccountLinkTransaction) Bytes() ([]byte, error) {
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

type NetworkConfigTransaction struct {
	AbstractTransaction
	ApplyHeightDelta  Duration
	NetworkConfig     *NetworkConfig
	SupportedEntities *SupportedEntities
}

// returns NetworkConfigTransaction from passed ApplyHeightDelta, NetworkConfig and SupportedEntities
func NewNetworkConfigTransaction(deadline *Deadline, delta Duration, config *NetworkConfig, entities *SupportedEntities, networkType NetworkType) (*NetworkConfigTransaction, error) {
	if entities == nil {
		return nil, errors.New("Entities should not be nil")
	}
	if config == nil {
		return nil, errors.New("NetworkConfig should not be nil")
	}

	return &NetworkConfigTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        NetworkConfigEntityType,
			Version:     NetworkConfigVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		ApplyHeightDelta:  delta,
		NetworkConfig:     config,
		SupportedEntities: entities,
	}, nil
}

func (tx *NetworkConfigTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *NetworkConfigTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"ApplyHeightDelta": %s,
			"NetworkConfig": %s,
			"SupportedEntities": %s
		`,
		tx.AbstractTransaction.String(),
		tx.ApplyHeightDelta,
		tx.NetworkConfig,
		tx.SupportedEntities,
	)
}

func (tx *NetworkConfigTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	sup, err := tx.SupportedEntities.MarshalBinary()
	if err != nil {
		return nil, err
	}

	config, err := tx.NetworkConfig.MarshalBinary()
	if err != nil {
		return nil, err
	}
	deltaV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ApplyHeightDelta.toArray())
	configV := transactions.TransactionBufferCreateByteVector(builder, config)
	supportedV := transactions.TransactionBufferCreateByteVector(builder, sup)

	transactions.NetworkConfigTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)

	transactions.NetworkConfigTransactionBufferAddApplyHeightDelta(builder, deltaV)
	transactions.NetworkConfigTransactionBufferAddNetworkConfigSize(builder, uint16(len(config)))
	transactions.NetworkConfigTransactionBufferAddNetworkConfig(builder, configV)
	transactions.NetworkConfigTransactionBufferAddSupportedEntityVersionsSize(builder, uint16(len(sup)))
	transactions.NetworkConfigTransactionBufferAddSupportedEntityVersions(builder, supportedV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return networkConfigTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *NetworkConfigTransaction) Size() int {
	return NetworkConfigHeaderSize + len(tx.NetworkConfig.String()) + len(tx.SupportedEntities.String())
}

type NetworkConfigAbsoluteHeightTransaction struct {
	AbstractTransaction
	ApplyHeight       Height
	NetworkConfig     *NetworkConfig
	SupportedEntities *SupportedEntities
}

// returns NetworkConfigAbsoluteHeightTransaction from passed ApplyHeight, NetworkConfig and SupportedEntities
func NewNetworkConfigAbsoluteHeightTransaction(deadline *Deadline, height Height, config *NetworkConfig, entities *SupportedEntities, networkType NetworkType) (*NetworkConfigAbsoluteHeightTransaction, error) {
	if entities == nil {
		return nil, errors.New("Entities should not be nil")
	}
	if config == nil {
		return nil, errors.New("NetworkConfig should not be nil")
	}

	return &NetworkConfigAbsoluteHeightTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        NetworkConfigAbsoluteHeightEntityType,
			Version:     NetworkConfigAbsoluteHeightVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		ApplyHeight:       height,
		NetworkConfig:     config,
		SupportedEntities: entities,
	}, nil
}

func (tx *NetworkConfigAbsoluteHeightTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *NetworkConfigAbsoluteHeightTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"ApplyHeightDelta": %s,
			"NetworkConfig": %s,
			"SupportedEntities": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Height,
		tx.NetworkConfig,
		tx.SupportedEntities,
	)
}

func (tx *NetworkConfigAbsoluteHeightTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	sup, err := tx.SupportedEntities.MarshalBinary()
	if err != nil {
		return nil, err
	}

	config, err := tx.NetworkConfig.MarshalBinary()
	if err != nil {
		return nil, err
	}
	heightV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ApplyHeight.toArray())
	configV := transactions.TransactionBufferCreateByteVector(builder, config)
	supportedV := transactions.TransactionBufferCreateByteVector(builder, sup)

	transactions.NetworkConfigAbsoluteHeightTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)

	transactions.NetworkConfigAbsoluteHeightTransactionBufferAddApplyHeight(builder, heightV)
	transactions.NetworkConfigAbsoluteHeightTransactionBufferAddNetworkConfigSize(builder, uint16(len(config)))
	transactions.NetworkConfigAbsoluteHeightTransactionBufferAddNetworkConfig(builder, configV)
	transactions.NetworkConfigAbsoluteHeightTransactionBufferAddSupportedEntityVersionsSize(builder, uint16(len(sup)))
	transactions.NetworkConfigAbsoluteHeightTransactionBufferAddSupportedEntityVersions(builder, supportedV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return networkConfigTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *NetworkConfigAbsoluteHeightTransaction) Size() int {
	return NetworkConfigAbsoluteHeightHeaderSize + len(tx.NetworkConfig.String()) + len(tx.SupportedEntities.String())
}

type BlockchainUpgradeTransaction struct {
	AbstractTransaction
	UpgradePeriod        Duration
	NewBlockChainVersion BlockChainVersion
}

// returns NetworkConfigTransaction from passed ApplyHeightDelta, NetworkConfig and SupportedEntityVersions
func NewBlockchainUpgradeTransaction(deadline *Deadline, upgradePeriod Duration, newBlockChainVersion BlockChainVersion, networkType NetworkType) (*BlockchainUpgradeTransaction, error) {
	return &BlockchainUpgradeTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        BlockchainUpgrade,
			Version:     BlockchainUpgradeVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		UpgradePeriod:        upgradePeriod,
		NewBlockChainVersion: newBlockChainVersion,
	}, nil
}

func (tx *BlockchainUpgradeTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *BlockchainUpgradeTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"UpgradePeriod": %s,
			"NewBlockChainVersion": %s
		`,
		tx.AbstractTransaction.String(),
		tx.UpgradePeriod,
		tx.NewBlockChainVersion,
	)
}

func (tx *BlockchainUpgradeTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	upgradeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.UpgradePeriod.toArray())
	versionV := transactions.TransactionBufferCreateUint32Vector(builder, tx.NewBlockChainVersion.toArray())

	transactions.BlockchainUpgradeTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)

	transactions.BlockchainUpgradeTransactionBufferAddUpgradePeriod(builder, upgradeV)
	transactions.BlockchainUpgradeTransactionBufferAddNewBlockChainVersion(builder, versionV)
	t := transactions.NetworkConfigTransactionBufferEnd(builder)
	builder.Finish(t)

	return blockchainUpgradeTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *BlockchainUpgradeTransaction) Size() int {
	return BlockchainUpgradeTransactionSize
}

type AccountV2UpgradeTransaction struct {
	AbstractTransaction
	NewAccountPublicKey *PublicAccount
}

// returns NetworkConfigTransaction from passed ApplyHeightDelta, NetworkConfig and SupportedEntityVersions
func NewAccountV2UpgradeTransaction(deadline *Deadline, newAccountPublicKey *PublicAccount, networkType NetworkType) (*AccountV2UpgradeTransaction, error) {
	return &AccountV2UpgradeTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        AccountV2Upgrade,
			Version:     AccountV2UpgradeVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		NewAccountPublicKey: newAccountPublicKey,
	}, nil
}

func (tx *AccountV2UpgradeTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountV2UpgradeTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"NewAccountPublicKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.NewAccountPublicKey,
	)
}

func (tx *AccountV2UpgradeTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, dV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := utils.HexDecodeStringOdd(tx.NewAccountPublicKey.PublicKey)
	if err != nil {
		return nil, err
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, b)

	transactions.AccountV2UpgradeTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, dV, fV)

	transactions.AccountV2UpgradeTransactionBufferAddNewaccountpublickey(builder, pV)
	t := transactions.AccountV2UpgradeTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountV2UpgradeTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountV2UpgradeTransaction) Size() int {
	return AccountV2UpgradeTransactionSize
}

type AggregateTransactionV1 struct {
	AbstractTransaction
	InnerTransactions []Transaction
	Cosignatures      []*AggregateTransactionCosignature
}

type AggregateTransactionV2 struct {
	AbstractTransaction
	InnerTransactions []Transaction
	Cosignatures      []*AggregateTransactionCosignature
}

// returns complete AggregateTransaction from passed array of own Transaction's to be included in
func NewCompleteAggregateV1Transaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransactionV1, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransactionV1{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateCompletedV1,
			Version:     AggregateCompletedV1Version,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

// returns bounded AggregateTransaction from passed array of transactions to be included in
func NewBondedAggregateV1Transaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransactionV1, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransactionV1{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateBondedV1,
			Version:     AggregateBondedV1Version,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

// returns complete AggregateTransaction from passed array of own Transaction's to be included in
func NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransactionV2, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransactionV2{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateCompletedV2,
			Version:     AggregateCompletedV2Version,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

// returns bounded AggregateTransaction from passed array of transactions to be included in
func NewBondedAggregateTransaction(deadline *Deadline, innerTxs []Transaction, networkType NetworkType) (*AggregateTransactionV2, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}
	return &AggregateTransactionV2{
		AbstractTransaction: AbstractTransaction{
			Type:        AggregateBondedV2,
			Version:     AggregateBondedV2Version,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		InnerTransactions: innerTxs,
	}, nil
}

func (tx *AggregateTransactionV1) UpdateUniqueAggregateHash(generationHash *Hash) (err error) {
	for _, innerTx := range tx.InnerTransactions {
		innerTx.GetAbstractTransaction().UniqueAggregateHash, err = UniqueAggregateHashV1(tx, innerTx, generationHash)
		if err != nil {
			break
		}
	}

	return err
}
func (tx *AggregateTransactionV2) UpdateUniqueAggregateHash(generationHash *Hash) (err error) {
	for _, innerTx := range tx.InnerTransactions {
		innerTx.GetAbstractTransaction().UniqueAggregateHash, err = UniqueAggregateHashV2(tx, innerTx, generationHash)
		if err != nil {
			break
		}
	}

	return err
}

func CompareInnerTransaction(left []Transaction, right []Transaction) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if !InnerTransactionHash(left[i]).Equal(InnerTransactionHash(right[i])) {
			return false
		}
	}

	return true
}

func (tx *AggregateTransactionV1) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AggregateTransactionV2) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AggregateTransactionV1) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"InnerTransactions": %+v,
			"Cosignatures": %s
		`,
		tx.AbstractTransaction.String(),
		tx.InnerTransactions,
		tx.Cosignatures,
	)
}

func (tx *AggregateTransactionV2) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"InnerTransactions": %+v,
			"Cosignatures": %s
		`,
		tx.AbstractTransaction.String(),
		tx.InnerTransactions,
		tx.Cosignatures,
	)
}

func (tx *AggregateTransactionV1) Bytes() ([]byte, error) {
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

func (tx *AggregateTransactionV2) Bytes() ([]byte, error) {
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

func (tx *AggregateTransactionV1) Size() int {
	sizeOfInnerTransactions := 0
	for _, itx := range tx.InnerTransactions {
		sizeOfInnerTransactions += itx.Size() - SignatureSize - MaxFeeSize - DeadLineSize
	}
	return AggregateBondedHeaderSize + sizeOfInnerTransactions
}

func (tx *AggregateTransactionV2) Size() int {
	sizeOfInnerTransactions := 0
	for _, itx := range tx.InnerTransactions {
		sizeOfInnerTransactions += itx.Size() - SignatureSize - MaxFeeSize - DeadLineSize
	}
	return AggregateBondedHeaderSize + sizeOfInnerTransactions
}

type BasicMetadataTransaction struct {
	AbstractTransaction
	TargetPublicAccount *PublicAccount
	ScopedMetadataKey   ScopedMetadataKey
	Value               []byte
	ValueDeltaSize      int16
}

func (tx *BasicMetadataTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"TargetPublicAccount": %s,
			"ScopedMetadataKey": %s,
			"Value": %s
		`,
		tx.AbstractTransaction.String(),
		tx.TargetPublicAccount,
		tx.ScopedMetadataKey,
		tx.Value,
	)
}

func (tx *BasicMetadataTransaction) Bytes(builder *flatbuffers.Builder, targetIdV flatbuffers.UOffsetT, size int) ([]byte, error) {

	targetKeyB, err := hex.DecodeString(tx.TargetPublicAccount.PublicKey)
	if err != nil {
		return nil, err
	}

	targetKeyV := transactions.TransactionBufferCreateByteVector(builder, targetKeyB)

	metadataKeyV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ScopedMetadataKey.toArray())
	valueV := transactions.TransactionBufferCreateByteVector(builder, tx.Value)

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(len(tx.Value)))
	valueSizeV := transactions.TransactionBufferCreateByteVector(builder, buf)

	buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(tx.ValueDeltaSize))
	valueDeltaSizeV := transactions.TransactionBufferCreateByteVector(builder, buf)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MetadataV2TransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, size)

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MetadataV2TransactionBufferAddTargetKey(builder, targetKeyV)
	transactions.MetadataV2TransactionBufferAddScopedMetadataKey(builder, metadataKeyV)
	transactions.MetadataV2TransactionBufferAddTargetId(builder, targetIdV)
	transactions.MetadataV2TransactionBufferAddValueSizeDelta(builder, valueDeltaSizeV)
	transactions.MetadataV2TransactionBufferAddValueSize(builder, valueSizeV)
	transactions.MetadataV2TransactionBufferAddValue(builder, valueV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return metadataTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *BasicMetadataTransaction) Size() int {
	return MetadataV2HeaderSize + len(tx.Value)
}

func (tx *BasicMetadataTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

type AccountMetadataTransaction struct {
	BasicMetadataTransaction
}

func NewAccountMetadataTransaction(deadline *Deadline,
	account *PublicAccount, scopedKey ScopedMetadataKey,
	newValue string, oldValue string, networkType NetworkType) (*AccountMetadataTransaction, error) {
	if newValue == oldValue {
		return nil, errors.New("new value is the same")
	}

	mmatx := AccountMetadataTransaction{
		BasicMetadataTransaction: BasicMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     AccountMetadataVersion,
				Deadline:    deadline,
				Type:        AccountMetadata,
				NetworkType: networkType,
			},
			TargetPublicAccount: account,
			ScopedMetadataKey:   scopedKey,
			ValueDeltaSize:      int16(len(newValue)) - int16(len(oldValue)),
		},
	}

	if len(newValue) < len(oldValue) {
		oldValue, newValue = newValue, oldValue
	}

	value := make([]byte, len(newValue))
	copy(value, newValue)
	for i, c := range []byte(oldValue) {
		value[i] ^= c
	}

	mmatx.Value = value

	return &mmatx, nil
}

func (tx *AccountMetadataTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	emptyB := make([]byte, 0)
	aV := transactions.TransactionBufferCreateByteVector(builder, emptyB)
	return tx.BasicMetadataTransaction.Bytes(builder, aV, tx.Size())
}

type MosaicMetadataTransaction struct {
	BasicMetadataTransaction
	TargetMosaicId *MosaicId
}

func NewMosaicMetadataTransaction(deadline *Deadline,
	mosaic *MosaicId, account *PublicAccount, scopedKey ScopedMetadataKey,
	newValue string, oldValue string, networkType NetworkType) (*MosaicMetadataTransaction, error) {
	if newValue == oldValue {
		return nil, errors.New("new value is the same")
	}

	mmatx := MosaicMetadataTransaction{
		BasicMetadataTransaction: BasicMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     MosaicMetadataVersion,
				Deadline:    deadline,
				Type:        MosaicMetadata,
				NetworkType: networkType,
			},
			TargetPublicAccount: account,
			ScopedMetadataKey:   scopedKey,
			ValueDeltaSize:      int16(len(newValue)) - int16(len(oldValue)),
		},
		TargetMosaicId: mosaic,
	}

	if len(newValue) < len(oldValue) {
		oldValue, newValue = newValue, oldValue
	}

	value := make([]byte, len(newValue))
	copy(value, newValue)
	for i, c := range []byte(oldValue) {
		value[i] ^= c
	}

	mmatx.Value = value

	return &mmatx, nil
}

func (tx *MosaicMetadataTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, MosaicIdSize)
	binary.LittleEndian.PutUint64(mosaicB, tx.TargetMosaicId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)
	return tx.BasicMetadataTransaction.Bytes(builder, mV, tx.Size())
}

func (tx *MosaicMetadataTransaction) Size() int {
	return MetadataV2HeaderSize + len(tx.Value) + MosaicIdSize
}

func (tx *MosaicMetadataTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"TargetPublicAccount": %s,
			"ScopedMetadataKey": %s,
			"TargetMosaicId": %s,
			"Value": %s
		`,
		tx.AbstractTransaction.String(),
		tx.TargetPublicAccount,
		tx.ScopedMetadataKey,
		tx.TargetMosaicId,
		tx.Value,
	)
}

type NamespaceMetadataTransaction struct {
	BasicMetadataTransaction
	TargetNamespaceId *NamespaceId
}

func NewNamespaceMetadataTransaction(deadline *Deadline,
	namespace *NamespaceId, account *PublicAccount, scopedKey ScopedMetadataKey,
	newValue string, oldValue string, networkType NetworkType) (*NamespaceMetadataTransaction, error) {
	if newValue == oldValue {
		return nil, errors.New("new value is the same")
	}

	mmatx := NamespaceMetadataTransaction{
		BasicMetadataTransaction: BasicMetadataTransaction{
			AbstractTransaction: AbstractTransaction{
				Version:     NamespaceMetadataVersion,
				Deadline:    deadline,
				Type:        NamespaceMetadata,
				NetworkType: networkType,
			},
			TargetPublicAccount: account,
			ScopedMetadataKey:   scopedKey,
			ValueDeltaSize:      int16(len(newValue)) - int16(len(oldValue)),
		},
		TargetNamespaceId: namespace,
	}

	if len(newValue) < len(oldValue) {
		oldValue, newValue = newValue, oldValue
	}

	value := make([]byte, len(newValue))
	copy(value, newValue)
	for i, c := range []byte(oldValue) {
		value[i] ^= c
	}

	mmatx.Value = value

	return &mmatx, nil
}

func (tx *NamespaceMetadataTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	namespaceB := make([]byte, NamespaceSize)
	binary.LittleEndian.PutUint64(namespaceB, tx.TargetNamespaceId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, namespaceB)
	return tx.BasicMetadataTransaction.Bytes(builder, mV, tx.Size())
}

func (tx *NamespaceMetadataTransaction) Size() int {
	return MetadataV2HeaderSize + len(tx.Value) + NamespaceSize
}

func (tx *NamespaceMetadataTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"TargetPublicAccount": %s,
			"ScopedMetadataKey": %s,
			"TargetNamespaceId": %s,
			"Value": %s
		`,
		tx.AbstractTransaction.String(),
		tx.TargetPublicAccount,
		tx.ScopedMetadataKey,
		tx.TargetNamespaceId,
		tx.Value,
	)
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

func (tx *ModifyMetadataTransaction) Bytes(builder *flatbuffers.Builder, metadataV flatbuffers.UOffsetT, sizeOfMetadata int) ([]byte, error) {

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

func (tx *ModifyMetadataAddressTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	a, err := tx.Address.Decode()
	if err != nil {
		return nil, err
	}

	aV := transactions.TransactionBufferCreateByteVector(builder, a)

	return tx.ModifyMetadataTransaction.Bytes(builder, aV, AddressSize)
}

func (tx *ModifyMetadataAddressTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + AddressSize
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

func (tx *ModifyMetadataMosaicTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	mosaicB := make([]byte, MosaicIdSize)
	binary.LittleEndian.PutUint64(mosaicB, tx.MosaicId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, mosaicB)

	return tx.ModifyMetadataTransaction.Bytes(builder, mV, MosaicIdSize)
}

func (tx *ModifyMetadataMosaicTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + MosaicIdSize
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

func (tx *ModifyMetadataNamespaceTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	namespaceB := make([]byte, NamespaceSize)
	binary.LittleEndian.PutUint64(namespaceB, tx.NamespaceId.Id())
	mV := transactions.TransactionBufferCreateByteVector(builder, namespaceB)

	return tx.ModifyMetadataTransaction.Bytes(builder, mV, NamespaceSize)
}

func (tx *ModifyMetadataNamespaceTransaction) Size() int {
	return tx.ModifyMetadataTransaction.Size() + NamespaceSize
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

	for _, p := range mosaicProps.OptionalProperties {
		if p.Value == 0 {
			return nil, errors.New("duration can't be zero")
		}
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

func (tx *MosaicDefinitionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)
	var f uint8 = 0
	if tx.MosaicProperties.SupplyMutable {
		f += Supply_Mutable
	}
	if tx.MosaicProperties.Transferable {
		f += Transferable
	}
	if tx.MosaicProperties.Restrictable {
		f += Restrictable
	}
	if tx.MosaicProperties.SupplyForceImmutable {
		f += Supply_Force_Immutable
	}
	if tx.MosaicProperties.DisableLocking {
		f += Disable_Locking
	}

	nonceB := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceB, tx.MosaicNonce)
	nonceV := transactions.TransactionBufferCreateByteVector(builder, nonceB)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())
	pV := mosaicPropertyArrayToBuffer(builder, tx.MosaicProperties.OptionalProperties)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.MosaicDefinitionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicDefinitionTransactionBufferAddMosaicNonce(builder, nonceV)
	transactions.MosaicDefinitionTransactionBufferAddMosaicId(builder, mV)
	transactions.MosaicDefinitionTransactionBufferAddFlags(builder, f)
	transactions.MosaicDefinitionTransactionBufferAddDivisibility(builder, tx.MosaicProperties.Divisibility)
	transactions.MosaicDefinitionTransactionBufferAddNumOptionalProperties(builder, byte(len(tx.MosaicProperties.OptionalProperties)))
	transactions.MosaicDefinitionTransactionBufferAddOptionalProperties(builder, pV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)
	return mosaicDefinitionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicDefinitionTransaction) Size() int {
	return MosaicDefinitionTransactionHeaderSize + len(tx.OptionalProperties)*MosaicOptionalPropertySize
}

type MosaicSupplyChangeTransaction struct {
	AbstractTransaction
	MosaicSupplyType
	AssetId
	Delta Amount
}

// returns MosaicSupplyChangeTransaction from passed AssetId, MosaicSupplyTypeand supply delta
func NewMosaicSupplyChangeTransaction(deadline *Deadline, assetId AssetId, supplyType MosaicSupplyType, delta Duration, networkType NetworkType) (*MosaicSupplyChangeTransaction, error) {
	if assetId == nil || assetId.Id() == 0 {
		return nil, ErrNilAssetId
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
		AssetId:          assetId,
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
			"AssetId": %s,
			"Delta": %d
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicSupplyType,
		tx.AssetId,
		tx.Delta,
	)
}

func (tx *MosaicSupplyChangeTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.AssetId.toArray())
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

// / region modify mosaic levy implementation
type MosaicModifyLevyTransaction struct {
	AbstractTransaction
	*MosaicId
	*MosaicLevy
}

func NewMosaicModifyLevyTransaction(deadline *Deadline, networkType NetworkType, mosaicId *MosaicId, levy *MosaicLevy) (*MosaicModifyLevyTransaction, error) {
	if levy.MosaicId == nil {
		levy.MosaicId = mosaicId
	}

	return &MosaicModifyLevyTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        MosaicModifyLevy,
			Version:     MosaicModifyLevyVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		MosaicId:   mosaicId,
		MosaicLevy: levy,
	}, nil
}

func (tx *MosaicModifyLevyTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *MosaicModifyLevyTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicId": %s,
			"MosaicLevy": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicId.String(),
		tx.MosaicLevy.String(),
	)
}

func (tx *MosaicModifyLevyTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	var r []byte
	if len(tx.MosaicLevy.Recipient.Address) > 0 {
		r, err = tx.MosaicLevy.Recipient.Decode()
		if err != nil {
			return nil, err
		}
	} else {
		r = make([]byte, AddressSize)
	}

	mL := tx.MosaicLevy.SetBuffers(builder, r)

	transactions.ModifyMosaicLevyTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ModifyMosaicLevyTransactionBufferAddMosaicId(builder, mV)
	transactions.ModifyMosaicLevyTransactionBufferAddLevy(builder, mL)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)
	return mosaicModifyLevyTransactionScheme().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicModifyLevyTransaction) Size() int {
	return MosaicModifyLevyTransactionSize
}

/// end region modify mosaic levy

// / region remove mosaic levy
type MosaicRemoveLevyTransaction struct {
	AbstractTransaction
	*MosaicId
}

func NewMosaicRemoveLevyTransaction(deadline *Deadline, networkType NetworkType, mosaicId *MosaicId) (*MosaicRemoveLevyTransaction, error) {
	return &MosaicRemoveLevyTransaction{
		AbstractTransaction: AbstractTransaction{
			Type:        MosaicRemoveLevy,
			Version:     MosaicRemoveLevyVersion,
			Deadline:    deadline,
			NetworkType: networkType,
		},
		MosaicId: mosaicId,
	}, nil
}

func (tx *MosaicRemoveLevyTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *MosaicRemoveLevyTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicId": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicId.String(),
	)
}

func (tx *MosaicRemoveLevyTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.RemoveMosaicLevyTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())

	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.RemoveMosaicLevyTransactionBufferAddMosaicId(builder, mV)
	t := transactions.TransactionBufferEnd(builder)

	builder.Finish(t)
	return mosaicRemoveLevyTransactionScheme().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicRemoveLevyTransaction) Size() int {
	return MosaicRemoveLevyTransactionSize
}

/// end region

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

func (tx *TransferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	ml := len(tx.Mosaics)
	mb := make([]flatbuffers.UOffsetT, ml)
	for i, mos := range tx.Mosaics {
		id := transactions.TransactionBufferCreateUint32Vector(builder, mos.AssetId.toArray())
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

	r, err := tx.Recipient.Decode()
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
	return TransferHeaderSize + ((MosaicIdSize + AmountSize) * len(tx.Mosaics)) + tx.MessageSize()
}

func (tx *TransferTransaction) MessageSize() int {
	// Message + MessageType
	return len(tx.Message.Payload()) + 1
}

type HarvesterTransaction struct {
	AbstractTransaction
}

type HarvesterTransactionType EntityType

const (
	AddHarvester    = HarvesterTransactionType(AddHarvesterEntityType)
	RemoveHarvester = HarvesterTransactionType(RemoveHarvesterEntityType)
)

// HarvesterTransaction creates new Harvester transaction
func NewHarvesterTransaction(deadline *Deadline, htt HarvesterTransactionType, networkType NetworkType) (*HarvesterTransaction, error) {
	return &HarvesterTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     HarvesterVersion,
			Deadline:    deadline,
			Type:        EntityType(htt),
			NetworkType: networkType,
		},
	}, nil
}

func (tx *HarvesterTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *HarvesterTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *HarvesterTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.HarvesterTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return harvesterTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *HarvesterTransaction) Size() int {
	return TransactionHeaderSize
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

func (tx *ModifyMultisigAccountTransaction) Bytes() ([]byte, error) {
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

type ModifyContractTransaction struct {
	AbstractTransaction
	DurationDelta Duration
	Hash          *Hash
	Customers     []*MultisigCosignatoryModification
	Executors     []*MultisigCosignatoryModification
	Verifiers     []*MultisigCosignatoryModification
}

// returns ModifyContractTransaction from passed duration delta in blocks, file hash, arrays of customers, replicators and verificators
func NewModifyContractTransaction(
	deadline *Deadline, durationDelta Duration, hash *Hash,
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

func (tx *ModifyContractTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	durationV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DurationDelta.toArray())
	hashV := hashToBuffer(builder, tx.Hash)

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
			"Duration": %d,
			"ParentId": %s
		`,
		tx.AbstractTransaction.String(),
		tx.NamspaceName,
		tx.NamespaceId,
		tx.Duration,
		tx.ParentId,
	)
}

func (tx *RegisterNamespaceTransaction) Bytes() ([]byte, error) {
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

	if signedTx.EntityType != AggregateBondedV1 && signedTx.EntityType != AggregateBondedV2 {
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
			"SignedTransaction": %s
		`,
		tx.AbstractTransaction.String(),
		tx.Mosaic,
		tx.Duration,
		tx.SignedTransaction,
	)
}

func (tx *LockFundsTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mv := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.AssetId.toArray())
	maV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.Amount.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())
	hV := transactions.TransactionBufferCreateByteVector(builder, tx.SignedTransaction.Hash[:])

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
		return nil, ErrNilSecret
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

func (tx *SecretLockTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	mV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.AssetId.toArray())
	maV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Mosaic.Amount.toArray())
	dV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())

	sV := transactions.TransactionBufferCreateByteVector(builder, tx.Secret.Hash[:])

	addr, err := tx.Recipient.Decode()
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

type SecretProofTransaction struct {
	AbstractTransaction
	HashType
	Proof     *Proof
	Recipient *Address
}

// returns a SecretProofTransaction from passed HashType and Proof
func NewSecretProofTransaction(deadline *Deadline, hashType HashType, proof *Proof, recipient *Address, networkType NetworkType) (*SecretProofTransaction, error) {
	if proof == nil {
		return nil, ErrNilProof
	}
	if recipient == nil {
		return nil, ErrNilAddress
	}

	return &SecretProofTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     SecretProofVersion,
			Deadline:    deadline,
			Type:        SecretProof,
			NetworkType: networkType,
		},
		HashType:  hashType,
		Proof:     proof,
		Recipient: recipient,
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
			"Recipient": %s
		`,
		tx.AbstractTransaction.String(),
		tx.HashType,
		tx.Proof,
		tx.Recipient,
	)
}

func (tx *SecretProofTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	secret, err := tx.Proof.Secret(tx.HashType)
	if err != nil {
		return nil, err
	}
	sV := transactions.TransactionBufferCreateByteVector(builder, secret.Hash[:])
	pV := transactions.TransactionBufferCreateByteVector(builder, tx.Proof.Data)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	addr, err := tx.Recipient.Decode()
	if err != nil {
		return nil, err
	}
	rV := transactions.TransactionBufferCreateByteVector(builder, addr)

	transactions.SecretProofTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.SecretProofTransactionBufferAddHashAlgorithm(builder, byte(tx.HashType))
	transactions.SecretProofTransactionBufferAddSecret(builder, sV)
	transactions.SecretProofTransactionBufferAddRecipient(builder, rV)
	transactions.SecretProofTransactionBufferAddProofSize(builder, uint16(tx.Proof.Size()))
	transactions.SecretProofTransactionBufferAddProof(builder, pV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return secretProofTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *SecretProofTransaction) Size() int {
	return SecretProofHeaderSize + tx.Proof.Size()
}

type CosignatureTransactionV1 struct {
	TransactionToCosign *AggregateTransactionV1
}

type CosignatureTransactionV2 struct {
	TransactionToCosign *AggregateTransactionV2
}

// returns a CosignatureTransaction from passed AggregateTransaction
func NewCosignatureTransactionV1(txToCosign *AggregateTransactionV1) (*CosignatureTransactionV1, error) {
	if txToCosign == nil {
		return nil, errors.New("txToCosign must not be nil")
	}
	return &CosignatureTransactionV1{txToCosign}, nil
}

// returns a CosignatureTransaction from passed hash of bounded AggregateTransaction
func NewCosignatureTransactionFromHashV1(hash *Hash) *CosignatureTransactionV1 {
	return &CosignatureTransactionV1{
		TransactionToCosign: &AggregateTransactionV1{
			AbstractTransaction: AbstractTransaction{
				TransactionInfo: TransactionInfo{
					TransactionHash: hash,
				},
			},
		},
	}
}

func (tx *CosignatureTransactionV1) String() string {
	return fmt.Sprintf(`"TransactionToCosign": %s`, tx.TransactionToCosign)
}

// returns a CosignatureTransaction from passed AggregateTransaction
func NewCosignatureTransaction(txToCosign *AggregateTransactionV2) (*CosignatureTransactionV2, error) {
	if txToCosign == nil {
		return nil, errors.New("txToCosign must not be nil")
	}
	return &CosignatureTransactionV2{txToCosign}, nil
}

// returns a CosignatureTransaction from passed hash of bounded AggregateTransaction
func NewCosignatureTransactionFromHash(hash *Hash) *CosignatureTransactionV2 {
	return &CosignatureTransactionV2{
		TransactionToCosign: &AggregateTransactionV2{
			AbstractTransaction: AbstractTransaction{
				TransactionInfo: TransactionInfo{
					TransactionHash: hash,
				},
			},
		},
	}
}

func (tx *CosignatureTransactionV2) String() string {
	return fmt.Sprintf(`"TransactionToCosign": %s`, tx.TransactionToCosign)
}

type signedTransactionDto struct {
	EntityType `json:"transactionType"`
	Payload    string `json:"payload"`
	Hash       string `json:"hash"`
}

type SignedTransaction struct {
	EntityType
	Payload string
	Hash    *Hash
}

func (tx *SignedTransaction) String() string {
	return fmt.Sprintf(
		`
			"EntityType": %d,
			"Payload": %s,
			"Hash": %s,
		`,
		tx.EntityType,
		tx.Payload,
		tx.Hash,
	)
}

type cosignatureSignedTransactionDto struct {
	ParentHash string `json:"parentHash"`
	Signature  string `json:"signature"`
	Scheme     string `json:"scheme"`
	Signer     string `json:"signer"`
}

type CosignatureSignedTransaction struct {
	ParentHash *Hash
	Signature  *Signature
	Scheme     crypto.DerivationScheme
	Signer     string
}

type AggregateTransactionCosignature struct {
	Signature string
	Signer    *PublicAccount
	Scheme    crypto.DerivationScheme
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

type TransactionStatus struct {
	Deadline *Deadline
	Group    TransactionGroup
	Status   string
	Hash     *Hash
	Height   Height
}

func (ts *TransactionStatus) String() string {
	return fmt.Sprintf(
		`
			"Group:" %s,
			"Status:" %s,
			"Hash": %s,
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

const (
	TransactionHeaderSize                        = SizeSize + SignerSize + SignatureSize + VersionSize + TypeSize + MaxFeeSize + DeadLineSize
	PropertyTypeSize                         int = 2
	PropertyModificationTypeSize             int = 1
	AccountPropertiesAddressModificationSize     = PropertyModificationTypeSize + AddressSize
	AccountPropertiesMosaicModificationSize      = PropertyModificationTypeSize + MosaicIdSize
	AccountPropertiesEntityModificationSize      = PropertyModificationTypeSize + TypeSize
	AccountPropertyAddressHeader                 = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyMosaicHeader                  = TransactionHeaderSize + PropertyTypeSize
	AccountPropertyEntityTypeHeader              = TransactionHeaderSize + PropertyTypeSize
	LinkActionSize                           int = 1
	AccountLinkTransactionSize                   = TransactionHeaderSize + KeySize + LinkActionSize
	NodeKeyLinkTransactionSize                   = TransactionHeaderSize + KeySize + LinkActionSize
	VrfKeyLinkTransactionSize                    = TransactionHeaderSize + KeySize + LinkActionSize
	AliasActionSize                          int = 1
	AliasTransactionHeaderSize                   = TransactionHeaderSize + NamespaceSize + AliasActionSize
	AggregateBondedHeaderSize                    = TransactionHeaderSize + SizeSize
	NetworkConfigHeaderSize                      = TransactionHeaderSize + BaseInt64Size + MaxStringSize + MaxStringSize
	NetworkConfigAbsoluteHeightHeaderSize        = NetworkConfigHeaderSize
	BlockchainUpgradeTransactionSize             = TransactionHeaderSize + DurationSize + BaseInt64Size
	AccountV2UpgradeTransactionSize              = TransactionHeaderSize + KeySize
	HashTypeSize                             int = 1
	LockSize                                     = TransactionHeaderSize + MosaicIdSize + AmountSize + DurationSize + Hash256
	MetadataTypeSize                         int = 1
	MetadataHeaderSize                           = TransactionHeaderSize + MetadataTypeSize
	MetadataV2HeaderSize                         = TransactionHeaderSize + KeySize + BaseInt64Size + 2 + 2
	ModificationsSizeSize                    int = 1
	ModifyContractHeaderSize                     = TransactionHeaderSize + DurationSize + Hash256 + 3*ModificationsSizeSize
	MinApprovalSize                          int = 1
	MinRemovalSize                           int = 1
	ModifyMultisigHeaderSize                     = TransactionHeaderSize + MinApprovalSize + MinRemovalSize + ModificationsSizeSize
	MosaicNonceSize                          int = 4
	MosaicPropertiesHeaderSize               int = 3
	MosaicPropertyIdSize                     int = 1
	MosaicOptionalPropertySize                   = MosaicPropertyIdSize + BaseInt64Size
	MosaicDefinitionTransactionHeaderSize        = TransactionHeaderSize + MosaicNonceSize + MosaicIdSize + MosaicPropertiesHeaderSize
	MosaicSupplyDirectionSize                int = 1
	MosaicSupplyChangeTransactionSize            = TransactionHeaderSize + MosaicIdSize + AmountSize + MosaicSupplyDirectionSize
	MosaicLevySize                               = 1 + AddressSize + MosaicIdSize + MaxFeeSize
	MosaicModifyLevyTransactionSize              = TransactionHeaderSize + MosaicLevySize + MosaicIdSize
	MosaicRemoveLevyTransactionSize              = TransactionHeaderSize + MosaicIdSize
	NamespaceTypeSize                        int = 1
	NamespaceNameSizeSize                    int = 1
	RegisterNamespaceHeaderSize                  = TransactionHeaderSize + NamespaceTypeSize + DurationSize + NamespaceSize + NamespaceNameSizeSize
	SecretLockSize                               = TransactionHeaderSize + MosaicIdSize + AmountSize + DurationSize + HashTypeSize + Hash256 + AddressSize
	ProofSizeSize                            int = 2
	SecretProofHeaderSize                        = TransactionHeaderSize + HashTypeSize + Hash256 + AddressSize + ProofSizeSize
	MosaicsSizeSize                          int = 1
	MessageSizeSize                          int = 2
	TransferHeaderSize                           = TransactionHeaderSize + AddressSize + MosaicsSizeSize + MessageSizeSize
	ReplicasSize                                 = 2
	MinReplicatorsSize                           = 2
	PercentApproversSize                         = 1
	PrepareDriveHeaderSize                       = TransactionHeaderSize + KeySize + DurationSize + DurationSize + AmountSize + StorageSizeSize + ReplicasSize + MinReplicatorsSize + PercentApproversSize
	JoinToDriveHeaderSize                        = TransactionHeaderSize + KeySize
	AddActionsSize                               = 2
	RemoveActionsSize                            = 2
	DriveFileSystemHeaderSize                    = TransactionHeaderSize + KeySize + Hash256 + Hash256 + AddActionsSize + RemoveActionsSize
	FilesSizeSize                                = 2
	FilesDepositHeaderSize                       = TransactionHeaderSize + KeySize + FilesSizeSize
	EndDriveHeaderSize                           = TransactionHeaderSize + KeySize
	StartDriveVerificationHeaderSize             = TransactionHeaderSize + KeySize
	OfferTypeSize                                = 1
	OffersCountSize                              = 1
	AddExchangeOfferSize                         = MosaicIdSize + DurationSize + 2*AmountSize + OfferTypeSize
	AddExchangeOfferHeaderSize                   = TransactionHeaderSize + OffersCountSize
	ExchangeOfferSize                            = DurationSize + 2*AmountSize + OfferTypeSize + KeySize
	ExchangeOfferHeaderSize                      = TransactionHeaderSize + OffersCountSize
	RemoveExchangeOfferSize                      = OfferTypeSize + MosaicIdSize
	RemoveExchangeOfferHeaderSize                = TransactionHeaderSize + OffersCountSize
	StartFileDownloadHeaderSize                  = TransactionHeaderSize + 2 + KeySize
	EndFileDownloadHeaderSize                    = TransactionHeaderSize + 2 + KeySize + Hash256
	OperationIdentifyHeaderSize                  = TransactionHeaderSize + Hash256
	EndOperationHeaderSize                       = TransactionHeaderSize + 1 + Hash256 + 2
	DeployHeaderSize                             = TransactionHeaderSize + KeySize + KeySize + Hash256 + BaseInt64Size
	StartExecuteHeaderSize                       = TransactionHeaderSize + KeySize + 1 + 1 + 2
	DeactivateHeaderSize                         = TransactionHeaderSize + KeySize + KeySize
	LockFundTransferHeaderSize                   = TransactionHeaderSize + DurationSize + 1 + 1
	LockFundCancelUnlockHeaderSize               = TransactionHeaderSize + BaseInt64Size
	AccountAddressRestrictionHeaderSize          = TransactionHeaderSize + CharCountSize + CharCountSize + HalfWordFlagsSize + IntPaddingSize
	AccountMosaicRestrictionHeaderSize           = TransactionHeaderSize + CharCountSize + CharCountSize + HalfWordFlagsSize + IntPaddingSize
	AccountOperationRestrictionHeaderSize        = TransactionHeaderSize + CharCountSize + CharCountSize + HalfWordFlagsSize + IntPaddingSize
	MosaicAddressRestrictionHeaderSize           = TransactionHeaderSize + MosaicIdSize + BaseInt64Size + BaseInt64Size + BaseInt64Size + AddressSize
	MosaicGlobalRestrictionHeaderSize            = TransactionHeaderSize + MosaicIdSize + MosaicIdSize + BaseInt64Size + BaseInt64Size + BaseInt64Size + ByteFlagsSize + ByteFlagsSize
)

type EntityType uint16

const (
	AccountPropertyAddress                EntityType = 0x4150
	AccountPropertyMosaic                 EntityType = 0x4250
	AccountPropertyEntityType             EntityType = 0x4350
	AddressAlias                          EntityType = 0x424e
	AggregateBondedV1                     EntityType = 0x4241
	AggregateCompletedV1                  EntityType = 0x4141
	AggregateBondedV2                     EntityType = 0x4441
	AggregateCompletedV2                  EntityType = 0x4341
	AddExchangeOffer                      EntityType = 0x415D
	AddHarvesterEntityType                EntityType = 0x4161
	ExchangeOffer                         EntityType = 0x425D
	RemoveExchangeOffer                   EntityType = 0x435D
	RemoveHarvesterEntityType             EntityType = 0x4261
	Block                                 EntityType = 0x8143
	NemesisBlock                          EntityType = 0x8043
	NetworkConfigEntityType               EntityType = 0x4159
	NetworkConfigAbsoluteHeightEntityType EntityType = 0x4259
	BlockchainUpgrade                     EntityType = 0x4158
	AccountV2Upgrade                      EntityType = 0x4258
	LinkAccount                           EntityType = 0x414c
	NodeKeyLink                           EntityType = 0x424c
	VrfKeyLink                            EntityType = 0x434c
	Lock                                  EntityType = 0x4148
	MetadataAddress                       EntityType = 0x413d
	MetadataMosaic                        EntityType = 0x423d
	MetadataNamespace                     EntityType = 0x433d
	AccountMetadata                       EntityType = 0x413f
	MosaicMetadata                        EntityType = 0x423f
	NamespaceMetadata                     EntityType = 0x433f
	ModifyContract                        EntityType = 0x4157
	ModifyMultisig                        EntityType = 0x4155
	MosaicAlias                           EntityType = 0x434e
	MosaicDefinition                      EntityType = 0x414d
	MosaicSupplyChange                    EntityType = 0x424d
	MosaicModifyLevy                      EntityType = 0x434d
	MosaicRemoveLevy                      EntityType = 0x444d
	RegisterNamespace                     EntityType = 0x414e
	SecretLock                            EntityType = 0x4152
	SecretProof                           EntityType = 0x4252
	Transfer                              EntityType = 0x4154
	PrepareDrive                          EntityType = 0x415A
	JoinToDrive                           EntityType = 0x425A
	DriveFileSystem                       EntityType = 0x435A
	FilesDeposit                          EntityType = 0x445A
	EndDrive                              EntityType = 0x455A
	DriveFilesReward                      EntityType = 0x465A
	StartDriveVerification                EntityType = 0x475A
	EndDriveVerification                  EntityType = 0x485A
	StartFileDownload                     EntityType = 0x495A
	EndFileDownload                       EntityType = 0x4A5A
	OperationIdentify                     EntityType = 0x415F
	StartOperation                        EntityType = 0x425F
	EndOperation                          EntityType = 0x435F
	Deploy                                EntityType = 0x4160
	StartExecute                          EntityType = 0x4260
	EndExecute                            EntityType = 0x4360
	SuperContractFileSystem               EntityType = 0x4460
	Deactivate                            EntityType = 0x4560
	LockFundTransfer                      EntityType = 0x4162
	LockFundCancelUnlock                  EntityType = 0x4262
	AccountAddressRestriction             EntityType = 0x4163
	AccountMosaicRestriction              EntityType = 0x4263
	AccountOperationRestriction           EntityType = 0x4363
	MosaicGlobalRestriction               EntityType = 0x4164
	MosaicAddressRestriction              EntityType = 0x4264
)

func (t EntityType) String() string {
	return fmt.Sprintf("0x%x", uint16(t))
}

type EntityVersion uint32

const (
	AccountPropertyAddressVersion      EntityVersion = 1
	AccountPropertyMosaicVersion       EntityVersion = 1
	AccountPropertyEntityTypeVersion   EntityVersion = 1
	AddressAliasVersion                EntityVersion = 1
	AggregateBondedV1Version           EntityVersion = 3
	AggregateCompletedV1Version        EntityVersion = 3
	AggregateBondedV2Version           EntityVersion = 1
	AggregateCompletedV2Version        EntityVersion = 1
	AddExchangeOfferVersion            EntityVersion = 4
	ExchangeOfferVersion               EntityVersion = 2
	RemoveExchangeOfferVersion         EntityVersion = 2
	NetworkConfigVersion               EntityVersion = 2
	NetworkConfigAbsoluteHeightVersion EntityVersion = 1
	BlockchainUpgradeVersion           EntityVersion = 1
	AccountV2UpgradeVersion            EntityVersion = 1
	LinkAccountVersion                 EntityVersion = 2
	NodeKeyLinkVersion                 EntityVersion = 1
	VrfKeyLinkVersion                  EntityVersion = 1
	LockVersion                        EntityVersion = 1
	AccountMetadataVersion             EntityVersion = 1
	MosaicMetadataVersion              EntityVersion = 1
	NamespaceMetadataVersion           EntityVersion = 1
	MetadataAddressVersion             EntityVersion = 1
	MetadataMosaicVersion              EntityVersion = 1
	MetadataNamespaceVersion           EntityVersion = 1
	ModifyContractVersion              EntityVersion = 3
	ModifyMultisigVersion              EntityVersion = 3
	MosaicAliasVersion                 EntityVersion = 1
	MosaicDefinitionVersion            EntityVersion = 4
	MosaicSupplyChangeVersion          EntityVersion = 3
	MosaicModifyLevyVersion            EntityVersion = 1
	MosaicRemoveLevyVersion            EntityVersion = 1
	RegisterNamespaceVersion           EntityVersion = 2
	SecretLockVersion                  EntityVersion = 1
	SecretProofVersion                 EntityVersion = 1
	TransferVersion                    EntityVersion = 4
	PrepareDriveVersion                EntityVersion = 3
	JoinToDriveVersion                 EntityVersion = 1
	DriveFileSystemVersion             EntityVersion = 1
	FilesDepositVersion                EntityVersion = 1
	EndDriveVersion                    EntityVersion = 1
	DriveFilesRewardVersion            EntityVersion = 1
	StartDriveVerificationVersion      EntityVersion = 1
	EndDriveVerificationVersion        EntityVersion = 1
	StartFileDownloadVersion           EntityVersion = 1
	EndFileDownloadVersion             EntityVersion = 1
	DeployVersion                      EntityVersion = 1
	StartExecuteVersion                EntityVersion = 1
	EndExecuteVersion                  EntityVersion = 1
	StartOperationVersion              EntityVersion = 1
	EndOperationVersion                EntityVersion = 1
	HarvesterVersion                   EntityVersion = 1
	OperationIdentifyVersion           EntityVersion = 1
	SuperContractFileSystemVersion     EntityVersion = 1
	DeactivateVersion                  EntityVersion = 1
	LockFundTransferVersion            EntityVersion = 1
	LockFundCancelUnlockVersion        EntityVersion = 1
	AccountAddressRestrictionVersion   EntityVersion = 1
	AccountMosaicRestrictionVersion    EntityVersion = 1
	AccountOperationRestrictionVersion EntityVersion = 1
	MosaicGlobalRestrictionVersion     EntityVersion = 1
	MosaicAddressRestrictionVersion    EntityVersion = 1
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

type MetadataV2Type uint8

func (t MetadataV2Type) String() string {
	return fmt.Sprintf("%d", t)
}

const (
	MetadataV2AddressType MetadataV2Type = iota
	MetadataV2MosaicType
	MetadataV2NamespaceType
)

func ExtractVersion(version int64) EntityVersion {
	return EntityVersion(uint32(version) & 0xFFFFFF)
}

func MapTransactions(b *bytes.Buffer, generationHash *Hash) ([]Transaction, error) {
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
			txs[i], errs[i] = MapTransaction(bytes.NewBuffer([]byte(t)), generationHash)
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

func dtoToTransaction(b *bytes.Buffer, dto transactionDto, generationHash *Hash) (Transaction, error) {
	if dto == nil {
		return nil, errors.New("dto can't be nil")
	}

	err := json.Unmarshal(b.Bytes(), dto)
	if err != nil {
		return nil, err
	}

	tx, err := dto.toStruct(generationHash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func MapTransaction(b *bytes.Buffer, generationHash *Hash) (Transaction, error) {
	rawT := struct {
		Transaction struct {
			Type EntityType
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
	case AggregateBondedV1, AggregateCompletedV1:
		dto = &aggregateTransactionV1DTO{}
	case AggregateBondedV2, AggregateCompletedV2:
		dto = &aggregateTransactionV2DTO{}
	case AddExchangeOffer:
		dto = &addExchangeOfferTransactionDTO{}
	case AddHarvesterEntityType:
		dto = &harvesterTransactionDTO{}
	case ExchangeOffer:
		dto = &exchangeOfferTransactionDTO{}
	case RemoveExchangeOffer:
		dto = &removeExchangeOfferTransactionDTO{}
	case NetworkConfigEntityType:
		dto = &networkConfigTransactionDTO{}
	case BlockchainUpgrade:
		dto = &blockchainUpgradeTransactionDTO{}
	case AccountV2Upgrade:
		dto = &accountV2UpgradeTransactionDTO{}
	case LinkAccount:
		dto = &accountLinkTransactionDTO{}
	case NodeKeyLink:
		dto = &nodeKeyLinkTransactionDTO{}
	case VrfKeyLink:
		dto = &vrfKeyLinkTransactionDTO{}
	case Lock:
		dto = &lockFundsTransactionDTO{}
	case AccountMetadata:
		dto = &accountMetadataTransactionDTO{}
	case MosaicMetadata:
		dto = &mosaicMetadataTransactionDTO{}
	case NamespaceMetadata:
		dto = &namespaceMetadataTransactionDTO{}
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
	case MosaicModifyLevy:
		dto = &mosaicModifyLevyTransactionDTO{}
	case MosaicRemoveLevy:
		dto = &mosaicRemoveLevyTransactionDTO{}
	case RegisterNamespace:
		dto = &registerNamespaceTransactionDTO{}
	case RemoveHarvesterEntityType:
		dto = &harvesterTransactionDTO{}
	case SecretLock:
		dto = &secretLockTransactionDTO{}
	case SecretProof:
		dto = &secretProofTransactionDTO{}
	case Transfer:
		dto = &transferTransactionDTO{}
	case PrepareDrive:
		dto = &prepareDriveTransactionDTO{}
	case JoinToDrive:
		dto = &joinToDriveTransactionDTO{}
	case DriveFileSystem:
		dto = &driveFileSystemTransactionDTO{}
	case FilesDeposit:
		dto = &filesDepositTransactionDTO{}
	case EndDrive:
		dto = &endDriveTransactionDTO{}
	case DriveFilesReward:
		dto = &driveFilesRewardTransactionDTO{}
	case StartDriveVerification:
		dto = &startDriveVerificationTransactionDTO{}
	case EndDriveVerification:
		dto = &endDriveVerificationTransactionDTO{}
	case StartFileDownload:
		dto = &startFileDownloadTransactionDTO{}
	case EndFileDownload:
		dto = &endFileDownloadTransactionDTO{}
	case OperationIdentify:
		dto = &operationIdentifyTransactionDTO{}
	case EndOperation:
		dto = &endOperationTransactionDTO{}
	case Deploy:
		dto = &deployTransactionDTO{}
	case StartExecute:
		dto = &startExecuteTransactionDTO{}
	case EndExecute:
		dto = &endOperationTransactionDTO{}
	case SuperContractFileSystem:
		dto = &driveFileSystemTransactionDTO{}
	case Deactivate:
		dto = &deactivateTransactionDTO{}
	case LockFundTransfer:
		dto = &lockFundTransferTransactionDto{}
	case LockFundCancelUnlock:
		dto = &lockFundCancelUnlockTransactionDto{}
	case AccountAddressRestriction:
		dto = &AccountAddressRestrictionTransactionDto{}
	case AccountMosaicRestriction:
		dto = &lockFundCancelUnlockTransactionDto{}
	case AccountOperationRestriction:
		dto = &lockFundTransferTransactionDto{}
	case MosaicGlobalRestriction:
		dto = &lockFundCancelUnlockTransactionDto{}
	case MosaicAddressRestriction:
		dto = &lockFundTransferTransactionDto{}
	}

	return dtoToTransaction(b, dto, generationHash)
}

func createTransactionHash(b []byte, generationHash *Hash) (*Hash, error) {
	var sizeOfGenerationHash = 0
	if generationHash != nil {
		sizeOfGenerationHash = len(generationHash)
	}

	sb := make([]byte, len(b)-SizeSize-HalfOfSignature+sizeOfGenerationHash)
	copy(sb[:HalfOfSignature], b[SizeSize:SizeSize+HalfOfSignature])
	copy(sb[HalfOfSignature:HalfOfSignature+SignerSize], b[SizeSize+SignatureSize:SizeSize+SignatureSize+SignerSize])

	if generationHash != nil {
		copy(sb[HalfOfSignature+SignerSize:], generationHash[:])
	}

	copy(sb[HalfOfSignature+SignerSize+sizeOfGenerationHash:], b[SizeSize+SignatureSize+SignerSize:])

	r, err := crypto.HashesSha3_256(sb)
	if err != nil {
		return nil, err
	}

	return bytesToHash(r)
}

func toAggregateTransactionBytes(tx Transaction) ([]byte, error) {
	if tx.GetAbstractTransaction().Signer == nil {
		return nil, fmt.Errorf("some of the transaction does not have a signer")
	}
	sb, err := hex.DecodeString(tx.GetAbstractTransaction().Signer.PublicKey)
	if err != nil {
		return nil, err
	}
	b, err := tx.Bytes()
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
	derivationScheme := GetDerivationSchemeForAccountVersion(a.Version).EngineDerivationScheme()
	// Embed signature derivation scheme to version field
	tx.GetAbstractTransaction().Version = EntityVersion(uint32(tx.GetAbstractTransaction().Version) | (uint32(derivationScheme) << 16))
	b, err := tx.Bytes()
	if err != nil {
		return nil, err
	}
	sb := make([]byte, len(b)-SizeSize-SignerSize-SignatureSize)
	copy(sb, b[SizeSize+SignerSize+SignatureSize:])

	if a.generationHash != nil {
		sb = append(a.generationHash[:], sb...)
	}
	signature, err := s.Sign(sb)
	if err != nil {
		return nil, err
	}
	p := make([]byte, len(b))

	copy(p[:SizeSize], b[:SizeSize])
	copy(p[SizeSize:SizeSize+SignatureSize], signature.Bytes())
	copy(p[SizeSize+SignatureSize:SizeSize+SignatureSize+SignerSize], a.KeyPair.PublicKey.Raw)
	copy(p[SizeSize+SignatureSize+SignerSize:], b[SizeSize+SignatureSize+SignerSize:])

	h, err := createTransactionHash(p, a.generationHash)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{tx.GetAbstractTransaction().Type, strings.ToUpper(hex.EncodeToString(p)), h}, nil
}

func InnerTransactionHash(tx Transaction) *Hash {
	b, err := toAggregateTransactionBytes(tx)
	if err != nil {
		panic(err)
	}
	sb := make([]byte, len(b)-SizeSize)
	copy(sb, b[SizeSize:SizeSize+SignerSize])
	copy(sb[SignerSize:], b[SizeSize+SignerSize:SizeSize+SignerSize+VersionSize+TypeSize])

	copy(
		sb[SignerSize+VersionSize+TypeSize:],
		b[SizeSize+SignerSize+VersionSize+TypeSize:],
	)

	r, err := crypto.HashesSha3_256(sb)
	if err != nil {
		panic(err)
	}

	result, err := bytesToHash(r)
	if err != nil {
		panic(err)
	}

	return result
}

func UniqueAggregateHashImpl(deadline *Deadline, tx Transaction, generationHash *Hash) (*Hash, error) {
	b, err := toAggregateTransactionBytes(tx)
	if err != nil {
		return nil, err
	}
	generationSize := len(generationHash)
	sb := make([]byte, len(b)-SizeSize+DeadLineSize+generationSize)
	copy(sb, b[SizeSize:SizeSize+SignerSize])
	copy(sb[SignerSize:], generationHash[:])
	copy(sb[SignerSize+generationSize:], b[SizeSize+SignerSize:SizeSize+SignerSize+VersionSize+TypeSize])

	// We are using dealine of aggregate transaction instead of deadline of transaction
	deadlineB := deadline.ToBlockchainTimestamp().toLittleEndian()
	copy(sb[SignerSize+generationSize+VersionSize+TypeSize:], deadlineB)
	copy(
		sb[SignerSize+generationSize+VersionSize+TypeSize+DeadLineSize:],
		b[SizeSize+SignerSize+VersionSize+TypeSize:],
	)

	r, err := crypto.HashesSha3_256(sb)
	if err != nil {
		return nil, err
	}

	return bytesToHash(r)
}
func UniqueAggregateHashV1(aggregateTx *AggregateTransactionV1, tx Transaction, generationHash *Hash) (*Hash, error) {
	return UniqueAggregateHashImpl(aggregateTx.Deadline, tx, generationHash)
}

func UniqueAggregateHashV2(aggregateTx *AggregateTransactionV2, tx Transaction, generationHash *Hash) (*Hash, error) {
	return UniqueAggregateHashImpl(aggregateTx.Deadline, tx, generationHash)
}

func signTransactionWithCosignaturesImpl(stx *SignedTransaction, txType EntityType, cosignatories []*Account, extended bool) (*SignedTransaction, error) {
	p := stx.Payload
	for _, cos := range cosignatories {
		s := crypto.NewSignerFromKeyPair(cos.KeyPair, cos.KeyPair.CryptoEngine)
		sb, err := s.Sign(stx.Hash[:])
		if err != nil {
			return nil, err
		}
		p += cos.KeyPair.PublicKey.String()
		if extended {
			sbe, err := crypto.NewExtendedSignatureFromSignature(sb, cos.KeyPair.EngineDerivationScheme())
			if err != nil {
				return nil, err
			}
			p += hex.EncodeToString(sbe.Bytes())
		} else {
			p += hex.EncodeToString(sb.Bytes())
		}
	}

	pb, err := hex.DecodeString(p)
	if err != nil {
		return nil, err
	}

	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, uint32(len(pb)))

	copy(pb[:len(s)], s)

	return &SignedTransaction{txType, hex.EncodeToString(pb), stx.Hash}, nil
}

func signTransactionWithCosignaturesV1(tx *AggregateTransactionV1, a *Account, cosignatories []*Account) (*SignedTransaction, error) {
	stx, err := signTransactionWith(tx, a)
	if err != nil {
		return nil, err
	}
	return signTransactionWithCosignaturesImpl(stx, tx.Type, cosignatories, false)
}
func signTransactionWithCosignaturesV2(tx *AggregateTransactionV2, a *Account, cosignatories []*Account) (*SignedTransaction, error) {
	stx, err := signTransactionWith(tx, a)
	if err != nil {
		return nil, err
	}
	return signTransactionWithCosignaturesImpl(stx, tx.Type, cosignatories, true)
}

func signCosignatureTransactionImpl(a *Account, tx *TransactionInfo) (*CosignatureSignedTransaction, error) {
	if tx.TransactionHash.Empty() {
		return nil, errors.New("cosignature transaction hash is nil")
	}

	s := crypto.NewSignerFromKeyPair(a.KeyPair, a.CryptoEngine)
	b := tx.TransactionHash[:]

	sb, err := s.Sign(b)
	if err != nil {
		return nil, err
	}

	signature, err := bytesToSignature(sb.Bytes())
	if err != nil {
		return nil, err
	}

	return &CosignatureSignedTransaction{tx.TransactionHash, signature, a.EngineDerivationScheme(), a.PublicAccount.PublicKey}, nil
}

func signCosignatureTransactionV1(a *Account, tx *CosignatureTransactionV1) (*CosignatureSignedTransaction, error) {
	return signCosignatureTransactionImpl(a, &tx.TransactionToCosign.TransactionInfo)
}

func signCosignatureTransactionV2(a *Account, tx *CosignatureTransactionV2) (*CosignatureSignedTransaction, error) {
	return signCosignatureTransactionImpl(a, &tx.TransactionToCosign.TransactionInfo)
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

func mosaicPropertyArrayToBuffer(builder *flatbuffers.Builder, properties []MosaicProperty) flatbuffers.UOffsetT {
	pBuffer := make([]flatbuffers.UOffsetT, len(properties))
	for i, p := range properties {
		valueV := transactions.TransactionBufferCreateUint32Vector(builder, p.Value.toArray())

		transactions.MosaicPropertyStart(builder)
		transactions.MosaicPropertyAddMosaicPropertyId(builder, byte(p.Id))
		transactions.MosaicPropertyAddValue(builder, valueV)

		pBuffer[i] = transactions.TransactionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, pBuffer)
}

func hashToBuffer(builder *flatbuffers.Builder, hash *Hash) flatbuffers.UOffsetT {
	pV := transactions.TransactionBufferCreateByteVector(builder, hash[:])

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
