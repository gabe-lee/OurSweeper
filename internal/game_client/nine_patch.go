package game_client

import (
	"image"

	C "github.com/gabe-lee/OurSweeper/internal/consts"
	"github.com/hajimehoshi/ebiten/v2"
)

type NinePatchSource struct {
	TopLeft  Vec2_Int
	BotRight Vec2_Int
	Margin   Vec2_Int
}

type NP_IDX uint16

const (
	NINE_P_DARK NP_IDX = iota
	NINE_P_LIGHT
)

var NINE_SOURCE = [...]NinePatchSource{
	NINE_P_DARK:  NinePatchSource{TopLeft: SPRITE[SPRITE_PANEL_DARK], BotRight: SPRITE[SPRITE_PANEL_DARK].Add(SPRITE_SIZE), Margin: SPRITE_PANEL_MARGIN},
	NINE_P_LIGHT: NinePatchSource{TopLeft: SPRITE[SPRITE_PANEL_LIGHT], BotRight: SPRITE[SPRITE_PANEL_LIGHT].Add(SPRITE_SIZE), Margin: SPRITE_PANEL_MARGIN},
}

type NinePatchPanel struct {
	SrcIdx     NP_IDX
	TopLeft    Vec2_Int
	Size       Vec2_Int
	Parent     UI_IDX
	Children   []UI_IDX
	FnOnClick  func()
	FnOnHover  func()
	FnOnLayout func()
	ZIndex     int
	Active     bool
}

func (n *NinePatchPanel) GetChild(relativeIdx int) UI_Element {
	return UI[n.Children[relativeIdx]]
}

// GetIndexInParent implements UI_Element.
func (n *NinePatchPanel) GetIndexInParent() int {
	return int(n.OwnIdx)
}

// GetParent implements UI_Element.
func (n *NinePatchPanel) GetParent() UI_Element {
	panic("unimplemented")
}

// IsActive implements UI_Element.
func (n *NinePatchPanel) IsActive() bool {
	panic("unimplemented")
}

// OnScroll implements UI_Element.
func (n *NinePatchPanel) OnScroll(dx float64, dy float64) {
	panic("unimplemented")
}

// ReLayout implements UI_Element.
func (n *NinePatchPanel) ReLayout(depth int) {
	panic("unimplemented")
}

// SetActive implements UI_Element.
func (n *NinePatchPanel) SetActive(state bool) {
	panic("unimplemented")
}

// Draw implements UI_Element.
func (n *NinePatchPanel) Draw(atlas *EbitImage, screen *EbitImage) {
	src := NINE_SOURCE[n.SrcIdx]
	srcSize := src.BotRight.Sub(src.TopLeft)
	srcRect := image.Rect(src.TopLeft.X, src.TopLeft.Y, src.TopLeft.X+C.TILE_SIZE, src.TopLeft.Y+C.TILE_SIZE)
	op := &ebiten.DrawImageOptions{}
	scale := n.Size.ToFloat64().Div(srcSize.ToFloat64())
	op.GeoM.Scale(scale.X, scale.Y)
	op.GeoM.Translate(float64(n.TopLeft.X), float64(n.TopLeft.Y))
	screen.DrawImage(atlas.SubImage(srcRect).(*EbitImage), op)
	for _, idx := range n.Children {
		UI[idx].Draw(atlas, screen)
	}
}

// GetPos implements UI_Element.
func (n *NinePatchPanel) GetPos() Vec2_Int {
	return n.TopLeft
}

// GetSize implements UI_Element.
func (n *NinePatchPanel) GetSize() Vec2_Int {
	return n.Size
}

// Layout implements UI_Element.
func (n *NinePatchPanel) Layout() {
	if n.FnOnLayout != nil {
		n.FnOnLayout()
	}
}

// OnClick implements UI_Element.
func (n *NinePatchPanel) OnClick() {
	if n.FnOnClick != nil {
		n.FnOnClick()
	} else {
		UI[n.Parent].OnClick()
	}
}

// OnHover implements UI_Element.
func (n *NinePatchPanel) OnHover() {
	if n.FnOnHover != nil {
		n.FnOnHover()
	}
}

var _ UI_Element = (*NinePatchPanel)(nil)
