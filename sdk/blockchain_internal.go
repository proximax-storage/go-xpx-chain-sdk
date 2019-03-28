package sdk

import "math/big"

type blockInfoDTO struct {
	BlockMeta struct {
		Hash            string    `json:"hash"`
		GenerationHash  string    `json:"generationHash"`
		TotalFee        uint64DTO `json:"totalFee"`
		NumTransactions uint64    `json:"numTransactions"`
		// MerkleTree      uint64DTO `json:"merkleTree"` is needed?
	} `json:"meta"`
	Block struct {
		Signature             string    `json:"signature"`
		Signer                string    `json:"signer"`
		Version               uint64    `json:"version"`
		Type                  uint64    `json:"type"`
		Height                uint64DTO `json:"height"`
		Timestamp             uint64DTO `json:"timestamp"`
		Difficulty            uint64DTO `json:"difficulty"`
		PreviousBlockHash     string    `json:"previousBlockHash"`
		BlockTransactionsHash string    `json:"blockTransactionsHash"`
	} `json:"block"`
}

// Chain Score
type chainScoreDTO struct {
	ScoreHigh uint64DTO `json:"scoreHigh"`
	ScoreLow  uint64DTO `json:"scoreLow"`
}

func (dto *chainScoreDTO) toStruct() *big.Int {
	return uint64DTO{uint32(dto.ScoreLow.toBigInt().Uint64()), uint32(dto.ScoreHigh.toBigInt().Uint64())}.toBigInt()
}

type blockInfoDTOs []*blockInfoDTO
