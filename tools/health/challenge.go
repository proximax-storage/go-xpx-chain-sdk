package main

import (
	"math/rand"
	"time"

	crypto "github.com/proximax-storage/go-xpx-crypto"
)

func GenerateServerChallengeRequest() *ServerChallengeRequest {
	request := NewServerChallengeRequest()
	request.Challenge = generateRandomChallenge()
	return request
}

func GenerateServerChallengeResponse(
	request *ServerChallengeRequest,
	keyPair *crypto.KeyPair,
	securityMode ConnectionSecurityMode) (*ServerChallengeResponse, error) {

	response := NewServerChallengeResponse()
	response.Challenge = generateRandomChallenge()

	buff := make([]byte, 0, ChallengeSize+1)
	buff = append(buff, request.Challenge[:]...)
	buff = append(buff, byte(securityMode))

	var err error
	response.Signature, err = signChallenge(keyPair, buff)
	if err != nil {
		return nil, err
	}

	response.PublicKey = keyPair.PublicKey
	response.SecurityMode = securityMode
	return response, nil
}

func GenerateClientChallengeResponse(request *ServerChallengeResponse, keyPair *crypto.KeyPair) (*ClientChallengeResponse, error) {
	response := NewClientChallengeResponse()

	var err error
	response.Signature, err = signChallenge(keyPair, request.Challenge[:])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func VerifyServerChallengeResponse(response *ServerChallengeResponse, challenge Challenge) bool {
	buff := make([]byte, 0, len(challenge)+SecurityModeSize)
	buff = append(buff, challenge[:]...)
	buff = append(buff, byte(response.SecurityMode))
	return verifyChallenge(response.PublicKey, buff, response.Signature)
}

func VerifyClientChallengeResponse(response *ClientChallengeResponse, serverPublicKey *crypto.PublicKey, challenge Challenge) bool {
	return verifyChallenge(serverPublicKey, challenge[:], response.Signature)
}

func generateRandomChallenge() Challenge {
	ch := Challenge{}
	rand.Seed(time.Now().UnixNano())
	rand.Read(ch[:])

	return ch
}

func signChallenge(keyPair *crypto.KeyPair, buffers []byte) (*crypto.Signature, error) {
	s := crypto.NewSignerFromKeyPair(keyPair, nil)
	sig, err := s.Sign(buffers)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func verifyChallenge(publicKey *crypto.PublicKey, buffers []byte, signature *crypto.Signature) bool {
	// TODO: the next struct is reduced. Would be to be able to verify just by pubKey
	keyPair := &crypto.KeyPair{PublicKey: publicKey}
	s := crypto.NewSignerFromKeyPair(keyPair, nil)
	return s.Verify(buffers, signature)
}
