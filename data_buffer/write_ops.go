package data_buffer

import (
	"hash"
	"unicode/utf8"
	"unsafe"
)

func (w *WriteBuffer) Skip(n int) {
	w.EnsureSpace(n)
	w.data = w.data[:len(w.data)+n]
}

// **************
// STATIC ORDER *
// **************

func (w *WriteBuffer) U8(src uint8) {
	w.EnsureSpace(1)
	w.data = append(w.data, src)
}

func (w *WriteBuffer) U8_Slice(src []uint8) {
	w.EnsureSpace(len(src))
	w.data = append(w.data, src...)
}

func (w *WriteBuffer) String(src string) {
	w.EnsureSpace(len(src))
	w.data = append(w.data, src...)
}

func (w *WriteBuffer) I8(src int8) {
	w.U8(*(*uint8)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) I8_Slice(src []int8) {
	w.U8_Slice(*(*[]uint8)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) UVar16(src uint16) {
	w.UVar64(uint64(src))
}

func (w *WriteBuffer) UVar16_Slice(src []uint16) {
	for i := range src {
		w.UVar16(src[i])
	}
}

func (w *WriteBuffer) UVar32(src uint32) {
	w.UVar64(uint64(src))
}

func (w *WriteBuffer) UVar32_Slice(src []uint32) {
	for i := range src {
		w.UVar32(src[i])
	}
}

func (w *WriteBuffer) UVar64(src uint64) {
	for src >= 0b10000000 {
		w.EnsureSpace(1)
		w.data = append(w.data, byte(src)|0b10000000)
		src >>= 7
	}
	w.EnsureSpace(1)
	w.data = append(w.data, byte(src))
}

func (w *WriteBuffer) UVar64_Slice(src []uint64) {
	for i := range src {
		w.UVar64(src[i])
	}
}

func (w *WriteBuffer) Writable(src Writable) {
	src.WriteToBuffer(w)
}

func (w *WriteBuffer) Writable_Slice(src []Writable) {
	for i := range src {
		src[i].WriteToBuffer(w)
	}
}

func (w *WriteBuffer) Rune(src rune) {
	var arr [4]byte
	n := utf8.EncodeRune(arr[:], src)
	w.EnsureSpace(n)
	w.data = append(w.data, arr[:n]...)
}

func (w *WriteBuffer) Rune_Slice(src []rune) {
	var arr [4]byte
	for _, r := range src {
		n := utf8.EncodeRune(arr[:], r)
		w.EnsureSpace(n)
		w.data = append(w.data, arr[:n]...)
	}
}

func (w *WriteBuffer) Hash(hasher hash.Hash) {
	w.EnsureSpace(hasher.Size())
	w.data = hasher.Sum(w.data)
}

// ***************
// LITTLE ENDIAN *
// ***************

func (w *WriteBuffer) U16_LE(src uint16) {
	w.EnsureSpace(2)
	w.data = append(w.data, byte(src), byte(src>>8))
}

func (w *WriteBuffer) I16_LE(src int16) {
	w.U16_LE(*(*uint16)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) U24_LE(src uint32) {
	w.EnsureSpace(3)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16))
}

func (w *WriteBuffer) U32_LE(src uint32) {
	w.EnsureSpace(4)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16), byte(src>>24))
}

func (w *WriteBuffer) I32_LE(src int32) {
	w.U32_LE(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) F32_LE(src float32) {
	w.U32_LE(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) U40_LE(src uint64) {
	w.EnsureSpace(5)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16), byte(src>>24), byte(src>>32))
}

func (w *WriteBuffer) U48_LE(src uint64) {
	w.EnsureSpace(6)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16), byte(src>>24), byte(src>>32), byte(src>>40))
}

func (w *WriteBuffer) U56_LE(src uint64) {
	w.EnsureSpace(7)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16), byte(src>>24), byte(src>>32), byte(src>>40), byte(src>>48))
}

func (w *WriteBuffer) U64_LE(src uint64) {
	w.EnsureSpace(8)
	w.data = append(w.data, byte(src), byte(src>>8), byte(src>>16), byte(src>>24), byte(src>>32), byte(src>>40), byte(src>>48), byte(src>>48))
}

func (w *WriteBuffer) I64_LE(src int64) {
	w.U64_LE(*(*uint64)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) F64_LE(src float64) {
	w.U64_LE(*(*uint64)(unsafe.Pointer(&src)))
}

// ************
// BIG ENDIAN *
// ************

func (w *WriteBuffer) U16_BE(src uint16) {
	w.EnsureSpace(2)
	w.data = append(w.data, byte(src>>8), byte(src))
}

func (w *WriteBuffer) I16_BE(src int16) {
	w.U16_BE(*(*uint16)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) U24_BE(src uint32) {
	w.EnsureSpace(3)
	w.data = append(w.data, byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) U32_BE(src uint32) {
	w.EnsureSpace(4)
	w.data = append(w.data, byte(src>>24), byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) I32_BE(src int32) {
	w.U32_BE(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) F32_BE(src float32) {
	w.U32_BE(*(*uint32)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) U40_BE(src uint64) {
	w.EnsureSpace(5)
	w.data = append(w.data, byte(src>>32), byte(src>>24), byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) U48_BE(src uint64) {
	w.EnsureSpace(6)
	w.data = append(w.data, byte(src>>40), byte(src>>32), byte(src>>24), byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) U56_BE(src uint64) {
	w.EnsureSpace(7)
	w.data = append(w.data, byte(src>>48), byte(src>>40), byte(src>>32), byte(src>>24), byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) U64_BE(src uint64) {
	w.EnsureSpace(8)
	w.data = append(w.data, byte(src>>56), byte(src>>48), byte(src>>40), byte(src>>32), byte(src>>24), byte(src>>16), byte(src>>8), byte(src))
}

func (w *WriteBuffer) I64_BE(src int64) {
	w.U64_BE(*(*uint64)(unsafe.Pointer(&src)))
}

func (w *WriteBuffer) F64_BE(src float64) {
	w.U64_BE(*(*uint64)(unsafe.Pointer(&src)))
}
