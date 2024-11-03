package screen

import (
	"fmt"

	"github.com/brunoga/deep"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
)

type SingleplayerScreen struct {
	win *turdgl.Window

	backend      *backend.Game
	arena        *common.Arena
	arenaInputCh chan func()

	heading    *turdgl.Text
	loseDialog *turdgl.Text
	logo2048   *turdgl.TextBox
	score      *common.ScoreBox
	highScore  *common.ScoreBox
	menu       *turdgl.Button
	newGame    *turdgl.Button
	guide      *turdgl.Text
	timer      *turdgl.Text

	debugGrid  *turdgl.Text
	debugTime  *turdgl.Text
	debugScore *turdgl.Text
}

// NewSingleplayerScreen constructs an uninitialised new singleplayer menu screen.
func NewSingleplayerScreen(win *turdgl.Window) *SingleplayerScreen {
	return &SingleplayerScreen{win: win}
}

// Enter initialises the screen.
func (s *SingleplayerScreen) Enter(_ InitData) {
	// Arena and supporting data structures
	{
		s.arena = common.NewArena(turdgl.Vec{X: 440, Y: 300})
		s.backend = backend.NewGame(nil)
		s.arenaInputCh = make(chan func(), 100)
	}

	// UI components
	{
		// Everything is sized relative to the tile size
		const unit = common.TileSizePx

		// Everything is positioned relative to the arena grid
		anchor := s.arena.Pos()

		s.heading = turdgl.NewText(
			"", // to be set and drawn when player loses
			turdgl.Vec{X: anchor.X + s.arena.Width()/2, Y: anchor.Y - 2.3*unit},
			common.FontPathBold,
		).SetSize(40).SetColour(common.GreyTextColour).SetAlignment(turdgl.AlignTopCentre)

		s.loseDialog = turdgl.NewText(
			"", // to be set and drawn when player loses
			turdgl.Vec{X: anchor.X + s.arena.Width()/2, Y: anchor.Y - 1.7*unit},
			common.FontPathBold,
		).SetSize(20).SetColour(common.GreyTextColour).SetAlignment(turdgl.AlignTopCentre)

		s.logo2048 = common.NewLogoBox(
			1.36*unit,
			turdgl.Vec{X: anchor.X, Y: anchor.Y - 2.58*unit},
		)

		const wScore = 90
		s.score = common.NewScoreBox(
			wScore, wScore,
			turdgl.Vec{X: anchor.X + s.arena.Width() - 2.74*unit, Y: anchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE")

		s.highScore = common.NewScoreBox(
			wScore, wScore,
			turdgl.Vec{X: anchor.X + s.arena.Width() - wScore, Y: anchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("BEST")

		const buttonWidth = unit * 1.27
		s.menu = common.NewGameButton(
			buttonWidth, 0.4*unit,
			turdgl.Vec{X: anchor.X + s.arena.Width() - buttonWidth, Y: anchor.Y - 1.21*unit},
			func() {
				SetScreen(Title, nil)
			},
		).SetLabelText("MENU")

		s.newGame = common.NewGameButton(
			buttonWidth, 0.4*unit,
			turdgl.Vec{X: anchor.X + s.arena.Width() - 2.74*unit, Y: anchor.Y - 1.21*unit},
			func() {
				s.arenaInputCh <- func() {
					s.backend.Reset()
					s.arena.Reset()
				}
			},
		).SetLabelText("NEW")

		s.guide = turdgl.NewText(
			"Join the numbers and get to the 2048 tile!",
			turdgl.Vec{X: anchor.X, Y: anchor.Y - 0.53*unit},
			common.FontPathBold,
		).SetSize(16).SetColour(common.GreyTextColour)

		s.timer = common.NewGameText("",
			turdgl.Vec{X: anchor.X + s.arena.Width(), Y: anchor.Y + s.arena.Height()*1.1},
		).SetSize(16).SetAlignment(turdgl.AlignBottomRight)
	}

	// Debug UI
	if config.Debug {
		s.debugGrid = turdgl.NewText("grid", turdgl.Vec{X: 930, Y: 600}, common.FontPathMedium).
			SetText(s.backend.Grid.Debug())
		s.debugTime = turdgl.NewText("time", turdgl.Vec{X: 1100, Y: 550}, common.FontPathMedium).
			SetText(s.backend.Timer.Time.String())
		s.debugScore = turdgl.NewText("score", turdgl.Vec{X: 950, Y: 550}, common.FontPathMedium).
			SetText(fmt.Sprint(s.backend.Score.Current))
	}

	// Set keybinds. User inputs are sent to the backend via a buffered channel
	// so the backend game cannot execute multiple moves before the frontend has
	// finished animating the first one
	{
		s.win.RegisterKeybind(turdgl.KeyUp, turdgl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirUp)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(turdgl.KeyDown, turdgl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirDown)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(turdgl.KeyLeft, turdgl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirLeft)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(turdgl.KeyRight, turdgl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirRight)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(turdgl.KeyR, turdgl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.Reset()
				s.arena.Reset()
			}
		})
		s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() {
			SetScreen(Title, nil)
		})
	}
}

// Exit deinitialises the screen.
func (s *SingleplayerScreen) Exit() {
	s.backend.Timer.Pause()

	if err := s.backend.Save(); err != nil {
		panic(err)
	}

	s.win.UnregisterKeybind(turdgl.KeyUp, turdgl.KeyPress)
	s.win.UnregisterKeybind(turdgl.KeyDown, turdgl.KeyPress)
	s.win.UnregisterKeybind(turdgl.KeyLeft, turdgl.KeyPress)
	s.win.UnregisterKeybind(turdgl.KeyRight, turdgl.KeyPress)
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyPress)

	s.arena.Destroy()
}

