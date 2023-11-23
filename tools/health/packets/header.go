package packets

import (
	"encoding/binary"
	"errors"
)

const (
	// Packet types
	ServerChallengePacketType        = PacketType(1)
	ClientChallengePacketType        = PacketType(2)
	ChainInfoPacketType              = PacketType(5)
	BlockHashesPacketType            = PacketType(7)
	NodeDiscoveryPullPeersPacketType = PacketType(603)

	// Sizes
	PacketHeaderSize = 4 + 4 // Size + PacketTypeSize
	ChallengeSize    = 64
	SignatureSize    = 64
	HashSize         = 32
	PublicKeySize    = 32
	SecurityModeSize = 1

	ServerChallengeRequestSize            = PacketHeaderSize + ChallengeSize
	ServerChallengeResponseSize           = PacketHeaderSize + ChallengeSize + SignatureSize + PublicKeySize + SecurityModeSize
	ClientChallengeResponseSize           = PacketHeaderSize + SignatureSize
	ChainInfoResponseSize                 = PacketHeaderSize + 8 + 8 + 8
	BlockHashesRequestSize                = PacketHeaderSize + 8 + 4
	MinNodeDiscoveryPullPeersResponseSize = PacketHeaderSize + 4 + PublicKeySize + 2 + PublicKeySize + 4 + 4 + 1 + 1
)

var (
	ErrPackerSizeTooSmall = errors.New("packer size too small")
	ErrMismatchSizes      = errors.New("read buffer size and parsed packer size are mismatched")
)

type (
	// PacketType is an enumeration of known packet types.
	PacketType uint32

	// PacketHeader that all transferable information is expected to have.
	PacketHeader struct {
		// Size of the packet.
		Size uint32
		// Type of the packet.
		Type PacketType
	}
)

func (ph *PacketHeader) Bytes() []byte {
	buff := make([]byte, PacketHeaderSize)
	binary.LittleEndian.PutUint32(buff[:4], ph.Size)
	binary.LittleEndian.PutUint32(buff[4:], uint32(ph.Type))
	return buff
}

func (ph *PacketHeader) Parse(buff []byte) error {
	if len(buff) < PacketHeaderSize {
		return ErrPackerSizeTooSmall
	}

	ph.Size = binary.LittleEndian.Uint32(buff[:4])
	ph.Type = PacketType(binary.LittleEndian.Uint32(buff[4:]))

	if ph.Size != uint32(len(buff)) {
		return ErrMismatchSizes
	}

	return nil
}

func NewPacketHeader(pt PacketType) PacketHeader {
	return PacketHeader{
		Size: PacketHeaderSize,
		Type: pt,
	}
}
