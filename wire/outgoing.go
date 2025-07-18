package wire

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type OutgoingWire struct {
	wtr io.Writer
	bin Order
	err error
	len int
}

func NewOutgoing(writer io.Writer, order Order) OutgoingWire {
	return OutgoingWire{
		wtr: writer,
		bin: order,
	}
}
func NewOutgoingSlice(initialCap int, order Order) OutgoingWire {
	buf := bytes.Buffer{}
	buf.Grow(initialCap)
	return OutgoingWire{
		wtr: &buf,
		bin: order,
	}
}

func (w *OutgoingWire) Len() int {
	return w.len
}

func (w *OutgoingWire) Err() error {
	return w.err
}

func (w *OutgoingWire) HasErr() bool {
	return w.err != nil
}

func (w *OutgoingWire) AddErr(err error) {
	if err != nil {
		if w.err == nil {
			w.err = err
		} else {
			w.err = errors.Join(w.err, err)
		}
	}
}

func (w *OutgoingWire) ClearErrs() {
	w.err = nil
}

func (w *OutgoingWire) GetOrder() Order {
	return w.bin
}

func (w *OutgoingWire) GetWriter() io.Writer {
	return w.wtr
}

func (w *OutgoingWire) TryWrite_U8(val uint8) {
	if w.err != nil {
		return
	}
	arr := (*[1]byte)(unsafe.Pointer(&val))
	var n int
	n, w.err = w.wtr.Write(arr[0:1])
	w.len += n
}

