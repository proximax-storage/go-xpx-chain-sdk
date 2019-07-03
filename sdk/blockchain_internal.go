package sdk

type blockInfoDTO struct {
	BlockMeta struct {
		BlockHash       hashDto   `json:"hash"`
		GenerationHash  hashDto   `json:"generationHash"`
		TotalFee        uint64DTO `json:"totalFee"`
		NumTransactions uint64    `json:"numTransactions"`
		// MerkleTree      uint64DTO `json:"merkleTree"` is needed?
	} `json:"meta"`
	Block struct {
		Signature              signatureDto           `json:"signature"`
		Signer                 string                 `json:"signer"`
		Version                uint64                 `json:"version"`
		Type                   uint64                 `json:"type"`
		Height                 uint64DTO              `json:"height"`
		Timestamp              blockchainTimestampDTO `json:"timestamp"`
		Difficulty             uint64DTO              `json:"difficulty"`
		FeeMultiplier          uint32                 `json:"feeMultiplier"`
		PreviousBlockHash      hashDto                `json:"previousBlockHash"`
		BlockTransactionsHash  hashDto                `json:"blockTransactionsHash"`
		BlockReceiptsHash      hashDto                `json:"blockReceiptsHash"`
		StateHash              hashDto                `json:"stateHash"`
		Beneficiary            string                 `json:"beneficiary"`
		FeeInterest            uint32                 `json:"feeInterest"`
		FeeInterestDenominator uint32                 `json:"feeInterestDenominator"`
	} `json:"block"`
}

func (dto *blockInfoDTO) toStruct() (*BlockInfo, error) {
	nt := ExtractNetworkType(dto.Block.Version)

	pa, err := NewAccountFromPublicKey(dto.Block.Signer, nt)
	if err != nil {
		return nil, err
	}

	v := ExtractVersion(dto.Block.Version)
	if err != nil {
		return nil, err
	}

	var bpa *PublicAccount = nil

	if dto.Block.Beneficiary != EmptyPublicKey {
		bpa, err = NewAccountFromPublicKey(dto.Block.Beneficiary, nt)
		if err != nil {
			return nil, err
		}
	}

	blockHash, err := dto.BlockMeta.BlockHash.Hash()
	if err != nil {
		return nil, err
	}

	generationHash, err := dto.BlockMeta.GenerationHash.Hash()
	if err != nil {
		return nil, err
	}

	signature, err := dto.Block.Signature.Signature()
	if err != nil {
		return nil, err
	}

	previousBlockHash, err := dto.Block.PreviousBlockHash.Hash()
	if err != nil {
		return nil, err
	}

	blockTransactionsHash, err := dto.Block.BlockTransactionsHash.Hash()
	if err != nil {
		return nil, err
	}

	blockReceiptsHash, err := dto.Block.BlockReceiptsHash.Hash()
	if err != nil {
		return nil, err
	}

	stateHash, err := dto.Block.StateHash.Hash()
	if err != nil {
		return nil, err
	}

	return &BlockInfo{
		NetworkType:            nt,
		BlockHash:              blockHash,
		GenerationHash:         generationHash,
		TotalFee:               dto.BlockMeta.TotalFee.toStruct(),
		NumTransactions:        dto.BlockMeta.NumTransactions,
		Signature:              signature,
		Signer:                 pa,
		Version:                v,
		Type:                   dto.Block.Type,
		Height:                 dto.Block.Height.toStruct(),
		Timestamp:              dto.Block.Timestamp.toStruct().ToTimestamp(),
		Difficulty:             dto.Block.Difficulty.toStruct(),
		FeeMultiplier:          dto.Block.FeeMultiplier,
		PreviousBlockHash:      previousBlockHash,
		BlockTransactionsHash:  blockTransactionsHash,
		BlockReceiptsHash:      blockReceiptsHash,
		StateHash:              stateHash,
		Beneficiary:            bpa,
		FeeInterest:            dto.Block.FeeInterest,
		FeeInterestDenominator: dto.Block.FeeInterestDenominator,
	}, nil
}

// Chain Score
type chainScoreDTO struct {
	ScoreHigh uint64DTO `json:"scoreHigh"`
	ScoreLow  uint64DTO `json:"scoreLow"`
}

func (dto *chainScoreDTO) toStruct() *ChainScore {
	return NewChainScore(dto.ScoreLow.toUint64(), dto.ScoreHigh.toUint64())
}

type blockInfoDTOs []*blockInfoDTO

func (b *blockInfoDTOs) toStruct() ([]*BlockInfo, error) {
	dtos := *b
	blocks := make([]*BlockInfo, 0, len(dtos))

	for _, dto := range dtos {
		block, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}
