package sdk

import (
	"encoding/binary"
	"github.com/proximax-storage/go-xpx-catapult-sdk/utils"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func bigIntToNamespaceId(bigInt *big.Int) *NamespaceId {
	if bigInt == nil {
		return nil
	}

	nsId := NamespaceId(*bigInt)

	return &nsId
}

func namespaceIdToBigInt(nsId *NamespaceId) *big.Int {
	if nsId == nil {
		return nil
	}

	return (*big.Int)(nsId)
}

type namespaceIdDTO uint64DTO

func (dto *namespaceIdDTO) toStruct() (*NamespaceId, error) {
	return NewNamespaceId(uint64DTO(*dto).toBigInt())
}

// namespaceNameDTO temporary struct for reading responce & fill NamespaceName
type namespaceNameDTO struct {
	NamespaceId uint64DTO
	Name        string
	ParentId    uint64DTO
}

func (ref *namespaceNameDTO) toStruct() (*NamespaceName, error) {
	nsId, err := NewNamespaceId(ref.NamespaceId.toBigInt())
	if err != nil {
		return nil, err
	}

	parentId, err := NewNamespaceId(ref.ParentId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &NamespaceName{
		nsId,
		ref.Name,
		parentId,
	}, nil
}

type namespaceNameDTOs []*namespaceNameDTO

func (n *namespaceNameDTOs) toStruct() ([]*NamespaceName, error) {
	dtos := *n
	nsNames := make([]*NamespaceName, 0, len(dtos))

	for _, dto := range dtos {
		nsName, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		nsNames = append(nsNames, nsName)
	}

	return nsNames, nil
}

// namespaceDTO temporary struct for reading responce & fill NamespaceInfo
type namespaceDTO struct {
	NamespaceId  uint64DTO
	FullName     string
	Type         int
	Depth        int
	Level0       *uint64DTO
	Level1       *uint64DTO
	Level2       *uint64DTO
	ParentId     uint64DTO
	MosaicIds    []uint64DTO
	Owner        string
	OwnerAddress string
	StartHeight  uint64DTO
	EndHeight    uint64DTO
}

// namespaceInfoDTO temporary struct for reading response & fill NamespaceInfo
type namespaceInfoDTO struct {
	Meta      namespaceMosaicMetaDTO
	Namespace namespaceDTO
}

func (ref *namespaceInfoDTO) extractLevels() ([]*NamespaceId, error) {
	levels := make([]*NamespaceId, 0)

	if ref.Namespace.Level0 != nil {
		nsName, err := NewNamespaceId(ref.Namespace.Level0.toBigInt())
		if err != nil {
			return nil, err
		}

		levels = append(levels, nsName)
	}

	if ref.Namespace.Level1 != nil {
		nsName, err := NewNamespaceId(ref.Namespace.Level1.toBigInt())
		if err != nil {
			return nil, err
		}

		levels = append(levels, nsName)
	}

	if ref.Namespace.Level2 != nil {
		nsName, err := NewNamespaceId(ref.Namespace.Level2.toBigInt())
		if err != nil {
			return nil, err
		}

		levels = append(levels, nsName)
	}

	return levels, nil
}

type namespaceInfoDTOs []*namespaceInfoDTO

func generateNamespaceId(name string, parentId *big.Int) (*big.Int, error) {
	b := make([]byte, 8)

	if parentId.Int64() != 0 {
		b = parentId.Bytes()
	}

	utils.ReverseByteArray(b)

	result := sha3.New256()

	if _, err := result.Write(b); err != nil {
		return nil, err
	}

	if _, err := result.Write([]byte(name)); err != nil {
		return nil, err
	}

	t := result.Sum(nil)

	return uint64DTO{binary.LittleEndian.Uint32(t[0:4]), binary.LittleEndian.Uint32(t[4:8]) | 0x80000000}.toBigInt(), nil
}
