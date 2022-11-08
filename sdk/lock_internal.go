// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type commonLockInfoDTO struct {
	Account  string         `json:"account"`
	MosaicId uint64DTO      `json:"mosaicId"`
	Amount   uint64DTO      `json:"amount"`
	Height   uint64DTO      `json:"height"`
	Status   LockStatusType `json:"status"`
}

type hashLockInfoDTO struct {
	Lock struct {
		commonLockInfoDTO
		Hash hashDto `json:"hash"`
	} `json:"lock"`
}

type secretLockInfoDTO struct {
	Lock struct {
		commonLockInfoDTO
		Recipient     string   `json:"recipient"`
		Secret        hashDto  `json:"secret"`
		HashAlgorithm HashType `json:"hashAlgorithm"`
		CompositeHash hashDto  `json:"compositeHash"`
	} `json:"lock"`
}

func (ref *commonLockInfoDTO) toStruct(networkType NetworkType) (*CommonLockInfo, error) {
	account, err := NewAccountFromPublicKey(ref.Account, networkType)
	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(ref.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}

	return &CommonLockInfo{
		Account:  account,
		MosaicId: mosaicId,
		Amount:   ref.Amount.toStruct(),
		Height:   ref.Height.toStruct(),
		Status:   ref.Status,
	}, nil
}

func (ref *hashLockInfoDTO) toStruct(networkType NetworkType) (*HashLockInfo, error) {
	commonLockInfo, err := ref.Lock.commonLockInfoDTO.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	hash, err := ref.Lock.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &HashLockInfo{
		CommonLockInfo: *commonLockInfo,
		Hash:           hash,
	}, nil
}

func (ref *secretLockInfoDTO) toStruct(networkType NetworkType) (*SecretLockInfo, error) {
	commonLockInfo, err := ref.Lock.commonLockInfoDTO.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	compositeHash, err := ref.Lock.CompositeHash.Hash()
	if err != nil {
		return nil, err
	}

	secret, err := ref.Lock.Secret.Hash()
	if err != nil {
		return nil, err
	}

	recipient, err := NewAddressFromHexString(ref.Lock.Recipient)
	if err != nil {
		return nil, err
	}

	return &SecretLockInfo{
		CommonLockInfo: *commonLockInfo,
		CompositeHash:  compositeHash,
		Secret:         secret,
		Recipient:      recipient,
		HashAlgorithm:  ref.Lock.HashAlgorithm,
	}, nil
}

type hashLockInfoDTOs []*hashLockInfoDTO
type secretLockInfoDTOs []*secretLockInfoDTO

func (ref *hashLockInfoDTOs) toStruct(networkType NetworkType) ([]*HashLockInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*HashLockInfo, 0, len(dtos))
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

func (ref *secretLockInfoDTOs) toStruct(networkType NetworkType) ([]*SecretLockInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*SecretLockInfo, 0, len(dtos))
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
