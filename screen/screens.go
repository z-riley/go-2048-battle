package screen

import "github.com/z-riley/turdgl"

// InitData can be passed to screens' Init function to share data between screens.
type InitData map[string]any

type Screen interface {
	// Init initialises the screen.
	Init(InitData)
	// Update updates and draws the screen.
	Update()
	// Deinit deinitialises the screen.
	Deinit()
}

// ID is the unique identifier for a screen.
type ID int

const (
	Title ID = iota
	Singleplayer
	MultiplayerMenu
	MultiplayerJoin
	MultiplayerHost
	Multiplayer
)

// currentScreen holds the game's current screen.
var currentScreen = Title

// screens contains every screen.
var screens map[ID]Screen

// Init populates the internal screens variable.
func Init(win *turdgl.Window) {
	screens = map[ID]Screen{
		Title:           NewTitleScreen(win),
		Singleplayer:    NewSingleplayerScreen(win),
		MultiplayerMenu: NewMultiplayerMenuScreen(win),
		MultiplayerJoin: NewMultiplayerJoinScreen(win),
		MultiplayerHost: NewMultiplayerHostScreen(win),
		Multiplayer:     NewMultiplayerScreen(win),
	}
}

// CurrentScreen returns the current screen.
func CurrentScreen() Screen {
	return screens[currentScreen]
}

// SetScreen changes the current screen to the given ID.
func SetScreen(s ID, data InitData) {
	CurrentScreen().Deinit()
	currentScreen = s
	CurrentScreen().Init(data)
}
