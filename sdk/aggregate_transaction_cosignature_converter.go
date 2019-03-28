package sdk

func newAggregateTransactionCosignatureConverter(factory AccountFactory) aggregateTransactionCosignatureConverter {
	return &aggregateTransactionCosignatureConverterImpl{
		accountFactory: factory,
	}
}

type aggregateTransactionCosignatureConverter interface {
	Convert(*aggregateTransactionCosignatureDTO, NetworkType) (*AggregateTransactionCosignature, error)
}

type aggregateTransactionCosignatureConverterImpl struct {
	accountFactory AccountFactory
}

func (c *aggregateTransactionCosignatureConverterImpl) Convert(dto *aggregateTransactionCosignatureDTO, networkType NetworkType) (*AggregateTransactionCosignature, error) {
	acc, err := c.accountFactory.NewAccountFromPublicKey(dto.Signer, networkType)
	if err != nil {
		return nil, err
	}
	return &AggregateTransactionCosignature{
		dto.Signature,
		acc,
	}, nil
}
