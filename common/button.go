package common

import (
	game "github.com/zac460/go-2048-battle"
	"github.com/zac460/turdgl"
)

var buttonStyleUnpressed = turdgl.Style{
	Colour:    buttonColourUnpressed,
	Thickness: 0,
	Bloom:     0,
}

var buttonStyleHovering = turdgl.Style{
	Colour:    buttonColourUnpressed,
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

// NewMenuButton constructs a new button with suitable defaults for a menu button.
func NewMenuButton(width, height float64, pos turdgl.Vec, cb func()) *MenuButton {
	r := turdgl.NewRect(width, height, pos, turdgl.WithStyle(buttonStyleUnpressed))
	b := turdgl.NewButton(r, game.FontPath)
	b.Behaviour = turdgl.OnRelease
	b.SetCallback(func(m turdgl.MouseState) { cb() })

	return &MenuButton{b}
}

func (b *MenuButton) Update(win *turdgl.Window) {
	// Adjust style if cursor hovering or button is pressed
	if b.Shape.IsWithin(win.MouseLocation()) {
		if win.MouseButtonState() == turdgl.LeftClick {
			b.Shape.SetStyle(buttonStylePressed)
		} else {
			b.Shape.SetStyle(buttonStyleHovering)
		}
	} else {
		b.Shape.SetStyle(buttonStyleUnpressed)
	}

	// Call underlying button update function
	b.Button.Update(win)
}
