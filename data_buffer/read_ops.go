package data_buffer

import (
	"io"
	"unicode/utf8"
	"unsafe"
)

func (r *ReadBuffer) Discard(n int) (discardedSlice []byte, err error) {
	nn, err := r.CheckReadSpace(n)
	dis := r.data[r.pos : r.pos+nn]
	r.pos += nn
	return dis, err
}

// **************
// STATIC ORDER *
// **************

func (r *ReadBuffer) U8(dst *uint8) error {
	src, err := r.Discard(1)
	if err != nil {
		return err
	}
	*dst = src[0]
	return nil
}

func (r *ReadBuffer) U8_Slice(dst []uint8) error {
	src, err := r.Discard(len(dst))
	if err != nil {
		return err
	}
	copy(dst, src)
	return nil
}

func (r *ReadBuffer) String(dst *string, len int) error {
	src, err := r.Discard(len)
	if err != nil {
		return err
	}
	*dst = string(src)
	return nil
}

func (r *ReadBuffer) I8(dst *int8) error {
	return r.U8((*uint8)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) I8_Slice(dst []int8) error {
	return r.U8_Slice(*(*[]uint8)(unsafe.Pointer(&dst)))
}

func (r *ReadBuffer) UVar16(dst *uint16) error {
	*dst = 0
	var s uint
	for i := range maxVarInt16Len {
		src, err := r.Discard(1)
		if err != nil {
			return err
		}
		if src[0] < 0b10000000 {
			*dst |= (uint16(src[0]) << s)
			if i == maxVarInt16LastIdx && src[0] > maxVarint16LastByte {
				return ErrVarintOverflow16
			}
			return nil
		}
		*dst |= uint16(src[0]&0b01111111) << s
		s += 7
	}
	return ErrVarintOverflow16
}

func (r *ReadBuffer) UVar16_Slice(src []uint16) error {
	for i := range src {
		err := r.UVar16(&src[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadBuffer) UVar32(dst *uint32) error {
	*dst = 0
	var s uint
	for i := range maxVarInt32Len {
		src, err := r.Discard(1)
		if err != nil {
			return err
		}
		if src[0] < 0b10000000 {
			*dst |= (uint32(src[0]) << s)
			if i == maxVarInt32LastIdx && src[0] > maxVarint32LastByte {
				return ErrVarintOverflow32
			}
			return nil
		}
		*dst |= uint32(src[0]&0b01111111) << s
		s += 7
	}
	return ErrVarintOverflow32
}

func (r *ReadBuffer) UVar32_Slice(src []uint32) error {
	for i := range src {
		err := r.UVar32(&src[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadBuffer) UVar64(dst *uint64) error {
	*dst = 0
	var s uint
	for i := range maxVarInt64Len {
		src, err := r.Discard(1)
		if err != nil {
			return err
		}
		if src[0] < 0b10000000 || i == maxVarInt64LastIdx {
			*dst |= (uint64(src[0]) << s)
			return nil
		}
		*dst |= uint64(src[0]&0b01111111) << s
		s += 7
	}
	return nil
}

func (r *ReadBuffer) UVar64_Slice(src []uint64) error {
	for i := range src {
		err := r.UVar64(&src[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadBuffer) Readable(dst Readable) error {
	return dst.ReadFromBuffer(r)
}

func (r *ReadBuffer) Readable_Slice(dst []Readable) error {
	for i := range dst {
		err := dst[i].ReadFromBuffer(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadBuffer) Rune(src *rune) error {
	rn, n := utf8.DecodeRune(r.data)
	*src = rn
	r.Discard(n)
	if rn == utf8.RuneError {
		if n == 0 {
			return io.ErrShortWrite
		} else {
			return ErrInvalidUTF8
		}
	}
	return nil
}

func (r *ReadBuffer) Rune_Slice(src []rune) error {
	for i := range src {
		err := r.Rune(&src[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// ***************
// LITTLE ENDIAN *
// ***************

func (r *ReadBuffer) U16_LE(dst *uint16) error {
	src, err := r.Discard(2)
	if err != nil {
		return err
	}
	*dst = uint16(src[0])
	*dst |= uint16(src[1]) << 8
	return nil
}

func (r *ReadBuffer) I16_LE(dst *int16) error {
	return r.U16_LE((*uint16)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) U24_LE(dst *uint32) error {
	src, err := r.Discard(3)
	if err != nil {
		return err
	}
	*dst = uint32(src[0])
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2]) << 16
	return nil
}

func (r *ReadBuffer) U32_LE(dst *uint32) error {
	src, err := r.Discard(4)
	if err != nil {
		return err
	}
	*dst = uint32(src[0])
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2]) << 16
	*dst |= uint32(src[3]) << 24
	return nil
}

func (r *ReadBuffer) I32_LE(dst *int32) error {
	return r.U32_LE((*uint32)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) F32_LE(dst *float32) error {
	return r.U32_LE((*uint32)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) U40_LE(dst *uint64) error {
	src, err := r.Discard(5)
	if err != nil {
		return err
	}
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	return nil
}

func (r *ReadBuffer) U48_LE(dst *uint64) error {
	src, err := r.Discard(6)
	if err != nil {
		return err
	}
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
	return nil
}

func (r *ReadBuffer) U56_LE(dst *uint64) error {
	src, err := r.Discard(7)
	if err != nil {
		return err
	}
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
	*dst |= uint64(src[6]) << 48
	return nil
}

func (r *ReadBuffer) U64_LE(dst *uint64) error {
	src, err := r.Discard(8)
	if err != nil {
		return err
	}
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
	*dst |= uint64(src[6]) << 48
	*dst |= uint64(src[7]) << 56
	return nil
}

func (r *ReadBuffer) I64_LE(dst *int64) error {
	return r.U64_LE((*uint64)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) F64_LE(dst *float64) error {
	return r.U64_LE((*uint64)(unsafe.Pointer(dst)))
}

// ************
// BIG ENDIAN *
// ************

func (r *ReadBuffer) U16_BE(dst *uint16) error {
	src, err := r.Discard(2)
	if err != nil {
		return err
	}
	*dst = uint16(src[0]) << 8
	*dst |= uint16(src[1])
	return nil
}

func (r *ReadBuffer) I16_BE(dst *int16) error {
	return r.U16_BE((*uint16)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) U24_BE(dst *uint32) error {
	src, err := r.Discard(3)
	if err != nil {
		return err
	}
	*dst = uint32(src[0]) << 16
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2])
	return nil
}

func (r *ReadBuffer) U32_BE(dst *uint32) error {
	src, err := r.Discard(4)
	if err != nil {
		return err
	}
	*dst = uint32(src[0]) << 24
	*dst |= uint32(src[1]) << 16
	*dst |= uint32(src[2]) << 8
	*dst |= uint32(src[3])
	return nil
}

func (r *ReadBuffer) I32_BE(dst *int32) error {
	return r.U32_BE((*uint32)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) F32_BE(dst *float32) error {
	return r.U32_BE((*uint32)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) U40_BE(dst *uint64) error {
	src, err := r.Discard(5)
	if err != nil {
		return err
	}
	*dst = uint64(src[0]) << 32
	*dst |= uint64(src[1]) << 24
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 8
	*dst |= uint64(src[4])
	return nil
}

func (r *ReadBuffer) U48_BE(dst *uint64) error {
	src, err := r.Discard(6)
	if err != nil {
		return err
	}
	*dst = uint64(src[0]) << 40
	*dst |= uint64(src[1]) << 32
	*dst |= uint64(src[2]) << 24
	*dst |= uint64(src[3]) << 16
	*dst |= uint64(src[4]) << 8
	*dst |= uint64(src[5])
	return nil
}

func (r *ReadBuffer) U56_BE(dst *uint64) error {
	src, err := r.Discard(7)
	if err != nil {
		return err
	}
	*dst = uint64(src[0]) << 48
	*dst |= uint64(src[1]) << 40
	*dst |= uint64(src[2]) << 32
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 16
	*dst |= uint64(src[5]) << 8
	*dst |= uint64(src[6])
	return nil
}

func (r *ReadBuffer) U64_BE(dst *uint64) error {
	src, err := r.Discard(8)
	if err != nil {
		return err
	}
	*dst = uint64(src[0]) << 56
	*dst |= uint64(src[1]) << 48
	*dst |= uint64(src[2]) << 40
	*dst |= uint64(src[3]) << 32
	*dst |= uint64(src[4]) << 24
	*dst |= uint64(src[5]) << 16
	*dst |= uint64(src[6]) << 8
	*dst |= uint64(src[7])
	return nil
}

func (r *ReadBuffer) I64_BE(dst *int64) error {
	return r.U64_BE((*uint64)(unsafe.Pointer(dst)))
}

func (r *ReadBuffer) F64_BE(dst *float64) error {
	return r.U64_BE((*uint64)(unsafe.Pointer(dst)))
}
