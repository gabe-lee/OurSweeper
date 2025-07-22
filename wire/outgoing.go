package wire

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type Outgoing struct {
	wtr     io.Writer
	bin     Order
	err     error
	OnWrite func(writtenBytes []byte)
}

func NewOutgoing(writer io.Writer, order Order) Outgoing {
	return Outgoing{
		wtr: writer,
		bin: order,
	}
}
func NewOutgoingBuffer(initialCap int, order Order) Outgoing {
	buf := bytes.Buffer{}
	buf.Grow(initialCap)
	return Outgoing{
		wtr: &buf,
		bin: order,
	}
}

func (w *Outgoing) Err() error {
	return w.err
}

func (w *Outgoing) JoinErr(err error) {
	w.err = errors.Join(w.err, err)
}

func (w *Outgoing) ReplaceErr(err error) {
	if err != nil {
		w.err = err
	}
}

func (w *Outgoing) HasErr() bool {
	return w.err != nil
}

func (w *Outgoing) ClearErrs() {
	w.err = nil
}

func (w *Outgoing) U8(val uint8) {
	if w.err != nil {
		return
	}
	arr := (*[1]byte)(unsafe.Pointer(&val))
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) I8(val int8) {
	w.U8(*(*uint8)(unsafe.Pointer(&val)))
}

func (w *Outgoing) Bool(val bool) {
	w.U8(*(*uint8)(unsafe.Pointer(&val)))
}

func (w *Outgoing) U16(val uint16) {
	if w.err != nil {
		return
	}
	var arr [2]byte
	w.bin.Write16(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) I16(val int16) {
	w.U16(*(*uint16)(unsafe.Pointer(&val)))
}

func (w *Outgoing) U24(val uint32) {
	if w.err != nil {
		return
	}
	var arr [3]byte
	w.bin.Write24(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) U32(val uint32) {
	if w.err != nil {
		return
	}
	var arr [4]byte
	w.bin.Write32(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) I32(val int32) {
	w.U32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *Outgoing) F32(val float32) {
	w.U32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *Outgoing) U40(val uint64) {
	if w.err != nil {
		return
	}
	var arr [5]byte
	w.bin.Write40(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}
func (w *Outgoing) U48(val uint64) {
	if w.err != nil {
		return
	}
	var arr [6]byte
	w.bin.Write48(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}
func (w *Outgoing) U54(val uint64) {
	if w.err != nil {
		return
	}
	var arr [7]byte
	w.bin.Write54(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) U64(val uint64) {
	if w.err != nil {
		return
	}
	var arr [8]byte
	w.bin.Write64(val, &arr)
	_, w.err = w.wtr.Write(arr[:])
	if w.OnWrite != nil {
		w.OnWrite(arr[:])
	}
}

func (w *Outgoing) I64(val int64) {
	w.U64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *Outgoing) F64(val float64) {
	w.U64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *Outgoing) UVar16(val uint16) {
	if w.err != nil {
		return
	}
	i := 0
	var arr [maxVarInt16Len]byte
	for val >= 0b10000000 {
		arr[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	arr[i] = byte(val)
	_, w.err = w.wtr.Write(arr[:i])
	if w.OnWrite != nil {
		w.OnWrite(arr[:i])
	}
}

func (w *Outgoing) IVar16(val int16) {
	w.UVar16(*(*uint16)(unsafe.Pointer(&val)))
}

func (w *Outgoing) UVar32(val uint32) {
	if w.err != nil {
		return
	}
	i := 0
	var arr [maxVarInt32Len]byte
	for val >= 0b10000000 {
		arr[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	arr[i] = byte(val)
	_, w.err = w.wtr.Write(arr[:i])
	if w.OnWrite != nil {
		w.OnWrite(arr[:i])
	}
}

func (w *Outgoing) IVar32(val int32) {
	w.UVar32(*(*uint32)(unsafe.Pointer(&val)))
}

func (w *Outgoing) UVar64(val uint64) {
	if w.err != nil {
		return
	}
	i := 0
	var arr [maxVarInt64Len]byte
	for val >= 0b10000000 && i < maxVarInt64LastIdx {
		arr[i] = byte(val) | 0b10000000
		val >>= 7
		i++
	}
	arr[i] = byte(val)
	_, w.err = w.wtr.Write(arr[:i])
	if w.OnWrite != nil {
		w.OnWrite(arr[:i])
	}
}

func (w *Outgoing) IVar64(val int64) {
	w.UVar64(*(*uint64)(unsafe.Pointer(&val)))
}

func (w *Outgoing) Struct(val WireWriter) {
	if w.err != nil {
		return
	}
	val.WireWrite(w)
}

func (w *Outgoing) U8_Slice(slice []uint8) {
	if w.err != nil {
		return
	}
	var n int
	n, w.err = w.wtr.Write(slice)
	if w.OnWrite != nil {
		w.OnWrite(slice[:n])
	}
}

func (w *Outgoing) String(str string) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.StringData(str)))
	uslice := unsafe.Slice(uptr, len(str))
	w.U8_Slice(uslice)
}

func (w *Outgoing) I8_Slice(slice []int8) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.U8_Slice(uslice)
}

func (w *Outgoing) Bool_Slice(slice []bool) {
	uptr := (*uint8)(unsafe.Pointer(unsafe.SliceData(slice)))
	uslice := unsafe.Slice(uptr, len(slice))
	w.U8_Slice(uslice)
}

func (w *Outgoing) U16_Slice(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U16(slice[i])
	}
}

func (w *Outgoing) I16_Slice(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I16(slice[i])
	}
}

func (w *Outgoing) U24_Slice(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U24(slice[i])
	}
}

func (w *Outgoing) U32_Slice(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U32(slice[i])
	}
}

func (w *Outgoing) I32_Slice(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I32(slice[i])
	}
}

func (w *Outgoing) F32_Slice(slice []float32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.F32(slice[i])
	}
}

