package sdk

func newBlockInfoConverter(accountFactory AccountFactory) blockInfoConverter {
	return &blockInfoConverterImpl{
		accountFactory: accountFactory,
	}
}

type blockInfoConverter interface {
	Convert(*blockInfoDTO) (*BlockInfo, error)
	ConvertMulti(*blockInfoDTOs) ([]*BlockInfo, error)
}

type blockInfoConverterImpl struct {
	accountFactory AccountFactory
}

func (c *blockInfoConverterImpl) Convert(dto *blockInfoDTO) (*BlockInfo, error) {
	nt := ExtractNetworkType(dto.Block.Version)

	pa, err := c.accountFactory.NewAccountFromPublicKey(dto.Block.Signer, nt)
	if err != nil {
		return nil, err
	}

	v := ExtractVersion(dto.Block.Version)
	if err != nil {
		return nil, err
	}

	return &BlockInfo{
		NetworkType:           nt,
		Hash:                  dto.BlockMeta.Hash,
		GenerationHash:        dto.BlockMeta.GenerationHash,
		TotalFee:              dto.BlockMeta.TotalFee.toBigInt(),
		NumTransactions:       dto.BlockMeta.NumTransactions,
		Signature:             dto.Block.Signature,
		Signer:                pa,
		Version:               v,
		Type:                  dto.Block.Type,
		Height:                dto.Block.Height.toBigInt(),
		Timestamp:             dto.Block.Timestamp.toBigInt(),
		Difficulty:            dto.Block.Difficulty.toBigInt(),
		PreviousBlockHash:     dto.Block.PreviousBlockHash,
		BlockTransactionsHash: dto.Block.BlockTransactionsHash,
	}, nil
}

func (c *blockInfoConverterImpl) ConvertMulti(b *blockInfoDTOs) ([]*BlockInfo, error) {
	dtos := *b
	blocks := make([]*BlockInfo, 0, len(dtos))

	for _, dto := range dtos {
		block, err := c.Convert(dto)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}
