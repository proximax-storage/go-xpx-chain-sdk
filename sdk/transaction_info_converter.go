package sdk

import "math/big"

func newTransactionInfoConverter() transactionInfoConverter {
	return &transactionInfoConverterImpl{}
}

type transactionInfoConverter interface {
	Convert(transactionInfoDTO) *TransactionInfo
}

type transactionInfoConverterImpl struct{}

func (*transactionInfoConverterImpl) Convert(dto transactionInfoDTO) *TransactionInfo {
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
