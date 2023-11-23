package packets

import (
	"encoding/binary"

	crypto "github.com/proximax-storage/go-xpx-crypto"
)

type (
	// NetworkNode information about a catapult node that is propagated through the network.
	NetworkNode struct {
		/// Size of the node.
		Size uint32

		/// Unique node identifier (public key).
		IdentityKey *crypto.PublicKey

		/// Port.
		Port uint16

		/// Network identifier.
		NetworkIdentifier uint8

		/// Version.
		Version uint32

		/// Role(s).
		NodeRoles uint32

		/// Size of the host in bytes.
		HostSize uint8

		/// Host of the node.
		Host string

		/// Size of the friendly name in bytes.
		FriendlyNameSize uint8

		/// FriendlyName of node.
		FriendlyName string
	}

	NodeDiscoveryPullPeersResponse struct {
		PacketHeader

		NetworkNodes []*NetworkNode
	}
)

func (n *NodeDiscoveryPullPeersResponse) Parse(buff []byte) error {
	n.NetworkNodes = make([]*NetworkNode, 0, 10)
	for len(buff) >= MinNodeDiscoveryPullPeersResponseSize-PacketHeaderSize {
		node := &NetworkNode{}

		node.Size = binary.LittleEndian.Uint32(buff[:4])
		buff = buff[4:]

		node.IdentityKey = crypto.NewPublicKey(buff[:PublicKeySize])
		buff = buff[PublicKeySize:]

		node.Port = binary.LittleEndian.Uint16(buff[:2])
		buff = buff[2:]

		node.NetworkIdentifier = buff[0]
		buff = buff[1:]

		node.Version = binary.LittleEndian.Uint32(buff[:4])
		buff = buff[4:]

		node.NodeRoles = binary.LittleEndian.Uint32(buff[:4])
		buff = buff[4:]

		node.HostSize = buff[0]
		buff = buff[1:]

		node.FriendlyNameSize = buff[0]
		buff = buff[1:]

		node.Host = string(buff[:node.HostSize])
		buff = buff[node.HostSize:]

		node.FriendlyName = string(buff[:node.FriendlyNameSize])
		buff = buff[node.FriendlyNameSize:]

		n.NetworkNodes = append(n.NetworkNodes, node)
	}

	return nil
}
