package common

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

// Menu colours
var (
	LightFontColour   = turdgl.RGB(221, 208, 201)
	LighterFontColour = turdgl.RGB(251, 238, 231)

	buttonColourUnpressed = turdgl.RGB(143, 122, 101)
	buttonColourPressed   = turdgl.RGB(143+20, 122+20, 101+20)
)

// UI colours
var (
	BackgroundColour     = turdgl.RGB(248, 248, 237) // official colour
	BackgroundColourWin  = turdgl.RGB(36, 59, 34)
	BackgroundColourLose = turdgl.RGB(38, 15, 15)

	LightGreyTextColour = turdgl.RGB(240, 229, 215) // official colour
	GreyTextColour      = turdgl.RGB(120, 110, 100) // official colour
	WhiteFontColour     = turdgl.RGB(255, 255, 255) // official colour

	buttonOrangeColour    = turdgl.RGB(235, 152, 91)  // official colour
	tileBackgroundColour  = turdgl.RGB(204, 192, 180) // official colour
	ArenaBackgroundColour = turdgl.RGB(187, 173, 160) // official colour

)

// Tile colours
var (
	tile2Colour    = turdgl.RGB(239, 229, 218) // official colour
	tile4Colour    = turdgl.RGB(236, 224, 198) // official colour
	tile8Colour    = turdgl.RGB(242, 176, 121) // official colour
	tile16Colour   = turdgl.RGB(235, 140, 83)  // official colour
	tile32Colour   = turdgl.RGB(245, 123, 93)  // official colour
	tile64Colour   = turdgl.RGB(233, 89, 55)   // official colour
	tile128Colour  = turdgl.RGB(242, 217, 107) // official colour
	tile256Colour  = turdgl.RGB(241, 208, 76)  // official colour
	tile512Colour  = turdgl.RGB(229, 192, 43)  // official colour
	tile1024Colour = turdgl.RGB(224, 192, 65)
	Tile2048Colour = turdgl.RGB(235, 196, 2) // official colour
	tile4096Colour = turdgl.RGB(255, 59, 59)
	tile8192Colour = turdgl.RGB(255, 32, 33)
)

const (
	FontPathMedium = "./assets/ClearSans/ClearSans-Medium.ttf"
	FontPathBold   = "./assets/ClearSans/ClearSans-Medium.ttf"
)

// tileColour returns the colour for a tile of a given value.
func tileColour(val int) color.Color {
	switch val {
	case 2:
		return tile2Colour
	case 4:
		return tile4Colour
	case 8:
		return tile8Colour
	case 16:
		return tile16Colour
	case 32:
		return tile32Colour
	case 64:
		return tile64Colour
	case 128:
		return tile128Colour
	case 256:
		return tile256Colour
	case 512:
		return tile512Colour
	case 1024:
		return tile1024Colour
	case 2048:
		return Tile2048Colour
	case 4096:
		return tile4096Colour
	case 8192:
		return tile8192Colour
	default:
		return turdgl.RGB(255, 0, 0)
	}
}

// tileTextColour returns the colour of the text for tile of a given value.
func tileTextColour(val int) color.Color {
	switch val {
	case 2, 4:
		return GreyTextColour
	default:
		return WhiteFontColour
	}
}
