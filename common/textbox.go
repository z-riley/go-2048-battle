package common

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

// EntryBox is a commonly used text box for the user to enter text.
type EntryBox struct{ *turdgl.TextBox }

// NewEntryBox constructs a new text box with suitable defaults.
func NewEntryBox(width, height float64, pos turdgl.Vec) *EntryBox {
	t := NewTextBox(width, height, pos)
	t.SetSelectedCB(func() { t.SetTextColour(turdgl.White) })
	t.SetDeselectedCB(func() { t.SetTextColour(LightFontColour) })

	return &EntryBox{t}
}

func (t *EntryBox) Update(win *turdgl.Window) {
	t.TextBox.Update(win)
}

func NewTextBox(width, height float64, pos turdgl.Vec) *turdgl.TextBox {
	r := turdgl.NewRect(
		width, height, pos,
		turdgl.WithStyle(turdgl.Style{Colour: color.RGBA{90, 65, 48, 255}, Thickness: 0}),
	)
	r.SetStyle(turdgl.Style{Colour: buttonColourUnpressed})

	tb := turdgl.NewTextBox(r, FontPathMedium).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetText("Click to edit").
		SetTextSize(36).
		SetTextColour(LightFontColour)

	return tb
}

// GameUIBox is a commonly used text box for displaying scores.
type GameUIBox struct {
	heading *turdgl.Text
	body    *turdgl.TextBox
}

// NewGameTextBox constructs a new text box for the game's UI.
func NewGameTextBox(width, height float64, pos turdgl.Vec, colour color.RGBA) *GameUIBox {
	headingPos := turdgl.Vec{
		X: pos.X + width/2,
		Y: pos.Y + 25,
	}
	heading := turdgl.NewText("Heading", headingPos, FontPathBold).
		SetColour(WhiteFontColour).
		SetSize(16).
		SetAlignment(turdgl.AlignTopCentre)

	r := turdgl.NewCurvedRect(
		width, height,
		3,
		pos,
		turdgl.WithStyle(turdgl.Style{Colour: colour}),
	)

	body := turdgl.NewTextBox(r, FontPathBold).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetTextSize(26).
		SetTextColour(WhiteFontColour)

	return &GameUIBox{heading, body}
}

// Draw draws the UI box to the window.
func (g *GameUIBox) Draw(win *turdgl.Window) {
	win.Draw(g.body)
	win.Draw(g.heading)
}

// SetHeading sets the heading text of the UI box.
func (g *GameUIBox) SetHeading(s string) *GameUIBox {
	g.heading.SetText(s)
	return g
}

// SetBody sets the body text of the UI box.
func (g *GameUIBox) SetBody(s string) *GameUIBox {
	g.body.SetText(s)
	return g
}

// NewGameText constructs text with sensible defaults.
func NewGameText(body string, pos turdgl.Vec) *turdgl.Text {
	return turdgl.NewText(body, pos, FontPathBold).
		SetColour(ArenaBackgroundColour).
		SetSize(18)
}
