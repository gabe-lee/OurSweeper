package wire

import (
	"io"
	"unsafe"
)

type Outgoing struct {
	Wtr     io.Writer
	Ord     Order
	Err     error
	Hasher  Hash
	XorMask []byte
	XorPos  int
	Rem     int
}

func NewOutgoing(writer io.Writer) Outgoing {
	return Outgoing{
		Wtr:     writer,
		Ord:     LE,
		Hasher:  nil,
		XorMask: nil,
		XorPos:  0,
		Rem:     maxInt,
	}
}

func NewOutgoingAdv(writer io.Writer, order Order, hasher Hash, xorMask []byte, maxLen int) Outgoing {
	return Outgoing{
		Wtr:     writer,
		Ord:     order,
		Hasher:  hasher,
		XorMask: xorMask,
		XorPos:  0,
		Rem:     maxLen,
	}
}

func (w *Outgoing) allAdv(err error, n int, targetN int, b []byte) {
	w.decrRem(n)
	w.coalescErrOrShort(err, n, targetN)
	w.addHash(b)
	w.xorData(b)
}

func (w *Outgoing) allAdvNoXor(err error, n int, targetN int, b []byte) {
	w.decrRem(n)
	w.coalescErrOrShort(err, n, targetN)
	w.addHash(b)
}

func (w *Outgoing) allAdvNoHashNoXor(err error, n int, targetN int) {
	w.decrRem(n)
	w.coalescErrOrShort(err, n, targetN)
}

func (w *Outgoing) addHash(b []byte) {
	if w.Hasher != nil {
		w.Hasher.Write(b)
	}
}

func (w *Outgoing) getMax(max int) int {
	return min(w.Rem, max)
}

func (w *Outgoing) decrRem(nn int) {
	w.Rem -= nn
}

func (w *Outgoing) coalescErrOrShort(err error, n, targetN int) {
	if n < targetN {
		w.Err = io.ErrShortBuffer
	}
	if w.Err == nil {
		w.Err = err
	}
}

func (w *Outgoing) xorSingle(b []byte) {
	if w.XorMask != nil {
		b[0] ^= w.XorMask[w.XorPos]
		w.XorPos += 1
		w.XorPos %= len(w.XorMask)
	}
}

func (w *Outgoing) xorData(b []byte) {
	if w.XorMask != nil {
		for i := range b {
			b[i] ^= w.XorMask[w.XorPos]
			w.XorPos += 1
			w.XorPos %= len(w.XorMask)
		}
	}
}

func (w *Outgoing) U8(src uint8) {
	var arr [1]byte
	arr[0] = src
	w.xorSingle(arr[:])
	n := w.getMax(1)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 1, arr[:])
}

func (w *Outgoing) I8(src int8) {
	w.U8(*(*uint8)(unsafe.Pointer(&src)))
}

func (w *Outgoing) Bool(src bool) {
	w.U8(*(*uint8)(unsafe.Pointer(&src)))
}

func (w *Outgoing) U16(src uint16) {
	var arr [2]byte
	w.Ord.Write16(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(2)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 2, arr[:])
}

func (w *Outgoing) I16(src int16) {
	w.U16(*(*uint16)(unsafe.Pointer(&src)))
}

func (w *Outgoing) U24(src uint32) {
	var arr [3]byte
	w.Ord.Write24(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(3)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 3, arr[:])
}

func (w *Outgoing) U32(src uint32) {
	var arr [4]byte
	w.Ord.Write32(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(4)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 4, arr[:])
}

func (w *Outgoing) I32(src int32) {
	w.U32(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *Outgoing) F32(src float32) {
	w.U32(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *Outgoing) U40(src uint64) {
	var arr [5]byte
	w.Ord.Write40(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(5)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 5, arr[:])
}

func (w *Outgoing) U48(src uint64) {
	var arr [6]byte
	w.Ord.Write48(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(6)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 6, arr[:])
}

func (w *Outgoing) U56(src uint64) {
	var arr [7]byte
	w.Ord.Write56(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(7)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 7, arr[:])
}

func (w *Outgoing) U64(src uint64) {
	var arr [8]byte
	w.Ord.Write64(src, &arr)
	w.xorData(arr[:])
	n := w.getMax(8)
	nn, err := w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 8, arr[:])
}

func (w *Outgoing) I64(src int64) {
	w.U64(*(*uint64)(unsafe.Pointer(&src)))
}

func (w *Outgoing) F64(src float64) {
	w.U64(*(*uint64)(unsafe.Pointer(&src)))
}

func (w *Outgoing) UVar16(src uint16) {
	var arr [1]byte
	var n int
	var nn int
	var err error
	for src >= 0b10000000 {
		arr[0] = byte(src) | 0b10000000
		w.xorSingle(arr[:])
		n = w.getMax(1)
		nn, err = w.Wtr.Write(arr[:n])
		w.allAdvNoXor(err, nn, 1, arr[:])
		src >>= 7
	}
	arr[0] = byte(src)
	n = w.getMax(1)
	nn, err = w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 1, arr[:])
}

func (w *Outgoing) IVar16(src int16) {
	w.UVar16(*(*uint16)(unsafe.Pointer(&src)))
}

