// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package transactions

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type AddExchangeOfferBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsAddExchangeOfferBuffer(buf []byte, offset flatbuffers.UOffsetT) *AddExchangeOfferBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AddExchangeOfferBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AddExchangeOfferBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AddExchangeOfferBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *AddExchangeOfferBuffer) MosaicId(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MosaicIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MutateMosaicId(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *AddExchangeOfferBuffer) MosaicAmount(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MosaicAmountLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MutateMosaicAmount(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *AddExchangeOfferBuffer) Cost(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) CostLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MutateCost(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *AddExchangeOfferBuffer) Type() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MutateType(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func (rcv *AddExchangeOfferBuffer) Duration(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) DurationLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferBuffer) MutateDuration(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func AddExchangeOfferBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func AddExchangeOfferBufferAddMosaicId(builder *flatbuffers.Builder, mosaicId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(mosaicId), 0)
}
func AddExchangeOfferBufferStartMosaicIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferBufferAddMosaicAmount(builder *flatbuffers.Builder, mosaicAmount flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(mosaicAmount), 0)
}
func AddExchangeOfferBufferStartMosaicAmountVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferBufferAddCost(builder *flatbuffers.Builder, cost flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(cost), 0)
}
func AddExchangeOfferBufferStartCostVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferBufferAddType(builder *flatbuffers.Builder, type_ byte) {
	builder.PrependByteSlot(3, type_, 0)
}
func AddExchangeOfferBufferAddDuration(builder *flatbuffers.Builder, duration flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(duration), 0)
}
func AddExchangeOfferBufferStartDurationVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type AddExchangeOfferTransactionBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsAddExchangeOfferTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *AddExchangeOfferTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &AddExchangeOfferTransactionBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *AddExchangeOfferTransactionBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *AddExchangeOfferTransactionBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *AddExchangeOfferTransactionBuffer) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *AddExchangeOfferTransactionBuffer) Signature(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) SignatureLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) SignatureBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateSignature(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *AddExchangeOfferTransactionBuffer) Signer(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) SignerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) SignerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateSigner(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *AddExchangeOfferTransactionBuffer) Version() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateVersion(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *AddExchangeOfferTransactionBuffer) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *AddExchangeOfferTransactionBuffer) MaxFee(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MaxFeeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateMaxFee(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *AddExchangeOfferTransactionBuffer) Deadline(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) DeadlineLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateDeadline(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *AddExchangeOfferTransactionBuffer) OffersCount() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *AddExchangeOfferTransactionBuffer) MutateOffersCount(n byte) bool {
	return rcv._tab.MutateByteSlot(18, n)
}

