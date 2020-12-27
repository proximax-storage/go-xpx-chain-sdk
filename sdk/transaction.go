// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type TransactionService struct {
	*service
	BlockchainService *BlockchainService
}

// returns Transaction for passed transaction group and hash
func (txs *TransactionService) getTransaction(ctx context.Context, group string, id string) (Transaction, error) {
	var b bytes.Buffer

	resp, err := txs.client.doNewRequest(ctx, http.MethodGet, fmt.Sprintf(transactionsByIdRoute, group, id), nil, &b)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return MapTransaction(&b, txs.client.GenerationHash())
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) getTransactionsByGroup(ctx context.Context, group TransactionGroup, tpOpts *TransactionsPageOptions) (*TransactionsPage, error) {
	tspDTO := &transactionsPageDTO{}

	u, err := addOptions(fmt.Sprintf(transactionsByGroupRoute, group.String()), tpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := txs.client.doNewRequest(ctx, http.MethodGet, u, nil, &tspDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return tspDTO.toStruct(txs.client.GenerationHash())
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) getTransactionsByIds(ctx context.Context, group TransactionGroup	, ids []string, tpOpts *TransactionsPageOptions) ([]Transaction, error) {
	var b bytes.Buffer
	txIds := &TransactionIdsDTO{
		ids,
	}

	u, err := addOptions(fmt.Sprintf(transactionsByGroupRoute, group.String()), tpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := txs.client.doNewRequest(ctx, http.MethodPost, u, txIds, &b)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return MapTransactions(&b, txs.client.GenerationHash())
}

// returns Transaction for passed transaction id or hash
func (txs *TransactionService) GetAnyTransactionById(ctx context.Context, id string) (Transaction, error) {
	trS, err := txs.GetTransactionStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	return txs.getTransaction(ctx, trS.Group, id)
}

// returns confirmed Transaction by id or hash
func (txs *TransactionService) GetConfirmedTransaction(ctx context.Context, id string) (Transaction, error) {
	return txs.getTransaction(ctx, confirmed.String(), id)
}

// returns confirmed Transactions
func (txs *TransactionService) GetConfirmedTransactions(ctx context.Context, tpOpts *TransactionsPageOptions) (*TransactionsPage, error) {
	return txs.getTransactionsByGroup(ctx, confirmed, tpOpts)
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetConfirmedTransactionsByIds(ctx context.Context, ids []string, tpOpts *TransactionsPageOptions) ([]Transaction, error) {
	return txs.getTransactionsByIds(ctx, confirmed, ids, tpOpts)
}

// returns unconfirmed Transaction by id or hash
func (txs *TransactionService) GetUnconfirmedTransaction(ctx context.Context, id string) (Transaction, error) {
	return txs.getTransaction(ctx, unconfirmed.String(), id)
}

// returns unconfirmed Transactions
func (txs *TransactionService) GetUnconfirmedTransactions(ctx context.Context, tpOpts *TransactionsPageOptions) (*TransactionsPage, error) {
	return txs.getTransactionsByGroup(ctx, unconfirmed, tpOpts)
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetUnconfirmedTransactionsByIds(ctx context.Context, ids []string, tpOpts *TransactionsPageOptions) ([]Transaction, error) {
	return txs.getTransactionsByIds(ctx, unconfirmed, ids, tpOpts)
}

// returns partial Transaction by id or hash
func (txs *TransactionService) GetPartialTransaction(ctx context.Context, id string) (Transaction, error) {
	return txs.getTransaction(ctx, partial.String(), id)
}

// returns partial Transactions
func (txs *TransactionService) GetPartialTransactions(ctx context.Context, tpOpts *TransactionsPageOptions) (*TransactionsPage, error) {
	return txs.getTransactionsByGroup(ctx, partial, tpOpts)
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetPartialTransactionsByIds(ctx context.Context, ids []string, tpOpts *TransactionsPageOptions) ([]Transaction, error) {
	return txs.getTransactionsByIds(ctx, partial, ids, tpOpts)
}

// returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetTransactions(ctx context.Context, ids []string) ([]Transaction, error) {
	var b bytes.Buffer
	txIds := &TransactionIdsDTO{
		ids,
	}

	resp, err := txs.client.doNewRequest(ctx, http.MethodPost, transactionsRoute, txIds, &b)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return MapTransactions(&b, txs.client.GenerationHash())
}

// returns transaction hash after announcing passed SignedTransaction
func (txs *TransactionService) Announce(ctx context.Context, tx *SignedTransaction) (string, error) {
	dto := signedTransactionDto{
		tx.EntityType,
		tx.Payload,
		tx.Hash.String(),
	}
	return txs.announceTransaction(ctx, &dto, transactionsRoute)
}

// returns transaction hash after announcing passed aggregate bounded SignedTransaction
func (txs *TransactionService) AnnounceAggregateBonded(ctx context.Context, tx *SignedTransaction) (string, error) {
	dto := signedTransactionDto{
		tx.EntityType,
		tx.Payload,
		tx.Hash.String(),
	}
	return txs.announceTransaction(ctx, &dto, fmt.Sprintf(transactionsByGroupRoute, partial.String()))
}

// returns transaction hash after announcing passed CosignatureSignedTransaction
func (txs *TransactionService) AnnounceAggregateBondedCosignature(ctx context.Context, c *CosignatureSignedTransaction) (string, error) {
	dto := cosignatureSignedTransactionDto{
		c.ParentHash.String(),
		c.Signature.String(),
		c.Signer,
	}
	return txs.announceTransaction(ctx, &dto, announceAggregateCosignatureRoute)
}

// returns TransactionStatus for passed transaction id or hash
func (txs *TransactionService) GetTransactionStatus(ctx context.Context, id string) (*TransactionStatus, error) {
	ts := &transactionStatusDTO{}

	resp, err := txs.client.doNewRequest(ctx, http.MethodGet, fmt.Sprintf(transactionStatusByIdRoute, id), nil, &ts)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return ts.toStruct()
}

// returns TransactionsStatuses for passed transactions id or hashes
func (txs *TransactionService) GetTransactionsStatuses(ctx context.Context, hashes []string) ([]*TransactionStatus, error) {
	txIds := &TransactionHashesDTO{
		hashes,
	}

	dtos := transactionStatusDTOs(make([]*transactionStatusDTO, len(hashes)))
	resp, err := txs.client.doNewRequest(ctx, http.MethodPost, transactionStatusRoute, txIds, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

func (txs *TransactionService) announceTransaction(ctx context.Context, tx interface{}, path string) (string, error) {
	m := struct {
		Message string `json:"message"`
	}{}

	resp, err := txs.client.doNewRequest(ctx, http.MethodPut, path, tx, &m)
	if err != nil {
		return "", err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return "", err
	}

	return m.Message, nil
}

// Gets a transaction's effective paid fee
func (txs *TransactionService) GetTransactionEffectiveFee(ctx context.Context, transactionId string) (int, error) {
	tx, err := txs.GetAnyTransactionById(ctx, transactionId)
	if err != nil {
		return -1, err
	}

	block, err := txs.BlockchainService.GetBlockByHeight(ctx, tx.GetAbstractTransaction().Height)
	if err != nil {
		return -1, err
	}

	return int(block.FeeMultiplier) * tx.Size(), nil
}
