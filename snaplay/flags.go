package snaplay

import "math"

type Snap_Flag uint32

const (
	__HIDDEN_MASK Snap_Flag = 0b1
	__HIDDEN_BITS Snap_Flag = 1

	__UPDATED_TOP_LEFT   Snap_Flag = 0b000000001
	__UPDATED_TOP_CENTER Snap_Flag = 0b000000010
	__UPDATED_TOP_RIGHT  Snap_Flag = 0b000000100
	__UPDATED_MID_LEFT   Snap_Flag = 0b000001000
	__UPDATED_MID_CENTER Snap_Flag = 0b000010000
	__UPDATED_MID_RIGHT  Snap_Flag = 0b000100000
	__UPDATED_BOT_LEFT   Snap_Flag = 0b001000000
	__UPDATED_BOT_CENTER Snap_Flag = 0b010000000
	__UPDATED_BOT_RIGHT  Snap_Flag = 0b100000000
	__UPDATED_MASK       Snap_Flag = 0b111111111
	__UPDATED_BITS       Snap_Flag = 9

	__PARAM_ABSOLUTE Snap_Flag = 0b0
	__PARAM_RELATIVE Snap_Flag = 0b1
	__PARAM_MASK     Snap_Flag = 0b1
	__PARAM_BITS     Snap_Flag = 1

	__ANCHOR_TOP_LEFT   Snap_Flag = 0b0000
	__ANCHOR_TOP_CENTER Snap_Flag = 0b0001
	__ANCHOR_TOP_RIGHT  Snap_Flag = 0b0010
	__ANCHOR_MID_LEFT   Snap_Flag = 0b0011
	__ANCHOR_MID_CENTER Snap_Flag = 0b0100
	__ANCHOR_MID_RIGHT  Snap_Flag = 0b0101
	__ANCHOR_BOT_LEFT   Snap_Flag = 0b0110
	__ANCHOR_BOT_CENTER Snap_Flag = 0b0111
	__ANCHOR_BOT_RIGHT  Snap_Flag = 0b1000
	__ANCHOR_MASK       Snap_Flag = 0b1111
	__ANCHOR_BITS       Snap_Flag = 4

	//FLAG LAYOUT

	_HIDDEN_SHIFT      Snap_Flag = 0
	_UPDATED_SHIFT     Snap_Flag = _HIDDEN_SHIFT + __HIDDEN_BITS
	_EXT_ANCHOR_SHIFT  Snap_Flag = _UPDATED_SHIFT + __UPDATED_BITS
	_SELF_ANCHOR_SHIFT Snap_Flag = _EXT_ANCHOR_SHIFT + __ANCHOR_BITS

	_X_SIZE_SHIFT                     = _SELF_ANCHOR_SHIFT + __ANCHOR_BITS
	_X_SIZE_ABSOLUTE        Snap_Flag = __PARAM_ABSOLUTE << _X_SIZE_SHIFT
	_X_SIZE_RELATIVE_ANCHOR Snap_Flag = __PARAM_RELATIVE << _X_SIZE_SHIFT
	_X_SIZE_MASK            Snap_Flag = __PARAM_MASK << _X_SIZE_SHIFT

	_Y_SIZE_SHIFT                     = _X_SIZE_SHIFT + __PARAM_BITS
	_Y_SIZE_ABSOLUTE        Snap_Flag = __PARAM_ABSOLUTE << _Y_SIZE_SHIFT
	_Y_SIZE_RELATIVE_ANCHOR Snap_Flag = __PARAM_RELATIVE << _Y_SIZE_SHIFT
	_Y_SIZE_MASK            Snap_Flag = __PARAM_MASK << _Y_SIZE_SHIFT

	_X_OFF_SHIFT                     = _Y_SIZE_SHIFT + __PARAM_BITS
	_X_OFF_ABSOLUTE        Snap_Flag = __PARAM_ABSOLUTE << _X_OFF_SHIFT
	_X_OFF_RELATIVE_ANCHOR Snap_Flag = __PARAM_RELATIVE << _X_OFF_SHIFT
	_X_OFF_MASK            Snap_Flag = __PARAM_MASK << _X_OFF_SHIFT

	_Y_OFF_SHIFT                     = _X_OFF_SHIFT + __PARAM_BITS
	_Y_OFF_ABSOLUTE        Snap_Flag = __PARAM_ABSOLUTE << _Y_OFF_SHIFT
	_Y_OFF_RELATIVE_ANCHOR Snap_Flag = __PARAM_RELATIVE << _Y_OFF_SHIFT
	_Y_OFF_MASK            Snap_Flag = __PARAM_MASK << _Y_OFF_SHIFT

	_FIT_X_SHIFT                = _Y_OFF_SHIFT + __PARAM_BITS
	_FIT_X_EXACT      Snap_Flag = __EXACT << _FIT_X_SHIFT
	_FIT_X_MIN_GROW   Snap_Flag = __MIN_GROW << _FIT_X_SHIFT
	_FIT_X_MAX_SHRINK Snap_Flag = __MAX_SHRINK << _FIT_X_SHIFT
	_FIT_X_MASK       Snap_Flag = __FIT_MASK << _FIT_X_SHIFT

	_FIT_Y_SHIFT                = _FIT_X_SHIFT + __FIT_BITS
	_FIT_Y_EXACT      Snap_Flag = __EXACT << _FIT_Y_SHIFT
	_FIT_Y_MIN_GROW   Snap_Flag = __MIN_GROW << _FIT_Y_SHIFT
	_FIT_Y_MAX_SHRINK Snap_Flag = __MAX_SHRINK << _FIT_Y_SHIFT
	_FIT_Y_MASK       Snap_Flag = __FIT_MASK << _FIT_Y_SHIFT

	_ANCHOR_IDX_SHIFT              = _FIT_Y_SHIFT + __FIT_BITS
	_ANCHOR_IDX_PARENT   Snap_Flag = __ANCHOR_IDX_PARENT << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_ABSOLUTE Snap_Flag = __ANCHOR_IDX_ABSOLUTE << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_RELATIVE Snap_Flag = __ANCHOR_IDX_RELATIVE << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_MASK     Snap_Flag = __ANCHOR_IDX_MASK << _ANCHOR_IDX_SHIFT

	_PLANE_SHIFT           = _ANCHOR_IDX_SHIFT + __ANCHOR_IDX_BITS
	_PLANE_0     Snap_Flag = __PLANE_0 << _PLANE_SHIFT
	_PLANE_1     Snap_Flag = __PLANE_1 << _PLANE_SHIFT
	_PLANE_2     Snap_Flag = __PLANE_2 << _PLANE_SHIFT
	_PLANE_3     Snap_Flag = __PLANE_3 << _PLANE_SHIFT
	_PLANE_4     Snap_Flag = __PLANE_4 << _PLANE_SHIFT
	_PLANE_5     Snap_Flag = __PLANE_5 << _PLANE_SHIFT
	_PLANE_6     Snap_Flag = __PLANE_6 << _PLANE_SHIFT
	_PLANE_7     Snap_Flag = __PLANE_7 << _PLANE_SHIFT
	_PLANE_MASK  Snap_Flag = __PLANE_MASK << _PLANE_SHIFT
)

