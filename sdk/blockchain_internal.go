package sdk

type blockInfoDTO struct {
	BlockMeta struct {
		Hash            string    `json:"hash"`
		GenerationHash  string    `json:"generationHash"`
		TotalFee        uint64DTO `json:"totalFee"`
		NumTransactions uint64    `json:"numTransactions"`
		// MerkleTree      uint64DTO `json:"merkleTree"` is needed?
	} `json:"meta"`
	Block struct {
		Signature             string                 `json:"signature"`
		Signer                string                 `json:"signer"`
		Version               uint64                 `json:"version"`
		Type                  uint64                 `json:"type"`
		Height                uint64DTO              `json:"height"`
		Timestamp             blockchainTimestampDTO `json:"timestamp"`
		Difficulty            uint64DTO              `json:"difficulty"`
		FeeMultiplier         uint32                 `json:"feeMultiplier"`
		PreviousBlockHash     string                 `json:"previousBlockHash"`
		BlockTransactionsHash string                 `json:"blockTransactionsHash"`
		BlockReceiptsHash     string                 `json:"blockReceiptsHash"`
		StateHash             string                 `json:"stateHash"`
		BeneficiaryPublicKey  string                 `json:"beneficiaryPublicKey"`
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

	if dto.Block.BeneficiaryPublicKey != EmptyPublicKey {
		bpa, err = NewAccountFromPublicKey(dto.Block.BeneficiaryPublicKey, nt)
		if err != nil {
			return nil, err
		}
	}

	return &BlockInfo{
		NetworkType:           nt,
		Hash:                  dto.BlockMeta.Hash,
		GenerationHash:        dto.BlockMeta.GenerationHash,
		TotalFee:              dto.BlockMeta.TotalFee.toStruct(),
		NumTransactions:       dto.BlockMeta.NumTransactions,
		Signature:             dto.Block.Signature,
		Signer:                pa,
		Version:               v,
		Type:                  dto.Block.Type,
		Height:                dto.Block.Height.toStruct(),
		Timestamp:             dto.Block.Timestamp.toStruct().ToTimestamp(),
		Difficulty:            dto.Block.Difficulty.toStruct(),
		FeeMultiplier:         dto.Block.FeeMultiplier,
		PreviousBlockHash:     dto.Block.PreviousBlockHash,
		BlockTransactionsHash: dto.Block.BlockTransactionsHash,
		BlockReceiptsHash:     dto.Block.BlockReceiptsHash,
		StateHash:             dto.Block.StateHash,
		Beneficiary:           bpa,
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
