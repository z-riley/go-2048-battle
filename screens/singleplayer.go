package screens

import (
	"fmt"
	"strconv"

	"github.com/brunoga/deep"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/common/backend"
	"github.com/z-riley/go-2048-battle/common/backend/grid"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/gogl"
)

type SingleplayerScreen struct {
	win *gogl.Window

	backend      *backend.Game
	arena        *common.Arena
	arenaInputCh chan func()

	heading    *gogl.Text
	loseDialog *gogl.Text
	logo2048   *gogl.TextBox
	score      *common.ScoreBox
	highScore  *common.ScoreBox
	menu       *gogl.Button
	newGame    *gogl.Button
	guide      *gogl.Text
	timer      *gogl.Text

	debugGrid  *gogl.Text
	debugTime  *gogl.Text
	debugScore *gogl.Text
}

// NewSingleplayerScreen constructs an uninitialised new singleplayer menu screen.
func NewSingleplayerScreen(win *gogl.Window) *SingleplayerScreen {
	return &SingleplayerScreen{win: win}
}

// Enter initialises the screen.
func (s *SingleplayerScreen) Enter(_ InitData) {
	// Arena and supporting data structures
	{
		s.arena = common.NewArena(gogl.Vec{X: 440, Y: 300})
		s.backend = backend.NewGame(nil)
		s.arenaInputCh = make(chan func(), 100)
	}

	// UI components
	{
		// Everything is sized relative to the tile size
		const unit = common.TileSizePx

		// Everything is positioned relative to the arena grid
		anchor := s.arena.Pos()

		s.heading = gogl.NewText(
			"", // to be set and drawn when player loses
			gogl.Vec{X: anchor.X + s.arena.Width()/2, Y: anchor.Y - 2.8*unit},
			common.FontPathBold,
		).SetSize(40).SetColour(common.GreyTextColour).SetAlignment(gogl.AlignTopCentre)

		s.loseDialog = gogl.NewText(
			"", // to be set and drawn when player loses
			gogl.Vec{X: anchor.X + s.arena.Width()/2, Y: anchor.Y - 1.9*unit},
			common.FontPathBold,
		).SetSize(20).SetColour(common.GreyTextColour).SetAlignment(gogl.AlignTopCentre)

		s.logo2048 = common.NewLogoBox(
			1.36*unit,
			gogl.Vec{X: anchor.X, Y: anchor.Y - 2.58*unit},
		)

		const wScore = 90
		s.score = common.NewScoreBox(
			wScore, wScore,
			gogl.Vec{X: anchor.X + s.arena.Width() - 2.74*unit, Y: anchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE")

		s.highScore = common.NewScoreBox(
			wScore, wScore,
			gogl.Vec{X: anchor.X + s.arena.Width() - wScore, Y: anchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("BEST")

		const buttonWidth = unit * 1.27
		s.menu = common.NewGameButton(
			buttonWidth, 0.4*unit,
			gogl.Vec{X: anchor.X + s.arena.Width() - buttonWidth, Y: anchor.Y - 1.21*unit},
			func() {
				SetScreen(Title, nil)
			},
		).SetLabelText("MENU")

		s.newGame = common.NewGameButton(
			buttonWidth, 0.4*unit,
			gogl.Vec{X: anchor.X + s.arena.Width() - 2.74*unit, Y: anchor.Y - 1.21*unit},
			func() {
				s.arenaInputCh <- func() {
					s.backend.Reset()
					s.arena.Reset()
				}
			},
		).SetLabelText("NEW")

		s.guide = gogl.NewText(
			"Join the numbers and get to the 2048 tile!",
			gogl.Vec{X: anchor.X, Y: anchor.Y - 0.60*unit},
			common.FontPathBold,
		).SetSize(16).SetColour(common.GreyTextColour)

		s.timer = common.NewGameText("",
			gogl.Vec{X: anchor.X + s.arena.Width(), Y: anchor.Y + s.arena.Height()*1.1},
		).SetSize(16).SetAlignment(gogl.AlignBottomRight)
	}

	// Debug UI
	s.debugGrid = gogl.NewText("grid", gogl.Vec{X: 930, Y: 600}, common.FontPathMedium).
		SetText(s.backend.Grid.Debug())
	s.debugTime = gogl.NewText("time", gogl.Vec{X: 1100, Y: 550}, common.FontPathMedium).
		SetText(s.backend.Timer.Time.String())
	s.debugScore = gogl.NewText("score", gogl.Vec{X: 950, Y: 550}, common.FontPathMedium).
		SetText(strconv.Itoa(s.backend.Score))

	// Set keybinds. User inputs are sent to the backend via a buffered channel
	// so the backend game cannot execute multiple moves before the frontend has
	// finished animating the first one
	{
		s.win.RegisterKeybind(gogl.KeyUp, gogl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirUp)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(gogl.KeyDown, gogl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirDown)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(gogl.KeyLeft, gogl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirLeft)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(gogl.KeyRight, gogl.KeyPress, func() {
			s.arenaInputCh <- func() {
				s.backend.ExecuteMove(grid.DirRight)
				s.debugGrid.SetText(s.backend.Grid.Debug())
			}
		})
		s.win.RegisterKeybind(gogl.KeyR, gogl.KeyRelease, func() {
			s.arenaInputCh <- func() {
				s.backend.Reset()
				s.arena.Reset()
			}
		})
		s.win.RegisterKeybind(gogl.KeyEscape, gogl.KeyRelease, func() {
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

	s.win.UnregisterKeybind(gogl.KeyUp, gogl.KeyPress)
	s.win.UnregisterKeybind(gogl.KeyDown, gogl.KeyPress)
	s.win.UnregisterKeybind(gogl.KeyLeft, gogl.KeyPress)
	s.win.UnregisterKeybind(gogl.KeyRight, gogl.KeyPress)
	s.win.UnregisterKeybind(gogl.KeyEscape, gogl.KeyRelease)

	s.arena.Destroy()
}

// Update updates and draws the singleplayer screen.
func (s *SingleplayerScreen) Update() {
	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
	default:
		// No user input; continue
	}

	// Deep copy so front-end has time to animate itself whilst allowing the
	// back-end to update
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

	// Draw debug grid
	if config.Debug {
		s.debugGrid.SetText(s.backend.Grid.Debug())
		s.debugTime.SetText(s.backend.Timer.Time.String())
		s.debugScore.SetText(
			fmt.Sprint(s.backend.Score, "|", s.backend.HighScore),
		)

		s.win.Draw(s.debugGrid)
		s.win.Draw(s.debugTime)
		s.win.Draw(s.debugScore)
	}
}

// Update updates and draws the singleplayer screen.
func (s *SingleplayerScreen) updateNormal(game backend.Game) {
	s.win.SetBackground(common.BackgroundColour)

	s.score.SetBody(strconv.Itoa(game.Score))
	s.menu.Update(s.win)
	s.highScore.SetBody(strconv.Itoa(game.HighScore))
	s.timer.SetText(game.Timer.Time.String())
	s.newGame.Update(s.win)

	s.arena.SetNormal()
	s.arena.Update(game)

	for _, d := range []gogl.Drawable{
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
		"You earned %d points in %v.", game.Score, game.Timer.Duration(),
	))

	s.menu.Update(s.win)
	s.newGame.Update(s.win)
	s.arena.Update(game)

	for _, d := range []gogl.Drawable{
		s.heading,
		s.loseDialog,
		s.menu,
		s.newGame,
		s.arena,
	} {
		s.win.Draw(d)
	}
}
