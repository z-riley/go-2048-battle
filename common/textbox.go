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

// Text returns the text content of the entry box.
func (e *EntryBox) SetText(text string) *EntryBox {
	e.TextBox.SetText(text)
	return e
}

func NewTextBox(width, height float64, pos turdgl.Vec) *turdgl.TextBox {
	r := turdgl.NewCurvedRect(
		width, height, 6, pos,
		turdgl.WithStyle(turdgl.Style{Colour: color.RGBA{90, 65, 48, 255}, Thickness: 0}),
	)
	r.SetStyle(turdgl.Style{Colour: buttonColourUnpressed})

	tb := turdgl.NewTextBox(r, FontPathMedium).
		SetTextOffset(turdgl.Vec{X: 0, Y: 15}).
		SetText("Click to edit").
		SetTextSize(36).
		SetTextColour(LightGreyTextColour)

	return tb
}

// ScoreBox is a commonly used text box for displaying scores.
type ScoreBox struct {
	heading *turdgl.Text
	body    *turdgl.TextBox
}

// NewScoreBox constructs a new text box for displaying a score.
func NewScoreBox(width, height float64, pos turdgl.Vec, colour color.RGBA) *ScoreBox {
	headingPos := turdgl.Vec{
		X: pos.X + width/2,
		Y: pos.Y + 25,
	}
	heading := turdgl.NewText("Heading", headingPos, FontPathBold).
		SetColour(LightGreyTextColour).
		SetSize(16).
		SetOffset(turdgl.Vec{Y: 3})

	r := turdgl.NewCurvedRect(
		width, height, 3,
		pos,
		turdgl.WithStyle(turdgl.Style{Colour: colour}),
	)

	body := turdgl.NewTextBox(r, FontPathBold).
		SetTextOffset(turdgl.Vec{X: 0, Y: 18}).
		SetTextSize(26).
		SetTextColour(WhiteFontColour)

	return &ScoreBox{heading, body}
}

// Draw draws the UI box to the window.
func (g *ScoreBox) Draw(buf *turdgl.FrameBuffer) {
	g.body.Draw(buf)
	g.heading.Draw(buf)
}

// SetHeading sets the heading text of the UI box.
func (g *ScoreBox) SetHeading(s string) *ScoreBox {
	g.heading.SetText(s)
	return g
}

// SetBody sets the body text of the UI box.
func (g *ScoreBox) SetBody(s string) *ScoreBox {
	g.body.SetText(s)
	return g
}

// NewGameText constructs text with sensible defaults.
func NewGameText(body string, pos turdgl.Vec) *turdgl.Text {
	return turdgl.NewText(body, pos, FontPathBold).
		SetColour(GreyTextColour).
		SetSize(17)
}

// NewLogoBox constructs a "2048" tile logo.
func NewLogoBox(size float64, pos turdgl.Vec, txt string) *turdgl.TextBox {
	logo := turdgl.NewTextBox(
		turdgl.NewCurvedRect(size, size, 3, pos),
		FontPathBold,
	).
		SetText(txt).
		SetTextSize(32).
		SetTextColour(WhiteFontColour)

	logo.Text.SetAlignment(turdgl.AlignCustom).SetOffset(turdgl.Vec{Y: 15})
	logo.Shape.SetStyle(turdgl.Style{Colour: Tile2048Colour})

	return logo
}