func (rcv *AddExchangeOfferTransactionBuffer) Offers(obj *AddExchangeOfferBuffer, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *AddExchangeOfferTransactionBuffer) OffersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func AddExchangeOfferTransactionBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(9)
}
func AddExchangeOfferTransactionBufferAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func AddExchangeOfferTransactionBufferAddSignature(builder *flatbuffers.Builder, signature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(signature), 0)
}
func AddExchangeOfferTransactionBufferStartSignatureVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func AddExchangeOfferTransactionBufferAddSigner(builder *flatbuffers.Builder, signer flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(signer), 0)
}
func AddExchangeOfferTransactionBufferStartSignerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func AddExchangeOfferTransactionBufferAddVersion(builder *flatbuffers.Builder, version uint32) {
	builder.PrependUint32Slot(3, version, 0)
}
func AddExchangeOfferTransactionBufferAddType(builder *flatbuffers.Builder, type_ uint16) {
	builder.PrependUint16Slot(4, type_, 0)
}
func AddExchangeOfferTransactionBufferAddMaxFee(builder *flatbuffers.Builder, maxFee flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(maxFee), 0)
}
func AddExchangeOfferTransactionBufferStartMaxFeeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferTransactionBufferAddDeadline(builder *flatbuffers.Builder, deadline flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(deadline), 0)
}
func AddExchangeOfferTransactionBufferStartDeadlineVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferTransactionBufferAddOffersCount(builder *flatbuffers.Builder, offersCount byte) {
	builder.PrependByteSlot(7, offersCount, 0)
}
func AddExchangeOfferTransactionBufferAddOffers(builder *flatbuffers.Builder, offers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(offers), 0)
}
func AddExchangeOfferTransactionBufferStartOffersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func AddExchangeOfferTransactionBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type ExchangeOfferBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsExchangeOfferBuffer(buf []byte, offset flatbuffers.UOffsetT) *ExchangeOfferBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ExchangeOfferBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ExchangeOfferBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ExchangeOfferBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ExchangeOfferBuffer) MosaicId(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MosaicIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MutateMosaicId(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *ExchangeOfferBuffer) MosaicAmount(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MosaicAmountLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MutateMosaicAmount(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *ExchangeOfferBuffer) Cost(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) CostLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MutateCost(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *ExchangeOfferBuffer) Type() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) MutateType(n byte) bool {
	return rcv._tab.MutateByteSlot(10, n)
}

func (rcv *ExchangeOfferBuffer) Owner(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) OwnerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferBuffer) OwnerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ExchangeOfferBuffer) MutateOwner(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func ExchangeOfferBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func ExchangeOfferBufferAddMosaicId(builder *flatbuffers.Builder, mosaicId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(mosaicId), 0)
}
func ExchangeOfferBufferStartMosaicIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferBufferAddMosaicAmount(builder *flatbuffers.Builder, mosaicAmount flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(mosaicAmount), 0)
}
func ExchangeOfferBufferStartMosaicAmountVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferBufferAddCost(builder *flatbuffers.Builder, cost flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(cost), 0)
}
func ExchangeOfferBufferStartCostVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferBufferAddType(builder *flatbuffers.Builder, type_ byte) {
	builder.PrependByteSlot(3, type_, 0)
}
func ExchangeOfferBufferAddOwner(builder *flatbuffers.Builder, owner flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(owner), 0)
}
func ExchangeOfferBufferStartOwnerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ExchangeOfferBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type ExchangeOfferTransactionBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsExchangeOfferTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *ExchangeOfferTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ExchangeOfferTransactionBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ExchangeOfferTransactionBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ExchangeOfferTransactionBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ExchangeOfferTransactionBuffer) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *ExchangeOfferTransactionBuffer) Signature(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) SignatureLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) SignatureBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ExchangeOfferTransactionBuffer) MutateSignature(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ExchangeOfferTransactionBuffer) Signer(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) SignerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) SignerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ExchangeOfferTransactionBuffer) MutateSigner(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *ExchangeOfferTransactionBuffer) Version() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateVersion(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *ExchangeOfferTransactionBuffer) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *ExchangeOfferTransactionBuffer) MaxFee(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MaxFeeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateMaxFee(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *ExchangeOfferTransactionBuffer) Deadline(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) DeadlineLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateDeadline(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *ExchangeOfferTransactionBuffer) OffersCount() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ExchangeOfferTransactionBuffer) MutateOffersCount(n byte) bool {
	return rcv._tab.MutateByteSlot(18, n)
}

