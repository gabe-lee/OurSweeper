package parameter_table

import "math"

type Op uint8

const (
	_Stop uint16 = math.MaxUint16 - iota
	// ### Patterns:
	//   - {V1, Add, V2} -> V1 + V2
	//   - {... , Add, VN} -> (result from all previous ops) + VN
	//   - {V1, Add, V2, V3, V4} -> V1 + V2 + V3 + V4
	Add
	// ### Patterns:
	//   - {V1, Subtract, V2} -> V1 - V2
	//   - {... , Subtract, VN} -> (result from all previous ops) - VN
	//   - {V1, Subtract, V2, V3, V4} -> V1 - V2 - V3 - V4
	Subtract
	// ### Patterns:
	//   - {V1, SubtractReverse, V2} -> V2 - V1
	//   - {... , SubtractReverse, VN} -> VN - (result from all previous ops)
	//   - {V1, SubtractReverse, V2, V3, V4} -> V4 - V3 - V2 - V1
	SubtractReverse
	// ### Patterns:
	//   - {V1, Multiply, V2} -> V1 * V2
	//   - {... , Multiply, VN} -> (result from all previous ops) * VN
	//   - {V1, Multiply, V2, V3, V4} -> V1 * V2 * V3 * V4
	Multiply
	// ### Patterns:
	//   - {V1, Divide, V2} -> V1 / V2
	//   - {... , Divide, VN} -> (result from all previous ops) / VN
	//   - {V1, Divide, V2, V3, V4} -> V1 / V2 / V3 / V4
	Divide
	// ### Patterns:
	//   - {V1, DivideReverse, V2} -> V2 / V1
	//   - {... , DivideReverse, VN} -> VN / (result from all previous ops)
	//   - {V1, DivideReverse, V2, V3, V4} -> V4 / V3 / V2 / V1
	DivideReverse
	// ### Patterns:
	//   - {V1, Power, V2} -> V1 ^ V2
	//   - {... , Power, VN} -> (result from all previous ops) ^ VN
	//   - {V1, Power, V2, V3, V4} -> V1 ^ V2 ^ V3 ^ V4
	Power
	// ### Patterns:
	//   - {V1, PowerReverse, V2} -> V2 ^ V1
	//   - {... , PowerReverse, VN} -> VN ^ (result from all previous ops)
	//   - {V1, PowerReverse, V2, V3, V4} -> V4 ^ V3 ^ V2 ^ V1
	PowerReverse
	// ### Patterns:
	//   - {V1, Root, V2} -> V1 √ V2
	//   - {... , Root, VN} -> (result from all previous ops) √ VN
	Root
	// ### Patterns:
	//   - {V1, RootReverse, V2} -> V2 √ V1
	//   - {... , RootReverse, VN} -> VN √ (result from all previous ops)
	RootReverse
	// ### Patterns:
	//   - {V1, LogBase, V2} -> logᵥ₂(V1)
	//   - {... , LogBase, VN} -> logᵥₙ(result from all previous ops)
	LogBase
	// ### Patterns:
	//   - {V1, LogBaseReverse, V2} -> logᵥ₁(V2)
	//   - {... , LogBaseReverse, VN} -> logᵥₓ(VN) where ᵥₓ == (result from all previous ops)
	LogBaseReverse
	// ### Patterns:
	//   - {V1, Modulo, V2} -> math.Mod(V1, V2)
	//   - {... , Modulo, VN} -> math.Mod((result from all previous ops), VN)
	Modulo
	// ### Patterns:
	//   - {V1, ModuloReverse, V2} -> V2 % V1
	//   - {... , ModuloReverse, VN} -> VN % (result from all previous ops)
	//   - {V1, ModuloReverse, V2, V3, V4} -> V4 % V3 % V2 % V1
	ModuloReverse
	// ### Patterns:
	//   - {V1, ModuloWrap, V2} -> V1 - (math.Floor(V1/V2)*V2)
	//   - {... , ModuloWrap, VN} -> VPrev - (math.Floor(VPrev/VN)*VN) where VPrev = (result from all previous ops)
	ModuloWrap
	// ### Patterns:
	//   - {V1, ModuloWrapReverse, V2} -> V2 - (math.Floor(V2/V1)*V1)
	//   - {... , ModuloWrapReverse, VN} -> VN - (math.Floor(VN/VPrev)*VPrev) where VPrev = (result from all previous ops)
	ModuloWrapReverse
	// ### Patterns:
	//   - {V1, Min, V2} -> min(V1, V2)
	//   - {... , Min, VN} -> min((result from all previous ops), VN)
	//   - {V1, Min, V2, V3, V4} -> min(V1, V2, V3, V4)
	Min
	// ### Patterns:
	//   - {V1, Max, V2} -> max(V1, V2)
	//   - {... , Max, VN} -> max((result from all previous ops), VN)
	//   - {V1, Max, V2, V3, V4} -> max(V1, V2, V3, V4)
	Max
	// ### Patterns:
	//   - {V1, LessThan, V2} -> 1.0 if V1 < V2, else 0.0
	//   - {... , LessThan, VN} -> 1.0 if (result from all previous ops) < VN, else 0.0
	LessThan
	// ### Patterns:
	//   - {V1, GreaterThan, V2} -> 1.0 if V1 > V2, else 0.0
	//   - {... , GreaterThan, VN} -> 1.0 if (result from all previous ops) > VN, else 0.0
	GreaterThan
	// ### Patterns:
	//   - {V1, LessThanOrEqual, V2} -> 1.0 if V1 <= V2, else 0.0
	//   - {... , LessThanOrEqual, VN} -> 1.0 if (result from all previous ops) <= VN, else 0.0
	LessThanOrEqual
	// ### Patterns:
	//   - {V1, GreaterThanOrEqual, V2} -> 1.0 if V1 >= V2, else 0.0
	//   - {... , GreaterThanOrEqual, VN} -> 1.0 if (result from all previous ops) >= VN, else 0.0
	GreaterThanOrEqual
	// ### Patterns:
	//   - {V1, EqualTo, V2} -> 1.0 if V1 == V2, else 0.0
	//   - {... , EqualTo, VN} -> 1.0 if (result from all previous ops) == VN, else 0.0
	EqualTo
	// ### Patterns:
	//   - {V1, NotEqualTo, V2} -> 1.0 if V1 != V2, else 0.0
	//   - {... , NotEqualTo, VN} -> 1.0 if (result from all previous ops) != VN, else 0.0
	NotEqualTo
	// ### Patterns:
	//   - {V1, Floor} -> math.Floor(V1)
	//   - {... , Floor} -> math.Floor((result from all previous ops))
	Floor
	// ### Patterns:
	//   - {V1, Ceil} -> math.Ceil(V1)
	//   - {... , Ceil} -> math.Ceil((result from all previous ops))
	Ceil
	// ### Patterns:
	//   - {V1, And, V2} -> 1.0 if V1 == 1.0 and V2 == 1.0, else 0.0
	//   - {... , GreaterThan, VN} -> 1.0 if (result from all previous ops) == 1.0 and VN == 1.0
	And
	// ### Patterns:
	//   - {V1, And, V2} -> 1.0 if V1 == 1.0 or V2 == 1.0, else 0.0
	//   - {... , GreaterThan, VN} -> 1.0 if (result from all previous ops) == 1.0 or VN == 1.0, else 0.0
	Or
	// ### Patterns:
	//   - {V1, And, V2} -> 1.0 if V1 == 1.0 or V2 == 1.0, but not both, else 0.0
	//   - {... , GreaterThan, VN} -> 1.0 if (result from all previous ops) == 1.0 or VN == 1.0, but not both, else 0.0
	Xor
	// ### Patterns:
	//   - {V1, BoolFlip} -> 1.0 if V1 == 0.0, else 0.0
	//   - {... , GreaterThan, VN} -> 1.0 if (result from all previous ops) == 0.0, else 0.0
	BoolFlip
	// ### Patterns:
	//   - {V1, Abs} -> math.Abs(V1)
	//   - {... , Ceil} -> math.Abs((result from all previous ops))
	Abs
	MaxParamIndex
)

type Param struct {
	val             float32
	operation_start uint32
	children_start  uint32
}
