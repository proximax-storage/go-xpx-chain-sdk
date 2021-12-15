package sdk

import (
	"encoding/hex"
	jsonLib "encoding/json"
	"errors"
)

type metadataV2InfoDTOs []*metadataV2InfoDTO

func (ref *metadataV2InfoDTOs) toStruct(networkType NetworkType) ([]*MetadataV2TupleInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*MetadataV2TupleInfo, 0, len(dtos))
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

type metadataV2InfoDTO struct {
	Metadata struct {
		CompositeHash     hashDto            `json:"compositeHash"`
		TargetKey         hashDto            `json:"targetKey"`
		ScopedMetadataKey uint64DTO          `json:"scopedMetadataKey"`
		SourceAddress     string             `json:"sourceAddress"`
		MetadataType      MetadataV2Type     `json:"metadataType"`
		Value             string             `json:"value"`
		TargetId          jsonLib.RawMessage `json:"targetId"`
	} `json:"metadataEntry"`
}

func (ref *metadataV2InfoDTO) toStruct(networkType NetworkType) (*MetadataV2TupleInfo, error) {
	metadataInfo := &MetadataV2TupleInfo{}
	var err error = nil

	commonMetadata := MetadataV2Info{}
	commonMetadata.Value, err = hex.DecodeString(ref.Metadata.Value)
	if err != nil {
		return nil, err
	}

	commonMetadata.Type = ref.Metadata.MetadataType
	commonMetadata.CompositeHash, err = ref.Metadata.CompositeHash.Hash()
	if err != nil {
		return nil, err
	}

	commonMetadata.TargetKey, err = ref.Metadata.TargetKey.Hash()
	if err != nil {
		return nil, err
	}

	commonMetadata.SourceAddress, err = NewAddressFromBase32(ref.Metadata.SourceAddress)
	if err != nil {
		return nil, err
	}

	commonMetadata.ScopedKey = ref.Metadata.ScopedMetadataKey.toStruct()

	if commonMetadata.Type == MetadataV2AddressType {
		metadataInfo.Address = &AddressMetadataV2Info{
			MetadataV2Info: commonMetadata,
			Address:        commonMetadata.SourceAddress,
		}

		return metadataInfo, nil
	}

	var targetId assetIdDTO
	err = json.Unmarshal(ref.Metadata.TargetId, &targetId)
	if err != nil {
		return nil, err
	}

	assetId, err := targetId.toStruct()
	if err != nil {
		return nil, err
	}

	if assetId.Type() == NamespaceAssetIdType {
		namespaceId, err := NewNamespaceId(assetId.Id())
		if err != nil {
			return nil, err
		}

		metadataInfo.Namespace = &NamespaceMetadataV2Info{
			MetadataV2Info: commonMetadata,
			NamespaceId:    namespaceId,
		}

		return metadataInfo, nil
	} else if assetId.Type() == MosaicAssetIdType {
		mosaicId, err := NewMosaicId(assetId.Id())
		if err != nil {
			return nil, err
		}

		metadataInfo.Mosaic = &MosaicMetadataV2Info{
			MetadataV2Info: commonMetadata,
			MosaicId:       mosaicId,
		}

		return metadataInfo, nil
	}

	return nil, errors.New("unknown type of asset id")
}

type computedHashes struct {
	Hashes []*Hash `json:"compositeHashes"`
}

func (ref *computedHashes) MarshalJSON() ([]byte, error) {
	buf := []byte(`{"compositeHashes": [`)

	for i, nsId := range ref.Hashes {
		if i > 0 {
			buf = append(buf, ',')
		}

		buf = append(buf, []byte(`"`+nsId.String()+`"`)...)
	}

	buf = append(buf, ']', '}')

	return buf, nil
}

type metadatasPageDTO struct {
	Metadatas metadataV2InfoDTOs `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (m *metadatasPageDTO) toStruct(networkType NetworkType) (*MetadatasPage, error) {
	metadatas, err := m.Metadatas.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	page := &MetadatasPage{
		Metadatas: metadatas,
		Pagination: Pagination{
			TotalEntries: m.Pagination.TotalEntries,
			PageNumber:   m.Pagination.PageNumber,
			PageSize:     m.Pagination.PageSize,
			TotalPages:   m.Pagination.TotalPages,
		},
	}

	return page, nil
}
