package screens

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/gogl"
)

type MultiplayerMenuScreen struct {
	win *gogl.Window

	title *gogl.Text

	hint             *gogl.Text
	buttonBackground *gogl.CurvedRect
	join             *gogl.Button
	host             *gogl.Button
	back             *gogl.Button
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerMenuScreen(win *gogl.Window) *MultiplayerMenuScreen {
	return &MultiplayerMenuScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerMenuScreen) Enter(_ InitData) {
	s.title = gogl.NewText("Versus", gogl.Vec{X: config.WinWidth / 2, Y: 260}, common.FontPathMedium).
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

	s.join = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(MultiplayerJoin, nil) },
	).SetLabelText("Join")
	s.join.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.join.Label.SetColour(common.WhiteFontColour)
			s.join.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Join a LAN game")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.join.Label.SetColour(common.WhiteFontColour)
			s.join.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.host = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(MultiplayerHost, nil) },
	).SetLabelText("Host")
	s.host.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.host.Label.SetColour(common.WhiteFontColour)
			s.host.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Host a LAN game")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.host.Label.SetColour(common.WhiteFontColour)
			s.host.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.back = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(2+3*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(Title, nil) },
	).SetLabelText("Back")
	s.back.SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnHold},
		func() {
			s.back.Label.SetColour(common.WhiteFontColour)
			s.back.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Go back to main menu")
		},
	).SetCallback(
		gogl.ButtonTrigger{State: gogl.NoClick, Behaviour: gogl.OnRelease},
		func() {
			s.back.Label.SetColour(common.WhiteFontColour)
			s.back.Shape.(*gogl.CurvedRect).SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.win.RegisterKeybind(gogl.Key1, gogl.KeyRelease, func() {
		SetScreen(MultiplayerJoin, nil)
	})
	s.win.RegisterKeybind(gogl.Key2, gogl.KeyRelease, func() {
		SetScreen(MultiplayerHost, nil)
	})
	s.win.RegisterKeybind(gogl.Key3, gogl.KeyRelease, func() {
		SetScreen(Title, nil)
	})
	s.win.RegisterKeybind(gogl.KeyEscape, gogl.KeyRelease, func() {
		SetScreen(Title, nil)
	})
}

// Exit deinitialises the screen.
func (s *MultiplayerMenuScreen) Exit() {
	s.win.UnregisterKeybind(gogl.Key1, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.Key2, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.Key3, gogl.KeyRelease)
	s.win.UnregisterKeybind(gogl.KeyEscape, gogl.KeyRelease)
}

// Update updates and draws multiplayer menu screen.
func (s *MultiplayerMenuScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)
	s.win.Draw(s.hint)
	s.win.Draw(s.buttonBackground)

	for _, b := range []*gogl.Button{
		s.join,
		s.host,
		s.back,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}
}
