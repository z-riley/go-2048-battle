package screen

import (
	"errors"
	"fmt"
	"image/color"
	"net"
	"os"
	"strings"

	"github.com/z-riley/go-2048-battle/common"
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
	labels      []*turdgl.TextBox
	buttons     []*common.MenuButton
	entries     []*common.EntryBox
	playerCards []*playerCard

	server *turdserve.Server
}

// NewTitle Screen constructs a new multiplayer host screen in the given window.
func NewMultiplayerHostScreen(win *turdgl.Window) *MultiplayerHostScreen {
	ipAddr := func() string {
		if isWSL() {
			return "check WSL host"
		}
		conn, err := LocalIP()
		if err != nil {
			panic(err)
		}
		return conn.String()
	}()

	title := turdgl.NewText("Host game", turdgl.Vec{X: 600, Y: 120}, common.FontPathMedium).
		SetColour(common.LightFontColour).
		SetAlignment(turdgl.AlignCentre).
		SetSize(40)

	nameHeading := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 200}, func() {})
	nameHeading.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Your name:")
	nameEntry := common.NewEntryBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 200})

	ipHeading := common.NewTextBox(400, 60, turdgl.Vec{X: 200 - 20, Y: 300})
	ipHeading.Body.SetOffset(turdgl.Vec{X: 0, Y: 32}).SetText("Your IP:")
	ipBody := common.NewTextBox(400, 60, turdgl.Vec{X: 600 + 20, Y: 300})
	ipBody.Body.SetOffset(turdgl.Vec{X: 0, Y: 32}).SetText(ipAddr)

	var cards []*playerCard
	for i := 0; i < 1; i++ {
		pos := turdgl.Vec{X: 300, Y: 450 + float64(i)*80}
		label := fmt.Sprintf("Waiting for player %d", i+2)
		cards = append(cards, newPlayerCard(pos, label))
	}

	start := common.NewMenuButton(400, 60, turdgl.Vec{X: 200 - 20, Y: 650}, func() { SetScreen(Title) })
	start.SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Start game")

	back := common.NewMenuButton(400, 60, turdgl.Vec{X: 600 + 20, Y: 650}, func() { SetScreen(MultiplayerMenu) })
	back.SetLabelAlignment(turdgl.AlignCustom).
		SetLabelOffset(turdgl.Vec{X: 0, Y: 32}).SetLabelText("Back")

	return &MultiplayerHostScreen{
		win:         win,
		title:       title,
		labels:      []*turdgl.TextBox{ipHeading, ipBody},
		buttons:     []*common.MenuButton{nameHeading, start, back},
		entries:     []*common.EntryBox{nameEntry},
		playerCards: cards,
		server:      turdserve.NewServer(maxPlayers - 1),
	}
}

// Init initialises the screen.
func (s *MultiplayerHostScreen) Init() {
	s.win.RegisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease, func() {
		SetScreen(MultiplayerMenu)
	})

	// Set up server
	s.server.SetCallback(func(id int, msg []byte) {
		fmt.Println("Received from client", id, ":", string(msg))
	})

	// Start server to allow other players to connect
	go func() {
		s.server.Run("0.0.0.0", serverPort)
	}()
}

// Deinit deinitialises the screen.
func (s *MultiplayerHostScreen) Deinit() {
	s.win.UnregisterKeybind(turdgl.KeyEscape, turdgl.KeyRelease)
	s.server.Destroy()
}

// Update updates and draws multiplayer host screen.
func (s *MultiplayerHostScreen) Update() {
	s.win.SetBackground(common.BackgroundColour)

	s.win.Draw(s.title)

	for _, l := range s.labels {
		s.win.Draw(l)
	}

	for _, b := range s.buttons {
		b.Update(s.win)
		s.win.Draw(b)
	}

	for _, e := range s.entries {
		s.win.Draw(e)
		e.Update(s.win)
	}

	for _, p := range s.playerCards {
		s.win.Draw(p)
	}

}

var (
	styleNotReady = turdgl.Style{Colour: color.RGBA{255, 0, 0, 255}, Thickness: 0, Bloom: 2}
	styleReady    = turdgl.Style{Colour: color.RGBA{0, 255, 0, 255}, Thickness: 0, Bloom: 10}
)

// playerCard displays a player in the lobby.
type playerCard struct {
	name  *common.EntryBox
	light *turdgl.Circle
}

func newPlayerCard(pos turdgl.Vec, txt string) *playerCard {
	const (
		width  = 400
		height = 60
	)

	name := common.NewEntryBox(width, height, pos)
	name.SetTextAlignment(turdgl.AlignCustom).
		SetTextOffset(turdgl.Vec{X: 0, Y: 32}).
		SetTextSize(30).
		SetTextColour(common.DarkerFontColour).
		SetText(txt)

	lightPos := turdgl.Vec{X: pos.X + width + 40, Y: pos.Y + height/2}
	light := turdgl.NewCircle(height*0.8, lightPos, turdgl.WithStyle(styleNotReady))

	return &playerCard{name, light}
}

func (p *playerCard) Draw(buf *turdgl.FrameBuffer) {
	p.name.Draw(buf)
	p.light.Draw(buf)
}

func (p *playerCard) setReady() {
	p.light.SetStyle(styleReady)
}

func (p *playerCard) setNotReady() {
	p.light.SetStyle(styleNotReady)
}

//////////////

// LocalIP returns the device's local IP address.
func LocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

// isPrivate IP returns true if the given IP address is reserved (private).
func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// isWSL returns true if the device is a WSL instance.
func isWSL() bool {
	data, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return false // unable to read, assume not WSL
	}
	return strings.Contains(string(data), "microsoft")
}
