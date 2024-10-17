package screen

import "github.com/z-riley/turdgl"

// InitData can be passed to screens' Init function to share data between screens.
type InitData map[string]any

type Screen interface {
	// Enter initialises the screen.
	Enter(InitData)
	// Update updates and draws the screen.
	Update()
	// Exit deinitialises the screen.
	Exit()
}

// ID is the unique identifier for a screen.
type ID string

const (
	Title           ID = "title"
	Singleplayer    ID = "singleplayer"
	MultiplayerMenu ID = "multiplayerMenu"
	MultiplayerJoin ID = "multiplayerJoin"
	MultiplayerHost ID = "multiplayerHost"
	Multiplayer     ID = "multiplayer"
)

func (id ID) String() string {
	return string(id)
}

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
	CurrentScreen().Exit()
	currentScreen = s
	CurrentScreen().Enter(data)
}
