package common

import (
	"github.com/z-riley/turdgl"
)

var (
	ButtonStyleUnpressed = turdgl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     0,
	}

	ButtonStyleHovering = turdgl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     8,
	}

	ButtonStylePressed = turdgl.Style{
		Colour:    buttonColourPressed,
		Thickness: 0,
		Bloom:     5,
	}
)

// NewMenuButton constructs a new menu button with sensible defaults.
func NewMenuButton(width, height float64, pos turdgl.Vec, callback func()) *turdgl.Button {
	r := turdgl.NewCurvedRect(width, height, 6, pos.Round()).SetStyle(ButtonStyleUnpressed)
	b := turdgl.NewButton(r, FontPathMedium).
		SetLabelText("SET ME").
		SetLabelSize(36).
		SetLabelColour(WhiteFontColour)

	b.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleHovering)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleUnpressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			r.SetStyle(ButtonStylePressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		callback,
	)

	return b
}

var gameButtonStyleHovering = turdgl.Style{
	Colour:    buttonColourUnpressed,
	Thickness: 0,
	Bloom:     4,
}

// NewGameButton constructs a new game button with sensible defaults.
func NewGameButton(width, height float64, pos turdgl.Vec, callback func()) *turdgl.Button {
	r := turdgl.NewCurvedRect(width, height, TileCornerRadius, pos.Round()).
		SetStyle(turdgl.Style{Colour: ButtonOrangeColour})
	b := turdgl.NewButton(r, FontPathBold).
		SetLabelText("BUTTON").
		SetLabelSize(14).
		SetLabelColour(WhiteFontColour)

	b.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(gameButtonStyleHovering)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleUnpressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			r.SetStyle(ButtonStylePressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		callback,
	)

	return b
}
