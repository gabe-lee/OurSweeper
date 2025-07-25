package wire

import (
	"hash"
	"io"
	"math"
	"unsafe"
)

type (
	Reader = io.Reader
	Hash   = hash.Hash
)

const (
	maxInt = math.MaxInt
)

type Incoming struct {
	Rdr     Reader
	Ord     Order
	Err     error
	Hasher  Hash
	XorMask []byte
	XorPos  int
	Rem     int
}

func NewIncoming(reader io.Reader) Incoming {
	return Incoming{
		Rdr:     reader,
		Ord:     LE,
		Hasher:  nil,
		XorMask: nil,
		XorPos:  0,
		Rem:     maxInt,
	}
}

func NewIncomingAdv(reader io.Reader, order Order, hasher Hash, xorMask []byte, maxLen int) Incoming {
	return Incoming{
		Rdr:     reader,
		Ord:     order,
		Hasher:  hasher,
		XorMask: xorMask,
		XorPos:  0,
		Rem:     maxLen,
	}
}

func (w *Incoming) allAdv(err error, n int, targetN int, b []byte) {
	w.decrRem(n)
	w.coalescErrOrShort(err, n, targetN)
	w.addHash(b)
	w.xorData(b)
}

func (w *Incoming) allAdvSingle(err error, n int, targetN int, b []byte) {
	w.decrRem(n)
	w.coalescErrOrShort(err, n, targetN)
	w.addHash(b)
	w.xorSingle(b)
}

func (w *Incoming) addHash(b []byte) {
	if w.Hasher != nil {
		w.Hasher.Write(b)
	}
}

func (w *Incoming) getMax(max int) int {
	return min(w.Rem, max)
}

func (w *Incoming) decrRem(nn int) {
	w.Rem -= nn
}

func (w *Incoming) coalescErrOrShort(err error, n, targetN int) {
	if n < targetN {
		w.Err = io.ErrShortBuffer
	}
	if w.Err == nil {
		w.Err = err
	}
}

func (w *Incoming) xorSingle(b []byte) {
	if w.XorMask != nil {
		b[0] ^= w.XorMask[w.XorPos]
		w.XorPos += 1
		w.XorPos %= len(w.XorMask)
	}
}

func (w *Incoming) xorData(b []byte) {
	if w.XorMask != nil {
		for i := range b {
			b[i] ^= w.XorMask[w.XorPos]
			w.XorPos += 1
			w.XorPos %= len(w.XorMask)
		}
	}
}

func (w *Incoming) U8(dst *uint8) {
	arr := (*[1]byte)(unsafe.Pointer(dst))
	n := min(w.Rem, 1)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.allAdv(err, nn, 1, arr[:])
}

func (w *Incoming) I8(dst *int8) {
	w.U8((*uint8)(unsafe.Pointer(dst)))
}

func (w *Incoming) Bool(dst *bool) {
	w.U8((*uint8)(unsafe.Pointer(dst)))
}

