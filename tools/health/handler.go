package health

import (
	"errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/tools/health/packets"
	crypto "github.com/proximax-storage/go-xpx-crypto"
)

var ErrClientChallengeResponseFailed = errors.New("client challenge response validation failed")

type NodeIo interface {
	Write(packets.Byter) (int, error)
	Read(packets.Parser, int) error
	Close() error
}

type Handler struct {
	nodeIo NodeIo
}

func NewHandler(io NodeIo) *Handler {
	return &Handler{nodeIo: io}
}

func (h *Handler) Close() error {
	return h.nodeIo.Close()
}

func (h *Handler) AuthHandle(clientKeyPair *crypto.KeyPair, serverPublicKey *crypto.PublicKey, securityMode packets.ConnectionSecurityMode) error {
	serverReq := &packets.ServerChallengeRequest{}
	err := h.nodeIo.Read(serverReq, packets.ServerChallengeRequestSize)
	if err != nil {
		return err
	}

	response, err := packets.GenerateServerChallengeResponse(serverReq, clientKeyPair, securityMode)
	if err != nil {
		return err
	}

	serverChallenge := packets.Challenge{}
	copy(serverChallenge[:], response.Challenge[:])

	_, err = h.nodeIo.Write(response)
	if err != nil {
		return err
	}

	clResp := &packets.ClientChallengeResponse{}
	err = h.nodeIo.Read(clResp, packets.ClientChallengeResponseSize)
	if err != nil {
		return err
	}

	if !packets.VerifyClientChallengeResponse(clResp, serverPublicKey, serverChallenge) {
		return ErrClientChallengeResponseFailed
	}

	return nil
}

func (h *Handler) CommonHandle(req packets.Byter, resp packets.ResponsePacket) error {
	_, err := h.nodeIo.Write(req)
	if err != nil {
		return err
	}

	err = h.nodeIo.Read(resp.Header(), packets.PacketHeaderSize)
	if err != nil {
		return errors.New("cannot read packet header: " + err.Error())
	}

	err = h.nodeIo.Read(resp, int(resp.Header().PacketSize()-packets.PacketHeaderSize))
	if err != nil {
		return errors.New("cannot read packet payload: " + err.Error())
	}

	return nil
}
