package packets

import (
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

type (
	Challenge [ChallengeSize]byte

	ConnectionSecurityMode uint8

	// ServerChallengeRequest Packet representing a challenge request from a server to a client.
	ServerChallengeRequest struct {
		PacketHeader

		// Challenge data that should be signed by the client.
		Challenge Challenge
	}

	// ServerChallengeResponse is a packet representing a challenge response and new challenge request from a client to a server.
	ServerChallengeResponse struct {
		PacketHeader

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
		PacketHeader

		// Server's signature on the client challenge.
		Signature *crypto.Signature
	}
)

const (
	// ConnectionSecurityModes
	NoneConnectionSecurity   = ConnectionSecurityMode(1)
	SignedConnectionSecurity = ConnectionSecurityMode(2)

	//Client = "client"
	//Server = "server"
)

//var headerInfos = map[string]PacketHeader{
//	Server: {Type: ServerChallengePacketType, Size: HeaderChallengeSize},
//	Client: {Type: ClientChallengePacketType, Size: HeaderChallengeSize},
//}

func NewServerChallengeRequest() *ServerChallengeRequest {
	ph := NewPacketHeader(ServerChallengePacketType)
	ph.Size += ChallengeSize
	return &ServerChallengeRequest{
		PacketHeader: ph,
		Challenge:    Challenge{},
	}
}

func (s *ServerChallengeRequest) Parse(buff []byte) error {
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

func (s *ServerChallengeResponse) Header() Header {
	return &s.PacketHeader
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
	copy(s.Challenge[:], buff[:ChallengeSize])

	buff = buff[ChallengeSize:]
	sig, err := crypto.NewSignatureFromBytes(buff[:SignatureSize])
	if err != nil {
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

func (c *ClientChallengeResponse) Parse(buff []byte) error {
	sig, err := crypto.NewSignatureFromBytes(buff[PacketHeaderSize:])
	if err != nil {
		return err
	}
	c.Signature = sig

	return nil
}
