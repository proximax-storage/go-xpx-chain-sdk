package sdk

import "math/big"

func newRegisterNamespaceTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) registerNamespaceTransactionConverter {
	return &registerNamespaceTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type registerNamespaceTransactionConverter interface {
	Convert(dto *registerNamespaceTransactionDTO) (*RegisterNamespaceTransaction, error)
}

type registerNamespaceTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *registerNamespaceTransactionConverterImpl) Convert(dto *registerNamespaceTransactionDTO) (*RegisterNamespaceTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	d := big.NewInt(0)
	n := &NamespaceId{}

	if dto.Tx.NamespaceType == Root {
		d = dto.Tx.Duration.toBigInt()
	} else {
		n, err = dto.Tx.ParentId.toStruct()
		if err != nil {
			return nil, err
		}
	}

	nsId, err := dto.Tx.Id.toStruct()
	if err != nil {
		return nil, err
	}

	return &RegisterNamespaceTransaction{
		*atx,
		nsId,
		dto.Tx.NamespaceType,
		dto.Tx.NamspaceName,
		d,
		n,
	}, nil
}
