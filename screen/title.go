package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
)

type TitleScreen struct {
	win *turdgl.Window

	title            *turdgl.Text
	hint             *turdgl.Text
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
	s.title = turdgl.NewText("2048 Battle", turdgl.Vec{X: config.WinWidth / 2, Y: 260}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(100)

	s.hint = turdgl.NewText("", turdgl.Vec{X: config.WinWidth / 2, Y: 375}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignBottomCentre).
		SetSize(20)

	// Adjustable settings for buttons
	const (
		TileSizePx        float64 = 170
		TileCornerRadius  float64 = 6
		TileBoundryFactor float64 = 0.15
	)

	// Background for buttons
	const w = TileSizePx * (3 + 4*TileBoundryFactor)
	s.buttonBackground = turdgl.NewCurvedRect(
		w, TileSizePx*(1+2*TileBoundryFactor), TileCornerRadius,
		turdgl.Vec{X: (config.WinWidth - w) / 2, Y: 400},
	)
	s.buttonBackground.SetStyle(turdgl.Style{Colour: common.ArenaBackgroundColour})

	// Menu buttons
	s.singleplayer = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			SetScreen(Singleplayer, nil)

		},
	).SetLabelText("Solo")
	s.singleplayer.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.singleplayer.Label.SetColour(common.WhiteFontColour)
			s.singleplayer.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Play alone")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.singleplayer.Label.SetColour(common.WhiteFontColour)
			s.singleplayer.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.multiplayer = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			SetScreen(MultiplayerMenu, nil)
		},
	).SetLabelText("Versus")
	s.multiplayer.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.multiplayer.Label.SetColour(common.WhiteFontColour)
			s.multiplayer.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Play against a friend")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.multiplayer.Label.SetColour(common.WhiteFontColour)
			s.multiplayer.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.quit = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(2+3*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			s.win.Quit()
		},
	).SetLabelText("Quit")
	s.quit.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.quit.Label.SetColour(common.WhiteFontColour)
			s.quit.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Exit to desktop")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.quit.Label.SetColour(common.WhiteFontColour)
			s.quit.Shape.(*turdgl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

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
	s.win.Draw(s.hint)
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
