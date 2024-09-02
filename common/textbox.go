package common

import (
	"image/color"

	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/turdgl"
)

// EntryBox is a commonly used text box for the user to enter text.
type EntryBox struct{ *turdgl.TextBox }

// NewEntryBox constructs a new text box with suitable defaults.
func NewEntryBox(width, height float64, pos turdgl.Vec) *EntryBox {
	r := turdgl.NewRect(
		width, height, pos,
		turdgl.WithStyle(turdgl.Style{Colour: color.RGBA{90, 65, 48, 255}, Thickness: 0}),
	)
	r.SetStyle(turdgl.Style{Colour: buttonColourUnpressed})

	t := turdgl.NewTextBox(r, game.FontPath).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetText("Click to edit").
		SetTextSize(36).
		SetTextColour(LightFontColour)
	t.SetSelectedCB(func() { t.SetTextColour(turdgl.White) })
	t.SetDeselectedCB(func() { t.SetTextColour(LightFontColour) })

	return &EntryBox{t}
}

func (t *EntryBox) Update(win *turdgl.Window) {
	t.TextBox.Update(win)
}
