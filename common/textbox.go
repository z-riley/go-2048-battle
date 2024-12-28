package common

import (
	"image/color"

	"github.com/z-riley/gogl"
)

// Entrybox is an interactive text box for data entry.
type EntryBox struct {
	TextBox *gogl.TextBox
	bloom   *gogl.CurvedRect
}

// NewEntryBox constructs a new text box with suitable defaults.
func NewEntryBox(width, height float64, pos gogl.Vec, txt string) *EntryBox {
	var (
		styleUnselected = gogl.Style{
			Colour:    gogl.LightGrey,
			Thickness: 1,
			Bloom:     0,
		}
		styleSelected = gogl.Style{
			Colour:    gogl.LightGrey,
			Thickness: 1,
			Bloom:     10,
		}
	)

	bloom := gogl.NewCurvedRect(width, height, 6, pos).SetStyle(styleUnselected)

	tb := NewTextBox(width, height, pos, txt).SetTextAlignment(gogl.AlignCentre)
	tb.SetSelectedCB(func() {
		tb.SetTextColour(gogl.White)
		bloom.SetStyle(styleSelected)
	}).SetDeselectedCB(func() {
		tb.SetTextColour(LightGreyTextColour)
		bloom.SetStyle(styleUnselected)
	})

	return &EntryBox{tb, bloom}
}

// Draw draws an entry box to the frame buffer.
func (e *EntryBox) Draw(buf *gogl.FrameBuffer) {
	e.bloom.Draw(buf)
	e.TextBox.Draw(buf)
}

// Update updates the entry box so it's interactive.
func (e *EntryBox) Update(win *gogl.Window) {
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
func NewTextBox(width, height float64, pos gogl.Vec, txt string) *gogl.TextBox {
	r := gogl.NewCurvedRect(width, height, 6, pos).
		SetStyle(gogl.Style{Colour: buttonColourUnpressed})

	return gogl.NewTextBox(r, txt, FontPathMedium).
		SetTextOffset(gogl.Vec{X: 0, Y: 15}).
		SetTextSize(36).
		SetTextColour(LightGreyTextColour)
}

// ScoreBox is a commonly used text box for displaying scores.
type ScoreBox struct {
	heading *gogl.Text
	body    *gogl.TextBox
}

// NewScoreBox constructs a new text box for displaying a score.
func NewScoreBox(width, height float64, pos gogl.Vec, colour color.RGBA) *ScoreBox {
	headingPos := gogl.Vec{
		X: pos.X + width/2,
		Y: pos.Y + 15,
	}
	heading := gogl.NewText("Heading", headingPos, FontPathBold).
		SetColour(LightGreyTextColour).
		SetSize(16).
		SetOffset(gogl.Vec{Y: 3})

	r := gogl.NewCurvedRect(
		width, height, 3,
		pos,
	).SetStyle(gogl.Style{Colour: colour})

	body := gogl.NewTextBox(r, "", FontPathBold).
		SetTextOffset(gogl.Vec{X: 0, Y: 10}).
		SetTextSize(26).
		SetTextColour(WhiteFontColour)

	return &ScoreBox{heading, body}
}

// Draw draws the UI box to the window.
func (g *ScoreBox) Draw(buf *gogl.FrameBuffer) {
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
func NewGameText(body string, pos gogl.Vec) *gogl.Text {
	return gogl.NewText(body, pos, FontPathBold).
		SetColour(GreyTextColour).
		SetSize(17)
}

// NewLogoBox constructs a "2048" tile logo.
func NewLogoBox(size float64, pos gogl.Vec) *gogl.TextBox {
	logo := gogl.NewTextBox(
		gogl.NewCurvedRect(size, size, 3, pos),
		"2048",
		FontPathBold,
	).
		SetTextSize(32).
		SetTextColour(WhiteFontColour)

	logo.Text.SetAlignment(gogl.AlignCustom)
	logo.Shape.(*gogl.CurvedRect).SetStyle(gogl.Style{Colour: Tile2048Colour})

	return logo
}

// NewLogoBox returns a new tooltip text box.
func NewTooltip() *gogl.TextBox {
	r := gogl.NewCurvedRect(110, 23, 2, gogl.Vec{}).
		SetStyle(gogl.Style{Colour: ArenaBackgroundColour})

	return gogl.NewTextBox(r, "Click to edit", FontPathMedium).
		SetTextSize(16).
		SetTextColour(LightGreyTextColour)
}
