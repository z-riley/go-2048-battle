package screen

import (
	"encoding/json"
	"fmt"

	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/turdgl"
)

type SingleplayerScreen struct {
	win     *turdgl.Window
	backend *backend.Game

	arena        *common.Arena
	arenaInputCh chan (func())

	debugGridText  *turdgl.Text
	debugTimeText  *turdgl.Text
	debugScoreText *turdgl.Text
}

// NewSingleplayerScreen constructs a new singleplayer menu screen.
func NewSingleplayerScreen(win *turdgl.Window) *SingleplayerScreen {
	return &SingleplayerScreen{
		win:     win,
		backend: backend.NewGame(),

		arena:        common.NewArena(turdgl.Vec{X: 700, Y: 80}),
		arenaInputCh: make(chan func(), 100),

		debugGridText:  turdgl.NewText("grid", turdgl.Vec{X: 100, Y: 100}, common.FontPathMedium),
		debugTimeText:  turdgl.NewText("time", turdgl.Vec{X: 500, Y: 100}, common.FontPathMedium),
		debugScoreText: turdgl.NewText("score", turdgl.Vec{X: 400, Y: 100}, common.FontPathMedium),
	}
}

// Init initialises the screen.
func (s *SingleplayerScreen) Init() {
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
		SetScreen(Title)
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
	s.win.SetBackground(common.BackgroundColour)

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
	b, err := s.backend.Serialise()
	if err != nil {
		panic(err)
	}
	var game backend.Game
	if err := json.Unmarshal(b, &game); err != nil {
		panic(err)
	}

	// Draw arena of tiles
	s.arena.Animate(game)
	s.arena.Draw(s.win)

	// Draw temporary debug grid
	s.win.Draw(s.debugGridText)
	s.win.Draw(s.debugTimeText)
	s.win.Draw(s.debugScoreText)
}
