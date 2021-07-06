// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type nodeRole uint8

const (
	None nodeRole = 0x00
	Peer nodeRole = 0x01
	Api  nodeRole = 0x02
)

type NodeUnlockedAccount struct {
	Account *PublicAccount
}

func (n *NodeUnlockedAccount) String() string {
	return fmt.Sprintf(
		`{ "Account": %s}`,
		n.Account,
	)
}

type NodeInfo struct {
	Account      *PublicAccount
	Host         string
	Port         int
	FriendlyName string
	NetworkType  NetworkType
	Roles        int
}

func (n *NodeInfo) isApiNode() bool {
	return n.Roles&int(Api) > 0
}

func (n *NodeInfo) isPeerNode() bool {
	return n.Roles&int(Peer) > 0
}

func (n *NodeInfo) String() string {
	return fmt.Sprintf(
		`{ "Account": %s, "Host": %s, "Port": %d, "FriendlyName": %s, "NetworkType": %s, "Roles": %d }`,
		n.Account,
		n.Host,
		n.Port,
		n.FriendlyName,
		n.NetworkType,
		n.Roles,
	)
}
