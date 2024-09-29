package screen

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type SingleplayerScreen struct {
	win     *turdgl.Window
	backend *backend.Game

	score     *common.GameUIBox
	highScore *common.GameUIBox
	newGame   *turdgl.Button

	arena        *common.Arena
	arenaInputCh chan (func())

	timer            *turdgl.Text
	backgroundColour color.RGBA

	debugGridText  *turdgl.Text
	debugTimeText  *turdgl.Text
	debugScoreText *turdgl.Text
}

// NewSingleplayerScreen constructs a new singleplayer menu screen.
func NewSingleplayerScreen(win *turdgl.Window) *SingleplayerScreen {
	s := SingleplayerScreen{
		win:     win,
		backend: backend.NewGame(&backend.Opts{SaveToDisk: true}),

		score: common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: 300, Y: 70},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE"),
		highScore: common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: 500, Y: 70},
			common.ArenaBackgroundColour,
		).SetHeading("BEST"),

		arena:        common.NewArena(turdgl.Vec{X: 250, Y: 250}),
		arenaInputCh: make(chan func(), 100),

		timer:            common.NewGameText("", turdgl.Vec{X: 520, Y: 620}),
		backgroundColour: common.BackgroundColour,

		debugGridText:  turdgl.NewText("grid", turdgl.Vec{X: 930, Y: 600}, common.FontPathMedium),
		debugTimeText:  turdgl.NewText("time", turdgl.Vec{X: 1100, Y: 550}, common.FontPathMedium),
		debugScoreText: turdgl.NewText("score", turdgl.Vec{X: 950, Y: 550}, common.FontPathMedium),
	}

	s.newGame = common.NewGameButton(
		200, 50,
		turdgl.Vec{X: 300, Y: 180},
		func() {
			s.arenaInputCh <- func() {
				s.backend.Reset()
				s.arena.Reset()
			}
		},
	).SetLabelText("NEW")

	return &s
}

// Init initialises the screen.
func (s *SingleplayerScreen) Init(_ InitData) {
	// Load debug UI
	s.debugGridText.SetText(s.backend.Grid.Debug())
	s.debugTimeText.SetText(s.backend.Timer.Time.String())
	s.debugScoreText.SetText(fmt.Sprint(s.backend.Score.Current))

	// Set keybinds. User inputs are sent to the backend via a buffered channel
	// so the backend game cannot execute multiple moves before the frontend has
	// finished animating the first one
	s.win.RegisterKeybind(turdgl.KeyUp, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirUp)
			s.debugGridText.SetText(s.backend.Grid.Debug())
		}
	})
	s.win.RegisterKeybind(turdgl.KeyDown, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirDown)
			s.debugGridText.SetText(s.backend.Grid.Debug())
		}
	})
	s.win.RegisterKeybind(turdgl.KeyLeft, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirLeft)
			s.debugGridText.SetText(s.backend.Grid.Debug())
		}
	})
	s.win.RegisterKeybind(turdgl.KeyRight, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirRight)
			s.debugGridText.SetText(s.backend.Grid.Debug())
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

// Deinit deinitialises the screen.
func (s *SingleplayerScreen) Deinit() {
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
	s.win.SetBackground(s.backgroundColour)

	// Temporary debug text
	s.debugGridText.SetText(s.backend.Grid.Debug())
	s.debugTimeText.SetText(s.backend.Timer.Time.String())
	s.debugScoreText.SetText(
		fmt.Sprint(s.backend.Score.Current, "|", s.backend.Score.High),
	)

	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
	default:
		// No user input; continue
	}

	// Serialise and deserialise grid to simulate receiving JSON from server
	// TODO: remove this once multiplayer is working
	b, err := s.backend.Serialise()
	if err != nil {
		panic(err)
	}
	var game backend.Game
	if err := json.Unmarshal(b, &game); err != nil {
		panic(err)
	}

	// Draw UI widgets
	s.score.Draw(s.win)
	s.score.SetBody(fmt.Sprint(game.Score.Current))

	s.highScore.Draw(s.win)
	s.highScore.SetBody(fmt.Sprint(game.Score.High))

	s.win.Draw(s.newGame)
	s.newGame.Update(s.win)

	s.timer.SetText(game.Timer.Time.String())
	s.win.Draw(s.timer)

	// Draw arena of tiles
	s.arena.Animate(game)
	s.arena.Draw(s.win)

	// Check for win or lose
	switch game.Outcome {
	case grid.None:
		s.backgroundColour = common.BackgroundColour
	case grid.Win:
		s.backgroundColour = common.BackgroundColourWin
	case grid.Lose:
		s.backgroundColour = common.BackgroundColourLose
	}

	// Draw temporary debug grid
	s.win.Draw(s.debugGridText)
	s.win.Draw(s.debugTimeText)
	s.win.Draw(s.debugScoreText)
}
