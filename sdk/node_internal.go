// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type nodeInfoDTO struct {
	PublicKey         string      `json:"publicKey"`
	Host              string      `json:"host"`
	FriendlyName      string      `json:"friendlyName"`
	Port              int         `json:"port"`
	Version           int         `json:"version"`
	Roles             int         `json:"roles"`
	NetworkIdentifier NetworkType `json:"networkIdentifier"`
}

func (ref *nodeInfoDTO) toStruct(networkType NetworkType) (*NodeInfo, error) {
	account, err := NewAccountFromPublicKey(ref.PublicKey, networkType)
	if err != nil {
		return nil, err
	}

	return &NodeInfo{
		Account:      account,
		Host:         ref.Host,
		Port:         ref.Port,
		Roles:        ref.Roles,
		FriendlyName: ref.FriendlyName,
		NetworkType:  ref.NetworkIdentifier,
	}, nil
}

type nodeInfoDTOs []*nodeInfoDTO

func (ref *nodeInfoDTOs) toStruct(networkType NetworkType) ([]*NodeInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*NodeInfo, 0, len(dtos))
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

type timeDTO struct {
	CommunicationTimestamps struct {
		SendTimestamp    uint64DTO `json:"sendTimestamp"`
		ReceiveTimestamp uint64DTO `json:"receiveTimestamp"`
	} `json:"communicationTimestamps"`
}

func (ref *timeDTO) toStruct(NetworkType) (*BlockchainTimestamp, error) {
	return NewBlockchainTimestamp(int64(ref.CommunicationTimestamps.ReceiveTimestamp.toUint64())), nil
}
