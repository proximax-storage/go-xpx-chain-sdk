package sdk

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

type namespaceIdDTO uint64DTO

func (dto *namespaceIdDTO) toStruct() (*NamespaceId, error) {
	return NewNamespaceId(uint64DTO(*dto).toUint64())
}

type namespaceIdDTOs []*namespaceIdDTO

func (dto *namespaceIdDTOs) toStruct() ([]*NamespaceId, error) {
	ids := make([]*NamespaceId, len(*dto))
	var err error

	for i, n := range *dto {
		ids[i], err = n.toStruct()
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}

// namespaceNameDTO temporary struct for reading responce & fill NamespaceName
type namespaceNameDTO struct {
	NamespaceId namespaceIdDTO `json:"namespaceId"`
	FullName    string         `json:"name"`
}

func (ref *namespaceNameDTO) toStruct() (*NamespaceName, error) {
	nsId, err := ref.NamespaceId.toStruct()
	if err != nil {
		return nil, err
	}

	return &NamespaceName{
		nsId,
		ref.FullName,
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

// namespaceAliasDTO
type namespaceAliasDTO struct {
	MosaicId *mosaicIdDTO
	Address  string
	Type     AliasType
}

func (dto *namespaceAliasDTO) toStruct() (*NamespaceAlias, error) {
	alias := NamespaceAlias{}

	alias.Type = dto.Type

	switch alias.Type {
	case AddressAliasType:
		a, err := NewAddressFromBase32(dto.Address)
		if err != nil {
			return nil, err
		}

		alias.address = a
	case MosaicAliasType:
		mosaicId, err := dto.MosaicId.toStruct()
		if err != nil {
			return nil, err
		}
		alias.mosaicId = mosaicId
	}

	return &alias, nil
}

// namespaceDTO temporary struct for reading responce & fill NamespaceInfo
type namespaceDTO struct {
	Type         int
	Depth        int
	Level0       *namespaceIdDTO
	Level1       *namespaceIdDTO
	Level2       *namespaceIdDTO
	Alias        *namespaceAliasDTO
	ParentId     namespaceIdDTO
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

//toStruct create & return new NamespaceInfo from namespaceInfoDTO
func (ref *namespaceInfoDTO) toStruct() (*NamespaceInfo, error) {
	address, err := NewAddressFromBase32(ref.Namespace.OwnerAddress)
	if err != nil {
		return nil, err
	}

	pubAcc, err := NewAccountFromPublicKey(ref.Namespace.Owner, address.Type)
	if err != nil {
		return nil, err
	}

	parentId, err := ref.Namespace.ParentId.toStruct()
	if err != nil {
		return nil, err
	}

	levels, err := ref.extractLevels()
	if err != nil {
		return nil, err
	}

	if len(levels) == 0 {
		return nil, ErrNilNamespaceId
	}

	alias, err := ref.Namespace.Alias.toStruct()
	if err != nil {
		return nil, err
	}

	ns := &NamespaceInfo{
		NamespaceId: levels[len(levels)-1],
		Active:      ref.Meta.Active,
		TypeSpace:   NamespaceType(ref.Namespace.Type),
		Depth:       ref.Namespace.Depth,
		Levels:      levels,
		Alias:       alias,
		Owner:       pubAcc,
		StartHeight: ref.Namespace.StartHeight.toStruct(),
		EndHeight:   ref.Namespace.EndHeight.toStruct(),
	}

	if parentId != nil && parentId.Id() != 0 {
		ns.Parent = &NamespaceInfo{NamespaceId: parentId}
	}

	return ns, nil
}

func (ref *namespaceInfoDTO) extractLevels() ([]*NamespaceId, error) {
	levels := make([]*NamespaceId, 0)
	extractLevel := func(level *namespaceIdDTO) error {
		if level != nil {
			nsName, err := level.toStruct()
			if err != nil {
				return err
			}

			levels = append(levels, nsName)
		}
		return nil
	}

	err := extractLevel(ref.Namespace.Level0)
	if err != nil {
		return nil, err
	}

	err = extractLevel(ref.Namespace.Level1)
	if err != nil {
		return nil, err
	}

	err = extractLevel(ref.Namespace.Level2)
	if err != nil {
		return nil, err
	}

	return levels, nil
}

type namespaceInfoDTOs []*namespaceInfoDTO

func (n *namespaceInfoDTOs) toStruct() ([]*NamespaceInfo, error) {
	dtos := *n
	nsInfos := make([]*NamespaceInfo, 0, len(dtos))

	for _, nsInfoDTO := range dtos {
		nsInfo, err := nsInfoDTO.toStruct()
		if err != nil {
			return nil, err
		}

		nsInfos = append(nsInfos, nsInfo)
	}

	return nsInfos, nil
}

func generateNamespaceId(name string, parentId *NamespaceId) (*NamespaceId, error) {
	b := parentId.toLittleEndian()

	result := sha3.New256()

	if _, err := result.Write(b); err != nil {
		return nil, err
	}

	if _, err := result.Write([]byte(name)); err != nil {
		return nil, err
	}

	t := result.Sum(nil)

	return NewNamespaceId(binary.LittleEndian.Uint64(t) | NamespaceBit)
}