// Update updates and draws the singleplayer screen.
func (s *SingleplayerScreen) Update() {
	// Temporary debug text
	if config.Debug {
		s.debugGrid.SetText(s.backend.Grid.Debug())
		s.debugTime.SetText(s.backend.Timer.Time.String())
		s.debugScore.SetText(
			fmt.Sprint(s.backend.Score.Current, "|", s.backend.Score.High),
		)
	}

	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
	default:
		// No user input; continue
	}

	// Deep copy so front-end hastime to animate itself whilst allowing the back
	// end to update
	game := deep.MustCopy(*s.backend)

	// Check for win or lose
	switch game.Grid.Outcome() {
	case grid.Win:
		s.updateWin(game)
	case grid.Lose:
		s.updateLose(game)
	default:
		s.updateNormal(game)
	}

	// Draw temporary debug grid
	if config.Debug {
		s.win.Draw(s.debugGrid)
		s.win.Draw(s.debugTime)
		s.win.Draw(s.debugScore)
	}
}

// Update updates and draws the singleplayer screen in a normal state.
func (s *SingleplayerScreen) updateNormal(game backend.Game) {
	s.win.SetBackground(common.BackgroundColour)

	s.score.SetBody(fmt.Sprint(game.Score.Current))
	s.menu.Update(s.win)
	s.highScore.SetBody(fmt.Sprint(game.Score.High))
	s.timer.SetText(game.Timer.Time.String())
	s.newGame.Update(s.win)

	s.arena.SetNormal()
	s.arena.Update(game)

	for _, d := range []turdgl.Drawable{
		s.logo2048,
		s.score,
		s.highScore,
		s.menu,
		s.newGame,
		s.guide,
		s.timer,
		s.arena,
	} {
		s.win.Draw(d)
	}
}

// updateWin updates and draws the singleplayer screen in a winning state.
func (s *SingleplayerScreen) updateWin(game backend.Game) {
	s.guide.SetText(
		fmt.Sprintf("Your next goal is to get to the %d tile!", game.Grid.HighestTile()*2),
	)
	s.updateNormal(game)
}

// updateLose updates and draws the singleplayer screen in a losing state.
func (s *SingleplayerScreen) updateLose(game backend.Game) {
	s.win.SetBackground(common.BackgroundColour)
	s.arena.SetLose()

	s.heading.SetText("Game over!")
	s.loseDialog.SetText(fmt.Sprintf(
		"You earned %d points in %v.", game.Score.Current, game.Timer.Duration(),
	))

	s.menu.Update(s.win)
	s.newGame.Update(s.win)
	s.arena.Update(game)

	for _, d := range []turdgl.Drawable{
		s.heading,
		s.loseDialog,
		s.menu,
		s.newGame,
		s.arena,
	} {
		s.win.Draw(d)
	}
}
