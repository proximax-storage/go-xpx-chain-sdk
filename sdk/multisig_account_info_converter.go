package sdk

func newMultisigAccountInfoDTOConverter(factory AccountFactory) multisigAccountInfoConverter {
	return &multisigAccountInfoDTOConverterImpl{
		accountFactory: factory,
	}
}

type multisigAccountInfoConverter interface {
	Convert(*multisigAccountInfoDTO, NetworkType) (*MultisigAccountInfo, error)
	ConvertMulti(*multisigAccountGraphInfoDTOS, NetworkType) (*MultisigAccountGraphInfo, error)
}

type multisigAccountInfoDTOConverterImpl struct {
	accountFactory AccountFactory
}

func (cnv *multisigAccountInfoDTOConverterImpl) Convert(dto *multisigAccountInfoDTO, networkType NetworkType) (*MultisigAccountInfo, error) {
	cs := make([]*PublicAccount, len(dto.Multisig.Cosignatories))
	ms := make([]*PublicAccount, len(dto.Multisig.MultisigAccounts))

	acc, err := cnv.accountFactory.NewAccountFromPublicKey(dto.Multisig.Account, networkType)
	if err != nil {
		return nil, err
	}

	for i, c := range dto.Multisig.Cosignatories {
		cs[i], err = cnv.accountFactory.NewAccountFromPublicKey(c, networkType)
		if err != nil {
			return nil, err
		}
	}

	for i, m := range dto.Multisig.MultisigAccounts {
		ms[i], err = cnv.accountFactory.NewAccountFromPublicKey(m, networkType)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &MultisigAccountInfo{
		Account:          *acc,
		MinApproval:      dto.Multisig.MinApproval,
		MinRemoval:       dto.Multisig.MinRemoval,
		Cosignatories:    cs,
		MultisigAccounts: ms,
	}, nil
}

func (cnv *multisigAccountInfoDTOConverterImpl) ConvertMulti(dto *multisigAccountGraphInfoDTOS, networkType NetworkType) (*MultisigAccountGraphInfo, error) {
	var (
		ms  = make(map[int32][]*MultisigAccountInfo)
		err error
	)

	for _, m := range *dto {
		mAccInfos := make([]*MultisigAccountInfo, len(m.Multisigs))

		for idx, c := range m.Multisigs {
			mAccInfos[idx], err = cnv.Convert(&c, networkType)
			if err != nil {
				return nil, err
			}
		}

		ms[m.Level] = mAccInfos
	}

	return &MultisigAccountGraphInfo{ms}, nil
}