func (w *Outgoing) U40_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U40(slice[i])
	}
}

func (w *Outgoing) U48_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U48(slice[i])
	}
}

func (w *Outgoing) U54_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U54(slice[i])
	}
}

func (w *Outgoing) U64_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.U64(slice[i])
	}
}

func (w *Outgoing) I64_Slice(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.I64(slice[i])
	}
}

func (w *Outgoing) F64_Slice(slice []float64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.F64(slice[i])
	}
}

func (w *Outgoing) UVar16_Slice(slice []uint16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar16(slice[i])
	}
}

func (w *Outgoing) IVar16_Slice(slice []int16) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar16(slice[i])
	}
}

func (w *Outgoing) UVar32_Slice(slice []uint32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar32(slice[i])
	}
}

func (w *Outgoing) IVar32_Slice(slice []int32) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar32(slice[i])
	}
}

func (w *Outgoing) UVar64_Slice(slice []uint64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.UVar64(slice[i])
	}
}

func (w *Outgoing) IVar64_Slice(slice []int64) {
	for i := range slice {
		if w.err != nil {
			return
		}
		w.IVar64(slice[i])
	}
}

func (w *Outgoing) Struct_Slice(slice []WireWriter) {
	for i := range slice {
		if w.err != nil {
			return
		}
		slice[i].WireWrite(w)
	}
}

func (w *Outgoing) Auto(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case bool:
		w.Bool(T)
	case int8:
		w.I8(T)
	case uint8:
		w.U8(T)
	case int16:
		w.I16(T)
	case uint16:
		w.U16(T)
	case int32:
		w.I32(T)
	case uint32:
		w.U32(T)
	case int64:
		w.I64(T)
	case uint64:
		w.U64(T)
	case float32:
		w.F32(T)
	case float64:
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
		if TT.Implements(reflect.TypeFor[WireWriter]()) {
			I := *(*WireWriter)(unsafe.Pointer(&T))
			w.Struct(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireWriter]()) {
				ISlice := *(*[]WireWriter)(unsafe.Pointer(&T))
				w.Struct_Slice(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryWrite_Auto: not a primitive type and does not implement WireWriter", TT.Name())
	}
}

func (w *Outgoing) AutoVarint(val any) {
	if w.err != nil {
		return
	}
	switch T := val.(type) {
	case bool:
		w.Bool(T)
	case int8:
		w.I8(T)
	case uint8:
		w.U8(T)
	case int16:
		w.IVar16(T)
	case uint16:
		w.UVar16(T)
	case int32:
		w.IVar32(T)
	case uint32:
		w.UVar32(T)
	case int64:
		w.IVar64(T)
	case uint64:
		w.UVar64(T)
	case float32:
		w.F32(T)
	case float64:
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
		if TT.Implements(reflect.TypeFor[WireWriter]()) {
			I := *(*WireWriter)(unsafe.Pointer(&T))
			w.Struct(I)
			return
		}
		if TT.Kind() == reflect.Slice {
			TTT := TT.Elem()
			if TTT.Implements(reflect.TypeFor[WireWriter]()) {
				ISlice := *(*[]WireWriter)(unsafe.Pointer(&T))
				w.Struct_Slice(ISlice)
			}
		}
		w.err = fmt.Errorf("invalid type `%s` for TryWrite_AutoVarint: not a primitive type and does not implement WireWriter", TT.Name())
	}
}
