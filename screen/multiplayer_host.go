package screen

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/go-2048-battle/comms"
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/turdgl"
	"github.com/z-riley/turdserve"
)

const (
	maxPlayers = 2
	serverPort = 8080
)

type MultiplayerHostScreen struct {
	win *turdgl.Window

	title       *turdgl.Text
	ipHeading   *turdgl.TextBox
	ipBody      *turdgl.TextBox
	nameHeading *common.MenuButton
	nameEntry   *common.EntryBox
	start       *common.MenuButton
	back        *common.MenuButton
	playerCards []*playerCard

	server *turdserve.Server
}

// NewMultiplayerHostScreen constructs an uninitialised multiplayer host screen.
func NewMultiplayerHostScreen(win *turdgl.Window) *MultiplayerHostScreen {
	return &MultiplayerHostScreen{win: win}
}

// Enter initialises the screen.
func (s *MultiplayerHostScreen) Enter(_ InitData) {
	s.title = turdgl.NewText("Host game", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.ArenaBackgroundColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	s.ipHeading = common.NewTextBox(400, 60, turdgl.Vec{X: 200 - 20, Y: 300}).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).SetText("Your IP:")
	s.ipBody = common.NewTextBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 300}).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).SetText(s.getIPAddr())

	s.nameHeading = common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	s.nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Your name:")
	s.nameEntry = common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})

	s.start = common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 650}, func() {})
	s.start.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Start game")
	s.start.SetCallback(func(_ turdgl.MouseState) {
		if err := s.startGame(); err != nil {
			fmt.Println("Failed to start game:", err)
		}
	})

	s.back = common.NewMenuButton(400, 60, turdgl.Vec{X: 600 + 20, Y: 650},
		func() { SetScreen(MultiplayerMenu, nil) },
	)
	s.back.SetLabelAlignment(turdgl.AlignCustom).SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	s.playerCards = make([]*playerCard, maxPlayers-1)
	for i := range s.playerCards {
		pos := turdgl.Vec{X: 300, Y: 450 + float64(i)*80}
		s.playerCards[i] = newPlayerCard(pos, i)
	}

	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu, nil)
	})

	// Set up server
	{
		s.server = turdserve.NewServer(maxPlayers - 1).
			SetCallback(func(id int, b []byte) {
				if err := s.handleClientData(id, b); err != nil {
					fmt.Println("Failed to handle data from client:", err)
				}
			}).SetDisconnectCallback(func(id int) {
			s.handleClientDisconnect(id)
		})

		// Start server to allow other players to connect
		errCh := make(chan error)
		go func() {
			for err := range errCh {
				if err != nil {
					panic("server crashed: " + err.Error())
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
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
}

// Update updates and draws multiplayer host screen.
func (s *MultiplayerHostScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, l := range []*turdgl.TextBox{
		s.ipHeading,
		s.ipBody,
	} {
		s.win.Draw(l)
	}

	for _, b := range []*common.MenuButton{
		s.nameHeading,
		s.start,
		s.back,
	} {
		b.Update(s.win)
		s.win.Draw(b)
	}

	s.win.Draw(s.nameEntry)
	s.nameEntry.Update(s.win)

	for _, p := range s.playerCards {
		s.win.Draw(p)
	}
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

	// Update player card with new data
	s.playerCards[id].setReady(data.Username)
	return nil
}

func (s *MultiplayerHostScreen) handleClientDisconnect(id int) {
	s.playerCards[id].setNotReady()
}

// getIPAddr returns the IP address of the host.
func (s *MultiplayerHostScreen) getIPAddr() string {
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
	// Check all players are connected
	for i := range maxPlayers - 1 {
		if !s.playerCards[i].isReady() {
			return fmt.Errorf("player %d not ready", i)
		}
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

// Styles for playerCards.
var (
	styleNotReady = turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 2}
	styleReady    = turdgl.Style{Colour: color.RGBA{0, 255, 0, 255}, Thickness: 0, Bloom: 10}
)

// playerCard displays a player in the lobby.
type playerCard struct {
	id    int
	ready bool

	name  *common.EntryBox
	light *turdgl.Circle
}

// newPlayerCard constructs a new player card. The ID must be unique.
func newPlayerCard(pos turdgl.Vec, id int) *playerCard {
	const (
		width  = 400
		height = 60
	)

	name := common.NewEntryBox(width, height, pos)
	name.SetTextAlignment(turdgl.AlignCustom).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetTextSize(30).
		SetTextColour(common.GreyTextColour).
		SetText(fmt.Sprintf("Waiting for player %d", id+2))

	lightPos := turdgl.Vec{X: pos.X + width + 40, Y: pos.Y + height/2}
	light := turdgl.NewCircle(height*0.8, lightPos, turdgl.WithStyle(styleNotReady))

	return &playerCard{
		id:    id,
		ready: false,
		name:  name,
		light: light,
	}
}

// Draw draws a player card to the supplied framebuffer.
func (p *playerCard) Draw(buf *turdgl.FrameBuffer) {
	p.name.Draw(buf)
	p.light.Draw(buf)
}

// setReady sets a player's card to ready.
func (p *playerCard) setReady(username string) {
	p.ready = true
	p.name.Text.SetText(username)
	p.light.SetStyle(styleReady)
}

// setNotReady removes player data from card and sets the status to not connected.
func (p *playerCard) setNotReady() {
	p.ready = false
	p.name.Text.SetText(fmt.Sprintf("Waiting for player %d", p.id+2))
	p.light.SetStyle(styleNotReady)
}

// isReady returns whether the player is ready.
func (p *playerCard) isReady() bool {
	return p.ready
}
