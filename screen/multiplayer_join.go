package screen

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/z-riley/go-2048-battle/backend/store"
	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/log"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerJoinScreen struct {
	win *turdgl.Window

	title            *turdgl.Text
	nameHeading      *turdgl.Text
	nameEntry        *common.EntryBox
	ipHeading        *turdgl.Text
	ipStore          *store.Store
	ipEntry          *common.EntryBox
	opponentName     string
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
	s.title = turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(100)

	s.nameHeading = turdgl.NewText(
		"Your name:",
		turdgl.Vec{X: 600, Y: 230},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(30)

	s.nameEntry = common.NewEntryBox(
		400, 60,
		turdgl.Vec{X: 600 - 400/2, Y: s.nameHeading.Pos().Y + 15},
		namesgenerator.GetRandomName(0),
	).
		SetModifiedCB(func() {
			// Update host with new username
			if s.client != nil {
				if err := s.sendPlayerData(); err != nil {
					log.Println("Failed to send username update to host:", err)
				}
			}
		})

	s.ipHeading = turdgl.NewText(
		"Host IP:",
		turdgl.Vec{X: 600, Y: 370},
		common.FontPathMedium,
	).
		SetColour(common.GreyTextColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(30)

	s.ipStore = store.NewStore(".ip.bruh")
	b, err := s.ipStore.ReadBytes()
	if err != nil {
		log.Println("Failed to read IP address store:", err)
		b = []byte("Enter IP address")
	}

	s.ipEntry = common.NewEntryBox(
		400, 60,
		turdgl.Vec{X: 600 - 400/2, Y: s.ipHeading.Pos().Y + 15},
		string(b),
	).SetModifiedCB(func() {
		if err := s.ipStore.SaveBytes([]byte(s.ipEntry.Text())); err != nil {
			log.Println("Failed to save IP address to store")
		}
	})

	s.opponentStatus = turdgl.NewText(
		"",
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

	// Set up client
	s.hostIsReady = make(chan bool)
	s.join = common.NewMenuButton(
		TileSizePx, TileSizePx,
		turdgl.Vec{
			X: s.buttonBackground.Pos.X + TileSizePx*TileBoundryFactor,
			Y: s.buttonBackground.Pos.Y + TileSizePx*TileBoundryFactor,
		},
		s.joinButtonHandler,
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

// joinButtonHandler handles presses of the join button.
func (s *MultiplayerJoinScreen) joinButtonHandler() {

	// Handle asynchronous errors from client
	errCh := make(chan error)
	go func() {
		for err := range errCh {
			if err != nil {
				log.Println("Client error:", err)

				// Re-enable button
				s.join.SetCallback(
					turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
					s.joinButtonHandler,
				)

				// Display error to user
				s.opponentStatus.SetText("Lost connection with host")
				go func() {
					timer := time.NewTimer(2 * time.Second)
					<-timer.C
					s.opponentStatus.SetText("")
				}()

			}
		}
	}()

	err := s.joinGame(errCh)
	if err != nil {
		s.opponentStatus.SetText("Failed to connect to host")
		go func() {
			time.Sleep(time.Second)
			s.opponentStatus.SetText("")
		}()
		log.Println("Failed to join game:", err)
		return
	}

	// Disable the button so user can't connect again
	s.join.SetCallback(
		turdgl.ButtonTrigger{State: turdgl.LeftClick, Behaviour: turdgl.OnRelease},
		func() {},
	)

	go func() {
		if <-s.hostIsReady {
			SetScreen(Multiplayer, InitData{
				clientKey:           s.client,
				usernameKey:         s.nameEntry.Text(),
				opponentUsernameKey: s.opponentName,
			})
			return
		}
	}()
}

// joinGame attempts to join a multiplayer game.
func (s *MultiplayerJoinScreen) joinGame(errCh chan error) error {

	// Connect using the user-specified IP address
	if err := s.client.Connect(context.Background(), s.ipEntry.Text(), serverPort, errCh); err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	s.client.SetCallback(func(b []byte) {
		if err := s.handleServerData(b); err != nil {
			log.Println("Join screen failed to handle data from server:", err)
		}
	})

	// Send player data to host
	if err := s.sendPlayerData(); err != nil {
		return fmt.Errorf("failed to send player data: %w", err)
	}

	return nil
}

// handleServerData handles all data received from the server.
func (s *MultiplayerJoinScreen) handleServerData(data []byte) error {
	msg, err := comms.ParseMessage(data)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	switch msg.Type {
	case comms.TypeEventData:
		eventData, err := comms.ParseEventData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse event data: %w", err)
		}
		return s.handleEventData(eventData)

	case comms.TypePlayerData:
		playerData, err := comms.ParsePlayerData(msg.Content)
		if err != nil {
			return fmt.Errorf("failed to parse player data: %w", err)
		}
		return s.handlePlayerData(playerData)

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// handleEventData handles incoming event data.
func (s *MultiplayerJoinScreen) handleEventData(data comms.EventData) error {
	if data.Event == comms.EventHostStartGame {
		s.hostIsReady <- true
	}
	return nil
}

// sendPlayerData sends the player data to the host.
func (s *MultiplayerJoinScreen) sendPlayerData() error {
	msg, err := comms.PlayerData{
		Version:  config.Version,
		Username: s.nameEntry.Text(),
	}.Serialise()
	if err != nil {
		return fmt.Errorf("failed to serialise player data: %w", err)
	}

	return s.client.Write(msg)
}

// handleEventData handles incoming player data.
func (s *MultiplayerJoinScreen) handlePlayerData(data comms.PlayerData) error {
	// Make sure versions are compatible
	if data.Version != config.Version {
		return fmt.Errorf("incompatible versions (peer %s, local %s)", data.Version, config.Version)
	}

	s.opponentName = data.Username
	s.opponentStatus.SetText(fmt.Sprintf("Waiting for \"%s\" to start the game", s.opponentName))

	return nil
}
