package game_client

import (
	ui "github.com/gabe-lee/OurSweeper/snaplay"
	"github.com/gabe-lee/OurSweeper/vec2"
)

type (
	UI_Element = ui.UI_Element[UI_UserData]
	UI_Flag    = ui.Snap_Flag
	UI_Idx     = ui.Snap_Idx
	Vec2_F32   = ui.Vec2_F32
)

const (
	UI_ROOT UI_Idx = iota
	UI_LOGIN_PANEL
	UI_LOGIN_PANEL_BUTTON_ANON
	UI_LOGIN_PANEL_BUTTON_ANON_TEXT
	UI_LOGIN_PANEL_BUTTON_LOGIN
	UI_LOGIN_PANEL_OR_SEPARATOR
	UI_LOGIN_PANEL_USERNAME_LABEL
	UI_LOGIN_PANEL_USERNAME_INPUT
	UI_LOGIN_PANEL_PASSWORD_LABEL
	UI_LOGIN_PANEL_PASSWORD_INPUT
)

const (
	PADDING float32 = 4.0
)

type UI_UserDataIdx uint16
type UI_UserDataKind uint16

const (
	USER_DATA_TEXT UI_UserDataKind = iota
	USER_DATA_IMAGE
)

type UI_UserData struct {
	Kind    UI_UserDataKind
	Param_1 UI_UserDataIdx
	Param_2 UI_UserDataIdx
	Param_3 UI_UserDataIdx
	Param_4 UI_UserDataIdx
}

var UI = [...]UI_Element{
	UI_ROOT: UI_Element{
		Children: []UI_Idx{UI_LOGIN_PANEL},
		Flags:    ui.FLAG.Anchor(ui.TOP_LEFT, ui.TOP_LEFT, ui.ANCHOR_TO_PARENT).XYFit(ui.EXACT).XYSizeMode(ui.RELATIVE_ANCHOR_SIZE).XYOffset(ui.ABSOLUTE),
		Size:     vec2.New[float32](1.0, 1.0),
	},
	UI_LOGIN_PANEL: UI_Element{
		Children: []UI_Idx{UI_LOGIN_PANEL_BUTTON_ANON, UI_LOGIN_PANEL_BUTTON_LOGIN, UI_LOGIN_PANEL_OR_SEPARATOR, UI_LOGIN_PANEL_PASSWORD_INPUT,
			UI_LOGIN_PANEL_PASSWORD_LABEL, UI_LOGIN_PANEL_USERNAME_INPUT, UI_LOGIN_PANEL_USERNAME_LABEL},
		Flags:        ui.FLAG.Anchor(ui.MID_CENTER, ui.MID_CENTER, ui.ANCHOR_TO_PARENT).XFit(ui.MIN_GROW).YFit(ui.MAX_SHRINK).XYSizeMode(ui.RELATIVE_ANCHOR_SIZE).Plane(ui.PLANE_5),
		Size:         vec2.New[float32](0.33, 0.80),
		Padding:      PADDING,
		AnchorRefIdx: UI_ROOT,
	},
	UI_LOGIN_PANEL_BUTTON_ANON: UI_Element{
		Children:     []UI_Idx{UI_LOGIN_PANEL_BUTTON_ANON_TEXT},
		Flags:        ui.FLAG.Anchor(ui.TOP_CENTER, ui.TOP_CENTER, ui.ANCHOR_TO_PARENT).XYFit(ui.MIN_GROW).XSizeMode(ui.RELATIVE_ANCHOR_SIZE).YSizeMode(ui.ABSOLUTE),
		Size:         vec2.New[float32](0.25, 24.0),
		Padding:      PADDING,
		AnchorRefIdx: UI_LOGIN_PANEL,
	},
	UI_LOGIN_PANEL_BUTTON_ANON_TEXT: UI_Element{
		Flags:        ui.FLAG.Anchor(ui.TOP_CENTER, ui.TOP_CENTER, ui.ANCHOR_TO_PARENT),
		GetSize:      GetUiTextSize,
		Padding:      PADDING,
		AnchorRefIdx: UI_LOGIN_PANEL_BUTTON_ANON,
	},
}
