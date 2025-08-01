package game_client

import (
	ui "github.com/gabe-lee/OurSweeper/snapui"
	"github.com/gabe-lee/OurSweeper/vec2"
)

type (
	UI_Element = ui.UI_Element[UI_UserData]
	UI_Flag    = ui.UI_Flag
	UI_IDX     = ui.UI_Idx
	Vec2_F32   = ui.Vec2_F32
)

const (
	UI_ROOT UI_IDX = iota
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
	USER_DATA_TEXT = iota
	USER_DATA_IMAGE
)

const (
	PADDING float32 = 4.0
)

type UI_UserData struct {
	Scale float32
	Idx   uint16
	Kind  uint8
	Red   uint8
	Blue  uint8
	Green uint8
	Alpha uint8
	Style uint8
}

var UI = [...]UI_Element{
	UI_ROOT: UI_Element{
		Children: []UI_IDX{UI_LOGIN_PANEL},
		Flags:    ui.FLAG.Anchor(ui.TOP_LEFT, ui.TOP_LEFT, ui.ANCHOR_TO_PARENT).XYFit(ui.EXACT).XYSizeMode(ui.RELATIVE_ANCHOR_SIZE).XYOffset(ui.ABSOLUTE),
		Size:     vec2.New[float32](1.0, 1.0),
	},
	UI_LOGIN_PANEL: UI_Element{
		Children: []UI_IDX{UI_LOGIN_PANEL_BUTTON_ANON, UI_LOGIN_PANEL_BUTTON_LOGIN, UI_LOGIN_PANEL_OR_SEPARATOR, UI_LOGIN_PANEL_PASSWORD_INPUT,
			UI_LOGIN_PANEL_PASSWORD_LABEL, UI_LOGIN_PANEL_USERNAME_INPUT, UI_LOGIN_PANEL_USERNAME_LABEL},
		Flags:        ui.FLAG.Anchor(ui.MID_CENTER, ui.MID_CENTER, ui.ANCHOR_TO_PARENT).XFit(ui.MIN_GROW).YFit(ui.MAX_SHRINK).XYSizeMode(ui.RELATIVE_ANCHOR_SIZE).Plane(ui.PLANE_5),
		Size:         vec2.New[float32](0.33, 0.80),
		Padding:      PADDING,
		AnchorRefIdx: UI_ROOT,
	},
	UI_LOGIN_PANEL_BUTTON_ANON: UI_Element{
		Children:     []UI_IDX{UI_LOGIN_PANEL_BUTTON_ANON_TEXT},
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
