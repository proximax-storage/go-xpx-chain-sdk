// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"errors"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewPlaceSdaExchangeOfferTransaction(deadline *Deadline, placeSdaOffers []*PlaceSdaOffer, networkType NetworkType) (*PlaceSdaExchangeOfferTransaction, error) {
	if len(placeSdaOffers) == 0 {
		return nil, errors.New("PlaceSdaOffers should be not empty")
	}

	tx := PlaceSdaExchangeOfferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     PlaceSdaExchangeOfferVersion,
			Deadline:    deadline,
			Type:        PlaceSdaExchangeOffer,
			NetworkType: networkType,
		},
		Offers: placeSdaOffers,
	}

	return &tx, nil
}

func (tx *PlaceSdaExchangeOfferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *PlaceSdaExchangeOfferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"PlaceSdaOffers": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Offers,
	)
}

func (tx *PlaceSdaExchangeOfferTransaction) Size() int {
	return PlaceSdaExchangeOfferHeaderSize + len(tx.Offers)*PlaceSdaExchangeOfferSize
}

func placeSdaExchangeOfferToArrayToBuffer(builder *flatbuffers.Builder, offers []*PlaceSdaOffer) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(offers))
	for i, offer := range offers {
		ob, err := hex.DecodeString(offer.Owner.PublicKey)
		if err != nil {
			return 0, err
		}

		mGiveV := transactions.TransactionBufferCreateUint32Vector(builder, offer.MosaicGive.AssetId.toArray())
		maGiveV := transactions.TransactionBufferCreateUint32Vector(builder, offer.MosaicGive.Amount.toArray())
		mGetV := transactions.TransactionBufferCreateUint32Vector(builder, offer.MosaicGet.AssetId.toArray())
		maGetV := transactions.TransactionBufferCreateUint32Vector(builder, offer.MosaicGet.Amount.toArray())
		oV := transactions.TransactionBufferCreateByteVector(builder, ob)
		dV := transactions.TransactionBufferCreateUint32Vector(builder, offer.Duration.toArray())

		transactions.PlaceSdaExchangeOfferBufferStart(builder)
		transactions.PlaceSdaExchangeOfferBufferAddMosaicIdGive(builder, mGiveV)
		transactions.PlaceSdaExchangeOfferBufferAddMosaicAmountGive(builder, maGiveV)
		transactions.PlaceSdaExchangeOfferBufferAddMosaicIdGet(builder, mGetV)
		transactions.PlaceSdaExchangeOfferBufferAddMosaicAmountGet(builder, maGetV)
		transactions.PlaceSdaExchangeOfferBufferAddOwner(builder, oV)
		transactions.PlaceSdaExchangeOfferBufferAddDuration(builder, dV)
		msb[i] = transactions.PlaceSdaExchangeOfferBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *PlaceSdaExchangeOfferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	offersV, err := placeSdaExchangeOfferToArrayToBuffer(builder, tx.Offers)
	if err != nil {
		return nil, err
	}

	transactions.PlaceSdaExchangeOfferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.PlaceSdaExchangeOfferTransactionBufferAddSdaOfferCount(builder, byte(len(tx.Offers)))
	transactions.PlaceSdaExchangeOfferTransactionBufferAddOffers(builder, offersV)
	t := transactions.PlaceSdaExchangeOfferTransactionBufferEnd(builder)
	builder.Finish(t)

	return placeSdaExchangeOfferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type sdaOfferDTO struct {
	AssetIdGive assetIdDTO `json:"mosaicIdGive"`
	AmountGive  uint64DTO  `json:"mosaicAmountGive"`
	AssetIdGet  assetIdDTO `json:"mosaicIdGet"`
	AmountGet   uint64DTO  `json:"mosaicAmountGet"`
}

func (dto *sdaOfferDTO) toStruct() (*SdaOffer, error) {
	hGive, err := dto.AssetIdGive.toStruct()
	if err != nil {
		return nil, err
	}

	hGet, err := dto.AssetIdGet.toStruct()
	if err != nil {
		return nil, err
	}

	return &SdaOffer{
		MosaicGive: newMosaicPanic(hGive, dto.AmountGive.toStruct()),
		MosaicGet:  newMosaicPanic(hGet, dto.AmountGet.toStruct()),
	}, nil
}

type placeSdaOfferDTO struct {
	sdaOfferDTO
	Owner    string    `json:"owner"`
	Duration uint64DTO `json:"duration"`
}

func placeSdaOfferDTOArrayToStruct(offers []*placeSdaOfferDTO, networkType NetworkType) ([]*PlaceSdaOffer, error) {
	offersResult := make([]*PlaceSdaOffer, len(offers))
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

		offersResult[i] = &PlaceSdaOffer{
			SdaOffer: *o,
			Owner:    a,
			Duration: offer.Duration.toStruct(),
		}
	}

	return offersResult, err
}

type placeSdaOfferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		SdaOffers []*placeSdaOfferDTO `json:"offers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *placeSdaOfferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	offers, err := placeSdaOfferDTOArrayToStruct(dto.Tx.SdaOffers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &PlaceSdaExchangeOfferTransaction{
		*atx,
		offers,
	}, nil
}

func NewRemoveSdaExchangeOfferTransaction(deadline *Deadline, removeSdaOffers []*RemoveSdaOffer, networkType NetworkType) (*RemoveSdaExchangeOfferTransaction, error) {
	if len(removeSdaOffers) == 0 {
		return nil, errors.New("RemoveSdaOffers should be not empty")
	}

	tx := RemoveSdaExchangeOfferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     RemoveSdaExchangeOfferVersion,
			Deadline:    deadline,
			Type:        RemoveSdaExchangeOffer,
			NetworkType: networkType,
		},
		Offers: removeSdaOffers,
	}

	return &tx, nil
}

func (tx *RemoveSdaExchangeOfferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *RemoveSdaExchangeOfferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RemoveSdaOffers": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Offers,
	)
}

func (tx *RemoveSdaExchangeOfferTransaction) Size() int {
	return RemoveSdaExchangeOfferHeaderSize + len(tx.Offers)*RemoveSdaExchangeOfferSize
}

func removeSdaExchangeOfferToArrayToBuffer(builder *flatbuffers.Builder, offers []*RemoveSdaOffer) flatbuffers.UOffsetT {
	msb := make([]flatbuffers.UOffsetT, len(offers))
	for i, offer := range offers {
		mGiveV := transactions.TransactionBufferCreateUint32Vector(builder, offer.AssetIdGive.toArray())
		mGetV := transactions.TransactionBufferCreateUint32Vector(builder, offer.AssetIdGet.toArray())

		transactions.RemoveSdaExchangeOfferBufferStart(builder)
		transactions.RemoveSdaExchangeOfferBufferAddMosaicIdGive(builder, mGiveV)
		transactions.RemoveSdaExchangeOfferBufferAddMosaicIdGet(builder, mGetV)
		msb[i] = transactions.RemoveSdaExchangeOfferBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb)
}

func (tx *RemoveSdaExchangeOfferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	offersV := removeSdaExchangeOfferToArrayToBuffer(builder, tx.Offers)

	transactions.RemoveSdaExchangeOfferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.RemoveSdaExchangeOfferTransactionBufferAddSdaOfferCount(builder, byte(len(tx.Offers)))
	transactions.RemoveSdaExchangeOfferTransactionBufferAddOffers(builder, offersV)
	t := transactions.RemoveSdaExchangeOfferTransactionBufferEnd(builder)
	builder.Finish(t)

	return removeSdaExchangeOfferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type removeSdaOfferDTO struct {
	AssetIdGive assetIdDTO `json:"mosaicIdGive"`
	AssetIdGet  assetIdDTO `json:"mosaicIdGet"`
}

func removeSdaOfferDTOArrayToStruct(offers []*removeSdaOfferDTO) ([]*RemoveSdaOffer, error) {
	offersResult := make([]*RemoveSdaOffer, len(offers))
	var err error = nil
	for i, offer := range offers {
		hGive, err := offer.AssetIdGive.toStruct()
		if err != nil {
			return nil, err
		}

		hGet, err := offer.AssetIdGet.toStruct()
		if err != nil {
			return nil, err
		}

		offersResult[i] = &RemoveSdaOffer{
			AssetIdGive: hGive,
			AssetIdGet:  hGet,
		}
	}

	return offersResult, err
}

type removeSdaExchangeOfferTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Offers []*removeSdaOfferDTO `json:"offers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *removeSdaExchangeOfferTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	offers, err := removeSdaOfferDTOArrayToStruct(dto.Tx.Offers)
	if err != nil {
		return nil, err
	}

	return &RemoveSdaExchangeOfferTransaction{
		*atx,
		offers,
	}, nil
}
