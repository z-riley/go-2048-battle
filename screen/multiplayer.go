package screen

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/brunoga/deep"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerScreen struct {
	win              *turdgl.Window
	backgroundColour color.RGBA
	logo2048         *turdgl.TextBox

	// Player's grid
	newGame      *turdgl.Button
	menu         *turdgl.Button
	score        *common.ScoreBox
	guide        *turdgl.Text
	timer        *turdgl.Text
	backend      *backend.Game
	arena        *common.Arena
	arenaInputCh chan func()
	debugGrid    *turdgl.Text

	// Opponent's grid
	opponentScore     *common.ScoreBox
	opponentGuide     *turdgl.Text
	opponentArena     *common.Arena
	opponentBackend   *backend.Game
	opponentDebugGrid *turdgl.Text

	// EITHER server or client will exist
	server *turdserve.Server
	client *turdserve.Client
}

// NewMultiplayerScreen constructs a new singleplayer menu screen.
func NewMultiplayerScreen(win *turdgl.Window) *MultiplayerScreen {
	return &MultiplayerScreen{
		win:              win,
		backgroundColour: common.BackgroundColour,
	}
}

const (
	// usernameKey is used for indentifying the player's username in InitData.
	usernameKey = "username"
	// usernameKey is used for indentifying the opponent's username in InitData.
	opponentUsernameKey = "opponentUsername"
)

// Enter initialises the screen.
func (s *MultiplayerScreen) Enter(initData InitData) {
	// UI widgets
	{
		s.arena = common.NewArena(
			turdgl.Vec{X: float64(s.win.Width())/3 - 249, Y: 300},
		)
		s.opponentArena = common.NewArena(
			turdgl.Vec{X: float64(s.win.Width())*2/3 - 71, Y: 300},
		)

		// Everything is sized relative to the tile size
		const unit = common.TileSizePx

		anchor := s.arena.Pos()

		const logoSize = 1.36 * unit
		s.logo2048 = common.NewLogoBox(
			logoSize,
			turdgl.Vec{X: (float64(s.win.Width()) - logoSize) / 2, Y: anchor.Y - 2.58*unit},
		)

		// Player's grid
		{
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
					SetScreen(MultiplayerMenu, nil)
				},
			).SetLabelText("MENU")

			const wScore = 90
			s.score = common.NewScoreBox(
				90, 90,
				turdgl.Vec{X: anchor.X + s.arena.Width() - wScore, Y: anchor.Y - 2.58*unit},
				common.ArenaBackgroundColour,
			).SetHeading("SCORE")

			s.guide = common.NewGameText(
				"Your grid",
				turdgl.Vec{X: anchor.X + s.arena.Width(), Y: anchor.Y - 0.53*unit},
			).SetAlignment(turdgl.AlignTopRight)

			s.backend = backend.NewGame(&backend.Opts{
				SaveToDisk: false,
			})
			s.arenaInputCh = make(chan func(), 100)

			s.timer = common.NewGameText("",
				turdgl.Vec{X: 600, Y: anchor.Y - 0.53*unit},
			).SetAlignment(turdgl.AlignTopCentre)

			s.debugGrid = turdgl.NewText(
				s.backend.Grid.Debug(),
				turdgl.Vec{X: 100, Y: 50},
				common.FontPathMedium,
			)
		}

		// Opponent's grid
		{
			// Everything is positioned relative to the arena grid
			opponentAnchor := s.opponentArena.Pos()

			s.opponentScore = common.NewScoreBox(
				90, 90,
				turdgl.Vec{X: opponentAnchor.X, Y: opponentAnchor.Y - 2.58*unit},
				common.ArenaBackgroundColour,
			).SetHeading("SCORE")

			s.opponentGuide = common.NewGameText(
				fmt.Sprintf("%s's grid", initData[opponentUsernameKey].(string)),
				turdgl.Vec{X: opponentAnchor.X, Y: opponentAnchor.Y - 0.53*unit},
			)

			s.opponentBackend = backend.NewGame(&backend.Opts{
				SaveToDisk: false,
			})

			s.opponentDebugGrid = turdgl.NewText(
				s.opponentBackend.Grid.Debug(),
				turdgl.Vec{X: 850, Y: 50},
				common.FontPathMedium,
			)
		}
	}

	// Initialise server/client
	{
		if server, ok := initData[serverKey]; ok {
			// Host mode - initialise server
			s.server = server.(*turdserve.Server)
			s.server.SetCallback(func(id int, b []byte) {
				if err := s.handleOpponentData(b); err != nil {
					fmt.Println("Failed to handle opponent data as server", err)
				}
			}).SetDisconnectCallback(func(id int) {
				fmt.Println("Opponent has left the game")
			})
		} else if client, ok := initData[clientKey]; ok {
			// Guest mode - initialise client
			s.client = client.(*turdserve.Client)
			s.client.SetCallback(func(b []byte) {
				if err := s.handleOpponentData(b); err != nil {
					fmt.Println("Failed to handle opponent data as client", err)
				}
			})
		} else {
			panic("neither server or client was passed to MultiplayerScreen Init")
		}

		// Tell the opponent that the local server/client is ready to receive data
		if err := s.sendScreenLoadedEvent(); err != nil {
			if config.Debug {
				fmt.Println("Failed to send game update", err)
			}
		}
	}

	// Set keybinds. User inputs are sent to the backend via a buffered channel
	// so the backend game cannot execute multiple moves before the frontend has
	// finished animating the first one
	{
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
}

