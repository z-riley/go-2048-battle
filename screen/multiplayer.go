package screen

import (
	"fmt"
	"image/color"

	"github.com/brunoga/deep"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/log"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerScreen struct {
	win              *turdgl.Window
	backgroundColour color.RGBA
	logo2048         *turdgl.TextBox

	newGame       *turdgl.Button
	menu          *turdgl.Button
	score         *common.ScoreBox
	guide         *turdgl.Text
	timer         *turdgl.Text
	backend       *backend.Game
	arena         *common.Arena
	arenaInputCh  chan func()
	endGameDialog *turdgl.Text
	debugGrid     *turdgl.Text

	opponentScore     *common.ScoreBox
	opponentName      string
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
			turdgl.Vec{X: config.WinWidth/3 - 249, Y: 300},
		)
		s.opponentArena = common.NewArena(
			turdgl.Vec{X: config.WinWidth*2/3 - 71, Y: 300},
		)

		// Everything is sized relative to the tile size and arena position
		const unit = common.TileSizePx
		anchor := s.arena.Pos()

		const logoSize = 1.36 * unit
		s.logo2048 = common.NewLogoBox(
			logoSize,
			turdgl.Vec{X: (config.WinWidth - logoSize) / 2, Y: anchor.Y - 2.58*unit},
		)

		s.endGameDialog = common.NewGameText(
			"Press MENU to\nplay again",
			turdgl.Vec{X: config.WinWidth / 2, Y: anchor.Y - 2.5*unit},
		).SetAlignment(turdgl.AlignTopCentre).SetSize(25)

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
				turdgl.Vec{X: config.WinWidth / 2, Y: anchor.Y - 0.53*unit},
			).SetAlignment(turdgl.AlignTopCentre)

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

			s.opponentName = initData[opponentUsernameKey].(string)
			s.opponentGuide = common.NewGameText(
				fmt.Sprintf("%s's grid", s.opponentName),
				turdgl.Vec{X: opponentAnchor.X, Y: opponentAnchor.Y - 0.53*unit},
			)

			s.opponentBackend = backend.NewGame(&backend.Opts{
				SaveToDisk: false,
			})
		}

		// Debug widgets
		if config.Debug {
			s.debugGrid = turdgl.NewText(
				s.backend.Grid.Debug(),
				turdgl.Vec{X: 100, Y: 50},
				common.FontPathMedium,
			)

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
					log.Println("Failed to handle opponent data as server", err)
				}
			}).SetDisconnectCallback(func(id int) {
				log.Println("Opponent has left the game")
			})
		} else if client, ok := initData[clientKey]; ok {
			// Guest mode - initialise client
			s.client = client.(*turdserve.Client)
			s.client.SetCallback(func(b []byte) {
				if err := s.handleOpponentData(b); err != nil {
					log.Println("Failed to handle opponent data as client", err)
				}
			})
		} else {
			panic("neither server or client was passed to MultiplayerScreen Init")
		}

		// Tell the opponent that the local server/client is ready to receive data
		if err := s.sendScreenLoadedEvent(); err != nil {
			log.Println("Failed to send game update", err)
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

// Update updates and draws the multiplayer screen.
func (s *MultiplayerScreen) Update() {
	s.win.SetBackground(s.backgroundColour)

	// Handle user inputs from user. Only 1 input must be sent per update cycle,
	// because the frontend can only animate one move at a time.
	select {
	case inputFunc := <-s.arenaInputCh:
		inputFunc()
		if err := s.sendGameData(); err != nil {
			log.Println("Failed to send game update:", err)
		}
	default:
		// No user input; continue
	}

	// Deep copy so front-end has time to animate itself whilst allowing the back
	// end to update
	s.arena.Update(deep.MustCopy(*s.backend))
	s.opponentArena.Update(deep.MustCopy(*s.opponentBackend))

	// Check for win or lose
	isLoss := s.backend.Grid.Outcome() == grid.Lose || s.opponentBackend.Grid.Outcome() == grid.Win
	isWin := s.backend.Grid.Outcome() == grid.Win || s.opponentBackend.Grid.Outcome() == grid.Lose
	switch {
	case isLoss:
		s.updateLose()
	case isWin:
		s.updateWin()
	default:
		s.updateNormal()
	}

	if config.Debug {
		s.debugGrid.SetText(s.backend.Grid.Debug())
		s.opponentDebugGrid.SetText(s.opponentBackend.Grid.Debug())
		s.win.Draw(s.debugGrid)
		s.win.Draw(s.opponentDebugGrid)
	}
}

// Update updates and draws the singleplayer screen in a normal state.
func (s *MultiplayerScreen) updateNormal() {
	s.newGame.Update(s.win)
	s.menu.Update(s.win)
	s.score.SetBody(fmt.Sprint(s.backend.Score.CurrentScore()))
	s.timer.SetText(s.backend.Timer.Time.String())
	s.opponentScore.SetBody(fmt.Sprint(s.opponentBackend.Score.CurrentScore()))

	for _, d := range []turdgl.Drawable{
		s.logo2048,
		s.newGame,
		s.menu,
		s.score,
		s.guide,
		s.timer,
		s.arena,
		s.opponentScore,
		s.opponentGuide,
		s.opponentArena,
	} {
		s.win.Draw(d)
	}
}

// updateWin updates and draws the singleplayer screen in a winning state.
func (s *MultiplayerScreen) updateWin() {
	s.arena.SetWin()
	s.opponentArena.SetLose()

	s.guide.SetText("You win!")
	s.opponentGuide.SetText(fmt.Sprintf("%s loses!", s.opponentName))

	s.menu.Update(s.win)
	s.score.SetBody(fmt.Sprint(s.backend.Score.CurrentScore()))
	s.timer.SetText(s.backend.Timer.Time.String())
	s.backend.Timer.Pause()
	s.opponentScore.SetBody(fmt.Sprint(s.opponentBackend.Score.CurrentScore()))

	for _, d := range []turdgl.Drawable{
		s.menu,
		s.score,
		s.guide,
		s.timer,
		s.arena,
		s.endGameDialog,
		s.opponentScore,
		s.opponentGuide,
		s.opponentArena,
	} {
		s.win.Draw(d)
	}
}

// updateLose updates and draws the singleplayer screen in a losing state.
func (s *MultiplayerScreen) updateLose() {
	s.arena.SetLose()
	s.opponentArena.SetWin()

	s.guide.SetText("You lose!")
	s.opponentGuide.SetText(fmt.Sprintf("%s wins!", s.opponentName))

	s.menu.Update(s.win)
	s.score.SetBody(fmt.Sprint(s.backend.Score.CurrentScore()))
	s.timer.SetText(s.backend.Timer.Time.String())
	s.backend.Timer.Pause()
	s.opponentScore.SetBody(fmt.Sprint(s.opponentBackend.Score.CurrentScore()))

	for _, d := range []turdgl.Drawable{
		s.menu,
		s.score,
		s.guide,
		s.timer,
		s.arena,
		s.endGameDialog,
		s.opponentScore,
		s.opponentGuide,
		s.opponentArena,
	} {
		s.win.Draw(d)
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
	msg, err := comms.ParseMessage(data)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	switch msg.Type {
	case comms.TypeGameData:
		gameData, err := comms.ParseGameData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse game data: %w", err)
		}
		return s.handleGameData(gameData)

	case comms.TypeEventData:
		eventData, err := comms.ParseEventData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse event data: %w", err)
		}
		return s.handleEventData(eventData)

	case comms.TypeRequestData:
		requestData, err := comms.ParseRequestData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse request data: %w", err)
		}
		return s.handleRequest(requestData)

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// sendGameData sends the local game state to the opponent.
func (s *MultiplayerScreen) sendGameData() error {
	msg, err := comms.GameData{
		Game: *s.backend,
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise game data: %w", err)
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
	msg, err := comms.EventData{
		Event: comms.EventScreenLoaded,
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise event data: %w", err)
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
	msg, err := comms.RequestData{
		Request: comms.TypeGameData,
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise request data: %w", err)
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
