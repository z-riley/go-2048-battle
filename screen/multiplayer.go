package screen

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerScreen struct {
	win              *turdgl.Window
	backgroundColour color.RGBA

	// Player's grid
	newGame      *turdgl.Button
	score        *common.GameUIBox
	menu         *turdgl.Button
	guide        *turdgl.Text
	backend      *backend.Game
	arena        *common.Arena
	arenaInputCh chan (func())

	// Opponent's grid
	opponentScore   *common.GameUIBox
	opponentGuide   *turdgl.Text
	opponentArena   *common.Arena
	opponentBackend *backend.Game

	// Shared widgets
	timer *turdgl.Text
	title *turdgl.Text

	// TODO: find a neater way of doing client/server polymorphism
	// but don't forget to account for 1 server -> multiple clients
	server *turdserve.Server
	client *turdserve.Client
}

// NewMultiplayerScreen constructs a new singleplayer menu screen.
func NewMultiplayerScreen(win *turdgl.Window) *MultiplayerScreen {
	s := MultiplayerScreen{
		win: win,
		title: turdgl.NewText("Head to Head", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
			SetColour(common.ArenaBackgroundColour).
			SetAlignment(turdgl.AlignCentre).
			SetSize(40),

		backend:      backend.NewGame(&backend.Opts{SaveToDisk: false}),
		arena:        common.NewArena(turdgl.Vec{X: 100, Y: 300}),
		arenaInputCh: make(chan func(), 100),

		opponentBackend: backend.NewGame(&backend.Opts{SaveToDisk: false}),
		opponentArena:   common.NewArena(turdgl.Vec{X: 700, Y: 300}),

		backgroundColour: common.BackgroundColour,
	}

	// Everything is sized relative to the tile size
	const unit = common.TileSizePx

	// Player's grid:
	{
		// Everything is positioned relative to the arena grid
		anchor := s.arena.Pos()

		const widgetWidth = unit * 1.27
		s.newGame = common.NewGameButton(
			widgetWidth, 0.4*unit,
			turdgl.Vec{X: anchor.X + s.arena.Width() - 2.74*unit, Y: anchor.Y - 1.21*unit},
			func() {
				s.arenaInputCh <- func() {
					s.backend.Reset()
					s.arena.Reset()
				}
			},
		).SetLabelText("NEW")

		s.menu = common.NewGameButton(
			widgetWidth, 0.4*unit,
			turdgl.Vec{X: anchor.X + s.arena.Width() - widgetWidth, Y: anchor.Y - 1.21*unit},
			func() {
				SetScreen(Multiplayer, nil)
			},
		).SetLabelText("MENU")

		const wScore = 90
		s.score = common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: anchor.X + s.arena.Width() - wScore, Y: anchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE")

		s.guide = turdgl.NewText(
			"Your grid", // TODO: Make this
			turdgl.Vec{X: anchor.X, Y: anchor.Y - 0.28*unit},
			common.FontPathBold,
		).SetSize(16).SetColour(common.GreyTextColour)

		s.timer = common.NewGameText("",
			turdgl.Vec{X: 600, Y: anchor.Y},
		).SetAlignment(turdgl.AlignBottomRight)
	}

	// Opponent's grid:
	{
		// Everything is positioned relative to the arena grid
		opponentAnchor := s.opponentArena.Pos()

		s.opponentScore = common.NewGameTextBox(
			90, 90,
			turdgl.Vec{X: opponentAnchor.X, Y: opponentAnchor.Y - 2.58*unit},
			common.ArenaBackgroundColour,
		).SetHeading("SCORE")

		s.opponentGuide = turdgl.NewText(
			"Opponent's grid",
			turdgl.Vec{X: opponentAnchor.X, Y: opponentAnchor.Y - 0.28*unit},
			common.FontPathBold,
		).SetSize(16).SetColour(common.GreyTextColour)
	}

	return &s
}

// Init initialises the screen.
func (s *MultiplayerScreen) Init(initData InitData) {
	if server, ok := initData[serverKey]; ok {
		// Host mode
		s.server = server.(*turdserve.Server)
		s.server.SetCallback(func(id int, b []byte) {
			if err := s.handleOpponentData(b); err != nil {
				fmt.Println("Failed to handle opponent data as server", err)
			}
		}).SetDisconnectCallback(func(id int) {
			fmt.Println("Opponent has left the game")
		})
	} else if client, ok := initData[clientKey]; ok {
		// Guest mode
		s.client = client.(*turdserve.Client)
		s.client.SetCallback(func(b []byte) {
			if err := s.handleOpponentData(b); err != nil {
				fmt.Println("Failed to handle opponent data as client", err)
			}
		})
	} else {
		panic("neither server or client was passed to MultiplayerScreen Init")
	}

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
		SetScreen(Title, nil)
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

	if s.server != nil {
		s.server.Destroy()
	} else if s.client != nil {
		s.client.Destroy()
	}

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
		if err := s.sendGameUpdate(); err != nil {
			fmt.Println("Failed to send game update:", err)
		}

	default:
		// No user input; continue
	}

	// Draw UI widgets
	s.win.Draw(s.title)
	s.win.Draw(s.guide)
	s.win.Draw(s.opponentGuide)

	s.win.Draw(s.menu)
	s.menu.Update(s.win)

	s.score.Draw(s.win)
	s.score.SetBody(fmt.Sprint(s.backend.Score.CurrentScore()))

	s.opponentScore.Draw(s.win)
	s.opponentScore.SetBody(fmt.Sprint(s.opponentBackend.Score.CurrentScore()))

	s.win.Draw(s.newGame)
	s.newGame.Update(s.win)

	s.timer.SetText(s.backend.Timer.Time.String())
	s.win.Draw(s.timer)

	// Draw arena of tiles
	s.arena.Animate(*s.backend)
	s.arena.Draw(s.win)

	// Draw opponent's arena
	s.opponentArena.Animate(*s.opponentBackend)
	s.opponentArena.Draw(s.win)

	// Check for win or lose
	switch s.backend.Outcome {
	case grid.None:
		s.backgroundColour = common.BackgroundColour
	case grid.Win:
		s.backgroundColour = common.BackgroundColourWin
	case grid.Lose:
		s.backgroundColour = common.BackgroundColourLose
	}

}

