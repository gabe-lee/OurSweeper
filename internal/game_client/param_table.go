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
	_U64_PARAMS_END PIdx_U64 = PIdx_U64(iota)
	// ... more uint64 param indexes
)

const (
	_I64_PARAMS_END PIdx_I64 = PIdx_I64(iota + _U64_PARAMS_END)
	// ... more uint64 param indexes
)

const (
	_F64_PARAMS_END PIdx_F64 = PIdx_F64(iota + _I64_PARAMS_END)
	// ... more float64 param indexes
)

const (
	_U32_PARAMS_END PIdx_U32 = PIdx_U32(iota + _F64_PARAMS_END)
	// ... more uint32 param indexes
)

const (
	WINDOW_WIDTH PIdx_I32 = PIdx_I32(iota + _U32_PARAMS_END)
	WINDOW_HEIGHT
	// ... more int32 param indexes
	_I32_PARAMS_END
)

const (
	_F32_PARAMS_END PIdx_F32 = PIdx_F32(iota + _I32_PARAMS_END)
	// ... more float32 param indexes
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
	_CALC_CENTER_OF_RECT
	_CALC_CENTER_OF_WINDOW
	_CALC_TEXT_SIZE
	// ... more calculation indexes
	_CALC_COUNT
)

const (
	_INPUT_CENTER_OF_WINDOW_WIDTH uint16 = iota
	_INPUT_CENTER_OF_WINDOW_HEIGHT
)

const (
	_INPUT_CENTER_OF_RECT_TOP_LEFT uint16 = iota
	_INPUT_CENTER_OF_RECT_TOP_RIGHT
	_INPUT_CENTER_OF_RECT_BOT_LEFT
	_INPUT_CENTER_OF_RECT_BOT_RIGHT
)

var _CENTER_OF_WINDOW_INPUTS = [...]uint16{
	_INPUT_CENTER_OF_WINDOW_WIDTH:  uint16(WINDOW_WIDTH),
	_INPUT_CENTER_OF_WINDOW_HEIGHT: uint16(WINDOW_HEIGHT),
}

const (
	_IN_TEXT_SIZE_SIZE uint16 = iota
	_IN_TEXT_SIZE_LANG
	_IN_TEXT_SIZE_STRING
)

var MyParamTable = paratable.NewParamTable(_U64_PARAMS_END, _I64_PARAMS_END, _F64_PARAMS_END, _U32_PARAMS_END, _I32_PARAMS_END, _F32_PARAMS_END, _U16_PARAMS_END, _I16_PARAMS_END, _U8_PARAMS_END, _I8_PARAMS_END, _BOOL_PARAMS_END, _CALC_COUNT)

func InitMyParamTable() {
	// Register all calculations first
	MyParamTable.RegisterCalc(_CALC_TEXT_SIZE, func(t *CalcInterface) {
		size := t.GetInput_F32(_IN_TEXT_SIZE_SIZE)
		lang := t.GetInput_U8(_IN_TEXT_SIZE_LANG)
		textIdx := t.GetInput_U16(_IN_TEXT_SIZE_STRING)
		w, h := GetUiTextSize(size, lang, textIdx)
		t.SetOutput_F32(val3)
	})
	// Init root values
	MyParamTable.SetRoot_U64(FIRST_U64_PARAM, 1)
	MyParamTable.SetRoot_I64(FIRST_I64_PARAM, -2)
	// Init derived values
	MyParamTable.InitDerived_I32(FIRST_I32_PARAM, _FIRST_CALC, _FIRST_CALC_1[:])
}
