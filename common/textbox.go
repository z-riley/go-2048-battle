package common

import (
	"image/color"

	"github.com/z-riley/turdgl"
)

// Entrybox is an interactive text box for data entry.
type EntryBox struct {
	TextBox *turdgl.TextBox
	bloom   *turdgl.CurvedRect
}

// NewEntryBox constructs a new text box with suitable defaults.
func NewEntryBox(width, height float64, pos turdgl.Vec) *EntryBox {
	var (
		styleUnselected = turdgl.Style{
			Colour:    turdgl.LightGrey,
			Thickness: 1,
			Bloom:     0,
		}
		styleSelected = turdgl.Style{
			Colour:    turdgl.LightGrey,
			Thickness: 1,
			Bloom:     10,
		}
	)

	bloom := turdgl.NewCurvedRect(width, height, 6, pos, turdgl.WithStyle(styleUnselected))

	tb := NewTextBox(width, height, pos)
	tb.SetSelectedCB(func() {
		tb.SetTextColour(turdgl.White)
		bloom.SetStyle(styleSelected)
	}).SetDeselectedCB(func() {
		tb.SetTextColour(LightGreyTextColour)
		bloom.SetStyle(styleUnselected)
	})

	return &EntryBox{tb, bloom}
}

// Draw draws an entry box to the frame buffer.
func (e *EntryBox) Draw(buf *turdgl.FrameBuffer) {
	e.bloom.Draw(buf)
	e.TextBox.Draw(buf)
}

// Update updates the entry box so it's interactive.
func (e *EntryBox) Update(win *turdgl.Window) {
	e.TextBox.Update(win)
}

// Text returns the text content of the entry box.
func (e *EntryBox) Text() string {
	return e.TextBox.Text.Text()
}

func NewTextBox(width, height float64, pos turdgl.Vec) *turdgl.TextBox {
	r := turdgl.NewCurvedRect(
		width, height, 6, pos,
		turdgl.WithStyle(turdgl.Style{Colour: color.RGBA{90, 65, 48, 255}, Thickness: 0}),
	)
	r.SetStyle(turdgl.Style{Colour: buttonColourUnpressed})

	tb := turdgl.NewTextBox(r, FontPathMedium).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetText("Click to edit").
		SetTextSize(36).
		SetTextColour(LightGreyTextColour)

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
		SetColour(LightGreyTextColour).
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
func (g *GameUIBox) Draw(buf *turdgl.FrameBuffer) {
	g.body.Draw(buf)
	g.heading.Draw(buf)
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
		SetColour(GreyTextColour).
		SetSize(17)
}

// NewLogoBox constructs a new "2048" tile logo.
func NewLogoBox(width, height float64, pos turdgl.Vec, txt string) *turdgl.TextBox {
	logo := turdgl.NewTextBox(
		turdgl.NewCurvedRect(width, height, 3, pos),
		FontPathBold,
	).
		SetText(txt).
		SetTextSize(32).
		SetTextColour(WhiteFontColour)

	logo.Text.SetAlignment(turdgl.AlignCustom).SetOffset(turdgl.Vec{Y: 27})
	logo.Shape.SetStyle(turdgl.Style{Colour: Tile2048Colour})

	return logo
}
