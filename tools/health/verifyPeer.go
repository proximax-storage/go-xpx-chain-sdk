package main

import (
	"errors"
	"net"

	crypto "github.com/proximax-storage/go-xpx-crypto"
)

const (
	HeaderChallengeSize = PacketHeaderSize + ChallengeSize

	Client = "client"
	Server = "server"
)

var (
	ErrPacketIsInvalid               = errors.New("packet is invalid")
	ErrClientChallengeResponseFailed = errors.New("client challenge response validation failed")

	headerInfos = map[string]PacketHeader{
		Server: {Type: ServerChallengePacketType, Size: HeaderChallengeSize},
		Client: {Type: ClientChallengePacketType, Size: HeaderChallengeSize},
	}
)

type AuthPacketHandler struct {
	clientKeyPair *crypto.KeyPair

	serverPublicKey *crypto.PublicKey
	securityMode    ConnectionSecurityMode // none (this is only mode currenty supported by upstream code)
	serverChallenge Challenge

	conn net.Conn
}

func NewAuthPacketHandler(
	clientKeyPair *crypto.KeyPair,
	serverPublicKey *crypto.PublicKey,
	securityMode ConnectionSecurityMode,
	conn net.Conn) *AuthPacketHandler {

	return &AuthPacketHandler{
		clientKeyPair:   clientKeyPair,
		serverPublicKey: serverPublicKey,
		securityMode:    securityMode,
		serverChallenge: Challenge{},
		conn:            conn,
	}
}

func (auth *AuthPacketHandler) Start() error {
	err := auth.HandleServerChallengeRequest()
	if err != nil {
		return err
	}

	v, err := auth.HandleClientChallengeResponse()
	if err != nil {
		return err
	}
	if !v {
		return ErrClientChallengeResponseFailed
	}

	return nil
}

func (auth *AuthPacketHandler) HandleServerChallengeRequest() error {
	buf, err := readFromConn(auth.conn, ServerChallengeRequestSize)
	if err != nil {
		return err
	}

	serverReq := &ServerChallengeRequest{}
	err = serverReq.Parse(buf)
	if err != nil {
		return err
	}

	response, err := GenerateServerChallengeResponse(serverReq, auth.clientKeyPair, auth.securityMode)
	if err != nil {
		return err
	}
	copy(auth.serverChallenge[:], response.Challenge[:])

	_, err = auth.conn.Write(response.Bytes())
	return err
}

func (auth *AuthPacketHandler) HandleServerChallengeResponse() (*ServerChallengeResponse, error) {
	buf := make([]byte, ServerChallengeResponseSize)
	_, err := auth.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	resp := &ServerChallengeResponse{}
	err = resp.Parse(buf)
	if err != nil {
		return nil, err
	}

	v := isPacketHeaderValid(Server, resp.PacketHeader)
	if !v {
		return nil, ErrPacketIsInvalid
	}

	return resp, nil
}

func (auth *AuthPacketHandler) HandleClientChallengeResponse() (bool, error) {
	buf, err := readFromConn(auth.conn, ClientChallengeResponseSize)
	if err != nil {
		return false, err
	}

	clResp := &ClientChallengeResponse{}
	err = clResp.Parse(buf)
	if err != nil {
		return false, err
	}

	return VerifyClientChallengeResponse(clResp, auth.serverPublicKey, auth.serverChallenge), nil
}

func isPacketHeaderValid(packetTypeName string, packet *PacketHeader) bool {
	inf, ok := headerInfos[packetTypeName]
	if !ok {
		return false
	}
	return inf.Type == packet.Type && inf.Size == packet.Size
}

func readFromConn(conn net.Conn, expectedSize uint64) ([]byte, error) {
	offset, n := 0, 0
	var err error
	buf := make([]byte, expectedSize)
	for offset < len(buf) {
		n, err = conn.Read(buf[offset:])
		if err != nil {
			return nil, err
		}
		offset += n
	}

	return buf, err
}
