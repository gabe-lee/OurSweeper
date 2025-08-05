package game_client

import (
	para "github.com/gabe-lee/OurSweeper/parameter_table"
)

const (
	FIRST_ROOT_U64_PARAM uint16 = iota
	// ... more uint64 param indexes
	_ROOT_U64_PARAMS_END
)

const (
	FIRST_ROOT_I64_PARAM uint16 = iota + _ROOT_U64_PARAMS_END
	// ... more uint64 param indexes
	_ROOT_I64_PARAMS_END
)

const (
	FIRST_ROOT_F64_PARAM uint16 = iota + _ROOT_I64_PARAMS_END
	// ... more float64 param indexes
	_ROOT_F64_PARAMS_END
)

const (
	FIRST_DERIVED_U64_PARAM uint16 = iota + _ROOT_F64_PARAMS_END
	// ... more uint64 param indexes
	_DERIVED_U64_PARAMS_END
)

const (
	FIRST_DERIVED_I64_PARAM uint16 = iota + _DERIVED_U64_PARAMS_END
	// ... more uint64 param indexes
	_DERIVED_I64_PARAMS_END
)

const (
	FIRST_DERIVED_F64_PARAM uint16 = iota + _DERIVED_I64_PARAMS_END
	// ... more float64 param indexes
	_DERIVED_F64_PARAMS_END
)

const (
	FIRST_ROOT_U32_PARAM uint16 = iota + _DERIVED_F64_PARAMS_END
	// ... more uint32 param indexes
	_ROOT_U32_PARAMS_END
)

const (
	FIRST_ROOT_I32_PARAM uint16 = iota + _ROOT_U32_PARAMS_END
	// ... more int32 param indexes
	_ROOT_I32_PARAMS_END
)

const (
	FIRST_ROOT_F32_PARAM uint16 = iota + _ROOT_I32_PARAMS_END
	// ... more float32 param indexes
	_ROOT_F32_PARAMS_END
)

const (
	FIRST_DERIVED_U32_PARAM uint16 = iota + _ROOT_F32_PARAMS_END
	// ... more uint32 param indexes
	_DERIVED_U32_PARAMS_END
)

const (
	FIRST_DERIVED_I32_PARAM uint16 = iota + _DERIVED_U32_PARAMS_END
	// ... more int32 param indexes
	_DERIVED_I32_PARAMS_END
)

const (
	FIRST_DERIVED_F32_PARAM uint16 = iota + _DERIVED_I32_PARAMS_END
	// ... more float32 param indexes
	_DERIVED_F32_PARAMS_END
)

const (
	FIRST_ROOT_U16_PARAM uint16 = iota + _DERIVED_F32_PARAMS_END
	// ... more uint16 param indexes
	_ROOT_U16_PARAMS_END
)

const (
	FIRST_ROOT_I16_PARAM uint16 = iota + _ROOT_U16_PARAMS_END
	// ... more uint16 param indexes
	_ROOT_I16_PARAMS_END
)

const (
	FIRST_DERIVED_U16_PARAM uint16 = iota + _ROOT_I16_PARAMS_END
	// ... more uint16 param indexes
	_DERIVED_U16_PARAMS_END
)

const (
	FIRST_DERIVED_I16_PARAM uint16 = iota + _DERIVED_U16_PARAMS_END
	// ... more uint16 param indexes
	_DERIVED_I16_PARAMS_END
)

const (
	FIRST_ROOT_U8_PARAM uint16 = iota + _DERIVED_I16_PARAMS_END
	// ... more uint8 param indexes
	_ROOT_U8_PARAMS_END
)

const (
	FIRST_ROOT_I8_PARAM uint16 = iota + _ROOT_U8_PARAMS_END
	// ... more uint8 param indexes
	_ROOT_I8_PARAMS_END
)

const (
	FIRST_ROOT_BOOL_PARAM uint16 = iota + _ROOT_I8_PARAMS_END
	// ... more bool param indexes
	_ROOT_BOOL_PARAMS_END
)

const (
	FIRST_DERIVED_U8_PARAM uint16 = iota + _ROOT_BOOL_PARAMS_END
	// ... more uint8 param indexes
	_DERIVED_U8_PARAMS_END
)

const (
	FIRST_DERIVED_I8_PARAM uint16 = iota + _DERIVED_U8_PARAMS_END
	// ... more uint8 param indexes
	_DERIVED_I8_PARAMS_END
)

const (
	FIRST_DERIVED_BOOL_PARAM uint16 = iota + _DERIVED_I8_PARAMS_END
	// ... more bool param indexes
	_DERIVED_BOOL_PARAMS_END
)

var MyParamTable = para.NewParamTable(_ROOT_U64_PARAMS_END, _ROOT_I64_PARAMS_END, _ROOT_F64_PARAMS_END, _DERIVED_U64_PARAMS_END, _DERIVED_I64_PARAMS_END, _DERIVED_F64_PARAMS_END, _ROOT_U32_PARAMS_END, _ROOT_I32_PARAMS_END, _ROOT_F32_PARAMS_END, _DERIVED_U32_PARAMS_END, _DERIVED_I32_PARAMS_END, _DERIVED_F32_PARAMS_END, _ROOT_U16_PARAMS_END, _ROOT_I16_PARAMS_END, _DERIVED_U16_PARAMS_END, _DERIVED_I16_PARAMS_END, _ROOT_U8_PARAMS_END, _ROOT_I8_PARAMS_END, _ROOT_BOOL_PARAMS_END, _DERIVED_U8_PARAMS_END, _DERIVED_I8_PARAMS_END, _DERIVED_BOOL_PARAMS_END)

func InitMyParamTable() {
	MyParamTable.Set_U64(FIRST_ROOT_U64_PARAM, 600)    // a radius
	MyParamTable.Set_F32(FIRST_ROOT_F32_PARAM, 3.1415) // pi
	MyParamTable.InitDerived_I32(FIRST_DERIVED_I32_PARAM, []uint16{FIRST_ROOT_U64_PARAM, FIRST_ROOT_F32_PARAM}, func(table *para.ParamTable, oldVal int32) int32 {
		// Area of a circle with radius 600, as an integer
		radius := table.Get_U64(FIRST_ROOT_U64_PARAM)
		pi := table.Get_F32(FIRST_ROOT_F32_PARAM)
		return int32(float32(radius*radius) * pi)
	})
}
