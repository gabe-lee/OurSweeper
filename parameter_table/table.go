package parameter_table

import (
	"fmt"
	"os"
	"slices"
	"unsafe"
)

// Recommended when first designing your parameter table
//
// Set this to true if your program is stuck in an infinite loop or parameters are not updating correctly
var EnableDebug bool

type ParamCalc func(recalc *RecalcInterface)

const (
	typeU64 = iota
	typeI64
	typeF64
	typeU32
	typeI32
	typeF32
	typeU16
	typeI16
	typeU8
	typeI8
	typeBool
	typeCount
)

const (
	size64 uint32 = 8
	size32 uint32 = 4
	size16 uint32 = 2
	size8  uint32 = 1
)

var sizeTable = [typeCount]uint32{
	typeU64:  size64,
	typeI64:  size64,
	typeF64:  size64,
	typeU32:  size32,
	typeI32:  size32,
	typeF32:  size32,
	typeU16:  size16,
	typeI16:  size16,
	typeU8:   size8,
	typeI8:   size8,
	typeBool: size8,
}

type paramHookups struct {
	parentsStart  uint32
	parentsEnd    uint32
	childrenStart uint32
	calculation   uint16
}

type ParamTable struct {
	values      []byte
	hookups     []paramHookups
	children    []uint16
	parents     []uint16
	calcs       []ParamCalc
	byteOffsets [typeCount]uint32
	idxOffsets  [typeCount]uint16
}

func NewParamTable(typeU64End, typeI64End, typeF64End, typeU32End, typeI32End, typeF32End, typeU16End, typeI16End, typeU8End, typeI8End, typeBoolEnd uint16, calcsCount uint16) ParamTable {
	if typeU64End > typeI64End || typeI64End > typeF64End || typeF64End > typeU32End ||
		typeU32End > typeI32End || typeI32End > typeF32End || typeF32End > typeU16End ||
		typeU16End > typeI16End || typeI16End > typeU8End ||
		typeU8End > typeI8End || typeI8End > typeBoolEnd {
		panic(`error: parameter table: indexes not in order: all parameter index ends MUST be in this EXACT order from smallest to largest:
	typeU64End <= typeI64End <= typeF64End <=
	typeU32End <= typeI32End <= typeF32End <=
	typeU16End <= typeI16End <=
	typeU8End <= typeI8End <= typeBoolEnd
For an example template that fulfills this requirement, see the doc-comment of var PARAM_TABLE_TEMPLATE_DOC_COMMENT`)
	}
	var idxOffsets = [typeCount]uint16{
		typeU64:  0,
		typeI64:  typeU64End,
		typeF64:  typeI64End,
		typeU32:  typeF64End,
		typeI32:  typeU32End,
		typeF32:  typeI32End,
		typeU16:  typeF32End,
		typeI16:  typeU16End,
		typeU8:   typeI16End,
		typeI8:   typeU8End,
		typeBool: typeI8End,
	}
	var valuesIdxLen = typeBoolEnd
	var byteOffsets = [typeCount]uint32{
		typeU64:  0,
		typeI64:  uint32(typeU64End) * 8,
		typeF64:  uint32(typeI64End) * 8,
		typeU32:  uint32(typeF64End) * 8,
		typeI32:  (uint32(typeF64End) * 8) + (uint32(typeU32End-typeF64End) * 4),
		typeF32:  (uint32(typeF64End) * 8) + (uint32(typeI32End-typeF64End) * 4),
		typeU16:  (uint32(typeF64End) * 8) + (uint32(typeF32End-typeF64End) * 4),
		typeI16:  (uint32(typeF64End) * 8) + (uint32(typeF32End-typeF64End) * 4) + (uint32(typeU16End-typeF32End) * 2),
		typeU8:   (uint32(typeF64End) * 8) + (uint32(typeF32End-typeF64End) * 4) + (uint32(typeI16End-typeF32End) * 2),
		typeI8:   (uint32(typeF64End) * 8) + (uint32(typeF32End-typeF64End) * 4) + (uint32(typeI16End-typeF32End) * 2) + uint32(typeU8End-typeI16End),
		typeBool: (uint32(typeF64End) * 8) + (uint32(typeF32End-typeF64End) * 4) + (uint32(typeI16End-typeF32End) * 2) + uint32(typeI8End-typeI16End),
	}
	var valuesByteLen = byteOffsets[typeBool] + uint32(typeBoolEnd-typeI16End)
	valuesSlice := make([]byte, valuesByteLen)
	hookupsSlice := make([]paramHookups, valuesIdxLen)
	calcsSlice := make([]ParamCalc, calcsCount)
	childrenSlice := make([]uint16, 1)
	parentsSlice := make([]uint16, 1)
	return ParamTable{
		values:      valuesSlice,
		hookups:     hookupsSlice,
		children:    childrenSlice,
		parents:     parentsSlice,
		calcs:       calcsSlice,
		byteOffsets: byteOffsets,
		idxOffsets:  idxOffsets,
	}
}

