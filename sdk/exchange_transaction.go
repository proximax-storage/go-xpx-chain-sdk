// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewAddExchangeOfferTransaction(deadline *Deadline, addOffers []*AddOffer, networkType NetworkType) (*AddExchangeOfferTransaction, error) {
	if len(addOffers) == 0 {
		return nil, errors.New("AddOffers should be not empty")
	}

	tx := AddExchangeOfferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AddExchangeOfferVersion,
			Deadline:    deadline,
			Type:        AddExchangeOffer,
			NetworkType: networkType,
		},
		Offers: addOffers,
	}

	return &tx, nil
}

func (tx *AddExchangeOfferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AddExchangeOfferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"AddOffers": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Offers,
	)
}

func (tx *AddExchangeOfferTransaction) Size() int {
	return AddExchangeOfferHeaderSize + len(tx.Offers)*AddExchangeOfferSize
}

func addExchangeOfferToArrayToBuffer(builder *flatbuffers.Builder, offers []*AddOffer) flatbuffers.UOffsetT {
	msb := make([]flatbuffers.UOffsetT, len(offers))
	for i, offer := range offers {

		mV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Mosaic.AssetId.toArray())
		maV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Mosaic.Amount.toArray())
		dV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Duration.toArray())
		cV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Cost.toArray())

		transactions.AddExchangeOfferBufferStart(builder)
		transactions.AddExchangeOfferBufferAddMosaicId(builder, mV)
		transactions.AddExchangeOfferBufferAddMosaicAmount(builder, maV)
		transactions.AddExchangeOfferBufferAddCost(builder, cV)
		transactions.AddExchangeOfferBufferAddDuration(builder, dV)
		transactions.AddExchangeOfferBufferAddType(builder, byte(offer.Type))
		msb[i] = transactions.AddExchangeOfferBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb)
}

func (tx *AddExchangeOfferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	offersV := addExchangeOfferToArrayToBuffer(builder, tx.Offers)

	transactions.AddExchangeOfferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AddExchangeOfferTransactionBufferAddOffersCount(builder, byte(len(tx.Offers)))
	transactions.AddExchangeOfferTransactionBufferAddOffers(builder, offersV)
	t := transactions.AddExchangeOfferTransactionBufferEnd(builder)
	builder.Finish(t)

	return addExchangeOfferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type offerDTO struct {
	AssetId assetIdDTO `json:"mosaicId"`
	Amount  uint64DTO  `json:"mosaicAmount"`
	Cost    uint64DTO  `json:"cost"`
	Type    OfferType  `json:"type"`
}

func (dto *offerDTO) toStruct() (*Offer, error) {
	h, err := dto.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	return &Offer{
		Type:   dto.Type,
		Mosaic: newMosaicPanic(h, dto.Amount.toStruct()),
		Cost:   dto.Cost.toStruct(),
	}, nil
}

type addOfferDTO struct {
	offerDTO
	Duration uint64DTO `json:"duration"`
}

func addOfferDTOArrayToStruct(offers []*addOfferDTO) ([]*AddOffer, error) {
	offersResult := make([]*AddOffer, len(offers))
	var err error = nil
	for i, offer := range offers {
		o, err := offer.toStruct()
		if err != nil {
			return nil, err
		}

		offersResult[i] = &AddOffer{
			Offer:    *o,
			Duration: offer.Duration.toStruct(),
		}

	}

	return offersResult, err
}

type addExchangeOfferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Offers []*addOfferDTO `json:"offers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *addExchangeOfferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	offers, err := addOfferDTOArrayToStruct(dto.Tx.Offers)
	if err != nil {
		return nil, err
	}

	return &AddExchangeOfferTransaction{
		*atx,
		offers,
	}, nil
}

func NewExchangeOfferTransaction(deadline *Deadline, confirmations []*ExchangeConfirmation, networkType NetworkType) (*ExchangeOfferTransaction, error) {
	if len(confirmations) == 0 {
		return nil, errors.New("Confirmations should be not empty")
	}

	tx := ExchangeOfferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ExchangeOfferVersion,
			Deadline:    deadline,
			Type:        ExchangeOffer,
			NetworkType: networkType,
		},
		Confirmations: confirmations,
	}

	return &tx, nil
}

func (tx *ExchangeOfferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ExchangeOfferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Confirmations": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Confirmations,
	)
}

func (tx *ExchangeOfferTransaction) Size() int {
	return ExchangeOfferHeaderSize + len(tx.Confirmations)*ExchangeOfferSize
}

func exchangeOfferToArrayToBuffer(builder *flatbuffers.Builder, offers []*ExchangeConfirmation) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(offers))
	for i, offer := range offers {

		mV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Mosaic.AssetId.toArray())
		maV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Mosaic.Amount.toArray())
		cV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Cost.toArray())

		ob, err := hex.DecodeString(offer.Owner.PublicKey)
		if err != nil {
			return 0, err
		}

		oV := transactions.TransactionBufferCreateByteVector(builder, ob)

		transactions.ExchangeOfferBufferStart(builder)
		transactions.ExchangeOfferBufferAddMosaicId(builder, mV)
		transactions.ExchangeOfferBufferAddMosaicAmount(builder, maV)
		transactions.ExchangeOfferBufferAddCost(builder, cV)
		transactions.ExchangeOfferBufferAddType(builder, byte(offer.Type))
		transactions.ExchangeOfferBufferAddOwner(builder, oV)
		msb[i] = transactions.ExchangeOfferBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *ExchangeOfferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	offersV, err := exchangeOfferToArrayToBuffer(builder, tx.Confirmations)
	if err != nil {
		return nil, err
	}

	transactions.ExchangeOfferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.ExchangeOfferTransactionBufferAddOffersCount(builder, byte(len(tx.Confirmations)))
	transactions.ExchangeOfferTransactionBufferAddOffers(builder, offersV)
	t := transactions.ExchangeOfferTransactionBufferEnd(builder)
	builder.Finish(t)

	return exchangeOfferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type confirmationOfferDTO struct {
	offerDTO
	Owner string   `json:"owner"`
}

func confirmationOfferDTOArrayToStruct(offers []*confirmationOfferDTO, networkType NetworkType) ([]*ExchangeConfirmation, error) {
	offersResult := make([]*ExchangeConfirmation, len(offers))
	var err error = nil
	for i, offer := range offers {
		o, err := offer.toStruct()
		if err != nil {
			return nil, err
		}

		a, err := NewAccountFromPublicKey(offer.Owner, networkType)
		if err != nil {
			return nil, err
		}

		offersResult[i] = &ExchangeConfirmation{
			Offer: *o,
			Owner: a,
		}

	}

	return offersResult, err
}

type exchangeOfferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Offers []*confirmationOfferDTO `json:"offers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *exchangeOfferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	offers, err := confirmationOfferDTOArrayToStruct(dto.Tx.Offers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ExchangeOfferTransaction{
		*atx,
		offers,
	}, nil
}

