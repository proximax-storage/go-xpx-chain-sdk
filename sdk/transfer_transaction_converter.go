package sdk

func newTransferTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) transferTransactionConverter {
	return &transferTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type transferTransactionConverter interface {
	Convert(*transferTransactionDTO) (*TransferTransaction, error)
}

type transferTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *transferTransactionConverterImpl) Convert(dto *transferTransactionDTO) (*TransferTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.Mosaics))

	for i, mosaic := range dto.Tx.Mosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	a, err := NewAddressFromEncoded(dto.Tx.Address)
	if err != nil {
		return nil, err
	}

	return &TransferTransaction{
		*atx,
		dto.Tx.Message.toStruct(),
		mosaics,
		a,
	}, nil
}
