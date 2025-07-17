package wire

import (
	"bytes"
	"io"
	"unsafe"
)

type IncomingWire struct {
	rdr io.Reader
	bin Order
	err error
	len int
}

func NewIncoming(reader io.Reader, order Order) IncomingWire {
	return IncomingWire{
		rdr: reader,
		bin: order,
	}
}
func NewIncomingSlice(data []byte, order Order) IncomingWire {
	return IncomingWire{
		rdr: bytes.NewReader(data),
		bin: order,
	}
}

func (w *IncomingWire) Len() int {
	return w.len
}

func (w *IncomingWire) Err() error {
	return w.err
}

func (w *IncomingWire) HasErr() bool {
	return w.err != nil
}

func (w *IncomingWire) GetOrder() Order {
	return w.bin
}

func (w *IncomingWire) GetReader() io.Reader {
	return w.rdr
}

func (w *IncomingWire) TryRead_U8(ptr *uint8) {
	if w.err != nil {
		return
	}
	arr := (*[1]byte)(unsafe.Pointer(ptr))
	var n int
	n, w.err = io.ReadFull(w.rdr, arr[0:1])
	w.len += n
}

func (w *IncomingWire) TryRead_I8(ptr *int8) {
	w.TryRead_U8((*uint8)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_U16(ptr *uint16) {
	if w.err != nil {
		return
	}
	var arr [2]byte
	var n int
	n, w.err = io.ReadFull(w.rdr, arr[0:2])
	w.bin.ReadU16(arr, ptr)
	w.len += n
}

func (w *IncomingWire) TryRead_I16(ptr *int16) {
	w.TryRead_U16((*uint16)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_U32(ptr *uint32) {
	if w.err != nil {
		return
	}
	var arr [4]byte
	var n int
	n, w.err = io.ReadFull(w.rdr, arr[0:4])
	w.bin.ReadU32(arr, ptr)
	w.len += n
}

func (w *IncomingWire) TryRead_I32(ptr *int32) {
	w.TryRead_U32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_F32(ptr *float32) {
	w.TryRead_U32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_U64(ptr *uint64) {
	if w.err != nil {
		return
	}
	var arr [8]byte
	var n int
	n, w.err = io.ReadFull(w.rdr, arr[0:8])
	w.bin.ReadU64(arr, ptr)
	w.len += n
}

func (w *IncomingWire) TryRead_I64(ptr *int64) {
	w.TryRead_U64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_F64(ptr *float64) {
	w.TryRead_U64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_UVar16(ptr *uint16) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt16Len {
		var n int
		n, w.err = io.ReadFull(w.rdr, arr[0:1])
		w.len += n
		if w.err != nil {
			return
		}
		if b < 0b10000000 {
			if i == maxVarInt16LastIdx && b > maxVarint16LastByte {
				w.err = ErrVarintOverflow16
				*ptr = *ptr | (uint16(b) << s)
				return
			}
			*ptr = *ptr | (uint16(b) << s)
			return
		}
		*ptr = *ptr | uint16(b&0b01111111)<<s
		s += 7
	}
	w.err = ErrVarintOverflow16
}

func (w *IncomingWire) TryRead_IVar16(ptr *int16) {
	w.TryRead_UVar16((*uint16)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_UVar32(ptr *uint32) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt32Len {
		var n int
		n, w.err = io.ReadFull(w.rdr, arr[0:1])
		w.len += n
		if w.err != nil {
			return
		}
		if b < 0b10000000 {
			if i == maxVarInt32LastIdx && b > maxVarint32LastByte {
				w.err = ErrVarintOverflow32
				*ptr = *ptr | (uint32(b) << s)
				return
			}
			*ptr = *ptr | (uint32(b) << s)
			return
		}
		*ptr = *ptr | uint32(b&0b01111111)<<s
		s += 7
	}
	w.err = ErrVarintOverflow32
}

func (w *IncomingWire) TryRead_IVar32(ptr *int32) {
	w.TryRead_UVar32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_UVar64(ptr *uint64) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt64Len {
		var n int
		n, w.err = io.ReadFull(w.rdr, arr[0:1])
		w.len += n
		if w.err != nil {
			return
		}
		if b < 0b10000000 || i == maxVarInt64LastIdx {
			*ptr = *ptr | (uint64(b) << s)
			return
		}
		*ptr = *ptr | uint64(b&0b01111111)<<s
		s += 7
	}
}

func (w *IncomingWire) TryRead_IVar64(ptr *int64) {
	w.TryRead_UVar64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *IncomingWire) TryRead_WireReader(x WireReader) {
	if w.err != nil {
		return
	}
	w.err = x.WireRead(w)
}

func (w *IncomingWire) TryRead_SliceU8(slice []uint8) {
	if w.err != nil {
		return
	}
	var n int
	n, w.err = io.ReadFull(w.rdr, slice)
	w.len += n
}

func (w *IncomingWire) TryRead_SliceI8(slice []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.TryRead_SliceU8(uslice)
}

func (w *IncomingWire) TryRead_SliceU16(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_U16(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceI16(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_I16(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceU32(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_U32(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceI32(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_I32(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceF32(slice []float32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_F32(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceU64(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_U64(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceI64(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_I64(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceF64(slice []float64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_F64(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceUVar16(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_UVar16(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceIVar16(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_IVar16(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceUVar32(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_UVar32(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceIVar32(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_IVar32(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceUVar64(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_UVar64(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceIVar64(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryRead_IVar64(&slice[i])
	}
}

func (w *IncomingWire) TryRead_SliceWireReader(slice []WireReader) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.err = slice[i].WireRead(w)

	}
}
