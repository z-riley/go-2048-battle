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
func NewEntryBox(width, height float64, pos turdgl.Vec, txt string) *EntryBox {
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

	bloom := turdgl.NewCurvedRect(width, height, 6, pos).SetStyle(styleUnselected)

	tb := NewTextBox(width, height, pos, txt).SetTextAlignment(turdgl.AlignCentre)
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

// SetModifiedCB sets a callback which is executed when the text in the entry
// box is modified.
func (e *EntryBox) SetModifiedCB(callback func()) *EntryBox {
	e.TextBox.SetCallback(callback)
	return e
}

// NewTextBox constructs a new text box.
func NewTextBox(width, height float64, pos turdgl.Vec, txt string) *turdgl.TextBox {
	r := turdgl.NewCurvedRect(width, height, 6, pos).
		SetStyle(turdgl.Style{Colour: buttonColourUnpressed})

	return turdgl.NewTextBox(r, txt, FontPathMedium).
		SetTextOffset(turdgl.Vec{X: 0, Y: 15}).
		SetTextSize(36).
		SetTextColour(LightGreyTextColour)
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
		Y: pos.Y + 15,
	}
	heading := turdgl.NewText("Heading", headingPos, FontPathBold).
		SetColour(LightGreyTextColour).
		SetSize(16).
		SetOffset(turdgl.Vec{Y: 3})

	r := turdgl.NewCurvedRect(
		width, height, 3,
		pos,
	).SetStyle(turdgl.Style{Colour: colour})

	body := turdgl.NewTextBox(r, "", FontPathBold).
		SetTextOffset(turdgl.Vec{X: 0, Y: 10}).
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
func NewLogoBox(size float64, pos turdgl.Vec) *turdgl.TextBox {
	logo := turdgl.NewTextBox(
		turdgl.NewCurvedRect(size, size, 3, pos),
		"2048",
		FontPathBold,
	).
		SetTextSize(32).
		SetTextColour(WhiteFontColour)

	logo.Text.SetAlignment(turdgl.AlignCustom)
	logo.Shape.(*turdgl.CurvedRect).SetStyle(turdgl.Style{Colour: Tile2048Colour})

	return logo
}

// NewLogoBox returns a new tooltip text box.
func NewTooltip() *turdgl.TextBox {
	r := turdgl.NewCurvedRect(110, 23, 2, turdgl.Vec{}).
		SetStyle(turdgl.Style{Colour: ArenaBackgroundColour})

	return turdgl.NewTextBox(r, "Click to edit", FontPathMedium).
		SetTextOffset(turdgl.Vec{Y: 4}).
		SetTextSize(16).
		SetTextColour(LightGreyTextColour)
}
