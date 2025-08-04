package parametric_ui

type UI_Element[UI_UserData any] struct {
	Children []uint16
	OnClick  func()
	OnHover  func()
	OnScroll func()
	OnDrag   func()
	P_Left   uint16
	P_Top    uint16
	P_Right  uint16
	P_Bot    uint16
	Zdepth   uint16
	Zplane   uint8
	Hidden   bool
	UserData UI_UserData
}

type testUsedData struct {
	Kind    uint8
	P_Scale uint16
	Color   uint32
}

type testElem = UI_Element[testUsedData]