func (t *ParamTable) checkIdxType(idx uint16, name string, validType int, final bool, canBeDerived bool) {
	if EnableDebug {
		if idx >= uint16(len(t.values)) {
			fmt.Fprintf(os.Stderr, "error: parameter table: index %d is outside bounds of parameter list (len %d)", idx, len(t.values))
		}
		if final {
			if idx < t.idxOffsets[validType] {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is not a %s value: %s values are in range [%d, %d)", idx, name, name, t.idxOffsets[validType], len(t.values))
			}
		} else {
			if idx < t.idxOffsets[validType] && idx >= t.idxOffsets[validType+1] {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is not a %s value: %s values are in range [%d, %d)", idx, name, name, t.idxOffsets[validType], t.idxOffsets[validType+1])
			}
		}
		if !canBeDerived {
			if t.hookups[idx].parentsStart != 0 || t.hookups[idx].parentsEnd != 0 || t.hookups[idx].calculation != 0 {
				fmt.Fprintf(os.Stderr, "error: parameter table: index %d is a derived value (has parents and calculation func), cannot update directly", idx)
			}
		}
	}
}

func (t *ParamTable) getBytePtr(idx uint16, typeIdx int) (ptr *byte, subIdx uint16) {
	subIdx = idx - t.idxOffsets[typeIdx]
	memOffset := t.byteOffsets[typeIdx] + (uint32(subIdx) * sizeTable[typeIdx])
	return &t.values[memOffset], subIdx
}

func (t *ParamTable) Get_U8(idx uint16) uint8 {
	t.checkIdxType(idx, "Uint8", typeU8, false, true)
	memPtr, _ := t.getBytePtr(idx, typeU8)
	return *memPtr
}