func (w *Incoming) U16(dst *uint16) {
	var arr [2]byte
	n := min(w.Rem, 2)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read16(&arr, dst)
	w.allAdv(err, nn, 2, (*[2]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) I16(dst *int16) {
	w.U16((*uint16)(unsafe.Pointer(dst)))
}

func (w *Incoming) U24(dst *uint32) {
	var arr [3]byte
	n := min(w.Rem, 3)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read24(&arr, dst)
	w.allAdv(err, nn, 3, (*[3]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) U32(dst *uint32) {
	var arr [4]byte
	n := min(w.Rem, 4)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read32(&arr, dst)
	w.allAdv(err, nn, 4, (*[4]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) I32(dst *int32) {
	w.U32((*uint32)(unsafe.Pointer(dst)))
}

func (w *Incoming) F32(dst *float32) {
	w.U32((*uint32)(unsafe.Pointer(dst)))
}

func (w *Incoming) U40(dst *uint64) {
	var arr [5]byte
	n := min(w.Rem, 5)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read40(&arr, dst)
	w.allAdv(err, nn, 5, (*[5]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) U48(dst *uint64) {
	var arr [6]byte
	n := min(w.Rem, 6)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read48(&arr, dst)
	w.allAdv(err, nn, 6, (*[6]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) U56(dst *uint64) {
	var arr [7]byte
	n := min(w.Rem, 7)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read56(&arr, dst)
	w.allAdv(err, nn, 7, (*[7]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) U64(dst *uint64) {
	var arr [8]byte
	n := min(w.Rem, 8)
	nn, err := io.ReadFull(w.Rdr, arr[:n])
	w.Ord.Read64(&arr, dst)
	w.allAdv(err, nn, 8, (*[8]byte)(unsafe.Pointer(dst))[:])
}

func (w *Incoming) I64(dst *int64) {
	w.U64((*uint64)(unsafe.Pointer(dst)))
}

func (w *Incoming) F64(dst *float64) {
	w.U64((*uint64)(unsafe.Pointer(dst)))
}

func (w *Incoming) UVar16(dst *uint16) {
	*dst = 0
	var s uint
	var arr [1]byte
	var nn int
	var err error
	for i := range maxVarInt16Len {
		n := w.getMax(1)
		nn, err = io.ReadFull(w.Rdr, arr[:n])
		w.allAdvSingle(err, nn, 1, arr[:])
		if arr[0] < 0b10000000 {
			if i == maxVarInt16LastIdx && arr[0] > maxVarint16LastByte {
				*dst |= (uint16(arr[0]) << s)
				return
			}
			*dst |= (uint16(arr[0]) << s)
			return
		}
		*dst |= uint16(arr[0]&0b01111111) << s
		s += 7
	}
}

func (w *Incoming) IVar16(dst *int16) {
	w.UVar16((*uint16)(unsafe.Pointer(dst)))
}

func (w *Incoming) UVar32(dst *uint32) {
	*dst = 0
	var s uint
	var arr [1]byte
	var nn int
	var err error
	for i := range maxVarInt32Len {
		n := w.getMax(1)
		nn, err = io.ReadFull(w.Rdr, arr[:n])
		w.allAdvSingle(err, nn, 1, arr[:])
		if arr[0] < 0b10000000 {
			if i == maxVarInt32LastIdx && arr[0] > maxVarint32LastByte {
				*dst |= (uint32(arr[0]) << s)
				return
			}
			*dst |= (uint32(arr[0]) << s)
			return
		}
		*dst |= uint32(arr[0]&0b01111111) << s
		s += 7
	}
}

func (w *Incoming) IVar32(dst *int32) {
	w.UVar32((*uint32)(unsafe.Pointer(dst)))
}

func (w *Incoming) UVar64(dst *uint64) {
	*dst = 0
	var s uint
	var arr [1]byte
	var nn int
	var err error
	for i := range maxVarInt64Len {
		n := w.getMax(1)
		nn, err = io.ReadFull(w.Rdr, arr[:n])
		w.allAdvSingle(err, nn, 1, arr[:])
		if arr[0] < 0b10000000 || i == maxVarInt64LastIdx {
			*dst = *dst | (uint64(arr[0]) << s)
			return
		}
		*dst = *dst | uint64(arr[0]&0b01111111)<<s
		s += 7
	}
}

func (w *Incoming) IVar64(dst *int64) {
	w.UVar64((*uint64)(unsafe.Pointer(dst)))
}

func (w *Incoming) Struct(dst WireReader) {
	dst.WireRead(w)
}

func (w *Incoming) U8_Slice(dst []uint8) {
	io.ReadFull(w.Rdr, dst)
}

func (w *Incoming) I8_Slice(dst []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(dst)))
	uslice := unsafe.Slice(uptr, len(dst))
	w.U8_Slice(uslice)
}

func (w *Incoming) Bool_Slice(dst []bool) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(dst)))
	uslice := unsafe.Slice(uptr, len(dst))
	w.U8_Slice(uslice)
}

func (w *Incoming) String(str *string, len int) {
	data := make([]byte, len)
	w.U8_Slice(data)
	*str = string(data)
}

func (w *Incoming) U16_Slice(dst []uint16) {
	for i := range dst {
		w.U16(&dst[i])
	}
}

func (w *Incoming) I16_Slice(dst []int16) {
	for i := range dst {
		w.I16(&dst[i])
	}
}

func (w *Incoming) U24_Slice(dst []uint32) {
	for i := range dst {
		w.U24(&dst[i])
	}
}

func (w *Incoming) U32_Slice(dst []uint32) {
	for i := range dst {
		w.U32(&dst[i])
	}
}

func (w *Incoming) I32_Slice(dst []int32) {
	for i := range dst {
		w.I32(&dst[i])
	}
}

func (w *Incoming) F32_Slice(dst []float32) {
	for i := range dst {
		w.F32(&dst[i])
	}
}

func (w *Incoming) U40_Slice(dst []uint64) {
	for i := range dst {
		w.U40(&dst[i])
	}
}

func (w *Incoming) U48_Slice(dst []uint64) {
	for i := range dst {
		w.U48(&dst[i])
	}
}

func (w *Incoming) U56_Slice(dst []uint64) {
	for i := range dst {
		w.U56(&dst[i])
	}
}

func (w *Incoming) U64_Slice(dst []uint64) {
	for i := range dst {
		w.U64(&dst[i])
	}
}

func (w *Incoming) I64_Slice(dst []int64) {
	for i := range dst {
		w.I64(&dst[i])
	}
}

func (w *Incoming) F64_Slice(dst []float64) {
	for i := range dst {
		w.F64(&dst[i])
	}
}

func (w *Incoming) UVar16_Slice(dst []uint16) {
	for i := range dst {
		w.UVar16(&dst[i])
	}
}

func (w *Incoming) IVar16_Slice(dst []int16) {
	for i := range dst {
		w.IVar16(&dst[i])
	}
}

func (w *Incoming) UVar32_Slice(dst []uint32) {
	for i := range dst {
		w.UVar32(&dst[i])
	}
}

func (w *Incoming) IVar32_Slice(dst []int32) {
	for i := range dst {
		w.IVar32(&dst[i])
	}
}

func (w *Incoming) UVar64_Slice(dst []uint64) {
	for i := range dst {
		w.UVar64(&dst[i])
	}
}

func (w *Incoming) IVar64_Slice(dst []int64) {
	for i := range dst {
		w.IVar64(&dst[i])
	}
}

func (w *Incoming) Struct_Slice(dst []WireReader) {
	for i := range dst {
		dst[i].WireRead(w)
	}
}
