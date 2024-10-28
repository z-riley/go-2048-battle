package screen

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

const (
	serverPort = 8080
)

type MultiplayerHostScreen struct {
	win *turdgl.Window

	title            *turdgl.Text
	nameHeading      *turdgl.Text
	nameEntry        *common.EntryBox
	opponentStatus   *turdgl.Text
	start            *turdgl.Button
	back             *turdgl.Button
	buttonBackground *turdgl.CurvedRect

	server              *turdserve.Server
	opponentIsConnected bool
}

// NewMultiplayerHostScreen constructs an uninitialised multiplayer host screen.
func NewMultiplayerHostScreen(win *turdgl.Window) *MultiplayerHostScreen {
	return &MultiplayerHostScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerHostScreen) Enter(_ InitData) {
	s.title = turdgl.NewText("Host Game", turdgl.Vec{X: 600, Y: 200}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(100)

	s.nameHeading = turdgl.NewText(
		"Your name:",
		turdgl.Vec{X: 600, Y: 320},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(30)

	s.nameEntry = common.NewEntryBox(
		400, 60,
		turdgl.Vec{X: 600 - 400/2, Y: 330},
	)

	s.opponentStatus = turdgl.NewText(
		fmt.Sprintf("Waiting for opponent to join \"%s\"", getIPAddr()),
		turdgl.Vec{X: 600, Y: 530},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(24)

	// Adjustable settings for buttons
	const (
		TileSizePx        float64 = 120
		TileCornerRadius  float64 = 6
		TileBoundryFactor float64 = 0.15
	)

	// Background for buttons
	const w = TileSizePx * (2 + 3*TileBoundryFactor)
	s.buttonBackground = turdgl.NewCurvedRect(
		w, TileSizePx*(1+2*TileBoundryFactor), TileCornerRadius,
		turdgl.Vec{X: (float64(s.win.Width()) - w) / 2, Y: 560},
	)
	s.buttonBackground.SetStyle(turdgl.Style{Colour: common.ArenaBackgroundColour})

	s.start = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			if !s.opponentIsConnected {

				// Make the opponent status text briefly change colour
				s.opponentStatus.SetColour(common.Tile64Colour)
				go func() {
					timer := time.NewTimer(200 * time.Millisecond)
					select {
					case <-timer.C:
						s.opponentStatus.SetColour(common.GreyTextColour)
					}
				}()

				return
			}

			if err := s.startGame(); err != nil {
				fmt.Println("Failed to start game:", err)
			}
		},
	).SetLabelText("Start")

	s.back = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() { SetScreen(MultiplayerMenu, nil) },
	).SetLabelText("Back")

	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})

	// Set up server
	{
		const maxClients = 1
		s.server = turdserve.NewServer(maxClients).
			SetCallback(func(id int, b []byte) {
				if err := s.handleClientData(id, b); err != nil {
					fmt.Println("Failed to handle data from client:", err)
				}
			}).SetDisconnectCallback(func(_ int) { s.handleOpponentDisconnect() })

		// Start server to allow other players to connect
		errCh := make(chan error)
		go func() {
			for err := range errCh {
				if err != nil {
					// Exit the loop if the server dies
					return
				}
			}
		}()
		if err := s.server.Start("0.0.0.0", serverPort, errCh); err != nil {
			panic(err)
		}
	}
}

// Exit deinitialises the screen.
func (s *MultiplayerHostScreen) Exit() {
	s.server.Destroy()
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer host screen.
func (s *MultiplayerHostScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)
	s.win.Draw(s.buttonBackground)

	for _, l := range []*turdgl.Text{
		s.opponentStatus,
		s.nameHeading,
	} {
		s.win.Draw(l)
	}

	for _, b := range []*turdgl.Button{
		s.start,
		s.back,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}

	s.win.Draw(s.nameEntry)
	s.nameEntry.Update(s.win)
}

// handleClientData handles all data received from a client.
func (s *MultiplayerHostScreen) handleClientData(id int, b []byte) error {
	fmt.Println("Received from client", id, ":", string(b))

	var msg comms.Message
	if err := json.Unmarshal(b, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal bytes from client: %w", err)
	}

	switch msg.Type {
	case comms.TypePlayerData:
		var data comms.PlayerData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal player data: %w", err)
		}
		return s.handlePlayerData(id, data)

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// handlePlayerData handles incoming player data.
func (s *MultiplayerHostScreen) handlePlayerData(id int, data comms.PlayerData) error {
	// Make sure versions are compatible
	if data.Version != config.Version {
		return fmt.Errorf("incompatible versions (peer %s, local %s)", data.Version, config.Version)
	}

	s.opponentStatus.SetText(
		fmt.Sprintf("\"%s\" has joined the game. Press Start to begin", data.Username),
	)
	s.opponentIsConnected = true

	// Send host player data to the client
	username := s.nameEntry.Text()
	playerData, err := json.Marshal(
		comms.PlayerData{
			Version:  config.Version,
			Username: username,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal player data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypePlayerData,
			Content: playerData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal joining message: %w", err)
	}

	// Send data to client
	if err := s.server.WriteToClient(id, msg); err != nil {
		return fmt.Errorf("failed to send message to server: %w", err)
	}

	return nil
}

func (s *MultiplayerHostScreen) handleOpponentDisconnect() {
	s.opponentStatus.SetText(fmt.Sprintf("Waiting for opponent to join \"%s\"", getIPAddr()))
	s.opponentIsConnected = false
}

// getIPAddr returns the IP address of the host.
func getIPAddr() string {
	if comms.IsWSL() {
		// If using WSL, the host IP address must be used
		return "check WSL host"
	}
	conn, err := comms.LocalIP()
	if err != nil {
		panic(err)
	}
	return conn.String()
}

// serverKey is used for indentifying the server in InitData.
const serverKey = "server"

// startGame attempts to start a multiplayer game.
func (s *MultiplayerHostScreen) startGame() error {
	// Check opponent is connected
	if !s.opponentIsConnected {
		return errors.New("opponent is not connected")
	}

	// Inform other players that game is starting
	eventData, err := json.Marshal(comms.EventData{Event: comms.EventHostStartGame})
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	msg, err := json.Marshal(
		comms.Message{
			Type:    comms.TypeEventData,
			Content: eventData,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal event message: %w", err)
	}
	for _, id := range s.server.GetClientIDs() {
		fmt.Println("Writing to client ID", id)
		if err := s.server.WriteToClient(id, msg); err != nil {
			return fmt.Errorf("failed to send message to server: %w", err)
		}
	}

	// Pass server to next screen
	SetScreen(Multiplayer, InitData{serverKey: s.server})
	return nil
}
