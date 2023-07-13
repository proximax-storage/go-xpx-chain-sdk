package main

import (
	"encoding/binary"
	"errors"

	crypto "github.com/proximax-storage/go-xpx-crypto"
)

const (
	// ConnectionSecurityModes
	NoneConnectionSecurity   = ConnectionSecurityMode(1)
	SignedConnectionSecurity = ConnectionSecurityMode(2)

	// Packet types
	ServerChallengePacketType = PacketType(1)
	ClientChallengePacketType = PacketType(2)
	ChainInfoPacketType       = PacketType(5)

	// Sizes
	PacketHeaderSize = 4 + 4 // Size + PacketTypeSize
	ChallengeSize    = 64
	SignatureSize    = 64
	PublicKeySize    = 32
	SecurityModeSize = 1

	ServerChallengeRequestSize  = PacketHeaderSize + ChallengeSize
	ServerChallengeResponseSize = PacketHeaderSize + ChallengeSize + SignatureSize + PublicKeySize + SecurityModeSize
	ClientChallengeResponseSize = PacketHeaderSize + SignatureSize
	ChainInfoSize               = PacketHeaderSize + 8 + 8 + 8
)

var (
	ErrPackerSizeTooSmall = errors.New("packer size too small")
	ErrMismatchSizes      = errors.New("read buffer size and parsed packer size are mismatched")
)

type (
	Packet interface {
		Bytes() []byte
		Parse([]byte) error
	}

	// PacketType is an enumeration of known packet types.
	PacketType uint32

	Challenge [ChallengeSize]byte

	ConnectionSecurityMode uint8

	// PacketHeader that all transferable information is expected to have.
	PacketHeader struct {
		// Size of the packet.
		Size uint32
		// Type of the packet.
		Type PacketType
	}

	// ServerChallengeRequest Packet representing a challenge request from a server to a client.
	ServerChallengeRequest struct {
		*PacketHeader

		// Challenge data that should be signed by the client.
		Challenge Challenge
	}

	// ServerChallengeResponse is a packet representing a challenge response and new challenge request from a client to a server.
	ServerChallengeResponse struct {
		*PacketHeader

		// Challenge data that should be signed by the server.
		Challenge Challenge
		// Client's signature on the server challenge and any additional request information.
		Signature *crypto.Signature
		// Client's public key.
		PublicKey *crypto.PublicKey
		// Security mode requested by the client.
		SecurityMode ConnectionSecurityMode
	}

	// ClientChallengeResponse is a packet representing a challenge response from a server to a client.
	ClientChallengeResponse struct {
		*PacketHeader

		// Server's signature on the client challenge.
		Signature *crypto.Signature
	}

	ChainInfoResponse struct {
		*PacketHeader

		// Chain height.
		Height uint64
		// High part of the score.
		ScoreHigh uint64
		// Low part of the score.
		ScoreLow uint64
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

func NewPacketHeader(pt PacketType) *PacketHeader {
	return &PacketHeader{
		Size: PacketHeaderSize,
		Type: pt,
	}
}

func NewServerChallengeRequest() *ServerChallengeRequest {
	ph := NewPacketHeader(ServerChallengePacketType)
	ph.Size += ChallengeSize
	return &ServerChallengeRequest{
		PacketHeader: ph,
		Challenge:    Challenge{},
	}
}

func (s *ServerChallengeRequest) Bytes() []byte {
	return append(s.PacketHeader.Bytes(), s.Challenge[:]...)
}

func (s *ServerChallengeRequest) Parse(buff []byte) error {
	ph := &PacketHeader{}
	err := ph.Parse(buff)
	if err != nil {
		return err
	}

	s.PacketHeader = ph
	copy(s.Challenge[:], buff[PacketHeaderSize:])

	return nil
}

func NewServerChallengeResponse() *ServerChallengeResponse {
	ph := NewPacketHeader(ServerChallengePacketType)
	ph.Size = ServerChallengeResponseSize
	return &ServerChallengeResponse{
		PacketHeader: ph,
		Challenge:    Challenge{},
	}
}

func (s *ServerChallengeResponse) Bytes() []byte {
	buf := make([]byte, 0, s.Size)
	buf = append(buf, s.PacketHeader.Bytes()...)
	buf = append(buf, s.Challenge[:]...)
	buf = append(buf, s.Signature.Bytes()...)
	buf = append(buf, s.PublicKey.Raw...)
	buf = append(buf, byte(s.SecurityMode))

	return buf
}

func (s *ServerChallengeResponse) Parse(buff []byte) error {
	ph := &PacketHeader{}
	err := ph.Parse(buff)
	if err != nil {
		return err
	}

	s.PacketHeader = ph
	buff = buff[PacketHeaderSize:]
	copy(s.Challenge[:], buff[:ChallengeSize])

	buff = buff[ChallengeSize:]
	sig, err := crypto.NewSignatureFromBytes(buff[:SignatureSize])
	if err != nil {
		s.PacketHeader = nil
		s.Challenge = Challenge{}
		return err
	}
	s.Signature = sig

	buff = buff[SignatureSize:]
	s.PublicKey = crypto.NewPublicKey(buff[:PublicKeySize])
	s.SecurityMode = ConnectionSecurityMode(buff[len(buff)-1])

	return nil
}

func NewClientChallengeResponse() *ClientChallengeResponse {
	ph := NewPacketHeader(ClientChallengePacketType)
	ph.Size += SignatureSize
	return &ClientChallengeResponse{
		PacketHeader: ph,
		Signature:    &crypto.Signature{},
	}
}

func (c *ClientChallengeResponse) Bytes() []byte {
	return append(c.PacketHeader.Bytes(), c.Signature.Bytes()...)
}

func (c *ClientChallengeResponse) Parse(buff []byte) error {
	ph := &PacketHeader{}
	err := ph.Parse(buff)
	if err != nil {
		return err
	}

	sig, err := crypto.NewSignatureFromBytes(buff[PacketHeaderSize:])
	if err != nil {
		return err
	}
	c.Signature = sig
	c.PacketHeader = ph

	return nil
}

func (cir *ChainInfoResponse) Bytes() []byte {
	buff := make([]byte, ChainInfoSize)
	copy(buff, cir.PacketHeader.Bytes())

	offset := PacketHeaderSize
	binary.LittleEndian.PutUint64(buff[offset:], cir.Height)
	offset += 8
	binary.LittleEndian.PutUint64(buff[offset:], cir.ScoreLow)
	offset += 8
	binary.LittleEndian.PutUint64(buff[offset:], cir.ScoreHigh)

	return buff
}

func (cir *ChainInfoResponse) Parse(buff []byte) error {
	ph := &PacketHeader{}
	err := ph.Parse(buff)
	if err != nil {
		return err
	}

	offset := PacketHeaderSize
	cir.Height = binary.LittleEndian.Uint64(buff[offset : offset+8])
	offset += 8
	cir.ScoreHigh = binary.LittleEndian.Uint64(buff[offset : offset+8])
	offset += 8
	cir.ScoreLow = binary.LittleEndian.Uint64(buff[offset : offset+8])

	return nil
}
