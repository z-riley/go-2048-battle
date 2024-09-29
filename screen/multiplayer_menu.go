package screen

import (
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
	title := turdgl.NewText("Multiplayer", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, func() { SetScreen(MultiplayerJoin, nil) })
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join game")

	host := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(MultiplayerHost, nil) })
	host.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Host game")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() { SetScreen(Title, nil) })
	back.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	// TODO: poke around 2048 app UI and copy that menu style

	return &MultiplayerMenuScreen{
		win,
		title,
		[]*common.MenuButton{join, host, back},
	}
}

// Init initialises the screen.
func (s *MultiplayerMenuScreen) Init(_ InitData) {
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

// Deinit deinitialises the screen.
func (s *MultiplayerMenuScreen) Deinit() {
	s.win.UnregisterKeybind(turdgl.Key1, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key2, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key3, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer menu screen.
func (s *MultiplayerMenuScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range s.buttons {
		s.win.Draw(b)
		b.Update(s.win)
	}
}
