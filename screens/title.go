package screens

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/gogl"
)

type TitleScreen struct {
	win *gogl.Window

	title            *gogl.Text
	hint             *gogl.Text
	buttonBackground *gogl.CurvedRect
	singleplayer     *gogl.Button
	multiplayer      *gogl.Button
	quit             *gogl.Button
}

// NewTitle Screen constructs a new title screen for the given window.
func NewTitleScreen(win *gogl.Window) *TitleScreen {
	return &TitleScreen{win: win}
}

// Enter initialises the screen.
func (s *TitleScreen) Enter(_ InitData) {
	s.title = gogl.NewText("2048 Battle", gogl.Vec{X: config.WinWidth / 2, Y: 260}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(gogl.AlignCentre).
		SetSize(100)

	s.hint = gogl.NewText("", gogl.Vec{X: config.WinWidth / 2, Y: 375}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(gogl.AlignBottomCentre).
		SetSize(20)

	// Adjustable settings for buttons
	const (
		TileSizePx        float64 = 170
		TileCornerRadius  float64 = 6
		TileBoundryFactor float64 = 0.15
	)

	// Background for buttons
	const w = TileSizePx * (3 + 4*TileBoundryFactor)
	s.buttonBackground = gogl.NewCurvedRect(
		w, TileSizePx*(1+2*TileBoundryFactor), TileCornerRadius,
		gogl.Vec{X: (config.WinWidth - w) / 2, Y: 400},
	)
	s.buttonBackground.SetStyle(gogl.Style{Colour: common.ArenaBackgroundColour})

	// Menu buttons
	s.singleplayer = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			SetScreen(Singleplayer, nil)
		},
	).SetLabelText("Solo")
	s.singleplayer.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.singleplayer.Label.SetColour(common.WhiteFontColour)
			s.singleplayer.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Play alone")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.singleplayer.Label.SetColour(common.WhiteFontColour)
			s.singleplayer.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.multiplayer = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			SetScreen(MultiplayerMenu, nil)
		},
	).SetLabelText("Versus")
	s.multiplayer.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.multiplayer.Label.SetColour(common.WhiteFontColour)
			s.multiplayer.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Play against a friend")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.multiplayer.Label.SetColour(common.WhiteFontColour)
			s.multiplayer.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.quit = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(2+3*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			s.win.Quit()
		},
	).SetLabelText("Quit")
	s.quit.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.quit.Label.SetColour(common.WhiteFontColour)
			s.quit.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Exit to desktop")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.quit.Label.SetColour(common.WhiteFontColour)
			s.quit.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	// Keybinds
	s.win.RegisterKeybind(gogl.Key1, gogl.KeyRelease, func() {
		SetScreen(Singleplayer, nil)
	})
	s.win.RegisterKeybind(gogl.Key2, gogl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})
	s.win.RegisterKeybind(gogl.Key3, gogl.KeyRelease, s.win.Quit)
	s.win.RegisterKeybind(gogl.KeyEscape, gogl.KeyRelease, s.win.Quit)
}

// Exit deinitialises the screen.
func (s *TitleScreen) Exit() {
	s.win.UnregisterKeybind(gogl.Key1, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.Key2, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.Key3, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.KeyEscape, gogl.KeyRelease)
}

// Update draws the title screen and updates its components.
func (s *TitleScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)
	s.win.Draw(s.hint)
	s.win.Draw(s.buttonBackground)

	for _, b := range []*gogl.Button{
		s.singleplayer,
		s.multiplayer,
		s.quit,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}
}
