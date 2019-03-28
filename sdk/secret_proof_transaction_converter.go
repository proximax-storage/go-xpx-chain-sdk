package sdk

func newSecretProofTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) secretProofTransactionConverter {
	return &secretProofTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type secretProofTransactionConverter interface {
	Convert(*secretProofTransactionDTO) (*SecretProofTransaction, error)
}

type secretProofTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *secretProofTransactionConverterImpl) Convert(dto *secretProofTransactionDTO) (*SecretProofTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	return &SecretProofTransaction{
		*atx,
		dto.Tx.HashType,
		dto.Tx.Secret,
		dto.Tx.Proof,
	}, nil
}
