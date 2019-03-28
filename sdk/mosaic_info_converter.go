package sdk

import "errors"

func newMosaicInfoConverter(factory AccountFactory) mosaicInfoConverter {
	return &MosaicInfoConverterImpl{
		accountFactory: factory,
	}
}

type mosaicInfoConverter interface {
	Convert(*mosaicInfoDTO, NetworkType) (*MosaicInfo, error)
	ConvertMulti(mosaicInfoDTOs, NetworkType) ([]*MosaicInfo, error)
}

type MosaicInfoConverterImpl struct {
	accountFactory AccountFactory
}

func (c *MosaicInfoConverterImpl) Convert(dto *mosaicInfoDTO, networkType NetworkType) (*MosaicInfo, error) {
	publicAcc, err := c.accountFactory.NewAccountFromPublicKey(dto.Mosaic.Owner, networkType)
	if err != nil {
		return nil, err
	}

	if len(dto.Mosaic.Properties) < 3 {
		return nil, errors.New("mosaic Properties is not valid")
	}

	mosaicId, err := NewMosaicId(dto.Mosaic.MosaicId.toBigInt())

	mscInfo := &MosaicInfo{
		MosaicId:   mosaicId,
		Supply:     dto.Mosaic.Supply.toBigInt(),
		Height:     dto.Mosaic.Height.toBigInt(),
		Owner:      publicAcc,
		Revision:   dto.Mosaic.Revision,
		Properties: dto.Mosaic.Properties.toStruct(),
	}

	return mscInfo, nil
}

func (c *MosaicInfoConverterImpl) ConvertMulti(dtos mosaicInfoDTOs, networkType NetworkType) ([]*MosaicInfo, error) {
	mscInfos := make([]*MosaicInfo, 0, len(dtos))

	for _, dto := range dtos {
		mscInfo, err := c.Convert(dto, networkType)
		if err != nil {
			return nil, err
		}

		mscInfos = append(mscInfos, mscInfo)
	}

	return mscInfos, nil
}
