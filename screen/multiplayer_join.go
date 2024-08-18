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
	entries []*common.EntryBox
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerJoinScreen(win *turdgl.Window) *MultiplayerJoinScreen {
	title := turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, game.FontPath).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	ipHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	ipHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).SetLabelText("Host IP:")

	ipEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})

	nameHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 300}, func() {})
	nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).SetLabelText("Your name:")

	nameEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 300})

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(Title) })
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).SetLabelText("Join")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(MultiplayerMenu) })
	back.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 30}).SetLabelText("Back")

	return &MultiplayerJoinScreen{
		win,
		title,
		[]*common.MenuButton{ipHeading, nameHeading, join, back},
		[]*common.EntryBox{nameEntry, ipEntry},
	}
}

// Update updates and draws multiplayer join screen.
func (t *MultiplayerJoinScreen) Update() {
	t.win.SetBackground(color.RGBA{46, 36, 27, 255})

	t.win.Draw(t.title)

	for _, b := range t.buttons {
		b.Update(t.win)
		t.win.Draw(b)
	}

	for _, e := range t.entries {
		t.win.Draw(e)
		e.Update(t.win)
	}

}
