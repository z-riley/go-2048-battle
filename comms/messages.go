package comms

import (
	"encoding/json"

	"github.com/z-riley/go-2048-battle/backend"
)

// Message contains data for multiplayer mode communication. The Type field
// should first be decoded, then the appropriate action can be taken to
// decode the Content field.
type Message struct {
	Type    MessageType `json:"type"`
	Content []byte      `json:"content"`
}

// ParseMessage returns a message from a byte slice.
func ParseMessage(b []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(b, &msg)
	return msg, err
}

// MessageType defines the type of message.
type MessageType string

const (
	TypePlayerData MessageType = "playerData"
	TypeGameData               = "gameData"
	TypeEventData              = "eventData"
	TypeRequest                = "request"
)

// PlayerData contains data about a player.
type PlayerData struct {
	Version  string `json:"version"`
	Username string `json:"username"`
}

// ParsePlayerData returns player data from a byte slice.
func ParsePlayerData(b []byte) (PlayerData, error) {
	var data PlayerData
	err := json.Unmarshal(b, &data)
	return data, err
}

// GameData contains a game's current state.
type GameData struct {
	Game backend.Game `json:"game"`
}

// ParseGameData returns game data from a byte slice.
func ParseGameData(b []byte) (GameData, error) {
	var data GameData
	err := json.Unmarshal(b, &data)
	return data, err
}

// EventData contains an event which has occurred.
type EventData struct {
	Event Event `json:"event"`
}

// Event is a stand-alone occurrence.
type Event string

const (
	// EventHostStartGame signifies that the host is starting the game.
	EventHostStartGame Event = "host started game"
	// EventScreenLoaded signifies that the screen has finished initialising.
	EventScreenLoaded Event = "screen loaded"
)

// ParseEventData returns event data from a byte slice.
func ParseEventData(b []byte) (EventData, error) {
	var data EventData
	err := json.Unmarshal(b, &data)
	return data, err
}

// RequestData contains a request for data of a certain type.
type RequestData struct {
	Request MessageType `json:"request"`
}

// ParseRequestData returns request data from a byte slice.
func ParseRequestData(b []byte) (RequestData, error) {
	var data RequestData
	err := json.Unmarshal(b, &data)
	return data, err
}
