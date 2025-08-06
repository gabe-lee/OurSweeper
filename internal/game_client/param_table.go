package game_client

import (
	"github.com/gabe-lee/OurSweeper/paratable"
)

type (
	ParamCalc     = paratable.ParamCalc
	PIdx_U64      = paratable.PIdx_U64
	PIdx_I64      = paratable.PIdx_I64
	PIdx_F64      = paratable.PIdx_F64
	PIdx_U32      = paratable.PIdx_U32
	PIdx_I32      = paratable.PIdx_I32
	PIdx_F32      = paratable.PIdx_F32
	PIdx_U16      = paratable.PIdx_U16
	PIdx_I16      = paratable.PIdx_I16
	PIdx_U8       = paratable.PIdx_U8
	PIdx_I8       = paratable.PIdx_I8
	PIdx_Bool     = paratable.PIdx_Bool
	PIdx_Calc     = paratable.PIdx_Calc
	ParamTable    = paratable.ParamTable
	CalcInterface = paratable.CalcInterface
)

const (
	FIRST_U64_PARAM PIdx_U64 = PIdx_U64(iota)
	// ... more uint64 param indexes
	_U64_PARAMS_END
)

const (
	FIRST_I64_PARAM PIdx_I64 = PIdx_I64(iota + _U64_PARAMS_END)
	// ... more uint64 param indexes
	_I64_PARAMS_END
)

const (
	FIRST_F64_PARAM PIdx_F64 = PIdx_F64(iota + _I64_PARAMS_END)
	// ... more float64 param indexes
	_F64_PARAMS_END
)

const (
	FIRST_U32_PARAM PIdx_U32 = PIdx_U32(iota + _F64_PARAMS_END)
	// ... more uint32 param indexes
	_U32_PARAMS_END
)

const (
	FIRST_I32_PARAM PIdx_I32 = PIdx_I32(iota + _U32_PARAMS_END)
	// ... more int32 param indexes
	_I32_PARAMS_END
)

const (
	FIRST_F32_PARAM PIdx_F32 = PIdx_F32(iota + _I32_PARAMS_END)
	// ... more float32 param indexes
	_F32_PARAMS_END
)

const (
	FIRST_U16_PARAM PIdx_U16 = PIdx_U16(iota + _F32_PARAMS_END)
	// ... more uint16 param indexes
	_U16_PARAMS_END
)

const (
	FIRST_I16_PARAM PIdx_I16 = PIdx_I16(iota + _U16_PARAMS_END)
	// ... more uint16 param indexes
	_I16_PARAMS_END
)

const (
	FIRST_U8_PARAM PIdx_U8 = PIdx_U8(iota + _I16_PARAMS_END)
	// ... more uint8 param indexes
	_U8_PARAMS_END
)

const (
	FIRST_I8_PARAM PIdx_I8 = PIdx_I8(iota + _U8_PARAMS_END)
	// ... more uint8 param indexes
	_I8_PARAMS_END
)

const (
	FIRST_BOOL_PARAM PIdx_Bool = PIdx_Bool(iota + _I8_PARAMS_END)
	// ... more bool param indexes
	_BOOL_PARAMS_END
)

const (
	_CALC_INVALID PIdx_Calc = PIdx_Calc(iota) // recommended to leave idx 0 as a nil func for debug purposes
	_FIRST_CALC
	// ... more calculation indexes
	_CALC_COUNT
)

const (
	_INPUT_FIRST_CALC_1 uint16 = iota
	_INPUT_FIRST_CALC_2
)

var _FIRST_CALC_1 = [...]uint16{
	_INPUT_FIRST_CALC_1: uint16(FIRST_U64_PARAM),
	_INPUT_FIRST_CALC_2: uint16(FIRST_I64_PARAM),
}

var MyParamTable = paratable.NewParamTable(_U64_PARAMS_END, _I64_PARAMS_END, _F64_PARAMS_END, _U32_PARAMS_END, _I32_PARAMS_END, _F32_PARAMS_END, _U16_PARAMS_END, _I16_PARAMS_END, _U8_PARAMS_END, _I8_PARAMS_END, _BOOL_PARAMS_END, _CALC_COUNT)

func InitMyParamTable() {
	// Register all calculations first
	MyParamTable.RegisterCalc(_FIRST_CALC, func(t *CalcInterface) {
		val1 := t.GetInput_U64(_INPUT_FIRST_CALC_1) // first input
		val2 := t.GetInput_I64(_INPUT_FIRST_CALC_2) // second input
		val3 := int32(val1) + int32(val2)
		t.SetOutput_I32(val3)
	})
	// Init root values
	MyParamTable.SetRoot_U64(FIRST_U64_PARAM, 1)
	MyParamTable.SetRoot_I64(FIRST_I64_PARAM, -2)
	// Init derived values
	MyParamTable.InitDerived_I32(FIRST_I32_PARAM, _FIRST_CALC, _FIRST_CALC_1[:])
}
