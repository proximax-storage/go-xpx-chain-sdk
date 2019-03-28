package sdk

func newLockFundsTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) lockFundsTransactionConverter {
	return &lockFundsTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type lockFundsTransactionConverter interface {
	Convert(*lockFundsTransactionDTO) (*LockFundsTransaction, error)
}

type lockFundsTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *lockFundsTransactionConverterImpl) Convert(dto *lockFundsTransactionDTO) (*LockFundsTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	mosaic, err := dto.Tx.Mosaic.toStruct()
	if err != nil {
		return nil, err
	}

	return &LockFundsTransaction{
		*atx,
		mosaic,
		dto.Tx.Duration.toBigInt(),
		&SignedTransaction{Lock, "", dto.Tx.Hash},
	}, nil
}
