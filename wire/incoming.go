package wire

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type Incoming struct {
	rdr    io.Reader
	bin    Order
	err    error
	OnRead func(readBytes []byte)
}

func NewIncoming(reader io.Reader, order Order) Incoming {
	return Incoming{
		rdr: reader,
		bin: order,
	}
}
func NewIncomingSlice(data []byte, order Order) Incoming {
	return Incoming{
		rdr: bytes.NewReader(data),
		bin: order,
	}
}

func (w *Incoming) Err() error {
	return w.err
}

func (w *Incoming) HasErr() bool {
	return w.err != nil
}

func (w *Incoming) JoinErr(err error) {
	w.err = errors.Join(w.err, err)
}

func (w *Incoming) ReplaceErr(err error) {
	if err != nil {
		w.err = err
	}
}

func (w *Incoming) ClearErrs() {
	w.err = nil
}

func (w *Incoming) U8(ptr *uint8) {
	if w.err != nil {
		return
	}
	arr := (*[1]byte)(unsafe.Pointer(ptr))
	_, w.err = io.ReadFull(w.rdr, arr[:])
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) I8(ptr *int8) {
	w.U8((*uint8)(unsafe.Pointer(ptr)))
}

func (w *Incoming) Bool(ptr *bool) {
	w.U8((*uint8)(unsafe.Pointer(ptr)))
}

func (w *Incoming) U16(ptr *uint16) {
	if w.err != nil {
		return
	}
	var arr [2]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read16(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) I16(ptr *int16) {
	w.U16((*uint16)(unsafe.Pointer(ptr)))
}

func (w *Incoming) U24(ptr *uint32) {
	if w.err != nil {
		return
	}
	var arr [3]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read24(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) U32(ptr *uint32) {
	if w.err != nil {
		return
	}
	var arr [4]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read32(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) I32(ptr *int32) {
	w.U32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *Incoming) F32(ptr *float32) {
	w.U32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *Incoming) U40(ptr *uint64) {
	if w.err != nil {
		return
	}
	var arr [5]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read40(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) U48(ptr *uint64) {
	if w.err != nil {
		return
	}
	var arr [6]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read48(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) U54(ptr *uint64) {
	if w.err != nil {
		return
	}
	var arr [7]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read54(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) U64(ptr *uint64) {
	if w.err != nil {
		return
	}
	var arr [8]byte
	_, w.err = io.ReadFull(w.rdr, arr[:])
	w.bin.Read64(arr, ptr)
	if w.OnRead != nil {
		w.OnRead(arr[:])
	}
}

func (w *Incoming) I64(ptr *int64) {
	w.U64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *Incoming) F64(ptr *float64) {
	w.U64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *Incoming) UVar16(ptr *uint16) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt16Len {
		_, w.err = io.ReadFull(w.rdr, arr[0:1])
		if w.err != nil {
			return
		}
		if b < 0b10000000 {
			if i == maxVarInt16LastIdx && b > maxVarint16LastByte {
				w.err = ErrVarintOverflow16
				*ptr = *ptr | (uint16(b) << s)
				if w.OnRead != nil {
					w.OnRead(arr[:i+1])
				}
				return
			}
			*ptr = *ptr | (uint16(b) << s)
			if w.OnRead != nil {
				w.OnRead(arr[:i+1])
			}
			return
		}
		*ptr = *ptr | uint16(b&0b01111111)<<s
		s += 7
	}
	w.err = ErrVarintOverflow16
}

func (w *Incoming) IVar16(ptr *int16) {
	w.UVar16((*uint16)(unsafe.Pointer(ptr)))
}

func (w *Incoming) UVar32(ptr *uint32) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt32Len {
		_, w.err = io.ReadFull(w.rdr, arr[0:1])
		if w.err != nil {
			return
		}
		if b < 0b10000000 {
			if i == maxVarInt32LastIdx && b > maxVarint32LastByte {
				w.err = ErrVarintOverflow32
				*ptr = *ptr | (uint32(b) << s)
				if w.OnRead != nil {
					w.OnRead(arr[:i+1])
				}
				return
			}
			*ptr = *ptr | (uint32(b) << s)
			if w.OnRead != nil {
				w.OnRead(arr[:i+1])
			}
			return
		}
		*ptr = *ptr | uint32(b&0b01111111)<<s
		s += 7
	}
	w.err = ErrVarintOverflow32
}

func (w *Incoming) IVar32(ptr *int32) {
	w.UVar32((*uint32)(unsafe.Pointer(ptr)))
}

func (w *Incoming) UVar64(ptr *uint64) {
	if w.err != nil {
		return
	}
	*ptr = 0
	var s uint
	var b byte
	arr := (*[1]byte)(unsafe.Pointer(&b))
	for i := range maxVarInt64Len {
		_, w.err = io.ReadFull(w.rdr, arr[0:1])
		if w.err != nil {
			return
		}
		if b < 0b10000000 || i == maxVarInt64LastIdx {
			*ptr = *ptr | (uint64(b) << s)
			if w.OnRead != nil {
				w.OnRead(arr[:i+1])
			}
			return
		}
		*ptr = *ptr | uint64(b&0b01111111)<<s
		s += 7
	}
}

func (w *Incoming) IVar64(ptr *int64) {
	w.UVar64((*uint64)(unsafe.Pointer(ptr)))
}

func (w *Incoming) Struct(impl WireReader) {
	if w.err != nil {
		return
	}
	impl.WireRead(w)
}

func (w *Incoming) U8_Slice(slice []uint8) {
	if w.err != nil {
		return
	}
	var n int
	n, w.err = io.ReadFull(w.rdr, slice)
	if w.OnRead != nil {
		w.OnRead(slice[:n])
	}
}

func (w *Incoming) I8_Slice(slice []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.U8_Slice(uslice)
}

func (w *Incoming) Bool_Slice(slice []bool) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.U8_Slice(uslice)
}

func (w *Incoming) String(str *string, len int) {
	data := make([]byte, len)
	w.U8_Slice(data)
	*str = string(data)
}

func (w *Incoming) U16_Slice(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U16(&slice[i])
	}
}

func (w *Incoming) I16_Slice(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I16(&slice[i])
	}
}

func (w *Incoming) U32_Slice(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U32(&slice[i])
	}
}

