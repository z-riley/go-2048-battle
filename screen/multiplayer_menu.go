package screen

import (
	"image/color"

	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type MultiplayerMenuScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*common.MenuButton
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerMenuScreen(win *turdgl.Window) *MultiplayerMenuScreen {
	title := turdgl.NewText("Multiplayer", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, func() { SetScreen(MultiplayerJoin) })
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join game")

	host := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(MultiplayerMenu) })
	host.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Host game")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(Title) })
	back.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	return &MultiplayerMenuScreen{
		win,
		title,
		[]*common.MenuButton{join, host, back},
	}
}

// Update updates and draws multiplayer menu screen.
func (t *MultiplayerMenuScreen) Update() {
	t.win.SetBackground(color.RGBA{46, 36, 27, 255})

	t.win.Draw(t.title)

	for _, b := range t.buttons {
		t.win.Draw(b)
		b.Update(t.win)
	}
}
