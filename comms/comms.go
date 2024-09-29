package comms

import "github.com/z-riley/go-2048-battle/backend"

// MessageType defines the type of message.
type MessageType string

const (
	TypePlayerData = "playerData"
	TypeGameData   = "gameData"
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

// GameData contains a game state.
type GameData struct {
	Game backend.Game `json:"game"`
}
