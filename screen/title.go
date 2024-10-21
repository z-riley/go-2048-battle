package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type TitleScreen struct {
	win *turdgl.Window

	title        *turdgl.Text
	singleplayer *turdgl.Button
	multiplayer  *turdgl.Button
	quit         *turdgl.Button
}

// NewTitle Screen constructs a new title screen for the given window.
func NewTitleScreen(win *turdgl.Window) *TitleScreen {
	return &TitleScreen{win: win}
}

// Enter initialises the screen.
func (s *TitleScreen) Enter(_ InitData) {
	// Main title
	s.title = turdgl.NewText("2048 Battle", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.ArenaBackgroundColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	// Menu buttons
	s.singleplayer = common.NewMenuButton(
		400, 60,
		turdgl.Vec{X: 400, Y: 300},
		func() { SetScreen(Singleplayer, nil) },
	).SetLabelText("Singleplayer")
	s.multiplayer = common.NewMenuButton(
		400, 60,
		turdgl.Vec{X: 400, Y: 400},
		func() { SetScreen(MultiplayerMenu, nil) },
	).SetLabelText("Multiplayer")
	s.quit = common.NewMenuButton(
		400, 60, turdgl.Vec{X: 400, Y: 500},
		s.win.Quit,
	).SetLabelText("Quit")

	// Keybinds
	s.win.RegisterKeybind(turdgl.Key1, turdgl.KeyRelease, func() {
		SetScreen(Singleplayer, nil)
	})
	s.win.RegisterKeybind(turdgl.Key2, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})
	s.win.RegisterKeybind(turdgl.Key3, turdgl.KeyRelease, s.win.Quit)
}

// Exit deinitialises the screen.
func (s *TitleScreen) Exit() {
	s.win.UnregisterKeybind(turdgl.Key1, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key2, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key3, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update draws the title screen and updates its components.
func (s *TitleScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range []*turdgl.Button{
		s.singleplayer,
		s.multiplayer,
		s.quit,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}
}
