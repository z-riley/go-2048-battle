package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type TitleScreen struct {
	win *turdgl.Window

	title *turdgl.Text

	buttonBackground *turdgl.CurvedRect
	singleplayer     *turdgl.Button
	multiplayer      *turdgl.Button
	quit             *turdgl.Button
}

// NewTitle Screen constructs a new title screen for the given window.
func NewTitleScreen(win *turdgl.Window) *TitleScreen {
	return &TitleScreen{win: win}
}

// Enter initialises the screen.
func (s *TitleScreen) Enter(_ InitData) {
	// Main title
	s.title = turdgl.NewText("2048 Battle", turdgl.Vec{X: 600, Y: 220}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(80)

	// Background for buttons
	const buttonSize = 180
	var w float64 = 600
	s.buttonBackground = turdgl.NewCurvedRect(
		w, buttonSize+30, common.TileCornerRadius,
		turdgl.Vec{X: (float64(s.win.Width()) - w) / 2, Y: 300},
	)
	s.buttonBackground.SetStyle(turdgl.Style{Colour: common.ArenaBackgroundColour})

	// Menu buttons
	s.singleplayer = common.NewMenuButton(
		buttonSize, buttonSize,
		turdgl.Vec{X: s.buttonBackground.Pos.X + 15, Y: s.buttonBackground.Pos.X + 15},
		func() { SetScreen(Singleplayer, nil) },
	).SetLabelText("Solo")

	s.multiplayer = common.NewMenuButton(
		buttonSize, buttonSize,
		turdgl.Vec{X: 600 - buttonSize/2, Y: s.buttonBackground.Pos.X + 15},
		func() { SetScreen(MultiplayerMenu, nil) },
	).SetLabelText("Versus")

	s.quit = common.NewMenuButton(
		buttonSize, buttonSize,
		turdgl.Vec{X: 900 - buttonSize - 15, Y: s.buttonBackground.Pos.X + 15},
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
	s.win.Draw(s.buttonBackground)

	for _, b := range []*turdgl.Button{
		s.singleplayer,
		s.multiplayer,
		s.quit,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}
}
