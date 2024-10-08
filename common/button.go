package common

import (
	"github.com/z-riley/turdgl"
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
	b := turdgl.NewButton(r, FontPathMedium).
		SetLabelSize(36).
		SetLabelColour(LightFontColour)
	b.Behaviour = turdgl.OnRelease
	b.SetCallback(func(m turdgl.MouseState) { cb() })

	return &MenuButton{b}
}

func (b *MenuButton) Update(win *turdgl.Window) {
	// Adjust style if cursor hovering or button is pressed
	if b.Shape.IsWithin(win.MouseLocation()) {
		if win.MouseButtonState() == turdgl.LeftClick {
			b.Label.SetColour(LighterFontColour)
			b.Shape.SetStyle(buttonStylePressed)
		} else {
			b.Label.SetColour(LighterFontColour)
			b.Shape.SetStyle(buttonStyleHovering)
		}
	} else {
		b.Label.SetColour(LightFontColour)
		b.Shape.SetStyle(buttonStyleUnpressed)
	}

	// Call underlying button update function
	b.Button.Update(win)
}

// GameButton is a button that's used in the main game's UI.
type GameButton turdgl.Button

// NewGameButton constructs a new game button with sensible defaults.
func NewGameButton(width, height float64, pos turdgl.Vec, cb func()) *turdgl.Button {
	r := turdgl.NewCurvedRect(width, height, 3, pos)
	r.SetStyle(turdgl.Style{Colour: buttonOrangeColour})

	b := turdgl.NewButton(r, FontPathBold).
		SetLabelText("BUTTON").
		SetLabelSize(14).
		SetLabelColour(WhiteFontColour).
		SetLabelAlignment(turdgl.AlignCustom).
		SetCallback(func(turdgl.MouseState) { cb() })
	b.Behaviour = turdgl.OnRelease
	b.Label.SetOffset(turdgl.Vec{X: width / 2, Y: height / 1.2})

	return b
}
