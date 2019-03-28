package sdk

func newMultisigCosignatoryModificationConverter(factory AccountFactory) multisigCosignatoryModificationConverter {
	return &multisigCosignatoryModificationConverterImpl{
		accountFactory: factory,
	}
}

type multisigCosignatoryModificationConverter interface {
	Convert(*multisigCosignatoryModificationDTO, NetworkType) (*MultisigCosignatoryModification, error)
}

type multisigCosignatoryModificationConverterImpl struct {
	accountFactory AccountFactory
}

func (c *multisigCosignatoryModificationConverterImpl) Convert(dto *multisigCosignatoryModificationDTO, networkType NetworkType) (*MultisigCosignatoryModification, error) {
	acc, err := c.accountFactory.NewAccountFromPublicKey(dto.PublicAccount, networkType)
	if err != nil {
		return nil, err
	}

	return &MultisigCosignatoryModification{
		dto.Type,
		acc,
	}, nil
}