func (w *OutgoingWire) TryWrite_I8(val int8) {
	w.TryWrite_U8(*(*uint8)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_Bool(val bool) {
	w.TryWrite_U8(*(*uint8)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_U16(val uint16) {
	if w.err != nil {
		return
	}
	var arr [2]byte
	w.bin.WriteU16(val, &arr)
	var n int
	n, w.err = w.wtr.Write(arr[0:2])
	w.len += n
}

func (w *OutgoingWire) TryWrite_I16(val int16) {
	w.TryWrite_U16(*(*uint16)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_U32(val uint32) {
	if w.err != nil {
		return
	}
	var arr [4]byte
	w.bin.WriteU32(val, &arr)
	var n int
	n, w.err = w.wtr.Write(arr[0:4])
	w.len += n
}

func (w *OutgoingWire) TryWrite_I32(val int32) {
	w.TryWrite_U32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_F32(val float32) {
	w.TryWrite_U32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_U64(val uint64) {
	if w.err != nil {
		return
	}
	var arr [8]byte
	w.bin.WriteU64(val, &arr)
	var n int
	n, w.err = w.wtr.Write(arr[0:8])
	w.len += n
}

func (w *OutgoingWire) TryWrite_I64(val int64) {
	w.TryWrite_U64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_F64(val float64) {
	w.TryWrite_U64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_UVar16(val uint16) {
	if w.err != nil {
		return
	}
	i := 0
	var dst [maxVarInt16Len]byte
	for val >= 0b10000000 {
		dst[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	dst[i] = byte(val)
	var n int
	n, w.err = w.wtr.Write(dst[0:i])
	w.len += n
}

func (w *OutgoingWire) TryWrite_IVar16(val int16) {
	w.TryWrite_UVar16(*(*uint16)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_UVar32(val uint32) {
	if w.err != nil {
		return
	}
	i := 0
	var dst [maxVarInt32Len]byte
	for val >= 0b10000000 {
		dst[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	dst[i] = byte(val)
	var n int
	n, w.err = w.wtr.Write(dst[0:i])
	w.len += n
}

func (w *OutgoingWire) TryWrite_IVar32(val int32) {
	w.TryWrite_UVar32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_UVar64(val uint64) {
	if w.err != nil {
		return
	}
	i := 0
	var dst [maxVarInt64Len]byte
	for val >= 0b10000000 && i < maxVarInt64LastIdx {
		dst[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	dst[i] = byte(val)
	var n int
	n, w.err = w.wtr.Write(dst[0:i])
	w.len += n
}

func (w *OutgoingWire) TryWrite_IVar64(val int64) {
	w.TryWrite_UVar64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *OutgoingWire) TryWrite_WireWriter(val WireWriter) {
	if w.err != nil {
		return
	}
	val.WireWrite(w)
}

func (w *OutgoingWire) TryWrite_SliceU8(slice []uint8) {
	if w.err != nil {
		return
	}
	var n int
	n, w.err = w.wtr.Write(slice)
	w.len += n
}

func (w *OutgoingWire) TryWrite_SliceI8(slice []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.TryWrite_SliceU8(uslice)
}

func (w *OutgoingWire) TryWrite_SliceBool(slice []bool) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.TryWrite_SliceU8(uslice)
}

func (w *OutgoingWire) TryWrite_SliceU16(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_U16(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceI16(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_I16(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceU32(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_U32(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceI32(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_I32(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceF32(slice []float32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_F32(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceU64(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_U64(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceI64(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_I64(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceF64(slice []float64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_F64(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceUVar16(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_UVar16(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceIVar16(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_IVar16(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceUVar32(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_UVar32(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceIVar32(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_IVar32(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceUVar64(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_UVar64(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceIVar64(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.TryWrite_IVar64(slice[i])
	}
}

func (w *OutgoingWire) TryWrite_SliceWireWriter(slice []WireWriter) {
	for i := range slice {
		if w.err != nil {
			return
		}
		slice[i].WireWrite(w)
	}
}

func (w *OutgoingWire) TryWrite_Auto(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case bool:
		w.TryWrite_Bool(T)
	case int8:
		w.TryWrite_I8(T)
	case uint8:
		w.TryWrite_U8(T)
	case int16:
		w.TryWrite_I16(T)
	case uint16:
		w.TryWrite_U16(T)
	case int32:
		w.TryWrite_I32(T)
	case uint32:
		w.TryWrite_U32(T)
	case int64:
		w.TryWrite_I64(T)
	case uint64:
		w.TryWrite_U64(T)
	case float32:
		w.TryWrite_F32(T)
	case float64:
		w.TryWrite_F64(T)
	case []bool:
		w.TryWrite_SliceBool(T)
	case []int8:
		w.TryWrite_SliceI8(T)
	case []uint8:
		w.TryWrite_SliceU8(T)
	case []int16:
		w.TryWrite_SliceI16(T)
	case []uint16:
		w.TryWrite_SliceU16(T)
	case []int32:
		w.TryWrite_SliceI32(T)
	case []uint32:
		w.TryWrite_SliceU32(T)
	case []int64:
		w.TryWrite_SliceI64(T)
	case []uint64:
		w.TryWrite_SliceU64(T)
	case []float32:
		w.TryWrite_SliceF32(T)
	case []float64:
		w.TryWrite_SliceF64(T)
	default:
		TT := reflect.TypeOf(T)
		if TT.Implements(reflect.TypeFor[WireWriter]()) {
			I := *(*WireWriter)(unsafe.Pointer(&T))
			w.TryWrite_WireWriter(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireWriter]()) {
				ISlice := *(*[]WireWriter)(unsafe.Pointer(&T))
				w.TryWrite_SliceWireWriter(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryWrite_Auto: not a primitive type and does not implement WireWriter", TT.Name())
	}
}

func (w *OutgoingWire) TryWrite_AutoVarint(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case bool:
		w.TryWrite_Bool(T)
	case int8:
		w.TryWrite_I8(T)
	case uint8:
		w.TryWrite_U8(T)
	case int16:
		w.TryWrite_IVar16(T)
	case uint16:
		w.TryWrite_UVar16(T)
	case int32:
		w.TryWrite_IVar32(T)
	case uint32:
		w.TryWrite_UVar32(T)
	case int64:
		w.TryWrite_IVar64(T)
	case uint64:
		w.TryWrite_UVar64(T)
	case float32:
		w.TryWrite_F32(T)
	case float64:
		w.TryWrite_F64(T)
	case []bool:
		w.TryWrite_SliceBool(T)
	case []int8:
		w.TryWrite_SliceI8(T)
	case []uint8:
		w.TryWrite_SliceU8(T)
	case []int16:
		w.TryWrite_SliceIVar16(T)
	case []uint16:
		w.TryWrite_SliceUVar16(T)
	case []int32:
		w.TryWrite_SliceIVar32(T)
	case []uint32:
		w.TryWrite_SliceUVar32(T)
	case []int64:
		w.TryWrite_SliceIVar64(T)
	case []uint64:
		w.TryWrite_SliceUVar64(T)
	case []float32:
		w.TryWrite_SliceF32(T)
	case []float64:
		w.TryWrite_SliceF64(T)
	default:
		TT := reflect.TypeOf(T)
		if TT.Implements(reflect.TypeFor[WireWriter]()) {
			I := *(*WireWriter)(unsafe.Pointer(&T))
			w.TryWrite_WireWriter(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireWriter]()) {
				ISlice := *(*[]WireWriter)(unsafe.Pointer(&T))
				w.TryWrite_SliceWireWriter(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryWrite_AutoVarint: not a primitive type and does not implement WireWriter", TT.Name())
	}
}
