package sdk

func newSecretLockTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) secretLockTransactionConverter {
	return &secretLockTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type secretLockTransactionConverter interface {
	Convert(*secretLockTransactionDTO) (*SecretLockTransaction, error)
}

type secretLockTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *secretLockTransactionConverterImpl) Convert(dto *secretLockTransactionDTO) (*SecretLockTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)
	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	a, err := NewAddressFromEncoded(dto.Tx.Recipient)
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(dto.Tx.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	mosaic, err := NewMosaic(mosaicId, dto.Tx.Amount.toBigInt())
	if err != nil {
		return nil, err
	}

	return &SecretLockTransaction{
		*atx,
		mosaic,
		dto.Tx.HashType,
		dto.Tx.Duration.toBigInt(),
		dto.Tx.Secret,
		a,
	}, nil
}
