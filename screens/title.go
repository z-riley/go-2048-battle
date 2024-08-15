package screens

import (
	"image/color"

	game "github.com/zac460/go-2048-battle"
	"github.com/zac460/go-2048-battle/common"
	"github.com/zac460/turdgl"
)

type TitleScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*turdgl.Button
}

func NewTitleScreen(win *turdgl.Window) *TitleScreen {
	title := turdgl.NewText("2048 Battle", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.ButtonFont).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	btnSingleplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, win.Quit).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Singleplayer")

	btnMultiplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, win.Quit).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Multiplayer")

	btnQuit := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, win.Quit).
		SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).
		SetLabelSize(36).
		SetLabelColour(common.ButtonFont).
		SetLabelText("Quit")

	return &TitleScreen{
		win,
		title,
		[]*turdgl.Button{btnSingleplayer, btnMultiplayer, btnQuit},
	}
}

// Update draws the title screen and updates its components.
func (t *TitleScreen) Update() {
	t.win.SetBackground(color.RGBA{46, 36, 27, 255})

	t.win.Draw(t.title)

	for _, b := range t.buttons {
		t.win.Draw(b)
		b.Update(t.win)
	}
}