func (w *Incoming) I32_Slice(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I32(&slice[i])
	}
}

func (w *Incoming) F32_Slice(slice []float32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.F32(&slice[i])
	}
}

func (w *Incoming) U64_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U64(&slice[i])
	}
}

func (w *Incoming) I64_Slice(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I64(&slice[i])
	}
}

func (w *Incoming) F64_Slice(slice []float64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.F64(&slice[i])
	}
}

func (w *Incoming) UVar16_Slice(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar16(&slice[i])
	}
}

func (w *Incoming) IVar16_Slice(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar16(&slice[i])
	}
}

func (w *Incoming) UVar32_Slice(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar32(&slice[i])
	}
}

func (w *Incoming) IVar32_Slice(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar32(&slice[i])
	}
}

func (w *Incoming) UVar64_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar64(&slice[i])
	}
}

func (w *Incoming) IVar64_Slice(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar64(&slice[i])
	}
}

func (w *Incoming) Struct_Slice(slice []WireReader) {
	for i := range slice {
		if w.err != nil {
			return
		}
		slice[i].WireRead(w)
	}
}

func (w *Incoming) Auto(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case *bool:
		w.Bool(T)
	case *int8:
		w.I8(T)
	case *uint8:
		w.U8(T)
	case *int16:
		w.I16(T)
	case *uint16:
		w.U16(T)
	case *int32:
		w.I32(T)
	case *uint32:
		w.U32(T)
	case *int64:
		w.I64(T)
	case *uint64:
		w.U64(T)
	case *float32:
		w.F32(T)
	case *float64:
		w.F64(T)
	case []bool:
		w.Bool_Slice(T)
	case []int8:
		w.I8_Slice(T)
	case []uint8:
		w.U8_Slice(T)
	case []int16:
		w.I16_Slice(T)
	case []uint16:
		w.U16_Slice(T)
	case []int32:
		w.I32_Slice(T)
	case []uint32:
		w.U32_Slice(T)
	case []int64:
		w.I64_Slice(T)
	case []uint64:
		w.U64_Slice(T)
	case []float32:
		w.F32_Slice(T)
	case []float64:
		w.F64_Slice(T)
	default:
		TT := reflect.TypeOf(T)
		if TT.Implements(reflect.TypeFor[WireReader]()) {
			I := *(*WireReader)(unsafe.Pointer(&T))
			w.Struct(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireReader]()) {
				ISlice := *(*[]WireReader)(unsafe.Pointer(&T))
				w.Struct_Slice(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryRead_Auto: not a pointer to a primitive type and does not implement WireReader", TT.Name())
	}
}

func (w *Incoming) AutoVarint(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case *bool:
		w.Bool(T)
	case *int8:
		w.I8(T)
	case *uint8:
		w.U8(T)
	case *int16:
		w.IVar16(T)
	case *uint16:
		w.UVar16(T)
	case *int32:
		w.IVar32(T)
	case *uint32:
		w.UVar32(T)
	case *int64:
		w.IVar64(T)
	case *uint64:
		w.UVar64(T)
	case *float32:
		w.F32(T)
	case *float64:
		w.F64(T)
	case []bool:
		w.Bool_Slice(T)
	case []int8:
		w.I8_Slice(T)
	case []uint8:
		w.U8_Slice(T)
	case []int16:
		w.IVar16_Slice(T)
	case []uint16:
		w.UVar16_Slice(T)
	case []int32:
		w.IVar32_Slice(T)
	case []uint32:
		w.UVar32_Slice(T)
	case []int64:
		w.IVar64_Slice(T)
	case []uint64:
		w.UVar64_Slice(T)
	case []float32:
		w.F32_Slice(T)
	case []float64:
		w.F64_Slice(T)
	default:
		TT := reflect.TypeOf(T)
		if TT.Implements(reflect.TypeFor[WireReader]()) {
			I := *(*WireReader)(unsafe.Pointer(&T))
			w.Struct(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireReader]()) {
				ISlice := *(*[]WireReader)(unsafe.Pointer(&T))
				w.Struct_Slice(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryRead_Auto: not a pointer to a primitive type and does not implement WireReader", TT.Name())
	}
}