func (w *Outgoing) UVar32(src uint32) {
	var arr [1]byte
	var n int
	var nn int
	var err error
	for src >= 0b10000000 {
		arr[0] = byte(src) | 0b10000000
		w.xorSingle(arr[:])
		n = w.getMax(1)
		nn, err = w.Wtr.Write(arr[:n])
		w.allAdvNoXor(err, nn, 1, arr[:])
		src >>= 7
	}
	arr[0] = byte(src)
	n = w.getMax(1)
	nn, err = w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 1, arr[:])
}

func (w *Outgoing) IVar32(src int32) {
	w.UVar32(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *Outgoing) UVar64(src uint64) {
	var arr [1]byte
	var n int
	var nn int
	var err error
	for src >= 0b10000000 {
		arr[0] = byte(src) | 0b10000000
		w.xorSingle(arr[:])
		n = w.getMax(1)
		nn, err = w.Wtr.Write(arr[:n])
		w.allAdvNoXor(err, nn, 1, arr[:])
		src >>= 7
	}
	arr[0] = byte(src)
	n = w.getMax(1)
	nn, err = w.Wtr.Write(arr[:n])
	w.allAdvNoXor(err, nn, 1, arr[:])
}

func (w *Outgoing) IVar64(src int64) {
	w.UVar64(*(*uint64)(unsafe.Pointer(&src)))
}

func (w *Outgoing) Struct(src WireWriter) {
	src.WireWrite(w)
}

// This does NOT write the current hash sum BACK into the hasher
func (w *Outgoing) OwnHash() {
	if w.Hasher != nil {
		size := w.Hasher.Size()
		h := make([]byte, 0, size)
		h = w.Hasher.Sum(h)
		w.xorData(h)
		n := w.getMax(size)
		nn, err := w.Wtr.Write(h[:n])
		w.allAdvNoHashNoXor(err, nn, size)
	}
}

func (w *Outgoing) U8_Slice(src []uint8) {
	arr := make([]byte, len(src))
	copy(arr, src)
	w.xorData(arr)
	n := w.getMax(len(src))
	nn, err := w.Wtr.Write(src[:n])
	w.allAdvNoXor(err, nn, len(src), arr)
}

func (w *Outgoing) String(src string) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.StringData(src)))
	uslice := unsafe.Slice(uptr, len(src))
	w.U8_Slice(uslice)
}

func (w *Outgoing) I8_Slice(src []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(src)))
	uslice := unsafe.Slice(uptr, len(src))
	w.U8_Slice(uslice)
}

func (w *Outgoing) Bool_Slice(src []bool) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(src)))
	uslice := unsafe.Slice(uptr, len(src))
	w.U8_Slice(uslice)
}

func (w *Outgoing) U16_Slice(src []uint16) {
	for i := range src {
		w.U16(src[i])
	}
}

func (w *Outgoing) I16_Slice(src []int16) {
	for i := range src {
		w.I16(src[i])
	}
}

func (w *Outgoing) U24_Slice(src []uint32) {
	for i := range src {
		w.U24(src[i])
	}
}

func (w *Outgoing) U32_Slice(src []uint32) {
	for i := range src {
		w.U32(src[i])
	}
}

func (w *Outgoing) I32_Slice(src []int32) {
	for i := range src {
		w.I32(src[i])
	}
}

func (w *Outgoing) F32_Slice(src []float32) {
	for i := range src {
		w.F32(src[i])
	}
}

func (w *Outgoing) U40_Slice(src []uint64) {
	for i := range src {
		w.U40(src[i])
	}
}

func (w *Outgoing) U48_Slice(src []uint64) {
	for i := range src {
		w.U48(src[i])
	}
}

func (w *Outgoing) U56_Slice(src []uint64) {
	for i := range src {
		w.U56(src[i])
	}
}

func (w *Outgoing) U64_Slice(src []uint64) {
	for i := range src {
		w.U64(src[i])
	}
}

func (w *Outgoing) I64_Slice(src []int64) {
	for i := range src {
		w.I64(src[i])
	}
}

func (w *Outgoing) F64_Slice(src []float64) {
	for i := range src {
		w.F64(src[i])
	}
}

func (w *Outgoing) UVar16_Slice(src []uint16) {
	for i := range src {
		w.UVar16(src[i])
	}
}

func (w *Outgoing) IVar16_Slice(src []int16) {
	for i := range src {
		w.IVar16(src[i])
	}
}

func (w *Outgoing) UVar32_Slice(src []uint32) {
	for i := range src {
		w.UVar32(src[i])
	}
}

func (w *Outgoing) IVar32_Slice(src []int32) {
	for i := range src {
		w.IVar32(src[i])
	}
}

func (w *Outgoing) UVar64_Slice(src []uint64) {
	for i := range src {
		w.UVar64(src[i])
	}
}

func (w *Outgoing) IVar64_Slice(src []int64) {
	for i := range src {
		w.IVar64(src[i])
	}
}

func (w *Outgoing) Struct_Slice(src []WireWriter) {
	for i := range src {
		src[i].WireWrite(w)
	}
}
