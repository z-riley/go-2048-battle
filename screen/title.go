package screen

import (
	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type TitleScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*common.MenuButton
}

// NewTitle Screen constructs a new title screen for the given window.
func NewTitleScreen(win *turdgl.Window) *TitleScreen {
	// Main title
	title := turdgl.NewText("2048 Battle", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	// Menu buttons
	singleplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, func() { SetScreen(Singleplayer) })
	singleplayer.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Singleplayer")
	multiplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(MultiplayerMenu) })
	multiplayer.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Multiplayer")
	quit := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, win.Quit)
	quit.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Quit")

	return &TitleScreen{
		win,
		title,
		[]*common.MenuButton{singleplayer, multiplayer, quit},
	}
}

// Update draws the title screen and updates its components.
func (s *TitleScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range s.buttons {
		s.win.Draw(b)
		b.Update(s.win)
	}
}
