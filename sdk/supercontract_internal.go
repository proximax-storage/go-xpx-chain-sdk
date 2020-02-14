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

type operationDTO struct {
	Operation struct {
		Initiator           string              `json:"account"`
		Height              uint64DTO           `json:"height"`
		Mosaics             []*mosaicDTO        `json:"mosaics"`
		Token               hashDto             `json:"token"`
		Status              OperationStatus     `json:"result"`
		Executors           []string            `json:"executors"`
		TransactionHashes   []*hashDto          `json:"transactionHashes"`
	}
}

func (ref *operationDTO) toStruct(networkType NetworkType) (*Operation, error) {
	operation := Operation{}

	var err error
	operation.Initiator, err = NewAccountFromPublicKey(ref.Operation.Initiator, networkType)
	if err != nil {
		return nil, err
	}

	operation.LockedMosaics = make([]*Mosaic, len(ref.Operation.Mosaics))
	for i, mosaic := range ref.Operation.Mosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		operation.LockedMosaics[i] = msc
	}

	operation.Token, err = ref.Operation.Token.Hash()
	if err != nil {
		return nil, err
	}

	operation.Executors = make([]*PublicAccount, len(ref.Operation.Executors))
	for i, executor := range ref.Operation.Executors {
		operation.Executors[i], err = NewAccountFromPublicKey(executor, networkType)
		if err != nil {
			return nil, err
		}
	}

	operation.AggregateHashes = make([]*Hash, len(ref.Operation.TransactionHashes))
	for i, hash := range ref.Operation.TransactionHashes {
		operation.AggregateHashes[i], err = hash.Hash()
		if err != nil {
			return nil, err
		}
	}

	operation.Height = ref.Operation.Height.toStruct()
	operation.Status = ref.Operation.Status

	return &operation, nil
}

type operationDTOs []*operationDTO

func (ref *operationDTOs) toStruct(networkType NetworkType) ([]*Operation, error) {
	var (
		dtos   = *ref
		objects = make([]*Operation, len(dtos))
	)

	for i, dto := range dtos {
		object, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		objects[i] = object
	}

	return objects, nil
}