// sendGameUpdate sends the local game state to the opponent.
func (s *MultiplayerScreen) sendGameUpdate() error {
	// Send game data update to opponent
	gameData, err := json.Marshal(comms.GameData{
		Game: *s.backend,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal player data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypeGameData,
			Content: gameData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal joining message: %w", err)
	}

	// Send data to opponent
	if s.server != nil {
		for _, id := range s.server.GetClientIDs() {
			if err := s.server.WriteToClient(id, msg); err != nil {
				return fmt.Errorf("failed to send message to server: %w", err)
			}
		}
	} else {
		if err := s.client.Write(msg); err != nil {
			return fmt.Errorf("failed to send message to server: %w", err)
		}
	}

	return nil
}

// handleOpponentData handles data from the opponent.
func (s *MultiplayerScreen) handleOpponentData(b []byte) error {
	// Note: in matches with 3 or more players, the server would need to forward
	// the incoming client data to the other clients. However, with 2 players, this
	// isn't necessary

	var msg comms.Message
	if err := json.Unmarshal(b, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal bytes from client: %w", err)
	}

	switch msg.Type {
	case comms.TypeGameData:
		var data comms.GameData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal game data: %w", err)
		}
		return s.handleGameData(data)

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// handlePlayerData handles incoming game data from the opponent.
func (s *MultiplayerScreen) handleGameData(data comms.GameData) error {
	fmt.Println("Received new data from opponent: ", data)

	s.opponentBackend = &data.Game

	return nil
}
