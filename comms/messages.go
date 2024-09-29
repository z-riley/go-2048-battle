package comms

import "github.com/z-riley/go-2048-battle/backend"

// MessageType defines the type of message.
type MessageType string

const (
	TypePlayerData = "playerData"
	TypeGameData   = "gameData"
	TypeEventData  = "eventData"
)

// Message contains data for multiplayer mode communication. The Type field
// should first be decoded, then the appropriate action can be taken to
// decode the Content field.
type Message struct {
	Type    MessageType `json:"type"`
	Content []byte      `json:"content"`
}

// PlayerData contains data about a player.
type PlayerData struct {
	Version  string `json:"version"`
	Username string `json:"username"`
}

// GameData contains a game's current state.
type GameData struct {
	Game backend.Game `json:"game"`
}

// EventData contains an event which has occurred.
type EventData struct {
	Event Event `json:"event"`
}

// Event is a type of stand-alone occurred.
type Event string

const (
	// EventHostStartGame signifies that the host is starting the game.
	EventHostStartGame Event = "host started game"
)