type UI_Anchor uint32

const (
	TOP_LEFT   UI_Anchor = UI_Anchor(__TOP_LEFT)
	TOP_CENTER UI_Anchor = UI_Anchor(__TOP_CENTER)
	TOP_RIGHT  UI_Anchor = UI_Anchor(__TOP_RIGHT)
	MID_LEFT   UI_Anchor = UI_Anchor(__MID_LEFT)
	MID_CENTER UI_Anchor = UI_Anchor(__MID_CENTER)
	MID_RIGHT  UI_Anchor = UI_Anchor(__MID_RIGHT)
	BOT_LEFT   UI_Anchor = UI_Anchor(__BOT_LEFT)
	BOT_CENTER UI_Anchor = UI_Anchor(__BOT_CENTER)
	BOT_RIGHT  UI_Anchor = UI_Anchor(__BOT_RIGHT)
)

type UI_DimMode uint32

const (
	ABSOLUTE             UI_DimMode = UI_DimMode(__PARAM_ABSOLUTE)
	RELATIVE_ANCHOR_SIZE UI_DimMode = UI_DimMode(__PARAM_RELATIVE)
)

type UI_Fit uint32

const (
	EXACT      UI_Fit = UI_Fit(__EXACT)
	MIN_GROW   UI_Fit = UI_Fit(__MIN_GROW)
	MAX_SHRINK UI_Fit = UI_Fit(__MAX_SHRINK)
)

