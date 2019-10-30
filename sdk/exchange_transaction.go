// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"errors"
	"fmt"
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
		Offers:	addOffers,
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
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *AddExchangeOfferTransaction) Size() int {
	return 0
}

func (tx *AddExchangeOfferTransaction) generateBytes() ([]byte, error) {
	return nil, nil
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
		Confirmations:	confirmations,
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
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *ExchangeOfferTransaction) Size() int {
	return 0
}

func (tx *ExchangeOfferTransaction) generateBytes() ([]byte, error) {
	return nil, nil
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
		Offers:	removeOffers,
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
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *RemoveExchangeOfferTransaction) Size() int {
	return 0
}

func (tx *RemoveExchangeOfferTransaction) generateBytes() ([]byte, error) {
	return nil, nil
}