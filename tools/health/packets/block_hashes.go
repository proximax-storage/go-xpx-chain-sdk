package packets

import (
	"encoding/binary"
	"errors"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

type (
	BlockHashesRequest struct {
		PacketHeader

		/// Requested block height.
		Height uint64
		/// Requested number of hashes.
		NumHashes uint32
	}

	BlockHashesResponse struct {
		PacketHeader

		/// Requested block height.
		Hashes []sdk.Hash
	}
)

var ErrBadHashesLength = errors.New("bad hashes length")

func NewBlockHashesRequest(height uint64, numHashes uint32) *BlockHashesRequest {
	ph := NewPacketHeader(BlockHashesPacketType)
	ph.Size = BlockHashesRequestSize
	return &BlockHashesRequest{
		PacketHeader: ph,
		Height:       height,
		NumHashes:    numHashes,
	}
}

func (bh *BlockHashesRequest) Bytes() []byte {
	buff := make([]byte, BlockHashesRequestSize)

	// copy header
	copy(buff[:PacketHeaderSize], bh.PacketHeader.Bytes())

	offset := PacketHeaderSize
	binary.LittleEndian.PutUint64(buff[offset:offset+8], bh.Height)

	offset += 8
	binary.LittleEndian.PutUint32(buff[offset:], bh.NumHashes)

	return buff
}

func (bh *BlockHashesResponse) Header() Header {
	return &bh.PacketHeader
}

func (bh *BlockHashesResponse) Parse(buff []byte) error {
	if len(buff)%HashSize != 0 {
		return ErrBadHashesLength
	}

	bh.Hashes = make([]sdk.Hash, 0, len(buff)/HashSize)
	for i := 0; i < len(buff); i += HashSize {
		h := sdk.Hash{}
		copy(h[:], buff[i:i+HashSize])
		bh.Hashes = append(bh.Hashes, h)
	}

	return nil
}
