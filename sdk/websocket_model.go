// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type StatusInfo struct {
	Status string `json:"status"`
	Hash   Hash   `json:"hash"`
}

type SignerInfo struct {
	Signer     string `json:"signer"`
	Signature  string `json:"signature"`
	ParentHash Hash   `json:"parentHash"`
}

type subscribeDTO struct {
	Uid       string `json:"uid"`
	Subscribe string `json:"subscribe"`
}

type unsubscribeDTO struct {
	Uid         string `json:"uid"`
	Unsubscribe string `json:"unsubscribe"`
}

type wsConnectionResponse struct {
	Uid string `json:"uid"`
}

type UnconfirmedRemoved struct {
	Meta *TransactionInfo
}

type unconfirmedRemovedDto struct {
	Meta *transactionInfoDTO `json:"meta"`
}

func (dto *unconfirmedRemovedDto) toStruct() *UnconfirmedRemoved {
	trInfo := dto.Meta.toStruct()
	return &UnconfirmedRemoved{
		Meta: trInfo,
	}
}

type partialRemovedInfoDTO struct {
	Meta *transactionInfoDTO `json:"meta"`
}

func (dto partialRemovedInfoDTO) toStruct() *PartialRemovedInfo {
	return &PartialRemovedInfo{
		Meta: dto.Meta.toStruct(),
	}
}

type PartialRemovedInfo struct {
	Meta *TransactionInfo
}

type wsMessageInfo struct {
	Address     *Address
	ChannelName string
}

type wsMessageInfoDTO struct {
	Meta wsMessageInfoMetaDTO `json:"meta"`
}

func (dto *wsMessageInfoDTO) toStruct() (*wsMessageInfo, error) {
	msg := &wsMessageInfo{
		ChannelName: dto.Meta.ChannelName,
	}

	if dto.Meta.ChannelName == pathBlock {
		return msg, nil
	}

	address, err := NewAddressFromEncoded(dto.Meta.Address)
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
