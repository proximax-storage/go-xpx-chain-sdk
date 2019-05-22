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
	Hash []byte
	Type HashType
}

func (s *Secret) String() string {
	return fmt.Sprintf(
		`"[ Type": %d, "Hash": %s ]`,
		s.Type,
		s.HashString(),
	)
}

func (s *Secret) HashString() string {
	return strings.ToUpper(hex.EncodeToString(s.Hash))
}

// returns Secret from passed hash and HashType
func NewSecret(hash []byte, hashType HashType) (*Secret, error) {
	l := len(hash)

	switch hashType {
	case SHA3_256, KECCAK_256, SHA_256:
		if l != 32 {
			return nil, errors.New("the length of Secret is wrong")
		}
	case HASH_160:
		if l != 20 && l != 32 {
			return nil, errors.New("the length of HASH_160 Secret is wrong")
		}
		if l == 20 {
			hash = append(hash, make([]byte, 12)...)
		}
	default:
		return nil, errors.New("Not supported HashType NewSecret")
	}

	secret := Secret{hash, hashType}
	return &secret, nil
}

// returns Secret from passed hex string hash and HashType
func NewSecretFromHexString(hash string, hashType HashType) (*Secret, error) {
	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return NewSecret(bytes, hashType)
}

type Proof struct {
	Data []byte
}

func (p *Proof) String() string {
	return fmt.Sprintf(
		`[ %s ]`,
		p.ProofString(),
	)
}

// bytes representation of Proof
func (p *Proof) ProofString() string {
	return strings.ToUpper(hex.EncodeToString(p.Data))
}

// bytes length of Proof
func (p *Proof) Size() int {
	return len(p.Data)
}

func NewProofFromBytes(proof []byte) *Proof {
	return &Proof{proof}
}

func NewProofFromString(proof string) *Proof {
	return &Proof{[]byte(proof)}
}

func NewProofFromHexString(hexProof string) (*Proof, error) {
	proofB, err := hex.DecodeString(hexProof)
	if err != nil {
		return nil, err
	}
	return &Proof{proofB}, nil
}

func NewProofFromUint8(number uint8) *Proof {
	numberB := []byte{number}
	return &Proof{numberB}
}

func NewProofFromUint16(number uint16) *Proof {
	numberB := make([]byte, 2)
	binary.LittleEndian.PutUint16(numberB, number)

	return &Proof{numberB}
}

func NewProofFromUint32(number uint32) *Proof {
	numberB := make([]byte, 4)
	binary.LittleEndian.PutUint32(numberB, number)

	return &Proof{numberB}
}

func NewProofFromUint64(number uint64) *Proof {
	numberB := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberB, number)

	return &Proof{numberB}
}

// returns Secret generated from Proof with passed HashType
func (p *Proof) Secret(hashType HashType) (*Secret, error) {
	secretB, err := generateSecret(p.Data, hashType)

	if err != nil {
		return nil, err
	}

	secret, err := NewSecret(secretB, hashType)

	if err != nil {
		return nil, err
	}

	return secret, nil
}

func generateSecret(proofB []byte, hashType HashType) ([]byte, error) {
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

	return nil, errors.New("Not supported HashType generateSecret")
}
