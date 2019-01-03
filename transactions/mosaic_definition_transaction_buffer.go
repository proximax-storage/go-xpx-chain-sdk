// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
// File is Auto-Generated

package transactions

import (
	"github.com/google/flatbuffers/go"
)

type MosaicDefinitionTransactionBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsMosaicDefinitionTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *MosaicDefinitionTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MosaicDefinitionTransactionBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *MosaicDefinitionTransactionBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MosaicDefinitionTransactionBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MosaicDefinitionTransactionBuffer) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Signature(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) SignatureLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) SignatureBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MosaicDefinitionTransactionBuffer) Signer(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) SignerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) SignerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MosaicDefinitionTransactionBuffer) Version() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateVersion(n uint16) bool {
	return rcv._tab.MutateUint16Slot(10, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Fee(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) FeeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) Deadline(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) DeadlineLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) ParentId(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) ParentIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MosaicId(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MosaicIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MosaicNameLength() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateMosaicNameLength(n byte) bool {
	return rcv._tab.MutateByteSlot(22, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) NumOptionalProperties() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateNumOptionalProperties(n byte) bool {
	return rcv._tab.MutateByteSlot(24, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Flags() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateFlags(n byte) bool {
	return rcv._tab.MutateByteSlot(26, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Divisibility() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateDivisibility(n byte) bool {
	return rcv._tab.MutateByteSlot(28, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) MosaicName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MosaicDefinitionTransactionBuffer) IndicateDuration() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) MutateIndicateDuration(n byte) bool {
	return rcv._tab.MutateByteSlot(32, n)
}

func (rcv *MosaicDefinitionTransactionBuffer) Duration(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MosaicDefinitionTransactionBuffer) DurationLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func MosaicDefinitionTransactionBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(16)
}
func MosaicDefinitionTransactionBufferAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func MosaicDefinitionTransactionBufferAddSignature(builder *flatbuffers.Builder, signature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(signature), 0)
}
func MosaicDefinitionTransactionBufferStartSignatureVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MosaicDefinitionTransactionBufferAddSigner(builder *flatbuffers.Builder, signer flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(signer), 0)
}
func MosaicDefinitionTransactionBufferStartSignerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MosaicDefinitionTransactionBufferAddVersion(builder *flatbuffers.Builder, version uint16) {
	builder.PrependUint16Slot(3, version, 0)
}
func MosaicDefinitionTransactionBufferAddType(builder *flatbuffers.Builder, type_ uint16) {
	builder.PrependUint16Slot(4, type_, 0)
}
func MosaicDefinitionTransactionBufferAddFee(builder *flatbuffers.Builder, fee flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(fee), 0)
}
func MosaicDefinitionTransactionBufferStartFeeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MosaicDefinitionTransactionBufferAddDeadline(builder *flatbuffers.Builder, deadline flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(deadline), 0)
}
func MosaicDefinitionTransactionBufferStartDeadlineVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MosaicDefinitionTransactionBufferAddParentId(builder *flatbuffers.Builder, parentId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(parentId), 0)
}
func MosaicDefinitionTransactionBufferStartParentIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MosaicDefinitionTransactionBufferAddMosaicId(builder *flatbuffers.Builder, mosaicId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(mosaicId), 0)
}
func MosaicDefinitionTransactionBufferStartMosaicIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MosaicDefinitionTransactionBufferAddMosaicNameLength(builder *flatbuffers.Builder, mosaicNameLength byte) {
	builder.PrependByteSlot(9, mosaicNameLength, 0)
}
func MosaicDefinitionTransactionBufferAddNumOptionalProperties(builder *flatbuffers.Builder, numOptionalProperties byte) {
	builder.PrependByteSlot(10, numOptionalProperties, 0)
}
func MosaicDefinitionTransactionBufferAddFlags(builder *flatbuffers.Builder, flags byte) {
	builder.PrependByteSlot(11, flags, 0)
}
func MosaicDefinitionTransactionBufferAddDivisibility(builder *flatbuffers.Builder, divisibility byte) {
	builder.PrependByteSlot(12, divisibility, 0)
}
func MosaicDefinitionTransactionBufferAddMosaicName(builder *flatbuffers.Builder, mosaicName flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(13, flatbuffers.UOffsetT(mosaicName), 0)
}
func MosaicDefinitionTransactionBufferAddIndicateDuration(builder *flatbuffers.Builder, indicateDuration byte) {
	builder.PrependByteSlot(14, indicateDuration, 0)
}
func MosaicDefinitionTransactionBufferAddDuration(builder *flatbuffers.Builder, duration flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(15, flatbuffers.UOffsetT(duration), 0)
}
func MosaicDefinitionTransactionBufferStartDurationVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MosaicDefinitionTransactionBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
