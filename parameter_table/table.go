package parameter_table

import (
	"fmt"
	"os"
	"slices"
	"unsafe"
)

// Set this to true if your program is stuck in an infinite loop or parameters are not updating correctly
var EnableDebug bool

type ParamCalcU64 func(table *ParamTable) uint64
type ParamCalcI64 func(table *ParamTable) int64
type ParamCalcF64 func(table *ParamTable) float64
type ParamCalcU32 func(table *ParamTable) uint32
type ParamCalcI32 func(table *ParamTable) int32
type ParamCalcF32 func(table *ParamTable) float32
type ParamCalcU16 func(table *ParamTable) uint16
type ParamCalcI16 func(table *ParamTable) int16
type ParamCalcU8 func(table *ParamTable) uint8
type ParamCalcI8 func(table *ParamTable) int8
type ParamCalcBool func(table *ParamTable) bool

const (
	rootU64 = iota
	rootI64
	rootF64
	derivedU64
	derivedI64
	derivedF64
	rootU32
	rootI32
	rootF32
	derivedU32
	derivedI32
	derivedF32
	rootU16
	rootI16
	derivedU16
	derivedI16
	rootU8
	rootI8
	rootBool
	derivedU8
	derivedI8
	derivedBool
	typeCount
)

const (
	dU64 = iota
	dI64
	dF64
	dU32
	dI32
	dF32
	dU16
	dI16
	dU8
	dI8
	dBool
	derivedCount
)

const (
	size64 uint32 = 8
	size32 uint32 = 4
	size16 uint32 = 2
	size8  uint32 = 1
)

var sizeTable = [typeCount]uint32{
	rootU64:     size64,
	rootI64:     size64,
	rootF64:     size64,
	derivedU64:  size64,
	derivedI64:  size64,
	derivedF64:  size64,
	rootU32:     size32,
	rootI32:     size32,
	rootF32:     size32,
	derivedU32:  size32,
	derivedI32:  size32,
	derivedF32:  size32,
	rootU16:     size16,
	rootI16:     size16,
	derivedU16:  size16,
	derivedI16:  size16,
	rootU8:      size8,
	rootI8:      size8,
	rootBool:    size8,
	derivedU8:   size8,
	derivedI8:   size8,
	derivedBool: size8,
}

var derivedTable = [typeCount]int{
	rootU64:     derivedCount,
	rootI64:     derivedCount,
	rootF64:     derivedCount,
	derivedU64:  dU64,
	derivedI64:  dI64,
	derivedF64:  dF64,
	rootU32:     derivedCount,
	rootI32:     derivedCount,
	rootF32:     derivedCount,
	derivedU32:  dU32,
	derivedI32:  dI32,
	derivedF32:  dF32,
	rootU16:     derivedCount,
	rootI16:     derivedCount,
	derivedU16:  dU16,
	derivedI16:  dI16,
	rootU8:      derivedCount,
	rootI8:      derivedCount,
	rootBool:    derivedCount,
	derivedU8:   dU8,
	derivedI8:   dI8,
	derivedBool: dBool,
}

type ParamTable struct {
	valuesPtr           *byte
	childStartsPtr      *uint32
	calcsPtr            *func()
	childrenPtr         *uint16
	childrenLen         uint32
	childrenCap         uint32
	byteOffsets         [typeCount]uint32
	valuesTotalByteLen  uint32
	idxOffsets          [typeCount]uint16
	valuesTotalIdxLen   uint16
	calcSubtractOffsets [derivedCount]uint16
	calcsTotalLen       uint16
}

