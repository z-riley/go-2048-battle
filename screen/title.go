package screen

import (
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
	title := turdgl.NewText("2048 Battle", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.ArenaBackgroundColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	// Menu buttons
	singleplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 300}, func() { SetScreen(Singleplayer, nil) })
	singleplayer.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Singleplayer")
	multiplayer := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() { SetScreen(MultiplayerMenu, nil) })
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

// Init initialises the screen.
func (s *TitleScreen) Init(_ InitData) {
	s.win.RegisterKeybind(turdgl.Key1, turdgl.KeyRelease, func() {
		SetScreen(Singleplayer, nil)
	})
	s.win.RegisterKeybind(turdgl.Key2, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})
	s.win.RegisterKeybind(turdgl.Key3, turdgl.KeyRelease, s.win.Quit)
	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(Title, nil)
	})
}

// Deinit deinitialises the screen.
func (s *TitleScreen) Deinit() {
	s.win.UnregisterKeybind(turdgl.Key1, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key2, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.Key3, turdgl.KeyRelease)
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
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
