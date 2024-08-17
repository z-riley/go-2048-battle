package screen

import (
	"image/color"

	game "github.com/zac460/go-2048-battle"
	"github.com/zac460/go-2048-battle/common"
	"github.com/zac460/turdgl"
)

type MultiplayerJoinScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*common.MenuButton
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerJoinScreen(win *turdgl.Window) *MultiplayerJoinScreen {
	title := turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.ButtonFont).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	a := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, win.Quit)
	a.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Host IP: (TODO textbox ->)")

	b := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, win.Quit)
	b.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Your name: (TODO textbox ->)")

	c := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(Title) })
	c.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Join")

	d := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(MultiplayerMenu) })
	d.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Back")

	return &MultiplayerJoinScreen{
		win,
		title,
		[]*common.MenuButton{a, b, c, d},
	}
}

// Update updates and draws multiplayer join screen.
func (t *MultiplayerJoinScreen) Update() {
	t.win.SetBackground(color.RGBA{46, 36, 27, 255})

	t.win.Draw(t.title)

	for _, b := range t.buttons {
		t.win.Draw(b)
		b.Update(t.win)
	}
}
