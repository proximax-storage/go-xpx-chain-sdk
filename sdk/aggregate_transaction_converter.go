package sdk

import "bytes"

func newAggregateTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter, cosignatureConverter aggregateTransactionCosignatureConverter) aggregateTransactionConverter {
	return &aggregateTransactionConverterImpl{
		abstractTransactionConverter:             converter,
		transactionInfoConverter:                 infoConverter,
		aggregateTransactionCosignatureConverter: cosignatureConverter,
	}
}

type aggregateTransactionConverter interface {
	Convert(*aggregateTransactionDTO) (*AggregateTransaction, error)
}

type aggregateTransactionConverterImpl struct {
	abstractTransactionConverter             abstractTransactionConverter
	transactionInfoConverter                 transactionInfoConverter
	aggregateTransactionCosignatureConverter aggregateTransactionCosignatureConverter
}

func (c *aggregateTransactionConverterImpl) Convert(dto *aggregateTransactionDTO) (*AggregateTransaction, error) {
	txsr, err := json.Marshal(dto.Tx.InnerTransactions)
	if err != nil {
		return nil, err
	}

	txs, err := MapTransactions(bytes.NewBuffer(txsr))
	if err != nil {
		return nil, err
	}

	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)
	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	as := make([]*AggregateTransactionCosignature, len(dto.Tx.Cosignatures))
	for i, a := range dto.Tx.Cosignatures {
		as[i], err = c.aggregateTransactionCosignatureConverter.Convert(a, atx.NetworkType)
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
