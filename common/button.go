package common

import (
	game "github.com/zac460/go-2048-battle"
	"github.com/zac460/turdgl"
)

var buttonStyleUnpressed = turdgl.Style{
	Colour:    buttonThemeUnpressed,
	Thickness: 0,
	Bloom:     0,
}

var buttonStyleHovering = turdgl.Style{
	Colour:    buttonThemeUnpressed,
	Thickness: 0,
	Bloom:     5,
}

var buttonStylePressed = turdgl.Style{
	Colour:    buttonColourPressed,
	Thickness: 0,
	Bloom:     5,
}

// MenuButton is a commonly used button for navigating menus.
type MenuButton struct{ *turdgl.Button }

// NewMenu button constructs a new button with defaults suitable for a menu button.
func NewMenuButton(width, height float64, pos turdgl.Vec, cb func()) *MenuButton {
	r := turdgl.NewRect(width, height, pos, turdgl.WithStyle(buttonStyleUnpressed))
	b := turdgl.NewButton(r, game.FontPath)
	b.SetCallback(func(m turdgl.MouseState) {
		// Callback executes after every update
		switch {
		case m == turdgl.LeftClick:
			if b.IsHovering() {
				r.SetStyle(buttonStylePressed)
				cb()
			}
		case b.IsHovering():
			r.SetStyle(buttonStyleHovering)
		case m == turdgl.NoClick:
			r.SetStyle(buttonStyleUnpressed)
		}
	})

	return &MenuButton{b}
}
