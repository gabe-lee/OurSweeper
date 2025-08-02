package parameter_table

import (
	"fmt"
	"math"
	"slices"
)

type ParamTable struct {
	params     []Param
	childLists []uint16
	opLists    []uint16
}

func NewParamTable(totalParams uint16) ParamTable {
	t := ParamTable{
		params:     make([]Param, totalParams),
		childLists: make([]uint16, 1),
		opLists:    make([]uint16, 1),
	}
	t.opLists[0] = _Stop
	return t
}

func (t *ParamTable) UpdateSourceParam(idx uint16, val float32) {
	t.params[idx].val = val
	t.updateChildren(idx)
}

func (t *ParamTable) GetParam(idx uint16) float32 {
	return t.params[idx].val
}

func (t *ParamTable) InitSourceParam(idx uint16, val float32) {
	if t.params[idx].operation_start > 0 || t.params[idx].children_start > 0 || t.params[idx].val != 0.0 {
		fmt.Printf("Parameter Heirarchy operation failure: parameter %d was already initialized", idx)
	}
	t.params[idx].val = val
	t.updateChildren(idx)
}

func (t *ParamTable) InitDerivedParam(idx uint16, operation ...uint16) {
	if t.params[idx].operation_start > 0 || t.params[idx].children_start > 0 || t.params[idx].val != 0.0 {
		fmt.Printf("Parameter Heirarchy operation failure: parameter %d was already initialized", idx)
	}
	parents := make([]uint16, 0)
	for _, token := range operation {
		if token <= MaxParamIndex {
			if !slices.Contains(parents, token) {
				parents = append(parents, token)
			}
		}
	}
	for _, parent := range parents {
		if t.params[parent].children_start == 0 {
			newStart := uint32(len(t.childLists))
			t.params[parent].children_start = newStart
			t.childLists = append(t.childLists, idx, 0)
		} else {
			parentChildrenEnd := t.params[parent].children_start
			for {
				if t.childLists[parentChildrenEnd] == 0 {
					break
				}
				parentChildrenEnd += 1
			}
			t.childLists = slices.Insert(t.childLists, int(parentChildrenEnd), idx)
			for p := range t.params {
				if t.params[p].children_start >= parentChildrenEnd {
					t.params[p].children_start += 1
				}
			}
		}
	}
	t.params[idx].operation_start = uint32(len(t.opLists))
	t.opLists = append(t.opLists, operation...)
	t.opLists = append(t.opLists, _Stop)
	t.doOp(idx, operation)
	t.updateChildren(idx)
}

func (t *ParamTable) updateChildren(idx uint16) {
	param := t.params[idx]
	if param.children_start == 0 {
		return
	}
	childIdxIdx := param.children_start
	if int(childIdxIdx) >= len(t.childLists) {
		return
	}
	childIdx := t.childLists[childIdxIdx]
	for childIdx != 0 {
		child := t.params[childIdx]
		didChange := t.doOp(childIdx, t.opLists[child.operation_start:])
		if didChange {
			t.updateChildren(childIdx)
		}
		childIdxIdx += 1
		if int(childIdxIdx) >= len(t.childLists) {
			return
		}
		childIdx = t.childLists[childIdxIdx]
	}
}

