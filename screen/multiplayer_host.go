package screen

import (
	"fmt"
	"image/color"

	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type MultiplayerHostScreen struct {
	win *turdgl.Window

	title       *turdgl.Text
	buttons     []*common.MenuButton
	entries     []*common.EntryBox
	playerCards []*playerCard
}

// NewTitle Screen constructs a new multiplayer host screen for the given window.
func NewMultiplayerHostScreen(win *turdgl.Window) *MultiplayerHostScreen {
	title := turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	nameHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Your name:")
	nameEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})

	var cards []*playerCard
	for i := 0; i < 4; i++ {
		pos := turdgl.Vec{X: 300, Y: 300 + float64(i)*80}
		label := fmt.Sprintf("Waiting for player %d", i+1)
		cards = append(cards, newPlayerCard(pos, label))
	}

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 650}, func() { SetScreen(Title) })
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 600 + 20, Y: 650}, func() { SetScreen(MultiplayerMenu) })
	back.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	return &MultiplayerHostScreen{
		win,
		title,
		[]*common.MenuButton{nameHeading, join, back},
		[]*common.EntryBox{nameEntry},
		cards,
	}
}

// Init initialises the screen.
func (s *MultiplayerHostScreen) Init() {}

// Deinit deinitialises the screen.
func (s *MultiplayerHostScreen) Deinit() {}

// Update updates and draws multiplayer host screen.
func (s *MultiplayerHostScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range s.buttons {
		b.Update(s.win)
		s.win.Draw(b)
	}

	for _, e := range s.entries {
		s.win.Draw(e)
		e.Update(s.win)
	}

	for _, p := range s.playerCards {
		s.win.Draw(p)
	}

}

var (
	styleNotReady = turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 2}
	styleReady    = turdgl.Style{Colour: color.RGBA{0, 255, 0, 255}, Thickness: 0, Bloom: 10}
)

type playerCard struct {
	name  *common.EntryBox
	light *turdgl.Circle
}

func newPlayerCard(pos turdgl.Vec, txt string) *playerCard {
	const (
		width  = 400
		height = 60
	)

	name := common.NewEntryBox(width, height, pos)
	name.SetTextAlignment(turdgl.AlignCustom).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetTextSize(30).
		SetTextColour(common.DarkerFontColour).
		SetText(txt)

	lightPos := turdgl.Vec{X: pos.X + width + 40, Y: pos.Y + height/2}
	light := turdgl.NewCircle(height*0.8, lightPos, turdgl.WithStyle(styleNotReady))

	return &playerCard{name, light}
}

func (p *playerCard) Draw(buf *turdgl.FrameBuffer) {
	p.name.Draw(buf)
	p.light.Draw(buf)
}

func (p *playerCard) setReady() {
	p.light.SetStyle(styleReady)
}

func (p *playerCard) setNotReady() {
	p.light.SetStyle(styleNotReady)
}