func NewRemoveExchangeOfferTransaction(deadline *Deadline, removeOffers []*RemoveOffer, networkType NetworkType) (*RemoveExchangeOfferTransaction, error) {
	if len(removeOffers) == 0 {
		return nil, errors.New("RemoveOffers should be not empty")
	}

	tx := RemoveExchangeOfferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     RemoveExchangeOfferVersion,
			Deadline:    deadline,
			Type:        RemoveExchangeOffer,
			NetworkType: networkType,
		},
		Offers: removeOffers,
	}

	return &tx, nil
}

func (tx *RemoveExchangeOfferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *RemoveExchangeOfferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RemoveOffers": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Offers,
	)
}

func (tx *RemoveExchangeOfferTransaction) Size() int {
	return RemoveExchangeOfferHeaderSize + len(tx.Offers)*RemoveExchangeOfferSize
}

func removeExchangeOfferToArrayToBuffer(builder *flatbuffers.Builder, offers []*RemoveOffer) flatbuffers.UOffsetT {
	msb := make([]flatbuffers.UOffsetT, len(offers))
	for i, offer := range offers {

		mV := transactions.TransactionBufferCreateUint32Vector(builder, offer.AssetId.toArray())

		transactions.RemoveExchangeOfferBufferStart(builder)
		transactions.RemoveExchangeOfferBufferAddMosaicId(builder, mV)
		transactions.RemoveExchangeOfferBufferAddType(builder, byte(offer.Type))
		msb[i] = transactions.RemoveExchangeOfferBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb)
}

func (tx *RemoveExchangeOfferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	offersV := removeExchangeOfferToArrayToBuffer(builder, tx.Offers)

	transactions.RemoveExchangeOfferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.RemoveExchangeOfferTransactionBufferAddOffersCount(builder, byte(len(tx.Offers)))
	transactions.RemoveExchangeOfferTransactionBufferAddOffers(builder, offersV)
	t := transactions.RemoveExchangeOfferTransactionBufferEnd(builder)
	builder.Finish(t)

	return removeExchangeOfferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type removeOfferDTO struct {
	AssetId assetIdDTO `json:"mosaicId"`
	Type    OfferType  `json:"offerType"`
}

func removeOfferDTOArrayToStruct(offers []*removeOfferDTO) ([]*RemoveOffer, error) {
	offersResult := make([]*RemoveOffer, len(offers))
	var err error = nil
	for i, offer := range offers {
		h, err := offer.AssetId.toStruct()
		if err != nil {
			return nil, err
		}

		offersResult[i] = &RemoveOffer{
			AssetId: h,
			Type:    offer.Type,
		}

	}

	return offersResult, err
}

type removeExchangeOfferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Offers []*removeOfferDTO `json:"offers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *removeExchangeOfferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	offers, err := removeOfferDTOArrayToStruct(dto.Tx.Offers)
	if err != nil {
		return nil, err
	}

	return &RemoveExchangeOfferTransaction{
		*atx,
		offers,
	}, nil
}
