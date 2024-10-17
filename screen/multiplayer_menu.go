package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type MultiplayerMenuScreen struct {
	win *turdgl.Window

	title *turdgl.Text
	join  *common.MenuButton
	host  *common.MenuButton
	back  *common.MenuButton
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerMenuScreen(win *turdgl.Window) *MultiplayerMenuScreen {
	return &MultiplayerMenuScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerMenuScreen) Enter(_ InitData) {
	s.title = turdgl.NewText("Multiplayer", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.ArenaBackgroundColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	s.join = common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, func() { SetScreen(MultiplayerJoin, nil) })
	s.join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join game")

	s.host = common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(MultiplayerHost, nil) })
	s.host.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Host game")

	s.back = common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(Title, nil) })
	s.back.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	s.win.RegisterKeybind(turdgl.Key1, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerJoin, nil)
	})
	s.win.RegisterKeybind(turdgl.Key2, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerHost, nil)
	})
	s.win.RegisterKeybind(turdgl.Key3, turdgl.KeyRelease, func() {
		SetScreen(Title, nil)
	})
	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(Title, nil)
	})
}

// Exit deinitialises the screen.
func (s *MultiplayerMenuScreen) Exit() {
	s.win.UnregisterKeybind(turdgl.Key1, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key2, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key3, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer menu screen.
func (s *MultiplayerMenuScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range []*common.MenuButton{
		s.join,
		s.host,
		s.back,
	} {
		s.win.Draw(b)
		b.Update(s.win)
	}
}
