package game_client

import (
	"math"

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
	// ... more int32 param indexes
	_I32_PARAMS_END = PIdx_I32(iota + _U32_PARAMS_END)
)

const (
	F32_0_00 = PIdx_F32(iota + _I32_PARAMS_END)
	F32_0_50
	WINDOW_WIDTH
	WINDOW_HEIGHT
	WINDOW_CENTER_X
	WINDOW_CENTER_Y
	LOGIN_WINDOW_WIDTH
	LOGIN_WINDOW_HEIGHT
	LOGIN_WINDOW_NUM_ROWS
	LOGIN_WINDOW_WIDEST_ROW
	LOGIN_WINDOW_PADDING
	LOGIN_WINDOW_ROW_HEIGHT_SUM
	LOGIN_WINDOW_X
	LOGIN_WINDOW_Y
	_F32_PARAMS_END
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
	_CALC_LERP_F32 = PIdx_Calc(iota)
	_CALC_MAX_F32
	_CALC_MIN_F32
	_CALC_SUM_F32
	_CALC_SUM_PLUS_MARGINS_F32
	_CALC_TEXT_SIZE
	_CALC_CENTER_LINE_AROUND_POINT
	// ... more calculation indexes
	_CALC_COUNT
)

const (
	_IN_CENTER_LINE_AROUND_POINT_CEN uint16 = iota
	_IN_CENTER_LINE_AROUND_POINT_LEN
)
const (
	_OUT_CENTER_LINE_AROUND_POINT_VAL uint16 = iota
)

var _INS_CENTER_LOGIN_PANEL_X_POS = [...]uint16{
	_IN_CENTER_LINE_AROUND_POINT_CEN: uint16(WINDOW_CENTER_X),
	_IN_CENTER_LINE_AROUND_POINT_LEN: uint16(LOGIN_WINDOW_WIDTH),
}
var _OUTS_CENTER_LOGIN_PANEL_X_POS = [...]uint16{
	_OUT_CENTER_LINE_AROUND_POINT_VAL: uint16(LOGIN_WINDOW_X),
}
var _INS_CENTER_LOGIN_PANEL_Y_POS = [...]uint16{
	_IN_CENTER_LINE_AROUND_POINT_CEN: uint16(WINDOW_CENTER_Y),
	_IN_CENTER_LINE_AROUND_POINT_LEN: uint16(LOGIN_WINDOW_HEIGHT),
}
var _OUTS_CENTER_LOGIN_PANEL_Y_POS = [...]uint16{
	_OUT_CENTER_LINE_AROUND_POINT_VAL: uint16(LOGIN_WINDOW_Y),
}

const (
	_IN_SUM_PLUS_MARGINS_F32_MARGIN uint16 = iota
	_IN_SUM_PLUS_MARGINS_F32_FIRST_VAL
)

const (
	_IN_LERP_F32_A uint16 = iota
	_IN_LERP_F32_B
	_IN_LERP_F32_DELTA
)
const (
	_OUT_LERP_F32 uint16 = iota
)

var _INS_WINDOW_CENTER_X = [...]uint16{
	_IN_LERP_F32_A:     uint16(F32_0_00),
	_IN_LERP_F32_B:     uint16(WINDOW_WIDTH),
	_IN_LERP_F32_DELTA: uint16(F32_0_50),
}
var _OUTS_WINDOW_CENTER_X = [...]uint16{
	_OUT_LERP_F32: uint16(WINDOW_CENTER_X),
}

var _INS_WINDOW_CENTER_Y = [...]uint16{
	_IN_LERP_F32_A:     uint16(F32_0_00),
	_IN_LERP_F32_B:     uint16(WINDOW_HEIGHT),
	_IN_LERP_F32_DELTA: uint16(F32_0_50),
}
var _OUTS_WINDOW_CENTER_Y = [...]uint16{
	_OUT_LERP_F32: uint16(WINDOW_CENTER_Y),
}

const (
	_IN_TEXT_SIZE_SIZE uint16 = iota
	_IN_TEXT_SIZE_LANG
	_IN_TEXT_SIZE_STRING
)

const (
	_OUT_TEXT_SIZE_X uint16 = iota
	_OUT_TEXT_SIZE_Y
)

