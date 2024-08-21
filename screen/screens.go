package screen

import "github.com/z-riley/turdgl"

type ID int

type Screen interface {
	// Init initialises the screen.
	Init()
	// Update updates and draws the screen.
	Update()
	// Deinit deinitialises the screen.
	Deinit()
}

const (
	Title ID = iota
	Singleplayer
	MultiplayerMenu
	MultiplayerJoin
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
	}
}

// CurrentScreen returns the current screen.
func CurrentScreen() Screen {
	return screens[currentScreen]
}

// SetScreen changes the current screen to the given ID.
func SetScreen(s ID) {
	CurrentScreen().Deinit()
	currentScreen = s
	CurrentScreen().Init()
}
