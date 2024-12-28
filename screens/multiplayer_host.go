package screens

import (
	"errors"
	"fmt"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/common/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/log"
	"github.com/z-riley/gogl"
	"github.com/z-riley/servesyouright"
)

const (
	serverPort = 8080
)

type MultiplayerHostScreen struct {
	win *gogl.Window

	title            *gogl.Text
	tooltip          *gogl.TextBox
	nameHeading      *gogl.Text
	nameEntry        *common.EntryBox
	opponentName     string
	opponentStatus   *gogl.Text
	start            *gogl.Button
	back             *gogl.Button
	buttonBackground *gogl.CurvedRect

	server            *servesyouright.Server
	opponentIsInLobby bool
}

// NewMultiplayerHostScreen constructs an uninitialised multiplayer host screen.
func NewMultiplayerHostScreen(win *gogl.Window) *MultiplayerHostScreen {
	return &MultiplayerHostScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerHostScreen) Enter(_ InitData) {
	s.title = gogl.NewText("Host Game", gogl.Vec{X: config.WinWidth / 2, Y: 150}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(gogl.AlignCentre).
		SetSize(100)

	s.tooltip = common.NewTooltip()

	s.nameHeading = gogl.NewText(
		"Your name:",
		gogl.Vec{X: config.WinWidth / 2, Y: 300},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(gogl.AlignCentre).
		SetSize(30)

	s.nameEntry = common.NewEntryBox(
		440, 60,
		gogl.Vec{X: (config.WinWidth - 440) / 2, Y: s.nameHeading.Pos().Y + 30},
		namesgenerator.GetRandomName(0),
	).
		SetModifiedCB(func() {
			// Update guest with new username
			if err := s.sendPlayerData(); err != nil {
				log.Println("Failed to send username update to guests:", err)
			}
		})

	s.opponentStatus = gogl.NewText(
		fmt.Sprintf("Waiting for opponent to join \"%s\"", getIPAddr()),
		gogl.Vec{X: config.WinWidth / 2, Y: 510},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(gogl.AlignCentre).
		SetSize(24)

	// Adjustable settings for buttons
	const (
		TileSizePx        float64 = 120
		TileCornerRadius  float64 = 6
		TileBoundryFactor float64 = 0.15
	)

	// Background for buttons
	const w = TileSizePx * (2 + 3*TileBoundryFactor)
	s.buttonBackground = gogl.NewCurvedRect(
		w, TileSizePx*(1+2*TileBoundryFactor), TileCornerRadius,
		gogl.Vec{X: (config.WinWidth - w) / 2, Y: 560},
	)
	s.buttonBackground.SetStyle(gogl.Style{Colour: common.ArenaBackgroundColour})

	s.start = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			if !s.opponentIsInLobby {
				// Make the opponent status text briefly change colour
				s.opponentStatus.SetColour(common.Tile64Colour)
				go func() {
					timer := time.NewTimer(200 * time.Millisecond)
					<-timer.C
					s.opponentStatus.SetColour(common.GreyTextColour)
				}()

				return
			}

			if err := s.startGame(); err != nil {
				log.Println("Failed to start game:", err)
			}
		},
	).SetLabelText("Start")

	s.back = common.NewMenuButton(
		TileSizePx, TileSizePx,
		gogl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			s.server.Destroy()
			SetScreen(MultiplayerMenu, nil)
		},
	).SetLabelText("Back")

	s.win.RegisterKeybind(gogl.KeyEscape, gogl.KeyRelease, func() {
		s.server.Destroy()
		SetScreen(MultiplayerMenu, nil)
	})

	// Set up server
	const maxClients = 1
	s.server = servesyouright.NewServer(maxClients).
		SetCallback(func(_ int, b []byte) {
			if err := s.handleClientData(b); err != nil {
				log.Println("Host screen failed to handle data from client:", err)
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

// Exit deinitialises the screen.
func (s *MultiplayerHostScreen) Exit() {
	s.win.UnregisterKeybind(gogl.KeyEscape, gogl.KeyRelease)
	s.opponentIsInLobby = false
}

// Update updates and draws multiplayer host screen.
func (s *MultiplayerHostScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)
	s.win.Draw(s.buttonBackground)

	for _, l := range []*gogl.Text{
		s.opponentStatus,
		s.nameHeading,
	} {
		s.win.Draw(l)
	}

	for _, b := range []*gogl.Button{
		s.start,
		s.back,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}

	s.win.Draw(s.nameEntry)
	s.nameEntry.Update(s.win)

	mouseLoc := s.win.MouseLocation()
	if s.nameEntry.TextBox.Shape.IsWithin(mouseLoc) && !s.nameEntry.TextBox.IsEditing() {
		s.tooltip.SetPos(gogl.Vec{X: mouseLoc.X, Y: mouseLoc.Y - s.tooltip.Shape.Height()})
		s.win.Draw(s.tooltip)
	}
}

// handleClientData handles all data received from a client.
func (s *MultiplayerHostScreen) handleClientData(data []byte) error {
	msg, err := comms.ParseMessage(data)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	switch msg.Type {
	case comms.TypePlayerData:
		data, err := comms.ParsePlayerData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse player data: %w", err)
		}
		return s.handlePlayerData(data)

	default:
		// Ignore other message type - don't return error
		return nil
	}
}

// sendPlayerData sends the player data to all connected guests.
func (s *MultiplayerHostScreen) sendPlayerData() error {
	msg, err := comms.PlayerData{
		Version:  config.Version,
		Username: s.nameEntry.Text(),
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise player data: %w", err)
	}

	// Send data to client
	for _, id := range s.server.GetClientIDs() {
		if err := s.server.WriteToClient(id, msg); err != nil {
			return fmt.Errorf("failed to send message to client: %w", err)
		}
	}

	return nil
}

// handlePlayerData handles incoming player data.
func (s *MultiplayerHostScreen) handlePlayerData(data comms.PlayerData) error {
	// Make sure versions are compatible
	if data.Version != config.Version {
		return fmt.Errorf("incompatible versions (peer %s, local %s)", data.Version, config.Version)
	}

	s.opponentName = data.Username
	s.opponentStatus.SetText(
		fmt.Sprintf("\"%s\" has joined the game. Press Start to begin", s.opponentName),
	)
	s.opponentIsInLobby = true

	// Send host player data to client
	if err := s.sendPlayerData(); err != nil {
		return fmt.Errorf("failed to send player data to client: %w", err)
	}

	return nil
}

// handleOpponentDisconnect handles the opponent disconnecting from the server.
func (s *MultiplayerHostScreen) handleOpponentDisconnect() {
	s.opponentStatus.SetText(fmt.Sprintf("Waiting for opponent to join \"%s\"", getIPAddr()))
	s.opponentIsInLobby = false
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
	if !s.opponentIsInLobby {
		return errors.New("opponent is not connected")
	}

	// Inform other players that game is starting
	msg, err := comms.EventData{
		Event: comms.EventHostStartGame,
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise event data: %w", err)
	}
	for _, id := range s.server.GetClientIDs() {
		if err := s.server.WriteToClient(id, msg); err != nil {
			return fmt.Errorf("failed to send message to server: %w", err)
		}
	}

	// Pass server to next screen
	SetScreen(Multiplayer, InitData{
		serverKey:           s.server,
		usernameKey:         s.nameEntry.Text(),
		opponentUsernameKey: s.opponentName,
	})
	return nil
}
