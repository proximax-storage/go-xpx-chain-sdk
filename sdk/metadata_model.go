package sdk

type MetadataInfo struct {
	MetadataType MetadataType
	Fields       map[string]string
}

type AddressMetadataInfo struct {
	MetadataInfo
	Address *Address
}

type MosaicMetadataInfo struct {
	MetadataInfo
	MosaicId *MosaicId
}

type NamespaceMetadataInfo struct {
	MetadataInfo
	NamespaceId *NamespaceId
}