// Exit deinitialises the screen.
func (s *MultiplayerScreen) Exit() {
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
	// Check for win or lose
	switch s.backend.Grid.Outcome() {
	case grid.None:
		s.backgroundColour = common.BackgroundColour
	case grid.Win:
		s.backgroundColour = common.BackgroundColourWin
	case grid.Lose:
		s.backgroundColour = common.BackgroundColourLose
	}

	s.win.SetBackground(s.backgroundColour)

	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
		if err := s.sendGameData(); err != nil {
			fmt.Println("Failed to send game update:", err)
		}

	default:
		// No user input; continue
	}

	s.menu.Update(s.win)
	s.score.SetBody(fmt.Sprint(s.backend.Score.CurrentScore()))
	s.opponentScore.SetBody(fmt.Sprint(s.opponentBackend.Score.CurrentScore()))
	s.newGame.Update(s.win)
	s.timer.SetText(s.backend.Timer.Time.String())

	// Deep copy so front-end has time to animate itself whilst allowing the back
	// end to update
	s.arena.Update(deep.MustCopy(*s.backend))
	s.opponentArena.Update(deep.MustCopy(*s.opponentBackend))

	for _, d := range []turdgl.Drawable{
		s.logo2048,
		s.guide,
		s.opponentGuide,
		s.menu,
		s.score,
		s.opponentScore,
		s.newGame,
		s.timer,
		s.arena,
		s.opponentArena,
	} {
		s.win.Draw(d)
	}

	if config.Debug {
		s.debugGrid.SetText(s.backend.Grid.Debug())
		s.opponentDebugGrid.SetText(s.opponentBackend.Grid.Debug())
		s.win.Draw(s.debugGrid)
		s.win.Draw(s.opponentDebugGrid)
	}
}

// sendToOpponent sends bytes to the opponent.
func (s *MultiplayerScreen) sendToOpponent(b []byte) error {
	if s.server != nil {
		for _, id := range s.server.GetClientIDs() {
			if err := s.server.WriteToClient(id, b); err != nil {
				return fmt.Errorf("failed to send message to server: %w", err)
			}
		}
	} else {
		if err := s.client.Write(b); err != nil {
			return fmt.Errorf("failed to send message to client: %w", err)
		}
	}
	return nil
}

// handleOpponentData handles data from the opponent.
func (s *MultiplayerScreen) handleOpponentData(data []byte) error {
	var msg comms.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	switch msg.Type {
	case comms.TypeGameData:
		var data comms.GameData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal game data: %w", err)
		}
		return s.handleGameData(data)

	case comms.TypeEventData:
		var data comms.EventData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		return s.handleEventData(data)

	case comms.TypeRequest:
		var data comms.RequestData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal request data: %w", err)
		}
		return s.handleRequest(data)

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// sendGameData sends the local game state to the opponent.
func (s *MultiplayerScreen) sendGameData() error {
	gameData, err := json.Marshal(
		comms.GameData{
			Game: *s.backend,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal game data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypeGameData,
			Content: gameData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.sendToOpponent(msg)
}

// handleGameData handles incoming game data from the opponent.
func (s *MultiplayerScreen) handleGameData(data comms.GameData) error {
	s.opponentBackend = &data.Game
	return nil
}

// sendScreenLoadedEvent sends the screen loaded event to the opponent.
func (s *MultiplayerScreen) sendScreenLoadedEvent() error {
	eventData, err := json.Marshal(
		comms.EventData{
			Event: comms.EventScreenLoaded,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypeEventData,
			Content: eventData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.sendToOpponent(msg)
}

// handleEventData handles incoming game data from the opponent.
func (s *MultiplayerScreen) handleEventData(data comms.EventData) error {
	switch data.Event {
	case comms.EventScreenLoaded:
		// Send game data to opponent
		if err := s.sendGameData(); err != nil {
			return fmt.Errorf("failed to send game data: %w", err)
		}

		// Request for opponent to send their game data
		if err := s.requestOpponentGameData(); err != nil {
			return fmt.Errorf("failed to request opponent's game data: %w", err)
		}
	}

	return nil
}

// requestOpponentData sends a request for the opponent to send their game data.
func (s *MultiplayerScreen) requestOpponentGameData() error {
	request, err := json.Marshal(
		comms.RequestData{
			Request: comms.TypeGameData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypeRequest,
			Content: request,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := s.sendToOpponent(msg); err != nil {
		return fmt.Errorf("failed to send data to opponent: %w", err)
	}

	return nil
}

// handleRequest handles an incoming request for data.
func (s *MultiplayerScreen) handleRequest(data comms.RequestData) error {
	switch data.Request {
	case comms.TypeGameData:
		if err := s.sendGameData(); err != nil {
			return fmt.Errorf("failed to send game data: %w", err)
		}
	}
	return nil
}