type UI_AnchorMode uint32

const (
	ANCHOR_TO_PARENT           UI_AnchorMode = UI_AnchorMode(__ANCHOR_IDX_PARENT)
	ANCHOR_TO_REF_IDX_ABSOLUTE UI_AnchorMode = UI_AnchorMode(__ANCHOR_IDX_ABSOLUTE)
	ANCHOR_TO_REF_IDX_RELATIVE UI_AnchorMode = UI_AnchorMode(__ANCHOR_IDX_RELATIVE)
)

// type UI_Plane uint32

// const (
// 	PLANE_0 UI_Plane = UI_Plane(__PLANE_0)
// 	PLANE_1 UI_Plane = UI_Plane(__PLANE_1)
// 	PLANE_2 UI_Plane = UI_Plane(__PLANE_2)
// 	PLANE_3 UI_Plane = UI_Plane(__PLANE_3)
// 	PLANE_4 UI_Plane = UI_Plane(__PLANE_4)
// 	PLANE_5 UI_Plane = UI_Plane(__PLANE_5)
// 	PLANE_6 UI_Plane = UI_Plane(__PLANE_6)
// 	PLANE_7 UI_Plane = UI_Plane(__PLANE_7)
// )

var FLAG = Snap_Flag(0)

func (f *Snap_Flag) SetHidden() {
	*f |= _IS_HIDDEN
}
func (f *Snap_Flag) SetVisible() {
	*f &= ^_IS_HIDDEN
}
func (f Snap_Flag) IsHidden() bool {
	return f&_IS_HIDDEN == _IS_HIDDEN
}
func (f Snap_Flag) IsVisible() bool {
	return f&_IS_HIDDEN == 0
}

func (f Snap_Flag) XOffsetMode(dim UI_DimMode) Snap_Flag {
	f &= ^_X_OFF_MASK
	return f | Snap_Flag(dim<<_X_OFF_SHIFT)
}
func (f Snap_Flag) GetXOffsetMode() UI_DimMode {
	return UI_DimMode((f & _X_OFF_MASK) >> _X_OFF_SHIFT)
}
func (f Snap_Flag) YOffsetMode(dim UI_DimMode) Snap_Flag {
	f &= ^_Y_OFF_MASK
	return f | Snap_Flag(dim<<_Y_OFF_SHIFT)
}
func (f Snap_Flag) GetYOffsetMode() UI_DimMode {
	return UI_DimMode((f & _Y_OFF_MASK) >> _Y_OFF_SHIFT)
}
func (f Snap_Flag) XYOffset(dim UI_DimMode) Snap_Flag {
	f &= ^_X_OFF_MASK
	f &= ^_Y_OFF_MASK
	f |= Snap_Flag(dim << _X_OFF_SHIFT)
	return f | Snap_Flag(dim<<_Y_OFF_SHIFT)
}

func (f Snap_Flag) XSizeMode(dim UI_DimMode) Snap_Flag {
	f &= ^_X_SIZE_MASK
	return f | Snap_Flag(dim<<_X_SIZE_SHIFT)
}
func (f Snap_Flag) GetXSizeMode() UI_DimMode {
	return UI_DimMode((f & _X_SIZE_MASK) >> _X_SIZE_SHIFT)
}
func (f Snap_Flag) YSizeMode(dim UI_DimMode) Snap_Flag {
	f &= ^_Y_SIZE_MASK
	return f | Snap_Flag(dim<<_Y_SIZE_SHIFT)
}
func (f Snap_Flag) GetYSizeMode() UI_DimMode {
	return UI_DimMode((f & _Y_SIZE_MASK) >> _Y_SIZE_SHIFT)
}
func (f Snap_Flag) XYSizeMode(dim UI_DimMode) Snap_Flag {
	f &= ^_X_SIZE_MASK
	f &= ^_Y_SIZE_MASK
	f |= Snap_Flag(dim << _X_SIZE_SHIFT)
	return f | Snap_Flag(dim<<_Y_SIZE_SHIFT)
}