func InitParamTable(table *ParamTable) {
	*table = paratable.NewParamTable(_U64_PARAMS_END, _I64_PARAMS_END, _F64_PARAMS_END, _U32_PARAMS_END, _I32_PARAMS_END, _F32_PARAMS_END, _U16_PARAMS_END, _I16_PARAMS_END, _U8_PARAMS_END, _I8_PARAMS_END, _BOOL_PARAMS_END, _CALC_COUNT)
	// Register all calculations first
	table.RegisterCalc(_CALC_TEXT_SIZE, func(t *CalcInterface) {
		size := t.GetInput_F32(_IN_TEXT_SIZE_SIZE)
		lang := t.GetInput_U8(_IN_TEXT_SIZE_LANG)
		textIdx := t.GetInput_U16(_IN_TEXT_SIZE_STRING)
		w, h := GetUiTextSize(size, lang, textIdx)
		t.SetOutput_F32(_OUT_TEXT_SIZE_X, w)
		t.SetOutput_F32(_OUT_TEXT_SIZE_Y, h)
	})
	table.RegisterCalc(_CALC_LERP_F32, func(t *CalcInterface) {
		a := t.GetInput_F32(_IN_LERP_F32_A)
		b := t.GetInput_F32(_IN_LERP_F32_B)
		delta := t.GetInput_F32(_IN_LERP_F32_DELTA)
		val := a + ((b - a) * delta)
		t.SetOutput_F32(_OUT_LERP_F32, val)
	})
	table.RegisterCalc(_CALC_MAX_F32, func(t *CalcInterface) {
		var val float32 = -math.MaxFloat32
		for _, idx := range t.GetAllInputs() {
			val = max(val, t.GetInput_F32(idx))
		}
		t.SetOutput_F32(0, val)
	})
	table.RegisterCalc(_CALC_MIN_F32, func(t *CalcInterface) {
		var val float32 = math.MaxFloat32
		for _, idx := range t.GetAllInputs() {
			val = min(val, t.GetInput_F32(idx))
		}
		t.SetOutput_F32(0, val)
	})
	table.RegisterCalc(_CALC_SUM_F32, func(t *CalcInterface) {
		var val float32 = 0
		for _, idx := range t.GetAllInputs() {
			val += t.GetInput_F32(idx)
		}
		t.SetOutput_F32(0, val)
	})
	table.RegisterCalc(_CALC_SUM_PLUS_MARGINS_F32, func(t *CalcInterface) {
		margin := t.GetInput_F32(_IN_SUM_PLUS_MARGINS_F32_MARGIN)
		var val float32 = margin
		for _, idx := range t.GetInputRangeStart(_IN_SUM_PLUS_MARGINS_F32_FIRST_VAL) {
			val += margin
			val += t.GetInput_F32(idx)
		}
		t.SetOutput_F32(0, val)
	})
	table.RegisterCalc(_CALC_CENTER_LINE_AROUND_POINT, func(t *CalcInterface) {
		cen := t.GetInput_F32(_IN_CENTER_LINE_AROUND_POINT_CEN)
		length := t.GetInput_F32(_IN_CENTER_LINE_AROUND_POINT_LEN)
		halfLen := length / 2.0
		val := cen - halfLen
		t.SetOutput_F32(0, val)
	})
	// Init root values
	table.SetRoot_F32(F32_0_00, 0.0)
	table.SetRoot_F32(F32_0_50, 0.5)
	table.SetRoot_F32(WINDOW_WIDTH, 800.0)
	table.SetRoot_F32(WINDOW_HEIGHT, 600.0)
	table.SetRoot_F32(LOGIN_WINDOW_WIDTH, 400.0)
	table.SetRoot_F32(LOGIN_WINDOW_HEIGHT, 300.0)
	// Init derived values
	table.InitDerived_F32(WINDOW_CENTER_X, _CALC_LERP_F32, _INS_WINDOW_CENTER_X[:], _OUTS_WINDOW_CENTER_X[:])
	table.InitDerived_F32(WINDOW_CENTER_Y, _CALC_LERP_F32, _INS_WINDOW_CENTER_Y[:], _OUTS_WINDOW_CENTER_Y[:])
	table.InitDerived_F32(LOGIN_WINDOW_X, _CALC_CENTER_LINE_AROUND_POINT, _INS_CENTER_LOGIN_PANEL_X_POS[:], _OUTS_CENTER_LOGIN_PANEL_X_POS[:])
	table.InitDerived_F32(LOGIN_WINDOW_Y, _CALC_CENTER_LINE_AROUND_POINT, _INS_CENTER_LOGIN_PANEL_Y_POS[:], _OUTS_CENTER_LOGIN_PANEL_Y_POS[:])
}