func NewParamTable(rootU64End, rootI64End, rootF64End, derivedU64End, derivedI64End, derivedF64End, rootU32End, rootI32End, rootF32End, derivedU32End, derivedI32End, derivedF32End, rootU16End, rootI16End, derivedU16End, derivedI16End, rootU8End, rootI8End, rootBoolEnd, derivedU8End, derivedI8End, derivedBoolEnd uint16) ParamTable {
	if rootU64End > rootI64End || rootI64End > rootF64End || rootF64End > derivedU64End ||
		derivedU64End > derivedI64End || derivedI64End > derivedF64End || derivedF64End > rootU32End ||
		rootU32End > rootI32End || rootF32End > derivedU32End || derivedU32End > derivedI32End ||
		derivedI32End > derivedF32End || derivedF32End > rootU16End || rootU16End > rootI16End ||
		rootI16End > derivedU16End || derivedU16End > derivedI16End || derivedI16End > rootU8End ||
		rootU8End > rootI8End || rootI8End > rootBoolEnd || rootBoolEnd > derivedU8End ||
		derivedU8End > derivedI8End || derivedI8End > derivedBoolEnd {
		panic(`error: parameter table: indexes not in order: all parameter index ends MUST be in this EXACT order from smallest to largest:
	rootU64End <=
	rootI64End <=
	rootF64End <=
	derivedU64End <=
	derivedI64End <=
	derivedF64End <=
	rootU32End <=
	rootI32End <=
	rootF32End <=
	derivedU32End <=
	derivedI32End <=
	derivedF32End <=
	rootU16End <=
	rootI16End <=
	derivedU16End <=
	derivedI16End <=
	rootU8End <=
	rootI8End <=
	rootBoolEnd <=
	derivedU8End <=
	derivedI8End <=
	derivedBoolEnd <=
For an example template that fulfills this requirement, see the doc-comment of var PARAM_TABLE_TEMPLATE_DOC_COMMENT
	`)
	}
	var idxOffsets = [typeCount]uint16{
		rootU64:     0,
		rootI64:     rootU64End,
		rootF64:     rootI64End,
		derivedU64:  rootF64End,
		derivedI64:  derivedU64End,
		derivedF64:  derivedI64End,
		rootU32:     derivedF64End,
		rootI32:     rootU32End,
		rootF32:     rootI32End,
		derivedU32:  rootF32End,
		derivedI32:  derivedU32End,
		derivedF32:  derivedI32End,
		rootU16:     derivedF32End,
		rootI16:     rootU16End,
		derivedU16:  rootI16End,
		derivedI16:  derivedU16End,
		rootU8:      derivedI16End,
		rootI8:      rootU8End,
		rootBool:    rootI8End,
		derivedU8:   rootBoolEnd,
		derivedI8:   derivedU8End,
		derivedBool: derivedI8End,
	}
	var valuesTotalIdxLen = derivedBoolEnd
	var byteOffsets = [typeCount]uint32{
		rootU64:     0,
		rootI64:     uint32(rootU64End) * 8,
		rootF64:     uint32(rootI64End) * 8,
		derivedU64:  uint32(rootF64End) * 8,
		derivedI64:  uint32(derivedU64End) * 8,
		derivedF64:  uint32(derivedI64End) * 8,
		rootU32:     uint32(derivedF64End) * 8,
		rootI32:     (uint32(derivedF64End) * 8) + (uint32(rootU32End-derivedF64End) * 4),
		rootF32:     (uint32(derivedF64End) * 8) + (uint32(rootI32End-derivedF64End) * 4),
		derivedU32:  (uint32(derivedF64End) * 8) + (uint32(rootF32End-derivedF64End) * 4),
		derivedI32:  (uint32(derivedF64End) * 8) + (uint32(derivedU32End-derivedF64End) * 4),
		derivedF32:  (uint32(derivedF64End) * 8) + (uint32(derivedI32End-derivedF64End) * 4),
		rootU16:     (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4),
		rootI16:     (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(rootU16End-derivedF32End) * 2),
		derivedU16:  (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(rootI16End-derivedF32End) * 2),
		derivedI16:  (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedU16End-derivedF32End) * 2),
		rootU8:      (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2),
		rootI8:      (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2) + (uint32(rootU8End - derivedI16End)),
		rootBool:    (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2) + (uint32(rootI8End - derivedI16End)),
		derivedU8:   (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2) + (uint32(rootBoolEnd - derivedI16End)),
		derivedI8:   (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2) + (uint32(derivedU8End - derivedI16End)),
		derivedBool: (uint32(derivedF64End) * 8) + (uint32(derivedF32End-derivedF64End) * 4) + (uint32(derivedI16End-derivedF32End) * 2) + (uint32(derivedI8End - derivedI16End)),
	}
	var valuesTotalByteLen = byteOffsets[3] + uint32(derivedBoolEnd-derivedI16End)
	var calcSubtractOffsets = [derivedCount]uint16{
		dU64:  rootF64End,
		dI64:  rootF64End,
		dF64:  rootF64End,
		dU32:  rootF64End + rootF32End - derivedF64End,
		dI32:  rootF64End + rootF32End - derivedF64End,
		dF32:  rootF64End + rootF32End - derivedF64End,
		dU16:  rootF64End + rootF32End - derivedF64End + rootI16End - derivedF32End,
		dI16:  rootF64End + rootF32End - derivedF64End + rootI16End - derivedF32End,
		dU8:   rootF64End + rootF32End - derivedF64End + rootI16End - derivedF32End + rootBoolEnd - derivedI16End,
		dI8:   rootF64End + rootF32End - derivedF64End + rootI16End - derivedF32End + rootBoolEnd - derivedI16End,
		dBool: rootF64End + rootF32End - derivedF64End + rootI16End - derivedF32End + rootBoolEnd - derivedI16End,
	}
	var calcsTotalLen = (derivedF64End - rootF64End) + (derivedF32End - rootF32End) + (derivedI16End - rootI16End) + (derivedBoolEnd - rootBoolEnd)
	valuesSlice := make([]byte, valuesTotalByteLen)
	valuesPtr := unsafe.SliceData(valuesSlice)
	childStartsSlice := make([]uint32, valuesTotalIdxLen)
	childStartsPtr := unsafe.SliceData(childStartsSlice)
	calcsSlice := make([]func(), calcsTotalLen)
	calcsPtr := unsafe.SliceData(calcsSlice)
	childrenSlice := make([]uint16, 1)
	childrenCap := cap(childrenSlice)
	childrenPtr := unsafe.SliceData(childrenSlice)
	return ParamTable{
		valuesPtr:           valuesPtr,
		childStartsPtr:      childStartsPtr,
		calcsPtr:            calcsPtr,
		childrenPtr:         childrenPtr,
		childrenLen:         1,
		childrenCap:         uint32(childrenCap),
		byteOffsets:         byteOffsets,
		valuesTotalByteLen:  valuesTotalByteLen,
		idxOffsets:          idxOffsets,
		valuesTotalIdxLen:   valuesTotalIdxLen,
		calcSubtractOffsets: calcSubtractOffsets,
		calcsTotalLen:       calcsTotalLen,
	}
}

