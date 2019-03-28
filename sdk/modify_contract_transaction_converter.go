package sdk

func newModifyContractTransactionConverter(converter abstractTransactionConverter, infoConverter transactionInfoConverter) modifyContractTransactionConverter {
	return &modifyContractTransactionConverterImpl{
		abstractTransactionConverter: converter,
		transactionInfoConverter:     infoConverter,
	}
}

type modifyContractTransactionConverter interface {
	Convert(*modifyContractTransactionDTO) (*ModifyContractTransaction, error)
}

type modifyContractTransactionConverterImpl struct {
	abstractTransactionConverter abstractTransactionConverter
	transactionInfoConverter     transactionInfoConverter
}

func (c *modifyContractTransactionConverterImpl) Convert(dto *modifyContractTransactionDTO) (*ModifyContractTransaction, error) {
	transactionInfo := c.transactionInfoConverter.Convert(dto.TDto)

	atx, err := c.abstractTransactionConverter.Convert(dto.Tx.abstractTransactionDTO, transactionInfo)
	if err != nil {
		return nil, err
	}

	customers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Customers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	executors, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Executors, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	verifiers, err := multisigCosignatoryDTOArrayToStruct(dto.Tx.Verifiers, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ModifyContractTransaction{
		*atx,
		dto.Tx.DurationDelta.toBigInt().Int64(),
		dto.Tx.Hash,
		customers,
		executors,
		verifiers,
	}, nil
}
