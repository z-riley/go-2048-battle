package common

import (
	"github.com/z-riley/gogl"
)

var (
	ButtonStyleUnpressed = gogl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     0,
	}

	ButtonStyleHovering = gogl.Style{
		Colour:    buttonColourUnpressed,
		Thickness: 0,
		Bloom:     8,
	}

	ButtonStylePressed = gogl.Style{
		Colour:    buttonColourPressed,
		Thickness: 0,
		Bloom:     5,
	}
)

// NewMenuButton constructs a new menu button with sensible defaults.
func NewMenuButton(width, height float64, pos gogl.Vec, callback func()) *gogl.Button {
	r := gogl.NewCurvedRect(width, height, 6, pos.Round()).SetStyle(ButtonStyleUnpressed)
	b := gogl.NewButton(r, FontPathMedium).
		SetLabelText("SET ME").
		SetLabelSize(36).
		SetLabelColour(WhiteFontColour)

	b.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleHovering)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleUnpressed)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			r.SetStyle(ButtonStylePressed)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnRelease},
		callback,
	)

	return b
}

var gameButtonStyleHovering = gogl.Style{
	Colour:    buttonColourUnpressed,
	Thickness: 0,
	Bloom:     4,
}

// NewGameButton constructs a new game button with sensible defaults.
func NewGameButton(width, height float64, pos gogl.Vec, callback func()) *gogl.Button {
	r := gogl.NewCurvedRect(width, height, TileCornerRadius, pos.Round()).
		SetStyle(gogl.Style{Colour: ButtonOrangeColour})
	b := gogl.NewButton(r, FontPathBold).
		SetLabelText("BUTTON").
		SetLabelSize(14).
		SetLabelColour(WhiteFontColour)

	b.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(gameButtonStyleHovering)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			b.Label.SetColour(WhiteFontColour)
			r.SetStyle(ButtonStyleUnpressed)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnPress},
		func() {
			b.Label.SetColour(GreyTextColour)
			r.SetStyle(ButtonStylePressed)
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.LeftClick, Behaviour: gogl.OnRelease},
		callback,
	)

	return b
}
