package sdk

import (
	"encoding/hex"
	jsonLib "encoding/json"
	"errors"
)

type metadataNemInfoDTOs []*metadataNemInfoDTO

func (ref *metadataNemInfoDTOs) toStruct(networkType NetworkType) ([]*MetadataNemTupleInfo, error) {
	var (
		dtos  = *ref
		infos = make([]*MetadataNemTupleInfo, 0, len(dtos))
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

type metadataNemInfoDTO struct {
	Metadata struct {
		CompositeHash     hashDto            `json:"compositeHash"`
		TargetKey         hashDto            `json:"targetKey"`
		ScopedMetadataKey uint64DTO          `json:"scopedMetadataKey"`
		SourceAddress     string             `json:"sourceAddress"`
		MetadataType      MetadataNemType    `json:"metadataType"`
		Value             string             `json:"value"`
		TargetId          jsonLib.RawMessage `json:"targetId"`
	} `json:"metadataEntry"`
}

func (ref *metadataNemInfoDTO) toStruct(networkType NetworkType) (*MetadataNemTupleInfo, error) {
	metadataInfo := &MetadataNemTupleInfo{}
	var err error = nil

	commonMetadata := MetadataNemInfo{}
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

	if commonMetadata.Type == MetadataNemAddressType {
		metadataInfo.Address = &AddressMetadataNemInfo{
			MetadataNemInfo: commonMetadata,
			Address:         commonMetadata.SourceAddress,
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

		metadataInfo.Namespace = &NamespaceMetadataNemInfo{
			MetadataNemInfo: commonMetadata,
			NamespaceId:     namespaceId,
		}

		return metadataInfo, nil
	} else if assetId.Type() == MosaicAssetIdType {
		mosaicId, err := NewMosaicId(assetId.Id())
		if err != nil {
			return nil, err
		}

		metadataInfo.Mosaic = &MosaicMetadataNemInfo{
			MetadataNemInfo: commonMetadata,
			MosaicId:        mosaicId,
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
