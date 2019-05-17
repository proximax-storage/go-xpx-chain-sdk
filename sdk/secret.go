package sdk

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/xpx-crypto-go"
	"strings"
)

type HashType uint8

func (ht HashType) String() string {
	return fmt.Sprintf("%d", ht)
}

const (
	/// Input is hashed using Sha-3-256.
	SHA3_256 HashType = iota
	/// Input is hashed using Keccak-256.
	KECCAK_256
	/// Input is hashed twice: first with SHA-256 and then with RIPEMD-160.
	HASH_160
	/// Input is hashed twice with SHA-256.
	SHA_256
)

type Secret struct {
	Hash string
	Type HashType
}

func (s *Secret) String() string {
	return fmt.Sprintf(
		`"[ Type": %d, "Hash": %s ]`,
		s.Type,
		s.Hash,
	)
}

func (s *Secret) HashBytes() ([]byte, error) {
	return hex.DecodeString(s.Hash)
}

// returns Secret from passed hash string and HashType
func NewSecret(hash string, hashType HashType) (*Secret, error) {
	l := len(hash)

	switch hashType {
	case SHA3_256, KECCAK_256, SHA_256:
		if l != 64 {
			return nil, errors.New("the length of Secret is wrong")
		}
	case HASH_160:
		if l != 40 && l != 64 {
			return nil, errors.New("the length of HASH_160 Secret is wrong")
		}
		if l == 40 {
			hash = hash + strings.Repeat("0", 24)
		}
	}

	secret := Secret{hash, hashType}
	return &secret, nil
}

type Proof struct {
	hexProof string
}

func (p *Proof) String() string {
	return fmt.Sprintf(
		`[ %s ]`,
		p.hexProof,
	)
}

func (p *Proof) Bytes() ([]byte, error) {
	return hex.DecodeString(p.hexProof)
}

func (p *Proof) Size() int {
	return len(p.hexProof) / 2
}

func NewProofFromBytes(proof []byte) *Proof {
	return &Proof{hex.EncodeToString(proof)}
}

func NewProofFromString(proof string) *Proof {
	return &Proof{hex.EncodeToString([]byte(proof))}
}

func NewProofFromHexString(hex string) *Proof {
	return &Proof{hex}
}

func NewProofFromUint8(number uint8) *Proof {
	numberB := []byte{number}
	return &Proof{hex.EncodeToString(numberB)}
}

func NewProofFromUint16(number uint16) *Proof {
	numberB := make([]byte, 2)
	binary.LittleEndian.PutUint16(numberB, number)

	return &Proof{hex.EncodeToString(numberB)}
}

func NewProofFromUint32(number uint32) *Proof {
	numberB := make([]byte, 4)
	binary.LittleEndian.PutUint32(numberB, number)

	return &Proof{hex.EncodeToString(numberB)}
}

func NewProofFromUint64(number uint64) *Proof {
	numberB := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberB, number)

	return &Proof{hex.EncodeToString(numberB)}
}

func (p *Proof) Secret(hashType HashType) (*Secret, error) {
	secretB, err := generateSecret(p.hexProof, hashType)

	if err != nil {
		return nil, err
	}

	secret, err := NewSecret(strings.ToUpper(hex.EncodeToString(secretB)), hashType)

	if err != nil {
		return nil, err
	}

	return secret, nil
}

func generateSecret(proof string, hashType HashType) ([]byte, error) {
	proofB, err := hex.DecodeString(proof)

	if err != nil {
		return nil, err
	}

	switch hashType {
	case SHA3_256:
		return crypto.HashesSha3_256(proofB)
	case KECCAK_256:
		return crypto.HashesKeccak_256(proofB)
	case HASH_160:
		secretFirstB, err := crypto.HashesSha_256(proofB)

		if err != nil {
			return nil, err
		}

		return crypto.HashesRipemd160(secretFirstB)
	case SHA_256:
		secretFirstB, err := crypto.HashesSha_256(proofB)

		if err != nil {
			return nil, err
		}

		return crypto.HashesSha_256(secretFirstB)
	}

	return nil, errors.New("Not supported HashType")
}
