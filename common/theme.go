package common

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

var (
	BackgroundColour      = turdgl.RGB(46, 36, 27)
	LightFontColour       = turdgl.RGB(221, 208, 201)
	LighterFontColour     = turdgl.RGB(251, 238, 231)
	buttonColourUnpressed = turdgl.RGB(143, 122, 101)
	buttonColourPressed   = turdgl.RGB(143+20, 122+20, 101+20)
)

// TileColour returns the colour for a tile of a given value.
func TileColour(val int) color.Color {
	switch val {
	case 0:
		return BackgroundColour
	case 2:
		return turdgl.RGB(236, 229, 219)
	case 4:
		return turdgl.RGB(235, 224, 202)
	case 8:
		return turdgl.RGB(232, 180, 130)
	case 16:
		return turdgl.RGB(232, 154, 108)
	case 32:
		return turdgl.RGB(230, 130, 102)
	case 64:
		return turdgl.RGB(228, 103, 71)
	case 128:
		return turdgl.RGB(234, 209, 127)
	case 256:
		return turdgl.RGB(232, 206, 113)
	case 512:
		return turdgl.RGB(238, 199, 68)
	case 1024:
		return turdgl.RGB(224, 192, 65)
	case 2048:
		return turdgl.RGB(238, 193, 48)
	case 4096:
		return turdgl.RGB(255, 59, 59)
	case 8192:
		return turdgl.RGB(255, 32, 33)
	default:
		return turdgl.RGB(255, 0, 0)
	}
}
