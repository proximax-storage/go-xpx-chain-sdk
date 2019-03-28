package sdk

func newMosaicSupplyChangeTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) mosaicSupplyChangeTransactionConverter {
	return &mosaicSupplyChangeTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type mosaicSupplyChangeTransactionConverter interface {
	Convert(*mosaicSupplyChangeTransactionDTO) (*MosaicSupplyChangeTransaction, error)
}

type mosaicSupplyChangeTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *mosaicSupplyChangeTransactionConverterImpl) Convert(dto *mosaicSupplyChangeTransactionDTO) (*MosaicSupplyChangeTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)
	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &MosaicSupplyChangeTransaction{
		*atx,
		dto.Tx.MosaicSupplyType,
		mosaicId,
		dto.Tx.Delta.toBigInt(),
	}, nil
}
