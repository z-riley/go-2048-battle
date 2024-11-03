package screen

import (
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type MultiplayerMenuScreen struct {
	win *turdgl.Window

	title *turdgl.Text

	hint             *turdgl.Text
	buttonBackground *turdgl.CurvedRect
	join             *turdgl.Button
	host             *turdgl.Button
	back             *turdgl.Button
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerMenuScreen(win *turdgl.Window) *MultiplayerMenuScreen {
	return &MultiplayerMenuScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerMenuScreen) Enter(_ InitData) {
	s.title = turdgl.NewText("Versus", turdgl.Vec{X: 600, Y: 260}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(100)

	s.hint = turdgl.NewText("", turdgl.Vec{X: 600, Y: 375}, common.FontPathMedium).
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
		turdgl.Vec{X: (float64(s.win.Width()) - w) / 2, Y: 400},
	)
	s.buttonBackground.SetStyle(turdgl.Style{Colour: common.ArenaBackgroundColour})

	s.join = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(MultiplayerJoin, nil) },
	).SetLabelText("Join")
	s.join.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.join.Label.SetColour(common.WhiteFontColour)
			s.join.Shape.SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Join a LAN game")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.join.Label.SetColour(common.WhiteFontColour)
			s.join.Shape.SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.host = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(MultiplayerHost, nil) },
	).SetLabelText("Host")
	s.host.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.host.Label.SetColour(common.WhiteFontColour)
			s.host.Shape.SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Host a LAN game")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.host.Label.SetColour(common.WhiteFontColour)
			s.host.Shape.SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

	s.back = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(2+3*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		}.Round(),
		func() { SetScreen(Title, nil) },
	).SetLabelText("Back")
	s.back.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnHold},
		func() {
			s.back.Label.SetColour(common.WhiteFontColour)
			s.back.Shape.SetStyle(common.ButtonStyleHovering)
			s.hint.SetText("Go back to main menu")
		},
	).SetCallback(
		turdgl.ButtonTrigger{State: turdgl.NoClick, Behaviour: turdgl.OnRelease},
		func() {
			s.back.Label.SetColour(common.WhiteFontColour)
			s.back.Shape.SetStyle(common.ButtonStyleUnpressed)
			s.hint.SetText("")
		},
	)

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
	s.win.Draw(s.hint)
	s.win.Draw(s.buttonBackground)

	for _, b := range []*turdgl.Button{
		s.join,
		s.host,
		s.back,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}
}
