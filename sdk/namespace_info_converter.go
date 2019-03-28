package sdk

func newNamespaceInfoConverterImpl(accountFactory AccountFactory) namespaceInfoConverter {
	return &NamespaceInfoConverterImpl{accountFactory: accountFactory}
}

type namespaceInfoConverter interface {
	Convert(*namespaceInfoDTO) (*NamespaceInfo, error)
	ConvertMulti(namespaceInfoDTOs) ([]*NamespaceInfo, error)
}

type NamespaceInfoConverterImpl struct {
	accountFactory AccountFactory
}

func (c *NamespaceInfoConverterImpl) Convert(ref *namespaceInfoDTO) (*NamespaceInfo, error) {
	nsId, err := NewNamespaceId(ref.Namespace.NamespaceId.toBigInt())
	if err != nil {
		return nil, err
	}

	pubAcc, err := c.accountFactory.NewAccountFromPublicKey(ref.Namespace.Owner, NetworkType(ref.Namespace.Type))
	if err != nil {
		return nil, err
	}

	parentId, err := NewNamespaceId(ref.Namespace.ParentId.toBigInt())
	if err != nil {
		return nil, err
	}

	levels, err := ref.extractLevels()
	if err != nil {
		return nil, err
	}

	mscIds := make([]*MosaicId, 0, len(ref.Namespace.MosaicIds))

	for _, mscIdDTO := range ref.Namespace.MosaicIds {
		mscId, err := NewMosaicId(mscIdDTO.toBigInt())
		if err != nil {
			return nil, err
		}

		mscIds = append(mscIds, mscId)
	}

	ns := &NamespaceInfo{
		NamespaceId: nsId,
		FullName:    ref.Namespace.FullName,
		Active:      ref.Meta.Active,
		Index:       ref.Meta.Index,
		MetaId:      ref.Meta.Id,
		TypeSpace:   NamespaceType(ref.Namespace.Type),
		Depth:       ref.Namespace.Depth,
		Levels:      levels,
		Owner:       pubAcc,
		StartHeight: ref.Namespace.StartHeight.toBigInt(),
		EndHeight:   ref.Namespace.EndHeight.toBigInt(),
	}

	if parentId != nil && namespaceIdToBigInt(parentId).Int64() != 0 {
		ns.Parent = &NamespaceInfo{NamespaceId: parentId}
	}

	return ns, nil
}

func (c *NamespaceInfoConverterImpl) ConvertMulti(dtos namespaceInfoDTOs) ([]*NamespaceInfo, error) {
	nsInfos := make([]*NamespaceInfo, 0, len(dtos))

	for _, nsInfoDTO := range dtos {
		nsInfo, err := c.Convert(nsInfoDTO)
		if err != nil {
			return nil, err
		}

		nsInfos = append(nsInfos, nsInfo)
	}

	return nsInfos, nil
}