func (t *ParamTable) doOp(idx uint16, operation []uint16) (didChange bool) {
	prevVal := t.params[idx].val
	if len(operation) == 0 {
		return
	}
	if len(operation) < 2 {
		fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: non-empty operation must have at least 2 codes (one paramter index followed by one op code), got operation %v", idx, operation)
		return
	}
	if operation[0] == _Stop {
		return
	}
	var currOp uint16
	var currVal float32
	var boolCollect = false
	var startedBoolCollect = false
	var opAllowed bool = true
	var indexAllowed bool = false
	if operation[0] > MaxParamIndex {
		fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: first code in operation list was not a parameter index, got op code %d", idx, operation[0])
		return
	}
	if int(operation[0]) > len(t.params) {
		fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: code index 0 in operation list is a parmeter index outside param table len %d, got parameter index %d", idx, len(t.params), operation[0])
		return
	}
	currVal = t.params[operation[0]].val
	i := 1
loop:
	for i < len(operation) {
		switch operation[i] {
		case _Stop:
			break loop
		case Floor:
			if startedBoolCollect {
				if boolCollect {
					currVal = 1.0
				} else {
					currVal = 0.0
				}
				startedBoolCollect = false
			}
			if opNotAllowed(opAllowed, idx, i, operation[i]) {
				return
			}
			currVal = float32(math.Floor(float64(currVal)))
			opAllowed = true
			indexAllowed = false
		case Abs:
			if startedBoolCollect {
				if boolCollect {
					currVal = 1.0
				} else {
					currVal = 0.0
				}
				startedBoolCollect = false
			}
			if opNotAllowed(opAllowed, idx, i, operation[i]) {
				return
			}
			currVal = float32(math.Abs(float64(currVal)))
			opAllowed = true
			indexAllowed = false
		case Ceil:
			if startedBoolCollect {
				if boolCollect {
					currVal = 1.0
				} else {
					currVal = 0.0
				}
				startedBoolCollect = false
			}
			if opNotAllowed(opAllowed, idx, i, operation[i]) {
				return
			}
			currVal = float32(math.Ceil(float64(currVal)))
			opAllowed = true
			indexAllowed = false
		case Add, Subtract, SubtractReverse, Multiply, Divide, DivideReverse, Power, PowerReverse,
			Root, RootReverse, LogBase, LogBaseReverse, Modulo, ModuloReverse, ModuloWrap, ModuloWrapReverse,
			Min, Max, LessThan, GreaterThan, LessThanOrEqual, GreaterThanOrEqual, EqualTo, NotEqualTo,
			And, Or, Xor, BoolFlip:
			if startedBoolCollect {
				if boolCollect {
					currVal = 1.0
				} else {
					currVal = 0.0
				}
				startedBoolCollect = false
			}
			if opNotAllowed(opAllowed, idx, i, operation[i]) {
				return
			}
			currOp = operation[i]
			opAllowed = false
			indexAllowed = true
		default:
			if indexNotAllowed(indexAllowed, idx, i, operation[i]) {
				return
			}
			nextParam := operation[i]
			if int(nextParam) > len(t.params) {
				fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: code index %d in operation list is a parmeter index outside param table len %d, got parameter index %d", idx, i, len(t.params), nextParam)
				return
			}
			nextVal := t.params[operation[i]].val
			switch currOp {
			case Add:
				currVal = currVal + nextVal
			case Subtract:
				currVal = currVal - nextVal
			case SubtractReverse:
				currVal = nextVal - currVal
			case Multiply:
				currVal = currVal * nextVal
			case Divide:
				currVal = currVal / nextVal
			case DivideReverse:
				currVal = nextVal / currVal
			case Power:
				currVal = float32(math.Pow(float64(currVal), float64(nextVal)))
			case PowerReverse:
				currVal = float32(math.Pow(float64(nextVal), float64(currVal)))
			case Root:
				currVal = float32(math.Pow(float64(currVal), float64(1.0/nextVal)))
			case RootReverse:
				currVal = float32(math.Pow(float64(1.0/nextVal), float64(currVal)))
			case LogBase:
				currVal = float32(math.Log2(float64(currVal)) / math.Log2(float64(nextVal)))
			case LogBaseReverse:
				currVal = float32(math.Log2(float64(nextVal)) / math.Log2(float64(currVal)))
			case Modulo:
				currVal = float32(math.Mod(float64(currVal), float64(nextVal)))
			case ModuloReverse:
				currVal = float32(math.Mod(float64(nextVal), float64(currVal)))
			case ModuloWrap:
				currVal = currVal - float32(math.Floor(float64(currVal)/float64(nextVal))*float64(nextVal))
			case ModuloWrapReverse:
				currVal = nextVal - float32(math.Floor(float64(nextVal)/float64(currVal))*float64(currVal))
			case Min:
				currVal = min(currVal, nextVal)
			case Max:
				currVal = max(currVal, nextVal)
			case LessThan:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal < nextVal
				} else {
					boolCollect = boolCollect && currVal < nextVal
				}
			case GreaterThan:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal > nextVal
				} else {
					boolCollect = boolCollect && currVal > nextVal
				}
			case LessThanOrEqual:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal <= nextVal
				} else {
					boolCollect = boolCollect && currVal <= nextVal
				}
			case GreaterThanOrEqual:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal >= nextVal
				} else {
					boolCollect = boolCollect && currVal >= nextVal
				}
			case EqualTo:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal == nextVal
				} else {
					boolCollect = boolCollect && currVal == nextVal
				}
			case NotEqualTo:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal != nextVal
				} else {
					boolCollect = boolCollect && currVal != nextVal
				}
			case And:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal == 1.0 && nextVal == 1.0
				} else {
					boolCollect = boolCollect && nextVal == 1.0
				}
			case Or:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal == 1.0 || nextVal == 1.0
				} else {
					boolCollect = boolCollect || nextVal == 1.0
				}
			case Xor:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = (currVal == 1.0 && nextVal != 1.0) || (nextVal == 1.0 && currVal != 1.0)
				} else {
					boolCollect = (boolCollect && nextVal != 1.0) || (nextVal == 1.0 && !boolCollect)
				}
			case BoolFlip:
				if !startedBoolCollect {
					startedBoolCollect = true
					boolCollect = currVal == 0.0
				} else {
					boolCollect = !boolCollect
				}
			}
		}
		opAllowed = true
		indexAllowed = true
		i += 1
	}
	if startedBoolCollect {
		if boolCollect {
			currVal = 1.0
		} else {
			currVal = 0.0
		}
		startedBoolCollect = false
	}
	t.params[idx].val = currVal
	return prevVal != currVal
}

func opNotAllowed(allowed bool, idx uint16, i int, code uint16) bool {
	if !allowed {
		fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: code index %d in operation list cannot be an op code (previous code was a binary op code), got op code %d", idx, i, code)
		return true
	}
	return false
}

func indexNotAllowed(allowed bool, idx uint16, i int, code uint16) bool {
	if !allowed {
		fmt.Printf("Parameter Heirarchy operation failure: operation for parameter %d: code index %d in operation list cannot be a parameter index (previous code was a unary op code, or current binary op code does not allow index chaining), got parameter index %d", idx, i, code)
		return true
	}
	return false
}
