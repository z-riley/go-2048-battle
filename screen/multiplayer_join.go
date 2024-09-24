package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type MultiplayerJoinScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*common.MenuButton
	entries []*common.EntryBox
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerJoinScreen(win *turdgl.Window) *MultiplayerJoinScreen {
	title := turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	ipHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	ipHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Host IP:")

	ipEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})

	nameHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 300}, func() {})
	nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Your name:")

	nameEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 300})

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(Title) })
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(MultiplayerMenu) })
	back.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	return &MultiplayerJoinScreen{
		win,
		title,
		[]*common.MenuButton{ipHeading, nameHeading, join, back},
		[]*common.EntryBox{nameEntry, ipEntry},
	}
}

// Init initialises the screen.
func (s *MultiplayerJoinScreen) Init() {}

// Deinit deinitialises the screen.
func (s *MultiplayerJoinScreen) Deinit() {}

// Update updates and draws multiplayer join screen.
func (s *MultiplayerJoinScreen) Update() {
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

}
