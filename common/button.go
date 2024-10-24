package common

import (
	"github.com/z-riley/turdgl"
)

var (
	buttonStyleUnpressed = turdgl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     0,
	}

	buttonStyleHovering = turdgl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     8,
	}

	buttonStylePressed = turdgl.Style{
		Colour:    buttonColourPressed,
		Thickness: 0,
		Bloom:     5,
	}
)

// NewMenuButton constructs a new menu button with sensible defaults.
func NewMenuButton(width, height float64, pos turdgl.Vec, cb func()) *turdgl.Button {
	b := turdgl.NewButton(
		turdgl.NewCurvedRect(width, height, 6, pos.Round(),
			turdgl.WithStyle(buttonStyleUnpressed)),
		FontPathMedium,
	).
		SetLabelText("SET ME").
		SetLabelSize(36).
		SetLabelColour(WhiteFontColour).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32})

	b.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			b.Shape.SetStyle(buttonStyleHovering)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			b.Shape.SetStyle(buttonStyleUnpressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			b.Shape.SetStyle(buttonStylePressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		cb,
	)

	return b
}

var gameButtonStyleHovering = turdgl.Style{
	Colour:    buttonColourUnpressed,
	Thickness: 0,
	Bloom:     4,
}

// NewGameButton constructs a new game button with sensible defaults.
func NewGameButton(width, height float64, pos turdgl.Vec, cb func()) *turdgl.Button {
	b := turdgl.NewButton(
		turdgl.NewCurvedRect(
			width, height, TileCornerRadius, pos.Round(),
			turdgl.WithStyle(turdgl.Style{Colour: buttonOrangeColour})),
		FontPathBold,
	).
		SetLabelText("BUTTON").
		SetLabelSize(14).
		SetLabelColour(WhiteFontColour).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 12})

	b.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			b.Shape.SetStyle(gameButtonStyleHovering)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			b.Shape.SetStyle(buttonStyleUnpressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			b.Shape.SetStyle(buttonStylePressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		cb,
	)

	return b
}
