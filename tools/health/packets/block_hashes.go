package packets

import (
	"encoding/binary"
	"errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	BlockHashesRequest struct {
		PacketHeader

		/// Requested number of hashes.
		NumHashes uint32
		/// Requested block height.
		Height uint64
	}

	BlockHashesResponse struct {
		PacketHeader

		/// Requested block height.
		Hashes []*sdk.Hash
	}
)

var ErrBadHashesLength = errors.New("bad hashes length")

func NewBlockHashesRequest(height uint64, numHashes uint32) *BlockHashesRequest {
	ph := NewPacketHeader(BlockHashesPacketType)
	ph.Size = BlockHashesRequestSize
	return &BlockHashesRequest{
		PacketHeader: ph,
		NumHashes:    numHashes,
		Height:       height,
	}
}

func (bh *BlockHashesRequest) Bytes() []byte {
	buff := make([]byte, BlockHashesRequestSize)

	// copy header
	copy(buff[:PacketHeaderSize], bh.PacketHeader.Bytes())

	binary.LittleEndian.PutUint32(buff[PacketHeaderSize:], bh.NumHashes)
	binary.LittleEndian.PutUint64(buff[PacketHeaderSize+4:], bh.Height)

	return buff
}

func (bh *BlockHashesResponse) Parse(buff []byte) error {
	if len(buff)%HashSize != 0 {
		return ErrBadHashesLength
	}

	bh.Hashes = make([]*sdk.Hash, 0, len(buff)/HashSize)
	for i := 0; i < len(buff); i += HashSize {
		h := &sdk.Hash{}
		copy(h[:], buff[i:i+HashSize])
		bh.Hashes = append(bh.Hashes, h)
	}

	return nil
}
