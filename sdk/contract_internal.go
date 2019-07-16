package sdk

type contractInfoDTOs []*contractInfoDTO

func (ref *contractInfoDTOs) toStruct(networkType NetworkType) ([]*ContractInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*ContractInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

type contractInfoDTO struct {
	Contract struct {
		Multisig        string
		MultisigAddress string
		Start           uint64DTO
		Duration        uint64DTO
		Hash            hashDto
		Customers       []string
		Executors       []string
		Verifiers       []string
	}
}

func (ref *contractInfoDTO) toStruct(networkType NetworkType) (*ContractInfo, error) {
	contract := ref.Contract

	hash, err := ref.Contract.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &ContractInfo{
		Multisig:        contract.Multisig,
		MultisigAddress: NewAddress(contract.MultisigAddress, networkType),
		Start:           contract.Start.toStruct(),
		Duration:        contract.Duration.toStruct(),
		Content:         hash,
		Customers:       contract.Customers,
		Executors:       contract.Executors,
		Verifiers:       contract.Verifiers,
	}, nil
}
