package packets

import (
	"encoding/binary"
)

type ChainInfoResponse struct {
	PacketHeader

	// Chain height.
	Height uint64
	// High part of the score.
	ScoreHigh uint64
	// Low part of the score.
	ScoreLow uint64
}

func (cir *ChainInfoResponse) Header() Header {
	return &cir.PacketHeader
}

func (cir *ChainInfoResponse) Parse(buff []byte) error {
	offset := 0
	cir.Height = binary.LittleEndian.Uint64(buff[offset : offset+8])
	offset += 8
	cir.ScoreHigh = binary.LittleEndian.Uint64(buff[offset : offset+8])
	offset += 8
	cir.ScoreLow = binary.LittleEndian.Uint64(buff[offset : offset+8])

	return nil
}
