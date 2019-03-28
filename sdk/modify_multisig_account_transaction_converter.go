package sdk

func newModifyMultisigAccountTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) modifyMultisigAccountTransactionConverter {
	return &modifyMultisigAccountTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type modifyMultisigAccountTransactionConverter interface {
	Convert(*modifyMultisigAccountTransactionDTO) (*ModifyMultisigAccountTransaction, error)
}

type modifyMultisigAccountTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *modifyMultisigAccountTransactionConverterImpl) Convert(dto *modifyMultisigAccountTransactionDTO) (*ModifyMultisigAccountTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	ms, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Modifications, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyMultisigAccountTransaction{
		*atx,
		dto.Tx.MinApprovalDelta,
		dto.Tx.MinRemovalDelta,
		ms,
	}, nil
}
