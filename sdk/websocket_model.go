// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type statusInfoDto struct {
	Status string  `json:"status"`
	Hash   hashDto `json:"hash"`
}

type StatusInfo struct {
	Status string
	Hash   *Hash
}

type signerInfoDto struct {
	Signer     string       `json:"signer"`
	Signature  signatureDto `json:"signature"`
	ParentHash hashDto      `json:"parentHash"`
}

type SignerInfo struct {
	Signer     string
	Signature  *Signature
	ParentHash *Hash
}

type driveStateDto struct {
	DriveKey    string        `json:"driveKey"`
	State       DriveState    `json:"state"`
}

func (dto *driveStateDto) toStruct() (*DriveStateInfo, error) {
	return &DriveStateInfo{
		DriveKey:   dto.DriveKey,
		State:      dto.State,
	}, nil
}

type DriveStateInfo struct {
	DriveKey    string
	State       DriveState
}

type UnconfirmedRemoved struct {
	Meta *TransactionInfo
}

type unconfirmedRemovedDto struct {
	Meta *transactionInfoDTO `json:"meta"`
}

func (dto *unconfirmedRemovedDto) toStruct() (*UnconfirmedRemoved, error) {
	info, err := dto.Meta.toStruct()
	if err != nil {
		return nil, err
	}

	return &UnconfirmedRemoved{
		Meta: info,
	}, nil
}

type partialRemovedInfoDTO struct {
	Meta *transactionInfoDTO `json:"meta"`
}

func (dto partialRemovedInfoDTO) toStruct() (*PartialRemovedInfo, error) {
	info, err := dto.Meta.toStruct()
	if err != nil {
		return nil, err
	}

	return &PartialRemovedInfo{
		Meta: info,
	}, nil
}

type PartialRemovedInfo struct {
	Meta *TransactionInfo
}

type WsMessageInfo struct {
	Address     *Address
	ChannelName string
}

type WsMessageInfoDTO struct {
	Meta wsMessageInfoMetaDTO `json:"meta"`
}

func (dto *WsMessageInfoDTO) ToStruct() (*WsMessageInfo, error) {
	msg := &WsMessageInfo{
		ChannelName: dto.Meta.ChannelName,
	}

	if dto.Meta.ChannelName == "block" {
		return msg, nil
	}

	address, err := NewAddressFromBase32(dto.Meta.Address)
	if err != nil {
		return nil, err
	}

	msg.Address = address

	return msg, nil
}

type wsMessageInfoMetaDTO struct {
	ChannelName string `json:"channelName"`
	Address     string `json:"address"`
}
