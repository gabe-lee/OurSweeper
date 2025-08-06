package game_client

import "image/color"

type (
	Color = color.RGBA
)

const (
	COLOR_TEXT_PRIMARY = iota
	_colCount
)

var PALETTE = [_colCount]Color{
	COLOR_TEXT_PRIMARY: Color{R: 225, G: 225, B: 225, A: 255},
	// COLOR_TEXT_PRIMARY: Color{R: 225, G: 225, B: 225, A: 255},
}
