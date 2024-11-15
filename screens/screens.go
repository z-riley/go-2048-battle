package screens

import (
	"github.com/z-riley/turdgl"
)

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

// screenChange contains data required to change the current screen.
type screenChange struct {
	id   ID
	data InitData
}

// currentScreen holds the game's current screen.
var currentScreen = Title

// screenChangeChan executes the screen change sequence on receipt.
var screenChangeChan = make(chan screenChange, 1)

// Update updates the current screen.
func Update() {
	select {
	case screen := <-screenChangeChan:
		CurrentScreen().Exit()
		currentScreen = screen.id
		CurrentScreen().Enter(screen.data)

	default:
		CurrentScreen().Update()
	}
}

// CurrentScreen returns the current screen.
func CurrentScreen() Screen {
	return screens[currentScreen]
}

// SetScreen changes the current screen to the given ID next time Update is called.
func SetScreen(id ID, data InitData) {
	switch id {
	case Title, Singleplayer, MultiplayerMenu, MultiplayerJoin, MultiplayerHost, Multiplayer:
		screenChangeChan <- screenChange{id, data}
	default:
		panic("invalid screen: " + id)
	}
}
