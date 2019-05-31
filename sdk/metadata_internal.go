package sdk

type addressMetadataInfoDTOs []*addressMetadataInfoDTO
type mosaicMetadataInfoDTOs []*mosaicMetadataInfoDTO
type namespaceMetadataInfoDTOs []*namespaceMetadataInfoDTO

func (ref *addressMetadataInfoDTOs) toStruct(networkType NetworkType) ([]*AddressMetadataInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*AddressMetadataInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (ref *mosaicMetadataInfoDTOs) toStruct(networkType NetworkType) ([]*MosaicMetadataInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*MosaicMetadataInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (ref *namespaceMetadataInfoDTOs) toStruct(networkType NetworkType) ([]*NamespaceMetadataInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*NamespaceMetadataInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

type metadataFieldDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type metadataInfoDTO struct {
	MetadataType MetadataType        `json:"metadataType"`
	Fields       []*metadataFieldDTO `json:"fields"`
}

type addressMetadataInfoDTO struct {
	Metadata struct {
		metadataInfoDTO
		Address string `json:"metadataId"`
	}
}

type mosaicMetadataInfoDTO struct {
	Metadata struct {
		metadataInfoDTO
		MosaicId uint64DTO `json:"metadataId"`
	}
}

type namespaceMetadataInfoDTO struct {
	Metadata struct {
		metadataInfoDTO
		NamespaceId uint64DTO `json:"metadataId"`
	}
}

func (ref *metadataInfoDTO) toStruct(networkType NetworkType) (*MetadataInfo, error) {
	metadataInfo := MetadataInfo{
		ref.MetadataType,
		make(map[string]string),
	}

	for _, f := range ref.Fields {
		metadataInfo.Fields[f.Key] = f.Value
	}

	return &metadataInfo, nil
}

func (ref *addressMetadataInfoDTO) toStruct(networkType NetworkType) (*AddressMetadataInfo, error) {
	metadata := ref.Metadata

	metadataInfo, err := metadata.toStruct(networkType)

	if err != nil {
		return nil, err
	}

	var a *Address = nil

	if len(metadata.Address) != 0 {
		a, err = NewAddressFromBase32(metadata.Address)
		if err != nil {
			return nil, err
		}
	}

	return &AddressMetadataInfo{
		MetadataInfo: *metadataInfo,
		Address:      a,
	}, nil
}

func (ref *mosaicMetadataInfoDTO) toStruct(networkType NetworkType) (*MosaicMetadataInfo, error) {
	metadata := ref.Metadata

	metadataInfo, err := metadata.toStruct(networkType)

	if err != nil {
		return nil, err
	}

	mosaicId, err := NewMosaicId(metadata.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &MosaicMetadataInfo{
		MetadataInfo: *metadataInfo,
		MosaicId:     mosaicId,
	}, nil
}

func (ref *namespaceMetadataInfoDTO) toStruct(networkType NetworkType) (*NamespaceMetadataInfo, error) {
	metadata := ref.Metadata

	metadataInfo, err := metadata.toStruct(networkType)

	if err != nil {
		return nil, err
	}

	namespaceId, err := NewNamespaceId(metadata.NamespaceId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &NamespaceMetadataInfo{
		MetadataInfo: *metadataInfo,
		NamespaceId:  namespaceId,
	}, nil
}
