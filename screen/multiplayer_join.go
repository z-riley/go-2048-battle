package screen

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerJoinScreen struct {
	win *turdgl.Window

	title            *turdgl.Text
	nameHeading      *turdgl.Text
	nameEntry        *common.EntryBox
	ipHeading        *turdgl.Text
	ipEntry          *common.EntryBox
	opponentStatus   *turdgl.Text
	join             *turdgl.Button
	back             *turdgl.Button
	buttonBackground *turdgl.CurvedRect

	client      *turdserve.Client
	hostIsReady chan bool
}

// NewTitle Screen constructs an uninitialised multiplayer join screen.
func NewMultiplayerJoinScreen(win *turdgl.Window) *MultiplayerJoinScreen {
	return &MultiplayerJoinScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerJoinScreen) Enter(_ InitData) {
	s.title = turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 200}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(100)

	s.nameHeading = turdgl.NewText(
		"Your name:",
		turdgl.Vec{X: 600, Y: 260},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(30)

	s.nameEntry = common.NewEntryBox(
		400, 60,
		turdgl.Vec{X: 600 - 400/2, Y: s.nameHeading.Pos().Y + 10},
	).SetText(namesgenerator.GetRandomName(0))

	s.opponentStatus = turdgl.NewText(
		"Click join......",
		turdgl.Vec{X: 600, Y: 530},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(24)

	s.ipHeading = turdgl.NewText(
		"Host IP:",
		turdgl.Vec{X: 600, Y: 420},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(30)

	s.ipEntry = common.NewEntryBox(
		400, 60,
		turdgl.Vec{X: 600 - 400/2, Y: s.ipHeading.Pos().Y + 10},
	)
	s.ipEntry.TextBox.SetText("127.0.0.1") // temporary for local testing

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

	s.hostIsReady = make(chan bool)
	s.join = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			if err := s.joinGame(); err != nil {
				fmt.Println("Failed to join game:", err)
				return
			}

			// Disable the button so user can't connect again
			s.join.SetCallback(
				turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
				func() {},
			)

			go func() {
				if <-s.hostIsReady {
					SetScreen(Multiplayer, InitData{clientKey: s.client})
					return
				}
			}()
		},
	).SetLabelText("Join")

	s.back = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*(1+2*TileBoundryFactor),
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		func() {
			s.join.SetLabelText("Join")
			s.client.Destroy()
			SetScreen(MultiplayerMenu, nil)
		}).SetLabelText("Back")

	s.client = turdserve.NewClient()

	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})
}

// Exit deinitialises the screen.
func (s *MultiplayerJoinScreen) Exit() {
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer join screen.
func (s *MultiplayerJoinScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)
	s.win.Draw(s.ipHeading)
	s.win.Draw(s.nameHeading)
	s.win.Draw(s.opponentStatus)
	s.win.Draw(s.buttonBackground)

	for _, b := range []*turdgl.Button{
		s.back,
		s.join,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}

	for _, e := range []*common.EntryBox{
		s.nameEntry,
		s.ipEntry,
	} {
		e.Update(s.win)
		s.win.Draw(e)
	}
}

// clientKey is used for indentifying the server in InitData.
const clientKey = "client"

// joinGame attempts to join a multiplayer game.
func (s *MultiplayerJoinScreen) joinGame() error {

	// Connect using the user-specified IP address
	ip := s.ipEntry.Text()
	errCh := make(chan error)
	go func() {
		for err := range errCh {
			if err != nil {
				fmt.Println(fmt.Errorf("client error: %w", err).Error())
				s.client.Destroy()
			}
		}
	}()
	if err := s.client.Connect(context.Background(), ip, serverPort, errCh); err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	s.client.SetCallback(func(b []byte) {
		if err := s.handleServerData(b); err != nil {
			fmt.Println("Failed to handle data from server:", err)
		}
	})

	// Construct message containing player data
	username := s.nameEntry.Text()
	playerData, err := json.Marshal(
		comms.PlayerData{
			Version:  config.Version,
			Username: username,
		})
	if err != nil {
		return fmt.Errorf("failed to marshal player data: %w", err)
	}
	msg, err := json.Marshal(comms.Message{
		Type:    comms.TypePlayerData,
		Content: playerData,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal joining message: %w", err)
	}

	// Send data to host
	if err := s.client.Write(msg); err != nil {
		return fmt.Errorf("failed to send message to server: %w", err)
	}

	return nil
}

// handleServerData handles all data received from the server.
func (s *MultiplayerJoinScreen) handleServerData(b []byte) error {
	fmt.Println("Received from server:", string(b))

	var msg comms.Message
	if err := json.Unmarshal(b, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal bytes from client: %w", err)
	}

	switch msg.Type {
	case comms.TypeEventData:
		var data comms.EventData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		s.handleEventData(data)
		return nil
	case comms.TypePlayerData:
		var data comms.PlayerData
		if err := json.Unmarshal(msg.Content, &data); err != nil {
			return fmt.Errorf("failed to unmarshal player data: %w", err)
		}
		return s.handlePlayerData(data)
	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// handleEventData handles incoming event data.
func (s *MultiplayerJoinScreen) handleEventData(data comms.EventData) {
	if data.Event == comms.EventHostStartGame {
		s.hostIsReady <- true
	}
}

// handleEventData handles incoming player data.
func (s *MultiplayerJoinScreen) handlePlayerData(data comms.PlayerData) error {
	// Make sure versions are compatible
	if data.Version != config.Version {
		return fmt.Errorf("incompatible versions (peer %s, local %s)", data.Version, config.Version)
	}

	s.opponentStatus.SetText(fmt.Sprintf("Waiting for \"%s\" to start the game", data.Username))

	return nil
}