func (t *ParamTable) getTypeIdx(idx uint16, name string, validRoot, validDerived int, final bool, canBeRoot bool, canBeDerived bool) (typeIdx int, isDerived bool) {
	if EnableDebug {
		if idx >= t.valuesTotalIdxLen {
			fmt.Fprintf(os.Stderr, "error: parameter table: index %d is outside bounds of parameter list (len %d)", idx, t.valuesTotalIdxLen)
		}
		if final {
			if idx < t.idxOffsets[validRoot] && idx >= t.idxOffsets[validRoot+1] && idx < t.idxOffsets[validDerived] {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is not a %s value: %s values are in range [%d, %d) (root) or [%d, %d) (derived)", idx, name, name, t.idxOffsets[validRoot], t.idxOffsets[validRoot+1], t.idxOffsets[validDerived], t.valuesTotalIdxLen)
			}
		} else {
			if idx < t.idxOffsets[validRoot] && idx >= t.idxOffsets[validRoot+1] && idx < t.idxOffsets[validDerived] && idx >= t.idxOffsets[validDerived+1] {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is not a %s value: %s values are in range [%d, %d) (root) or [%d, %d) (derived)", idx, name, name, t.idxOffsets[validRoot], t.idxOffsets[validRoot+1], t.idxOffsets[validDerived], t.idxOffsets[validDerived+1])
			}
		}
	}
	if idx < t.idxOffsets[validDerived] {
		typeIdx = validRoot
		if EnableDebug {
			if !canBeRoot {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is a root value (cannot be a root value)", idx)
			}
		}
	} else {
		typeIdx = validDerived
		if EnableDebug {
			if !canBeDerived {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is a derived value (cannot be a derived value)", idx)
			}
		}
	}
	return
}

