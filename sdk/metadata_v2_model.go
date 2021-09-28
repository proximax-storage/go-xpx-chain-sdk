package sdk

type MetadataV2Info struct {
	CompositeHash *Hash
	SourceAddress *Address
	TargetKey     *Hash
	ScopedKey     ScopedMetadataKey
	Type          MetadataV2Type
	Value         []byte
}

type AddressMetadataV2Info struct {
	MetadataV2Info
	Address *Address
}

type MosaicMetadataV2Info struct {
	MetadataV2Info
	MosaicId *MosaicId
}

type NamespaceMetadataV2Info struct {
	MetadataV2Info
	NamespaceId *NamespaceId
}

type MetadataV2TupleInfo struct {
	Address   *AddressMetadataV2Info
	Mosaic    *MosaicMetadataV2Info
	Namespace *NamespaceMetadataV2Info
}

type SearchMetadataV2Options struct {
	SourceAddress *Address          `url:"sourceAddress,omitempty"`
	TargetKey     *Hash             `url:"targetKey,omitempty"`
	ScopedKey     ScopedMetadataKey `url:"scopedMetadataKey,omitempty"`
	TargetId      baseInt64         `url:"targetId,omitempty"`
	Type          MetadataV2Type    `url:"metadataType,omitempty"`
}
