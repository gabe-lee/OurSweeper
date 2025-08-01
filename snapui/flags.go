package snapui

import "math"

type UI_Flag uint32

const (
	__HIDDEN_MASK UI_Flag = 0b1

	__HIDDEN_BITS = 1

	__DIM_ABSOLUTE        UI_Flag = 0b0
	__DIM_RELATIVE_ANCHOR UI_Flag = 0b1
	__DIM_MASK            UI_Flag = 0b1
	__DIM_BITS                    = 1

	__TOP_LEFT    UI_Flag = 0b0000
	__TOP_CENTER  UI_Flag = 0b0001
	__TOP_RIGHT   UI_Flag = 0b0010
	__MID_LEFT    UI_Flag = 0b0011
	__MID_CENTER  UI_Flag = 0b0100
	__MID_RIGHT   UI_Flag = 0b0101
	__BOT_LEFT    UI_Flag = 0b0110
	__BOT_CENTER  UI_Flag = 0b0111
	__BOT_RIGHT   UI_Flag = 0b1000
	__ANCHOR_MASK UI_Flag = 0b1111
	__ANCHOR_BITS         = 4

	__EXACT      UI_Flag = 0b00
	__MIN_GROW   UI_Flag = 0b01
	__MAX_SHRINK UI_Flag = 0b10
	__FIT_MASK   UI_Flag = 0b11
	__FIT_BITS           = 2

	__ANCHOR_IDX_PARENT   UI_Flag = 0b00
	__ANCHOR_IDX_ABSOLUTE UI_Flag = 0b01
	__ANCHOR_IDX_RELATIVE UI_Flag = 0b10
	__ANCHOR_IDX_MASK     UI_Flag = 0b11
	__ANCHOR_IDX_BITS             = 2

	__PLANE_0    UI_Flag = 0b000
	__PLANE_1    UI_Flag = 0b001
	__PLANE_2    UI_Flag = 0b010
	__PLANE_3    UI_Flag = 0b011
	__PLANE_4    UI_Flag = 0b100
	__PLANE_5    UI_Flag = 0b101
	__PLANE_6    UI_Flag = 0b110
	__PLANE_7    UI_Flag = 0b111
	__PLANE_MASK UI_Flag = 0b111
	__PLANE_BITS         = 3

	//LAYOUT

	_HIDDEN_SHIFT         = 0
	_IS_HIDDEN    UI_Flag = __HIDDEN_MASK << _HIDDEN_SHIFT

	_EXT_ANCHOR_SHIFT              = _HIDDEN_SHIFT + __HIDDEN_BITS
	_EXT_ANCHOR_TOP_LEFT   UI_Flag = __TOP_LEFT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_TOP_CENTER UI_Flag = __TOP_CENTER << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_TOP_RIGHT  UI_Flag = __TOP_RIGHT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_MID_LEFT   UI_Flag = __MID_LEFT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_MID_CENTER UI_Flag = __MID_CENTER << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_MID_RIGHT  UI_Flag = __MID_RIGHT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_BOT_LEFT   UI_Flag = __BOT_LEFT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_BOT_CENTER UI_Flag = __BOT_CENTER << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_BOT_RIGHT  UI_Flag = __BOT_RIGHT << _EXT_ANCHOR_SHIFT
	_EXT_ANCHOR_MASK       UI_Flag = __ANCHOR_MASK << _EXT_ANCHOR_SHIFT

	_SELF_ANCHOR_SHIFT              = _EXT_ANCHOR_SHIFT + __ANCHOR_BITS
	_SELF_ANCHOR_TOP_LEFT   UI_Flag = __TOP_LEFT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_TOP_CENTER UI_Flag = __TOP_CENTER << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_TOP_RIGHT  UI_Flag = __TOP_RIGHT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_MID_LEFT   UI_Flag = __MID_LEFT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_MID_CENTER UI_Flag = __MID_CENTER << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_MID_RIGHT  UI_Flag = __MID_RIGHT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_BOT_LEFT   UI_Flag = __BOT_LEFT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_BOT_CENTER UI_Flag = __BOT_CENTER << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_BOT_RIGHT  UI_Flag = __BOT_RIGHT << _SELF_ANCHOR_SHIFT
	_SELF_ANCHOR_MASK       UI_Flag = __ANCHOR_MASK << _SELF_ANCHOR_SHIFT

	_X_SIZE_SHIFT                   = _SELF_ANCHOR_SHIFT + __ANCHOR_BITS
	_X_SIZE_ABSOLUTE        UI_Flag = __DIM_ABSOLUTE << _X_SIZE_SHIFT
	_X_SIZE_RELATIVE_ANCHOR UI_Flag = __DIM_RELATIVE_ANCHOR << _X_SIZE_SHIFT
	_X_SIZE_MASK            UI_Flag = __DIM_MASK << _X_SIZE_SHIFT

	_Y_SIZE_SHIFT                   = _X_SIZE_SHIFT + __DIM_BITS
	_Y_SIZE_ABSOLUTE        UI_Flag = __DIM_ABSOLUTE << _Y_SIZE_SHIFT
	_Y_SIZE_RELATIVE_ANCHOR UI_Flag = __DIM_RELATIVE_ANCHOR << _Y_SIZE_SHIFT
	_Y_SIZE_MASK            UI_Flag = __DIM_MASK << _Y_SIZE_SHIFT

	_X_OFF_SHIFT                   = _Y_SIZE_SHIFT + __DIM_BITS
	_X_OFF_ABSOLUTE        UI_Flag = __DIM_ABSOLUTE << _X_OFF_SHIFT
	_X_OFF_RELATIVE_ANCHOR UI_Flag = __DIM_RELATIVE_ANCHOR << _X_OFF_SHIFT
	_X_OFF_MASK            UI_Flag = __DIM_MASK << _X_OFF_SHIFT

	_Y_OFF_SHIFT                   = _X_OFF_SHIFT + __DIM_BITS
	_Y_OFF_ABSOLUTE        UI_Flag = __DIM_ABSOLUTE << _Y_OFF_SHIFT
	_Y_OFF_RELATIVE_ANCHOR UI_Flag = __DIM_RELATIVE_ANCHOR << _Y_OFF_SHIFT
	_Y_OFF_MASK            UI_Flag = __DIM_MASK << _Y_OFF_SHIFT

	_FIT_X_SHIFT              = _Y_OFF_SHIFT + __DIM_BITS
	_FIT_X_EXACT      UI_Flag = __EXACT << _FIT_X_SHIFT
	_FIT_X_MIN_GROW   UI_Flag = __MIN_GROW << _FIT_X_SHIFT
	_FIT_X_MAX_SHRINK UI_Flag = __MAX_SHRINK << _FIT_X_SHIFT
	_FIT_X_MASK       UI_Flag = __FIT_MASK << _FIT_X_SHIFT

	_FIT_Y_SHIFT              = _FIT_X_SHIFT + __FIT_BITS
	_FIT_Y_EXACT      UI_Flag = __EXACT << _FIT_Y_SHIFT
	_FIT_Y_MIN_GROW   UI_Flag = __MIN_GROW << _FIT_Y_SHIFT
	_FIT_Y_MAX_SHRINK UI_Flag = __MAX_SHRINK << _FIT_Y_SHIFT
	_FIT_Y_MASK       UI_Flag = __FIT_MASK << _FIT_Y_SHIFT

	_ANCHOR_IDX_SHIFT            = _FIT_Y_SHIFT + __FIT_BITS
	_ANCHOR_IDX_PARENT   UI_Flag = __ANCHOR_IDX_PARENT << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_ABSOLUTE UI_Flag = __ANCHOR_IDX_ABSOLUTE << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_RELATIVE UI_Flag = __ANCHOR_IDX_RELATIVE << _ANCHOR_IDX_SHIFT
	_ANCHOR_IDX_MASK     UI_Flag = __ANCHOR_IDX_MASK << _ANCHOR_IDX_SHIFT

	_PLANE_SHIFT         = _ANCHOR_IDX_SHIFT + __ANCHOR_IDX_BITS
	_PLANE_0     UI_Flag = __PLANE_0 << _PLANE_SHIFT
	_PLANE_1     UI_Flag = __PLANE_1 << _PLANE_SHIFT
	_PLANE_2     UI_Flag = __PLANE_2 << _PLANE_SHIFT
	_PLANE_3     UI_Flag = __PLANE_3 << _PLANE_SHIFT
	_PLANE_4     UI_Flag = __PLANE_4 << _PLANE_SHIFT
	_PLANE_5     UI_Flag = __PLANE_5 << _PLANE_SHIFT
	_PLANE_6     UI_Flag = __PLANE_6 << _PLANE_SHIFT
	_PLANE_7     UI_Flag = __PLANE_7 << _PLANE_SHIFT
	_PLANE_MASK  UI_Flag = __PLANE_MASK << _PLANE_SHIFT
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
	ABSOLUTE             UI_DimMode = UI_DimMode(__DIM_ABSOLUTE)
	RELATIVE_ANCHOR_SIZE UI_DimMode = UI_DimMode(__DIM_RELATIVE_ANCHOR)
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

type UI_Plane uint32

const (
	PLANE_0 UI_Plane = UI_Plane(__PLANE_0)
	PLANE_1 UI_Plane = UI_Plane(__PLANE_1)
	PLANE_2 UI_Plane = UI_Plane(__PLANE_2)
	PLANE_3 UI_Plane = UI_Plane(__PLANE_3)
	PLANE_4 UI_Plane = UI_Plane(__PLANE_4)
	PLANE_5 UI_Plane = UI_Plane(__PLANE_5)
	PLANE_6 UI_Plane = UI_Plane(__PLANE_6)
	PLANE_7 UI_Plane = UI_Plane(__PLANE_7)
)

var FLAG = UI_Flag(0)

func (f *UI_Flag) SetHidden() {
	*f |= _IS_HIDDEN
}
func (f *UI_Flag) SetVisible() {
	*f &= ^_IS_HIDDEN
}
func (f UI_Flag) IsHidden() bool {
	return f&_IS_HIDDEN == _IS_HIDDEN
}
func (f UI_Flag) IsVisible() bool {
	return f&_IS_HIDDEN == 0
}

func (f UI_Flag) XOffsetMode(dim UI_DimMode) UI_Flag {
	f &= ^_X_OFF_MASK
	return f | UI_Flag(dim<<_X_OFF_SHIFT)
}
func (f UI_Flag) GetXOffsetMode() UI_DimMode {
	return UI_DimMode((f & _X_OFF_MASK) >> _X_OFF_SHIFT)
}
func (f UI_Flag) YOffsetMode(dim UI_DimMode) UI_Flag {
	f &= ^_Y_OFF_MASK
	return f | UI_Flag(dim<<_Y_OFF_SHIFT)
}
func (f UI_Flag) GetYOffsetMode() UI_DimMode {
	return UI_DimMode((f & _Y_OFF_MASK) >> _Y_OFF_SHIFT)
}
func (f UI_Flag) XYOffset(dim UI_DimMode) UI_Flag {
	f &= ^_X_OFF_MASK
	f &= ^_Y_OFF_MASK
	f |= UI_Flag(dim << _X_OFF_SHIFT)
	return f | UI_Flag(dim<<_Y_OFF_SHIFT)
}

func (f UI_Flag) XSizeMode(dim UI_DimMode) UI_Flag {
	f &= ^_X_SIZE_MASK
	return f | UI_Flag(dim<<_X_SIZE_SHIFT)
}
func (f UI_Flag) GetXSizeMode() UI_DimMode {
	return UI_DimMode((f & _X_SIZE_MASK) >> _X_SIZE_SHIFT)
}
func (f UI_Flag) YSizeMode(dim UI_DimMode) UI_Flag {
	f &= ^_Y_SIZE_MASK
	return f | UI_Flag(dim<<_Y_SIZE_SHIFT)
}
func (f UI_Flag) GetYSizeMode() UI_DimMode {
	return UI_DimMode((f & _Y_SIZE_MASK) >> _Y_SIZE_SHIFT)
}
func (f UI_Flag) XYSizeMode(dim UI_DimMode) UI_Flag {
	f &= ^_X_SIZE_MASK
	f &= ^_Y_SIZE_MASK
	f |= UI_Flag(dim << _X_SIZE_SHIFT)
	return f | UI_Flag(dim<<_Y_SIZE_SHIFT)
}

func (f UI_Flag) SelfAnchor(anchor UI_Anchor) UI_Flag {
	f &= ^_SELF_ANCHOR_MASK
	return f | UI_Flag(anchor<<_SELF_ANCHOR_SHIFT)
}
func (f UI_Flag) GetSelfAnchor() UI_Anchor {
	return UI_Anchor((f & _SELF_ANCHOR_MASK) >> _SELF_ANCHOR_SHIFT)
}

func (f UI_Flag) ExternalAnchor(anchor UI_Anchor, mode UI_AnchorMode) UI_Flag {
	f &= ^_EXT_ANCHOR_MASK
	f &= ^_ANCHOR_IDX_MASK
	f |= UI_Flag(mode << _ANCHOR_IDX_SHIFT)
	return f | UI_Flag(anchor<<_EXT_ANCHOR_SHIFT)
}
func (f UI_Flag) GetExternalAnchor() UI_Anchor {
	return UI_Anchor((f & _EXT_ANCHOR_MASK) >> _EXT_ANCHOR_SHIFT)
}
func (f UI_Flag) GetExternalAnchorMode() UI_AnchorMode {
	return UI_AnchorMode((f & _ANCHOR_IDX_MASK) >> _ANCHOR_IDX_SHIFT)
}

func (f UI_Flag) Anchor(selfAnchor UI_Anchor, externalAnchor UI_Anchor, mode UI_AnchorMode) UI_Flag {
	f &= ^_EXT_ANCHOR_MASK
	f &= ^_SELF_ANCHOR_MASK
	f &= ^_ANCHOR_IDX_MASK
	f |= UI_Flag(selfAnchor << _SELF_ANCHOR_SHIFT)
	f |= UI_Flag(mode << _ANCHOR_IDX_SHIFT)
	return f | UI_Flag(externalAnchor<<_EXT_ANCHOR_SHIFT)
}

func (f UI_Flag) XFit(fit UI_Fit) UI_Flag {
	f &= ^_FIT_X_MASK
	return f | UI_Flag(fit<<_FIT_X_SHIFT)
}
func (f UI_Flag) GetXFit() UI_Fit {
	return UI_Fit((f & _FIT_X_MASK) >> _FIT_X_SHIFT)
}
func (f UI_Flag) YFit(fit UI_Fit) UI_Flag {
	f &= ^_FIT_Y_MASK
	return f | UI_Flag(fit<<_FIT_Y_SHIFT)
}
func (f UI_Flag) GetYFit() UI_Fit {
	return UI_Fit((f & _FIT_Y_MASK) >> _FIT_Y_SHIFT)
}
func (f UI_Flag) XYFit(fit UI_Fit) UI_Flag {
	f &= ^_FIT_X_MASK
	f &= ^_FIT_Y_MASK
	f |= UI_Flag(fit << _FIT_X_SHIFT)
	return f | UI_Flag(fit<<_FIT_Y_SHIFT)
}

func (f UI_Flag) Plane(plane UI_Plane) UI_Flag {
	f &= ^_PLANE_MASK
	return f | UI_Flag(plane<<_PLANE_SHIFT)
}

func (f UI_Flag) GetPlane() UI_Plane {
	return UI_Plane((f & _PLANE_MASK) >> _PLANE_SHIFT)
}
func (f UI_Flag) GetPlaneAsZIndex() uint16 {
	return uint16((f&_PLANE_MASK)>>_PLANE_SHIFT) << (16 - __PLANE_BITS)
}

const (
	_PLANE_ZINDEX_SHIFT        = 16 - __PLANE_BITS
	_PLANE_IN_ZINDEX_MASK      = uint16((math.MaxUint16 << _PLANE_ZINDEX_SHIFT) & math.MaxUint16)
	_ZINDEX_WITHOUT_PLANE_MASK = math.MaxUint16 >> __PLANE_BITS
)
