package sdk

type MetadataNemInfo struct {
	CompositeHash *Hash
	SourceAddress *Address
	TargetKey     *Hash
	ScopedKey     ScopedMetadataKey
	Type          MetadataNemType
	Value         []byte
}

type AddressMetadataNemInfo struct {
	MetadataNemInfo
	Address *Address
}

type MosaicMetadataNemInfo struct {
	MetadataNemInfo
	MosaicId *MosaicId
}

type NamespaceMetadataNemInfo struct {
	MetadataNemInfo
	NamespaceId *NamespaceId
}

type MetadataNemTupleInfo struct {
	Address   *AddressMetadataNemInfo
	Mosaic    *MosaicMetadataNemInfo
	Namespace *NamespaceMetadataNemInfo
}
