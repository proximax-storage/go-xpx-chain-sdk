// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type TransactionService struct {
	*service
	BlockchainService *BlockchainService
}

// GetTransaction returns Transaction for passed transaction id or hash
func (txs *TransactionService) GetTransaction(ctx context.Context, group TransactionGroup, id string) (Transaction, error) {
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

// GetAnyTransaction returns Transaction for passed transaction id or hash
func (txs *TransactionService) GetAnyTransaction(ctx context.Context, id string) (Transaction, error) {
	trS, err := txs.GetTransactionStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	return txs.GetTransaction(ctx, trS.Group, id)
}

// GetTransactionsByGroup returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetTransactionsByGroup(ctx context.Context, group TransactionGroup, tpOpts *TransactionsPageOptions) (*TransactionsPage, error) {
	tspDTO := &transactionsPageDTO{}

	u, err := addOptions(fmt.Sprintf(transactionsByGroupRoute, group), tpOpts)
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

// GetTransactionsByIds returns an array of Transaction's for passed array of transaction ids or hashes
func (txs *TransactionService) GetTransactionsByIds(ctx context.Context, group TransactionGroup, ids []string, tpOpts *TransactionsPageOptions) ([]Transaction, error) {
	var b bytes.Buffer
	txIds := &TransactionIdsDTO{
		ids,
	}

	u, err := addOptions(fmt.Sprintf(transactionsByGroupRoute, group), tpOpts)
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

// Announce returns transaction hash after announcing passed SignedTransaction
func (txs *TransactionService) Announce(ctx context.Context, tx *SignedTransaction) (string, error) {
	dto := signedTransactionDto{
		tx.EntityType,
		tx.Payload,
		tx.Hash.String(),
	}

	msg, err := txs.announceTransaction(ctx, &dto, transactionsRoute)
	if err != nil {
		return "", errors.Errorf("%s. Message: %s", err, msg)
	}

	return tx.Hash.String(), nil
}

// AnnounceAggregateBonded returns transaction hash after announcing passed aggregate bounded SignedTransaction
func (txs *TransactionService) AnnounceAggregateBonded(ctx context.Context, tx *SignedTransaction) (string, error) {
	dto := signedTransactionDto{
		tx.EntityType,
		tx.Payload,
		tx.Hash.String(),
	}

	msg, err := txs.announceTransaction(ctx, &dto, announceAggregateRoute)
	if err != nil {
		return "", errors.Errorf("%s. Message: %s", err, msg)
	}

	return tx.Hash.String(), nil
}

// AnnounceAggregateBondedCosignature returns transaction hash after announcing passed CosignatureSignedTransaction
func (txs *TransactionService) AnnounceAggregateBondedCosignature(ctx context.Context, c *CosignatureSignedTransaction) (string, error) {
	dto := cosignatureSignedTransactionDto{
		c.ParentHash.String(),
		c.Signature.String(),
		c.Scheme,
		c.Signer,
	}

	msg, err := txs.announceTransaction(ctx, &dto, announceAggregateCosignatureRoute)
	if err != nil {
		return "", errors.Errorf("%s. Message: %s", err, msg)
	}

	return c.ParentHash.String(), nil
}

// GetTransactionStatus returns TransactionStatus for passed transaction id or hash
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

// GetTransactionsStatuses returns TransactionsStatuses for passed transactions id or hashes
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

// GetTransactionEffectiveFee gets a transaction's effective paid fee
func (txs *TransactionService) GetTransactionEffectiveFee(ctx context.Context, transactionId string) (int, error) {
	tx, err := txs.GetTransaction(ctx, Confirmed, transactionId)
	if err != nil {
		return -1, err
	}

	block, err := txs.BlockchainService.GetBlockByHeight(ctx, tx.GetAbstractTransaction().Height)
	if err != nil {
		return -1, err
	}

	return int(block.FeeMultiplier) * tx.Size(), nil
}
