// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type superContractDTO struct {
	SuperContract struct {
		SuperContractKey string             `json:"multisig"`
		Start            uint64DTO          `json:"start"`
		End              uint64DTO          `json:"end"`
		MainDriveKey     string             `json:"mainDriveKey"`
		FileHash         hashDto            `json:"fileHash"`
		Version          uint64DTO          `json:"vmVersion"`
	}
}

func (ref *superContractDTO) toStruct(networkType NetworkType) (*SuperContract, error) {
	contract := SuperContract{}

	driveAccount, err := NewAccountFromPublicKey(ref.SuperContract.MainDriveKey, networkType)
	if err != nil {
		return nil, err
	}

	scAcc, err := NewAccountFromPublicKey(ref.SuperContract.SuperContractKey, networkType)
	if err != nil {
		return nil, err
	}

	fileHash, err := ref.SuperContract.FileHash.Hash()
	if err != nil {
		return nil, err
	}

	contract.Account = scAcc
	contract.Drive = driveAccount
	contract.FileHash = fileHash
	contract.VMVersion = ref.SuperContract.Version.toUint64()
	contract.Start = ref.SuperContract.Start.toStruct()
	contract.End = ref.SuperContract.End.toStruct()

	return &contract, nil
}

type superContractDTOs []*superContractDTO

func (ref *superContractDTOs) toStruct(networkType NetworkType) ([]*SuperContract, error) {
	var (
		dtos   = *ref
		contracts = make([]*SuperContract, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, info)
	}

	return contracts, nil
}
