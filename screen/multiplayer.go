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

type MultiplayerScreen struct {
	win     *turdgl.Window
	backend *backend.Game

	title *turdgl.Text

	score     *common.GameUIBox
	highScore *common.GameUIBox
	newGame   *turdgl.Button

	arena         *common.Arena
	arenaInputCh  chan (func())
	opponentArena *common.Arena

	timer            *turdgl.Text
	backgroundColour color.RGBA
}

// NewMultiplayerScreen constructs a new singleplayer menu screen.
func NewMultiplayerScreen(win *turdgl.Window) *MultiplayerScreen {
	s := MultiplayerScreen{
		win:     win,
		backend: backend.NewGame(&backend.Opts{SaveToDisk: false}),

		title: turdgl.NewText("Head to Head", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
			SetColour(common.LightFontColour).
			SetAlignment(turdgl.AlignCentre).
			SetSize(40),

		score: common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: 100, Y: 70},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE"),
		highScore: common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: 220, Y: 70},
			common.ArenaBackgroundColour,
		).SetHeading("BEST"),

		arena:        common.NewArena(turdgl.Vec{X: 100, Y: 250}),
		arenaInputCh: make(chan func(), 100),

		opponentArena: common.NewArena(turdgl.Vec{X: 700, Y: 250}),

		timer:            common.NewGameText("", turdgl.Vec{X: 370, Y: 620}),
		backgroundColour: common.BackgroundColour,
	}

	s.newGame = common.NewGameButton(
		200, 50,
		turdgl.Vec{X: 100, Y: 180},
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
func (s *MultiplayerScreen) Init() {
	// Set keybinds. User inputs are sent to the backend via a buffered channel
	// so the backend game cannot execute multiple moves before the frontend has
	// finished animating the first one
	s.win.RegisterKeybind(turdgl.KeyUp, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirUp)
		}
	})
	s.win.RegisterKeybind(turdgl.KeyDown, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirDown)
		}
	})
	s.win.RegisterKeybind(turdgl.KeyLeft, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirLeft)
		}
	})
	s.win.RegisterKeybind(turdgl.KeyRight, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.ExecuteMove(grid.DirRight)
		}
	})
	s.win.RegisterKeybind(turdgl.KeyR, turdgl.KeyPress, func() {
		s.arenaInputCh <- func() {
			s.backend.Reset()
			s.arena.Reset()
		}
	})
	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyPress, func() {
		SetScreen(Title)
	})
}

// Deinit deinitialises the screen.
func (s *MultiplayerScreen) Deinit() {
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
func (s *MultiplayerScreen) Update() {
	s.win.SetBackground(s.backgroundColour)

	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
	default:
		// No user input; continue
	}

	// Serialise and deserialise grid to simulate receiving JSON from server
	b, err := s.backend.Serialise()
	if err != nil {
		panic(err)
	}
	var game backend.Game
	if err := json.Unmarshal(b, &game); err != nil {
		panic(err)
	}

	// Draw UI widgets
	s.win.Draw(s.title)

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

	// Draw opponent's arena
	s.opponentArena.Animate(game)
	s.opponentArena.Draw(s.win)

	// Check for win or lose
	switch game.Outcome {
	case grid.None:
		s.backgroundColour = common.BackgroundColour
	case grid.Win:
		s.backgroundColour = common.BackgroundColourWin
	case grid.Lose:
		s.backgroundColour = common.BackgroundColourLose
	}

}