func (t *ParamTable) Get_I8(idx uint16) int8 {
	t.checkIdxType(idx, "Int8", typeI8, false, true)
	memPtr, _ := t.getBytePtr(idx, typeI8)
	return *(*int8)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_Bool(idx uint16) bool {
	t.checkIdxType(idx, "Bool", typeBool, true, true)
	memPtr, _ := t.getBytePtr(idx, typeBool)
	return *(*bool)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U16(idx uint16) uint16 {
	t.checkIdxType(idx, "Uint16", typeU16, false, true)
	memPtr, _ := t.getBytePtr(idx, typeU16)
	return *(*uint16)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I16(idx uint16) int16 {
	t.checkIdxType(idx, "Int16", typeI16, false, true)
	memPtr, _ := t.getBytePtr(idx, typeI16)
	return *(*int16)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U32(idx uint16) uint32 {
	t.checkIdxType(idx, "Uint32", typeU32, false, true)
	memPtr, _ := t.getBytePtr(idx, typeU32)
	return *(*uint32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I32(idx uint16) int32 {
	t.checkIdxType(idx, "Int32", typeI32, false, true)
	memPtr, _ := t.getBytePtr(idx, typeI32)
	return *(*int32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_F32(idx uint16) float32 {
	t.checkIdxType(idx, "Float32", typeF32, false, true)
	memPtr, _ := t.getBytePtr(idx, typeF32)
	return *(*float32)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_U64(idx uint16) uint64 {
	t.checkIdxType(idx, "Uint64", typeU64, false, true)
	memPtr, _ := t.getBytePtr(idx, typeU64)
	return *(*uint64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_I64(idx uint16) int64 {
	t.checkIdxType(idx, "Int64", typeI64, false, true)
	memPtr, _ := t.getBytePtr(idx, typeI64)
	return *(*int64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) Get_F64(idx uint16) float64 {
	t.checkIdxType(idx, "Float64", typeF64, false, true)
	memPtr, _ := t.getBytePtr(idx, typeF64)
	return *(*float64)(unsafe.Pointer(memPtr))
}

func (t *ParamTable) set_U8(idx uint16, val uint8, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Uint8", typeU8, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeU8)
	oldVal := *memPtr
	*memPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_I8(idx uint16, val int8, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Int8", typeI8, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeI8)
	valPtr := (*int8)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_Bool(idx uint16, val bool, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Bool", typeBool, true, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeBool)
	valPtr := (*bool)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_U16(idx uint16, val uint16, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Uint16", typeU16, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeU16)
	valPtr := (*uint16)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_I16(idx uint16, val int16, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Int16", typeI16, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeI16)
	valPtr := (*int16)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_U32(idx uint16, val uint32, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Uint32", typeU32, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeU32)
	valPtr := (*uint32)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_I32(idx uint16, val int32, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Int32", typeI32, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeI32)
	valPtr := (*int32)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_F32(idx uint16, val float32, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Float32", typeF32, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeF32)
	valPtr := (*float32)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_U64(idx uint16, val uint64, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Uint64", typeU64, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeU64)
	valPtr := (*uint64)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_I64(idx uint16, val int64, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Int64", typeI64, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeI64)
	valPtr := (*int64)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) set_F64(idx uint16, val float64, canBeDerived bool, prevIdxs []uint16) (newPrevIdxs []uint16) {
	t.checkIdxType(idx, "Float64", typeF64, false, canBeDerived)
	memPtr, _ := t.getBytePtr(idx, typeF64)
	valPtr := (*float64)(unsafe.Pointer(memPtr))
	oldVal := *valPtr
	*valPtr = val
	newPrevIdxs = prevIdxs
	if oldVal != val {
		newPrevIdxs = t.updateChildren(idx, newPrevIdxs)
	}
	return
}

func (t *ParamTable) SetRoot_U8(idx uint16, val uint8) {
	prev := []uint16{idx}
	t.set_U8(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_I8(idx uint16, val int8) {
	prev := []uint16{idx}
	t.set_I8(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_Bool(idx uint16, val bool) {
	prev := []uint16{idx}
	t.set_Bool(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_U16(idx uint16, val uint16) {
	prev := []uint16{idx}
	t.set_U16(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_I16(idx uint16, val int16) {
	prev := []uint16{idx}
	t.set_I16(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_U32(idx uint16, val uint32) {
	prev := []uint16{idx}
	t.set_U32(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_I32(idx uint16, val int32) {
	prev := []uint16{idx}
	t.set_I32(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_F32(idx uint16, val float32) {
	prev := []uint16{idx}
	t.set_F32(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_U64(idx uint16, val uint64) {
	prev := []uint16{idx}
	t.set_U64(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_I64(idx uint16, val int64) {
	prev := []uint16{idx}
	t.set_I64(idx, val, false, prev)
}

func (t *ParamTable) SetRoot_F64(idx uint16, val float64) {
	prev := []uint16{idx}
	t.set_F64(idx, val, false, prev)
}

func (t *ParamTable) updateChildrenOfParents(idx uint16, parents []uint16) {
	t.hookups[idx].parentsStart = uint32(len(t.parents))
	t.hookups[idx].parentsEnd = uint32(len(t.parents)) + uint32(len(parents))
	t.parents = append(t.parents, parents...)
nextParent:
	for _, parent := range parents {
		if t.hookups[parent].childrenStart == 0 {
			newStart := uint32(len(t.children))
			t.hookups[parent].childrenStart = newStart
			t.children = append(t.children, idx, 0)
		} else {
			parentChildrenEnd := t.hookups[parent].childrenStart
			for {
				if t.children[parentChildrenEnd] == 0 {
					break
				}
				if t.children[parentChildrenEnd] == idx {
					continue nextParent
				}
				parentChildrenEnd += 1
			}
			t.children = slices.Insert(t.children, int(parentChildrenEnd), idx)
			for i := range t.hookups {
				if t.hookups[i].childrenStart >= parentChildrenEnd {
					t.hookups[i].childrenStart += 1
				}
			}
		}
	}
}

func (t *ParamTable) getCalc(calcIdx uint16) ParamCalc {
	if EnableDebug {
		if t.calcs[calcIdx] == nil {
			fmt.Fprintf(os.Stderr, "error: parameter table: calc index %d has not been registered", calcIdx)
		}
	}
	return t.calcs[calcIdx]
}

func (t *ParamTable) getCalcFromIdx(idx uint16) ParamCalc {
	calcIdx := t.hookups[idx].calculation
	return t.getCalc(calcIdx)
}

func (t *ParamTable) RegisterCalc(calcIdx uint16, calc ParamCalc) {
	if EnableDebug {
		if calcIdx > uint16(len(t.calcs)) {
			fmt.Fprintf(os.Stderr, "error: parameter table: calc index %d is outside bounds of calc list (len %d)", calcIdx, uint16(len(t.calcs)))
		}
		if t.calcs[calcIdx] != nil {
			fmt.Fprintf(os.Stderr, "error: parameter table: calc index %d is already registered", calcIdx)
		}
	}
	t.calcs[calcIdx] = calc
}

func (t *ParamTable) InitDerived_U8(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Uint8", typeU8, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_I8(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Int8", typeI8, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_Bool(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Bool", typeBool, true, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_U16(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Uint16", typeU16, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_I16(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Int16", typeI16, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_U32(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Uint32", typeU32, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_I32(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Int32", typeI32, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_F32(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Float32", typeF32, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_U64(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Uint64", typeU64, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_I64(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Int64", typeI64, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) InitDerived_F64(idx uint16, parents []uint16, calcIdx uint16) {
	t.checkIdxType(idx, "Float64", typeF64, false, true)
	calc := t.getCalc(calcIdx)
	t.hookups[idx].calculation = calcIdx
	t.updateChildrenOfParents(idx, parents)
	recalc := RecalcInterface{table: t, inputs: parents, output: idx}
	calc(&recalc)
	t.updateChildren(idx, []uint16{idx})
}

func (t *ParamTable) updateChildren(idx uint16, prevIdxs []uint16) (newPrevIdxs []uint16) {
	newPrevIdxs = prevIdxs
	childIdxIdx := t.hookups[idx].childrenStart
	if childIdxIdx == 0 {
		return
	}
	if childIdxIdx >= uint32(len(t.children)) {
		return
	}
	childIdx := t.children[childIdxIdx]
	for childIdx != 0 {
		if EnableDebug {
			for _, prevIdx := range prevIdxs {
				if childIdx == prevIdx {
					panic(fmt.Sprintf("error: parameter table: cyclic update loop: during update, idx %d was updated higher (previous) in the heirarchy, but idx %d had previous idx %d as a child, creating an infinite loop", prevIdx, idx, prevIdx))
				}
			}
			newPrevIdxs = append(newPrevIdxs, childIdx)
		}
		hookup := t.hookups[childIdx]
		calc := t.calcs[hookup.calculation]
		recalc := RecalcInterface{table: t, inputs: t.parents[hookup.parentsStart:hookup.parentsEnd], output: childIdx, prevIdxs: newPrevIdxs}
		calc(&recalc)
		newPrevIdxs = recalc.prevIdxs
		childIdxIdx += 1
		if childIdxIdx >= uint32(len(t.children)) {
			return
		}
		childIdx = t.children[childIdxIdx]
	}
	return
}

type RecalcInterface struct {
	table    *ParamTable
	inputs   []uint16
	prevIdxs []uint16
	output   uint16
}

func (t RecalcInterface) GetInput_U8(inputIdx uint16) uint8 {
	idx := t.inputs[inputIdx]
	return t.table.Get_U8(idx)
}
func (t RecalcInterface) GetInput_I8(inputIdx uint16) int8 {
	idx := t.inputs[inputIdx]
	return t.table.Get_I8(idx)
}
func (t RecalcInterface) GetInput_Bool(inputIdx uint16) bool {
	idx := t.inputs[inputIdx]
	return t.table.Get_Bool(idx)
}
func (t RecalcInterface) GetInput_U16(inputIdx uint16) uint16 {
	idx := t.inputs[inputIdx]
	return t.table.Get_U16(idx)
}
func (t RecalcInterface) GetInput_I16(inputIdx uint16) int16 {
	idx := t.inputs[inputIdx]
	return t.table.Get_I16(idx)
}
func (t RecalcInterface) GetInput_U32(inputIdx uint16) uint32 {
	idx := t.inputs[inputIdx]
	return t.table.Get_U32(idx)
}
func (t RecalcInterface) GetInput_I32(inputIdx uint16) int32 {
	idx := t.inputs[inputIdx]
	return t.table.Get_I32(idx)
}
func (t RecalcInterface) GetInput_F32(inputIdx uint16) float32 {
	idx := t.inputs[inputIdx]
	return t.table.Get_F32(idx)
}
func (t RecalcInterface) GetInput_U64(inputIdx uint16) uint64 {
	idx := t.inputs[inputIdx]
	return t.table.Get_U64(idx)
}
func (t RecalcInterface) GetInput_I64(inputIdx uint16) int64 {
	idx := t.inputs[inputIdx]
	return t.table.Get_I64(idx)
}
func (t RecalcInterface) GetInput_F64(inputIdx uint16) float64 {
	idx := t.inputs[inputIdx]
	return t.table.Get_F64(idx)
}

func (t *RecalcInterface) SetOutput_U8(val uint8) {
	t.prevIdxs = t.table.set_U8(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_I8(val int8) {
	t.prevIdxs = t.table.set_I8(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_Bool(val bool) {
	t.prevIdxs = t.table.set_Bool(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_U16(val uint16) {
	t.prevIdxs = t.table.set_U16(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_I16(val int16) {
	t.prevIdxs = t.table.set_I16(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_U32(val uint32) {
	t.prevIdxs = t.table.set_U32(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_I32(val int32) {
	t.prevIdxs = t.table.set_I32(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_F32(val float32) {
	t.prevIdxs = t.table.set_F32(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_U64(val uint64) {
	t.prevIdxs = t.table.set_U64(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_I64(val int64) {
	t.prevIdxs = t.table.set_I64(t.output, val, true, t.prevIdxs)
}
func (t *RecalcInterface) SetOutput_F64(val float64) {
	t.prevIdxs = t.table.set_F64(t.output, val, true, t.prevIdxs)
}
