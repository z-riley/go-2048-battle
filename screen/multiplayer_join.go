package screen

import (
	"encoding/json"
	"fmt"

	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

type MultiplayerJoinScreen struct {
	win *turdgl.Window

	title   *turdgl.Text
	buttons []*common.MenuButton
	entries []*common.EntryBox

	hostIsReady chan bool
	client      *turdserve.Client
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerJoinScreen(win *turdgl.Window) *MultiplayerJoinScreen {
	title := turdgl.NewText("Join game", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.ArenaBackgroundColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	ipHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	ipHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Host IP:")

	ipEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})
	ipEntry.SetText("127.0.0.1") // temporary for local testing

	nameHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 300}, func() {})
	nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Your name:")

	nameEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 300})

	join := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 400}, func() {})
	join.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Join")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 400, Y: 500}, func() {})
	back.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	s := MultiplayerJoinScreen{
		win:         win,
		title:       title,
		buttons:     []*common.MenuButton{ipHeading, nameHeading, join, back},
		entries:     []*common.EntryBox{nameEntry, ipEntry},
		hostIsReady: make(chan bool),
		client:      turdserve.NewClient(),
	}

	back.SetCallback(func(_ turdgl.MouseState) {
		join.SetLabelText("Join")
		s.client.Destroy()
		SetScreen(MultiplayerMenu, nil)
	})

	join.SetCallback(func(_ turdgl.MouseState) {
		if err := s.joinGame(); err != nil {
			fmt.Println("Failed to join game:", err)
			return
		}

		join.SetLabelText("Waiting for host")

		// Disable the button
		join.SetCallback(func(_ turdgl.MouseState) {})

		go func() {
			if <-s.hostIsReady {
				SetScreen(Multiplayer, InitData{clientKey: s.client})
				return
			}
		}()
	})

	return &s
}

// Init initialises the screen.
func (s *MultiplayerJoinScreen) Init(_ InitData) {
	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})
}

// Deinit deinitialises the screen.
func (s *MultiplayerJoinScreen) Deinit() {
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer join screen.
func (s *MultiplayerJoinScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, b := range s.buttons {
		b.Update(s.win)
		s.win.Draw(b)
	}

	for _, e := range s.entries {
		s.win.Draw(e)
		e.Update(s.win)
	}

}

// clientKey is used for indentifying the server in InitData.
const clientKey = "client"

// joinGame attempts to join a multiplayer game.
func (s *MultiplayerJoinScreen) joinGame() error {

	// Connect using the user-specified IP address
	ipEntry := s.entries[1]
	ip := ipEntry.Text.Text()
	errCh := make(chan error)
	go func() {
		for err := range errCh {
			if err != nil {
				fmt.Println(fmt.Errorf("client error: %w", err).Error())
				s.client.Destroy()
			}
		}
	}()
	if err := s.client.Connect(ip, serverPort, errCh); err != nil {
		panic(err)
	}

	s.client.SetCallback(func(b []byte) {
		if err := s.handleServerData(b); err != nil {
			fmt.Println("Failed to handle data from server:", err)
		}
	})

	// Construct message containing player data
	usernameEntry := s.entries[0]
	username := usernameEntry.Text.Text()
	playerData, err := json.Marshal(comms.PlayerData{
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

	// Send data to host
	if err := s.client.Write(msg); err != nil {
		return fmt.Errorf("failed to send message to server: %w", err)
	}

	return nil
}

// handleServerData handles all data received from the server.
func (s *MultiplayerJoinScreen) handleServerData(b []byte) error {
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

	default:
		return fmt.Errorf("unsupported message type \"%s\"", msg.Type)
	}
}

// handleEventData handles incoming player data.
func (s *MultiplayerJoinScreen) handleEventData(data comms.EventData) {
	if data.Event == comms.EventHostStartGame {
		s.hostIsReady <- true
	}
}