func (f Snap_Flag) SelfAnchor(anchor UI_Anchor) Snap_Flag {
	f &= ^_SELF_ANCHOR_MASK
	return f | Snap_Flag(anchor<<_SELF_ANCHOR_SHIFT)
}
func (f Snap_Flag) GetSelfAnchor() UI_Anchor {
	return UI_Anchor((f & _SELF_ANCHOR_MASK) >> _SELF_ANCHOR_SHIFT)
}

func (f Snap_Flag) ExternalAnchor(anchor UI_Anchor, mode UI_AnchorMode) Snap_Flag {
	f &= ^_EXT_ANCHOR_MASK
	f &= ^_ANCHOR_IDX_MASK
	f |= Snap_Flag(mode << _ANCHOR_IDX_SHIFT)
	return f | Snap_Flag(anchor<<_EXT_ANCHOR_SHIFT)
}
func (f Snap_Flag) GetExternalAnchor() UI_Anchor {
	return UI_Anchor((f & _EXT_ANCHOR_MASK) >> _EXT_ANCHOR_SHIFT)
}
func (f Snap_Flag) GetExternalAnchorMode() UI_AnchorMode {
	return UI_AnchorMode((f & _ANCHOR_IDX_MASK) >> _ANCHOR_IDX_SHIFT)
}

func (f Snap_Flag) Anchor(selfAnchor UI_Anchor, externalAnchor UI_Anchor, mode UI_AnchorMode) Snap_Flag {
	f &= ^_EXT_ANCHOR_MASK
	f &= ^_SELF_ANCHOR_MASK
	f &= ^_ANCHOR_IDX_MASK
	f |= Snap_Flag(selfAnchor << _SELF_ANCHOR_SHIFT)
	f |= Snap_Flag(mode << _ANCHOR_IDX_SHIFT)
	return f | Snap_Flag(externalAnchor<<_EXT_ANCHOR_SHIFT)
}

func (f Snap_Flag) XFit(fit UI_Fit) Snap_Flag {
	f &= ^_FIT_X_MASK
	return f | Snap_Flag(fit<<_FIT_X_SHIFT)
}
func (f Snap_Flag) GetXFit() UI_Fit {
	return UI_Fit((f & _FIT_X_MASK) >> _FIT_X_SHIFT)
}
func (f Snap_Flag) YFit(fit UI_Fit) Snap_Flag {
	f &= ^_FIT_Y_MASK
	return f | Snap_Flag(fit<<_FIT_Y_SHIFT)
}
func (f Snap_Flag) GetYFit() UI_Fit {
	return UI_Fit((f & _FIT_Y_MASK) >> _FIT_Y_SHIFT)
}
func (f Snap_Flag) XYFit(fit UI_Fit) Snap_Flag {
	f &= ^_FIT_X_MASK
	f &= ^_FIT_Y_MASK
	f |= Snap_Flag(fit << _FIT_X_SHIFT)
	return f | Snap_Flag(fit<<_FIT_Y_SHIFT)
}

func (f Snap_Flag) Plane(plane UI_Plane) Snap_Flag {
	f &= ^_PLANE_MASK
	return f | Snap_Flag(plane<<_PLANE_SHIFT)
}

func (f Snap_Flag) GetPlane() UI_Plane {
	return UI_Plane((f & _PLANE_MASK) >> _PLANE_SHIFT)
}
func (f Snap_Flag) GetPlaneAsZIndex() Snap_Zindex {
	return Snap_Zindex((f&_PLANE_MASK)>>_PLANE_SHIFT) << (16 - __PLANE_BITS)
}

const (
	_PLANE_ZINDEX_SHIFT        = 16 - __PLANE_BITS
	_PLANE_IN_ZINDEX_MASK      = uint16((math.MaxUint16 << _PLANE_ZINDEX_SHIFT) & math.MaxUint16)
	_ZINDEX_WITHOUT_PLANE_MASK = math.MaxUint16 >> __PLANE_BITS
)
