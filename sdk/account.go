// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type AccountService service

// returns account info for passed address
func (a *AccountService) GetAccountInfo(ctx context.Context, address *Address) (*AccountInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	if len(address.Address) == 0 {
		return nil, ErrBlankAddress
	}

	url := net.NewUrl(fmt.Sprintf(accountRoute, address.Address))

	dto := &accountInfoDTO{}

	resp, err := a.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.reputationConfig)
}

// returns an array of account infos for passed addresses
func (a *AccountService) GetAccountsInfo(ctx context.Context, addresses []*Address) ([]*AccountInfo, error) {
	if len(addresses) == 0 {
		return nil, ErrEmptyAddressesIds
	}

	addrs := struct {
		Messages []string `json:"addresses"`
	}{
		Messages: make([]string, len(addresses)),
	}

	for i, address := range addresses {
		addrs.Messages[i] = address.Address
	}

	dtos := accountInfoDTOs(make([]*accountInfoDTO, 0))

	resp, err := a.client.DoNewRequest(ctx, http.MethodPost, accountsRoute, addrs, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct(a.client.config.reputationConfig)
}

// returns multisig account info for passed address
func (a *AccountService) GetMultisigAccountInfo(ctx context.Context, address *Address) (*MultisigAccountInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(multisigAccountRoute, address.Address))

	dto := &multisigAccountInfoDTO{}

	resp, err := a.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.NetworkType)
}

// returns multisig account info for passed address
func (a *AccountService) GetMultisigAccountGraphInfo(ctx context.Context, address *Address) (*MultisigAccountGraphInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(multisigAccountGraphInfoRoute, address.Address))

	dto := &multisigAccountGraphInfoDTOS{}

	resp, err := a.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.NetworkType)
}

// returns an array of confirmed transactions for which passed account is sender or receiver.
func (a *AccountService) Transactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error) {
	return a.findTransactions(ctx, account, opt, accountTransactionsRoute)
}

// returns an array of transactions for which passed account is receiver
func (a *AccountService) IncomingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error) {
	return a.findTransactions(ctx, account, opt, incomingTransactionsRoute)
}

// returns an array of transaction for which passed account is sender
func (a *AccountService) OutgoingTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error) {
	return a.findTransactions(ctx, account, opt, outgoingTransactionsRoute)
}


// returns an array of confirmed transactions for which passed account is sender or receiver.
// unconfirmed transactions are those transactions that have not yet been included in a block.
// unconfirmed transactions are not guaranteed to be included in any block.
func (a *AccountService) UnconfirmedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]Transaction, error) {
	return a.findTransactions(ctx, account, opt, unconfirmedTransactionsRoute)
}

// returns an array of aggregate bounded transactions where passed account is signer or cosigner
func (a *AccountService) AggregateBondedTransactions(ctx context.Context, account *PublicAccount, opt *AccountTransactionsOption) ([]*AggregateTransaction, error) {
	txs, err := a.findTransactions(ctx, account, opt, aggregateTransactionsRoute)
	if err != nil {
		return nil, err
	}

	atxs := make([]*AggregateTransaction, len(txs))
	for i, tx := range txs {
		atxs[i] = tx.(*AggregateTransaction)
	}

	return atxs, nil
}
