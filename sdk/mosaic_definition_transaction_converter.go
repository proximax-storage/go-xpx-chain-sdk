package sdk

func newMosaicDefinitionTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) mosaicDefinitionTransactionConverter {
	return &mosaicDefinitionTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type mosaicDefinitionTransactionConverter interface {
	Convert(*mosaicDefinitionTransactionDTO) (*MosaicDefinitionTransaction, error)
}

type mosaicDefinitionTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *mosaicDefinitionTransactionConverterImpl) Convert(dto *mosaicDefinitionTransactionDTO) (*MosaicDefinitionTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &MosaicDefinitionTransaction{
		*atx,
		dto.Tx.Properties.toStruct(),
		dto.Tx.MosaicNonce,
		mosaicId,
	}, nil
}
