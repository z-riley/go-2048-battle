package screen

type Screen int

const (
	Title Screen = iota
	Singleplayer
	MultiplayerMenu
	MultiplayerJoin
)

// Updater represents screens that can update themselves, including drawing themselves.
type Updater interface {
	Update()
}

// currentScreen holds the game's current screen.
var currentScreen = Title

func CurrentScreen() Screen {
	return currentScreen
}

func SetScreen(s Screen) {
	currentScreen = s
}
