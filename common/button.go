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

// NewMenuButton constructs a new menu button with sensible defaults.
func NewMenuButton(width, height float64, pos turdgl.Vec, cb func()) *turdgl.Button {
	b := turdgl.NewButton(
		turdgl.NewRect(width, height, pos, turdgl.WithStyle(buttonStyleUnpressed)),
		FontPathMedium,
	).
		SetLabelSize(36).
		SetLabelText("SET ME").
		SetLabelColour(LightFontColour).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32})

	b.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			b.Label.SetColour(LighterFontColour)
			b.Shape.SetStyle(buttonStyleHovering)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			b.Label.SetColour(LightFontColour)
			b.Shape.SetStyle(buttonStyleUnpressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnPress},
		func() {
			b.Label.SetColour(LighterFontColour)
			b.Shape.SetStyle(buttonStylePressed)
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		cb,
	)

	return b
}

// NewGameButton constructs a new game button with sensible defaults.
func NewGameButton(width, height float64, pos turdgl.Vec, cb func()) *turdgl.Button {
	r := turdgl.NewCurvedRect(width, height, 3, pos)
	r.SetStyle(turdgl.Style{Colour: buttonOrangeColour})

	b := turdgl.NewButton(r, FontPathBold).
		SetLabelText("BUTTON").
		SetLabelSize(14).
		SetLabelColour(WhiteFontColour).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: width / 2, Y: height / 1.2}).
		// TODO: add hover and press animations
		SetCallback(
			turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
			cb,
		)

	return b
}