func (t *ParamTable) getBytePtr(idx uint16, typeIdx int) (ptr *byte, subIdx uint16) {
	subIdx = idx - t.idxOffsets[typeIdx]
	memOffset := t.byteOffsets[typeIdx] + (uint32(subIdx) * sizeTable[typeIdx])
	memSlice := unsafe.Slice(t.valuesPtr, t.valuesTotalByteLen)
	return &memSlice[memOffset], subIdx
}

func (t *ParamTable) Get_U8(idx uint16) uint8 {
	typeIdx, _ := t.getTypeIdx(idx, "Uint8", rootU8, derivedU8, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *memPtr
}

func (t *ParamTable) Get_I8(idx uint16) int8 {
	typeIdx, _ := t.getTypeIdx(idx, "Int8", rootI8, derivedI8, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*int8)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_Bool(idx uint16) bool {
	typeIdx, _ := t.getTypeIdx(idx, "Bool", rootBool, derivedBool, true, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*bool)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U16(idx uint16) uint16 {
	typeIdx, _ := t.getTypeIdx(idx, "Uint16", rootU16, derivedU16, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*uint16)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I16(idx uint16) int16 {
	typeIdx, _ := t.getTypeIdx(idx, "Int16", rootI16, derivedI16, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*int16)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U32(idx uint16) uint32 {
	typeIdx, _ := t.getTypeIdx(idx, "Uint32", rootU32, derivedU32, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*uint32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I32(idx uint16) int32 {
	typeIdx, _ := t.getTypeIdx(idx, "Int32", rootI32, derivedI32, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*int32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_F32(idx uint16) float32 {
	typeIdx, _ := t.getTypeIdx(idx, "Float32", rootF32, derivedF32, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*float32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U64(idx uint16) uint64 {
	typeIdx, _ := t.getTypeIdx(idx, "Uint64", rootU64, derivedU64, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*uint64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I64(idx uint16) int64 {
	typeIdx, _ := t.getTypeIdx(idx, "Uint64", rootI64, derivedI64, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*int64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_F64(idx uint16) float64 {
	typeIdx, _ := t.getTypeIdx(idx, "Float64", rootF64, derivedF64, false, true, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	return *(*float64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Set_U8(idx uint16, val uint8) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint8", rootU8, derivedU8, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*memPtr = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_I8(idx uint16, val int8) {
	typeIdx, _ := t.getTypeIdx(idx, "Int8", rootI8, derivedI8, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*int8)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_Bool(idx uint16, val bool) {
	typeIdx, _ := t.getTypeIdx(idx, "Bool", rootBool, derivedBool, true, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*bool)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_U16(idx uint16, val uint16) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint16", rootU16, derivedU16, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*uint16)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_I16(idx uint16, val int16) {
	typeIdx, _ := t.getTypeIdx(idx, "Int16", rootI16, derivedI16, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*int16)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_U32(idx uint16, val uint32) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint32", rootU32, derivedU32, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*uint32)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_I32(idx uint16, val int32) {
	typeIdx, _ := t.getTypeIdx(idx, "Int32", rootI32, derivedI32, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*int32)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_F32(idx uint16, val float32) {
	typeIdx, _ := t.getTypeIdx(idx, "Float32", rootF32, derivedF32, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*float32)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_U64(idx uint16, val uint64) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint64", rootU64, derivedU64, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*uint64)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_I64(idx uint16, val int64) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint64", rootI64, derivedI64, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*int64)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) Set_F64(idx uint16, val float64) {
	typeIdx, _ := t.getTypeIdx(idx, "Float64", rootF64, derivedF64, false, false, true)
	memPtr, _ := t.getBytePtr(idx, typeIdx)
	*(*float64)(unsafe.Pointer(memPtr)) = val
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) updateChildrenOfParents(idx uint16, parents []uint16) {
	childrenStart := unsafe.Slice(t.childStartsPtr, t.valuesTotalIdxLen)
	children := unsafe.Slice(t.childrenPtr, t.childrenLen)
	for _, parent := range parents {
		if childrenStart[parent] == 0 {
			newStart := t.childrenLen
			childrenStart[parent] = newStart
			children = append(children, idx, 0)
		} else {
			parentChildrenEnd := childrenStart[parent]
			for {
				if children[parentChildrenEnd] == 0 {
					break
				}
				parentChildrenEnd += 1
			}
			children = slices.Insert(children, int(parentChildrenEnd), idx)
			for i, childStart := range childrenStart {
				if childStart >= parentChildrenEnd {
					childrenStart[i] += 1
				}
			}
		}
	}
	t.childrenLen = uint32(len(children))
	t.childrenPtr = unsafe.SliceData(children)
}

func (t *ParamTable) getFuncPtr(subIdx uint16, typeIdx int) *func() {
	dType := derivedTable[typeIdx]
	finalSubIdx := t.idxOffsets[typeIdx] + subIdx - t.calcSubtractOffsets[dType]
	funcSlice := unsafe.Slice(t.calcsPtr, t.calcsTotalLen)
	funcPtr := &funcSlice[finalSubIdx]
	return funcPtr
}

func (t *ParamTable) InitDerived_U8(idx uint16, parents []uint16, calc ParamCalcU8) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint8", rootU8, derivedU8, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*memPtr = calc(t)
	*(*ParamCalcU8)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_I8(idx uint16, parents []uint16, calc ParamCalcI8) {
	typeIdx, _ := t.getTypeIdx(idx, "Int8", rootI8, derivedI8, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*int8)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcI8)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_Bool(idx uint16, parents []uint16, calc ParamCalcBool) {
	typeIdx, _ := t.getTypeIdx(idx, "Bool", rootBool, derivedBool, true, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*bool)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcBool)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_U16(idx uint16, parents []uint16, calc ParamCalcU16) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint16", rootU16, derivedU16, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*uint16)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcU16)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_I16(idx uint16, parents []uint16, calc ParamCalcI16) {
	typeIdx, _ := t.getTypeIdx(idx, "Int16", rootI16, derivedI16, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*int16)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcI16)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_U32(idx uint16, parents []uint16, calc ParamCalcU32) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint32", rootU32, derivedU32, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*uint32)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcU32)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_I32(idx uint16, parents []uint16, calc ParamCalcI32) {
	typeIdx, _ := t.getTypeIdx(idx, "Int32", rootI32, derivedI32, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*int32)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcI32)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_F32(idx uint16, parents []uint16, calc ParamCalcF32) {
	typeIdx, _ := t.getTypeIdx(idx, "Float32", rootF32, derivedF32, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*float32)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcF32)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_U64(idx uint16, parents []uint16, calc ParamCalcU64) {
	typeIdx, _ := t.getTypeIdx(idx, "Uint64", rootU64, derivedU64, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*uint64)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcU64)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_I64(idx uint16, parents []uint16, calc ParamCalcI64) {
	typeIdx, _ := t.getTypeIdx(idx, "Int64", rootI64, derivedI64, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*int64)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcI64)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) InitDerived_F64(idx uint16, parents []uint16, calc ParamCalcF64) {
	typeIdx, _ := t.getTypeIdx(idx, "Float64", rootF64, derivedF64, false, false, true)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	t.updateChildrenOfParents(idx, parents)
	*(*float64)(unsafe.Pointer(memPtr)) = calc(t)
	*(*ParamCalcF64)(unsafe.Pointer(funcPtr)) = calc
	t.updateChildren(idx, typeIdx, idx)
}

func (t *ParamTable) getTypeIdxDerivedBlind(idx uint16) int {
	switch {
	case idx >= t.idxOffsets[derivedBool]:
		return derivedBool
	case idx >= t.idxOffsets[derivedI8]:
		return derivedI8
	case idx >= t.idxOffsets[derivedU8]:
		return derivedU8
	case idx >= t.idxOffsets[derivedI16]:
		return derivedI16
	case idx >= t.idxOffsets[derivedU16]:
		return derivedU16
	case idx >= t.idxOffsets[derivedF32]:
		return derivedF32
	case idx >= t.idxOffsets[derivedI32]:
		return derivedI32
	case idx >= t.idxOffsets[derivedU32]:
		return derivedU32
	case idx >= t.idxOffsets[derivedF64]:
		return derivedF64
	case idx >= t.idxOffsets[derivedI64]:
		return derivedI64
	case idx >= t.idxOffsets[derivedU64]:
		return derivedU64
	default:
		return typeCount
	}
}

func (t *ParamTable) recalculate(idx uint16) (typeIdx int, didChange bool) {
	typeIdx = t.getTypeIdxDerivedBlind(idx)
	memPtr, subIdx := t.getBytePtr(idx, typeIdx)
	funcPtr := t.getFuncPtr(subIdx, typeIdx)
	switch typeIdx {
	case derivedBool:
		ptr := (*bool)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcBool)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedI8:
		ptr := (*int8)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcI8)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedU8:
		ptr := (*uint8)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcU8)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedI16:
		ptr := (*int16)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcI16)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedU16:
		ptr := (*uint16)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcU16)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedF32:
		ptr := (*float32)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcF32)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedI32:
		ptr := (*int32)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcI32)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedU32:
		ptr := (*uint32)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcU32)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedF64:
		ptr := (*float64)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcF64)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedI64:
		ptr := (*int64)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcI64)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	case derivedU64:
		ptr := (*uint64)(unsafe.Pointer(memPtr))
		oldVal := *ptr
		newVal := (*(*ParamCalcU64)(unsafe.Pointer(funcPtr)))(t)
		didChange = oldVal != newVal
		*ptr = newVal
	}
	return
}

func (t *ParamTable) updateChildren(idx uint16, typeIdx int, triggerIdx uint16) {
	childrenStart := unsafe.Slice(t.childStartsPtr, t.valuesTotalIdxLen)
	children := unsafe.Slice(t.childrenPtr, t.childrenLen)
	childIdxIdx := childrenStart[idx]
	if childIdxIdx == 0 {
		return
	}
	if childIdxIdx >= t.childrenLen {
		return
	}
	childIdx := children[childIdxIdx]
	for childIdx != 0 {
		if EnableDebug {
			if childIdx == triggerIdx {
				panic(fmt.Sprintf("error: parameter table: cyclic update loop: update triggered by idx %d, idx %d is a child or descendant of %d, however %d is also registered as a direct child of %d, causing infinite loop", triggerIdx, idx, triggerIdx, triggerIdx, idx))
			}
		}
		childTypeIdx, didChange := t.recalculate(childIdx)
		if didChange {
			t.updateChildren(childIdx, childTypeIdx, triggerIdx)
		}
		childIdxIdx += 1
		if childIdxIdx >= t.childrenLen {
			return
		}
		childIdx = children[childIdxIdx]
	}
}
