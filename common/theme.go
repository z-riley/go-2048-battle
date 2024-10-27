package common

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

// Menu colours
var (
	LighterFontColour = turdgl.RGB(251, 238, 231)

	buttonColourUnpressed = Tile16Colour
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
	Tile2Colour    = turdgl.RGB(239, 229, 218) // official colour
	Tile4Colour    = turdgl.RGB(236, 224, 198) // official colour
	Tile8Colour    = turdgl.RGB(242, 176, 121) // official colour
	Tile16Colour   = turdgl.RGB(235, 140, 83)  // official colour
	Tile32Colour   = turdgl.RGB(245, 123, 93)  // official colour
	Tile64Colour   = turdgl.RGB(233, 89, 55)   // official colour
	Tile128Colour  = turdgl.RGB(242, 217, 107) // official colour
	Tile256Colour  = turdgl.RGB(241, 208, 76)  // official colour
	Tile512Colour  = turdgl.RGB(229, 192, 43)  // official colour
	Tile1024Colour = turdgl.RGB(224, 192, 65)
	Tile2048Colour = turdgl.RGB(235, 196, 2) // official colour
	Tile4096Colour = turdgl.RGB(255, 59, 59)
	Tile8192Colour = turdgl.RGB(255, 32, 33)
)

const (
	FontPathMedium = "./assets/ClearSans/ClearSans-Medium.ttf"
	FontPathBold   = "./assets/ClearSans/ClearSans-Medium.ttf"
)

// tileColour returns the colour for a tile of a given value.
func tileColour(val int) color.Color {
	switch val {
	case 2:
		return Tile2Colour
	case 4:
		return Tile4Colour
	case 8:
		return Tile8Colour
	case 16:
		return Tile16Colour
	case 32:
		return Tile32Colour
	case 64:
		return Tile64Colour
	case 128:
		return Tile128Colour
	case 256:
		return Tile256Colour
	case 512:
		return Tile512Colour
	case 1024:
		return Tile1024Colour
	case 2048:
		return Tile2048Colour
	case 4096:
		return Tile4096Colour
	case 8192:
		return Tile8192Colour
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