func (rcv *ExchangeOfferTransactionBuffer) Offers(obj *ExchangeOfferBuffer, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *ExchangeOfferTransactionBuffer) OffersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ExchangeOfferTransactionBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(9)
}
func ExchangeOfferTransactionBufferAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func ExchangeOfferTransactionBufferAddSignature(builder *flatbuffers.Builder, signature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(signature), 0)
}
func ExchangeOfferTransactionBufferStartSignatureVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ExchangeOfferTransactionBufferAddSigner(builder *flatbuffers.Builder, signer flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(signer), 0)
}
func ExchangeOfferTransactionBufferStartSignerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ExchangeOfferTransactionBufferAddVersion(builder *flatbuffers.Builder, version uint32) {
	builder.PrependUint32Slot(3, version, 0)
}
func ExchangeOfferTransactionBufferAddType(builder *flatbuffers.Builder, type_ uint16) {
	builder.PrependUint16Slot(4, type_, 0)
}
func ExchangeOfferTransactionBufferAddMaxFee(builder *flatbuffers.Builder, maxFee flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(maxFee), 0)
}
func ExchangeOfferTransactionBufferStartMaxFeeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferTransactionBufferAddDeadline(builder *flatbuffers.Builder, deadline flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(deadline), 0)
}
func ExchangeOfferTransactionBufferStartDeadlineVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferTransactionBufferAddOffersCount(builder *flatbuffers.Builder, offersCount byte) {
	builder.PrependByteSlot(7, offersCount, 0)
}
func ExchangeOfferTransactionBufferAddOffers(builder *flatbuffers.Builder, offers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(offers), 0)
}
func ExchangeOfferTransactionBufferStartOffersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ExchangeOfferTransactionBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type RemoveExchangeOfferBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsRemoveExchangeOfferBuffer(buf []byte, offset flatbuffers.UOffsetT) *RemoveExchangeOfferBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RemoveExchangeOfferBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *RemoveExchangeOfferBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RemoveExchangeOfferBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RemoveExchangeOfferBuffer) MosaicId(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *RemoveExchangeOfferBuffer) MosaicIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveExchangeOfferBuffer) MutateMosaicId(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *RemoveExchangeOfferBuffer) Type() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveExchangeOfferBuffer) MutateType(n byte) bool {
	return rcv._tab.MutateByteSlot(6, n)
}

func RemoveExchangeOfferBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func RemoveExchangeOfferBufferAddMosaicId(builder *flatbuffers.Builder, mosaicId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(mosaicId), 0)
}
func RemoveExchangeOfferBufferStartMosaicIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RemoveExchangeOfferBufferAddType(builder *flatbuffers.Builder, type_ byte) {
	builder.PrependByteSlot(1, type_, 0)
}
func RemoveExchangeOfferBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type RemoveExchangeOfferTransactionBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsRemoveExchangeOfferTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *RemoveExchangeOfferTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RemoveExchangeOfferTransactionBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Signature(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) SignatureLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) SignatureBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateSignature(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Signer(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) SignerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) SignerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateSigner(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Version() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateVersion(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MaxFee(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MaxFeeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateMaxFee(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Deadline(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) DeadlineLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateDeadline(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *RemoveExchangeOfferTransactionBuffer) OffersCount() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RemoveExchangeOfferTransactionBuffer) MutateOffersCount(n byte) bool {
	return rcv._tab.MutateByteSlot(18, n)
}

func (rcv *RemoveExchangeOfferTransactionBuffer) Offers(obj *RemoveExchangeOfferBuffer, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RemoveExchangeOfferTransactionBuffer) OffersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RemoveExchangeOfferTransactionBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(9)
}
func RemoveExchangeOfferTransactionBufferAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func RemoveExchangeOfferTransactionBufferAddSignature(builder *flatbuffers.Builder, signature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(signature), 0)
}
func RemoveExchangeOfferTransactionBufferStartSignatureVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func RemoveExchangeOfferTransactionBufferAddSigner(builder *flatbuffers.Builder, signer flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(signer), 0)
}
func RemoveExchangeOfferTransactionBufferStartSignerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func RemoveExchangeOfferTransactionBufferAddVersion(builder *flatbuffers.Builder, version uint32) {
	builder.PrependUint32Slot(3, version, 0)
}
func RemoveExchangeOfferTransactionBufferAddType(builder *flatbuffers.Builder, type_ uint16) {
	builder.PrependUint16Slot(4, type_, 0)
}
func RemoveExchangeOfferTransactionBufferAddMaxFee(builder *flatbuffers.Builder, maxFee flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(maxFee), 0)
}
func RemoveExchangeOfferTransactionBufferStartMaxFeeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RemoveExchangeOfferTransactionBufferAddDeadline(builder *flatbuffers.Builder, deadline flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(deadline), 0)
}
func RemoveExchangeOfferTransactionBufferStartDeadlineVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RemoveExchangeOfferTransactionBufferAddOffersCount(builder *flatbuffers.Builder, offersCount byte) {
	builder.PrependByteSlot(7, offersCount, 0)
}
func RemoveExchangeOfferTransactionBufferAddOffers(builder *flatbuffers.Builder, offers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(offers), 0)
}
func RemoveExchangeOfferTransactionBufferStartOffersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RemoveExchangeOfferTransactionBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
