// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewReplicatorOffboardingTransaction(
	deadline *Deadline,
	networkType NetworkType,
) (*ReplicatorOffboardingTransaction, error) {

	return &ReplicatorOffboardingTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     ReplicatorOffboardingVersion,
			Type:        ReplicatorOffboarding,
			NetworkType: networkType,
		},
	}, nil
}

func (tx *ReplicatorOffboardingTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ReplicatorOffboardingTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *ReplicatorOffboardingTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.ReplicatorOffboardingTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return replicatorOffboardingTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ReplicatorOffboardingTransaction) Size() int {
	return ReplicatorOffboardingHeaderSize
}

type replicatorOffboardingTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *replicatorOffboardingTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &ReplicatorOffboardingTransaction{
		*atx,
	}, nil
}