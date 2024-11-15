package comms

import (
	"encoding/json"

	"github.com/z-riley/go-2048-battle/common/backend"
)

// Message contains data for multiplayer mode communication. The Type field
// should first be decoded, then the appropriate action can be taken to
// decode the Content field.
type Message struct {
	Type    MessageType `json:"type"`
	Content []byte      `json:"content"`
}

// ParseMessage returns a message from a byte slice.
func ParseMessage(b []byte) (m Message, err error) {
	err = json.Unmarshal(b, &m)
	return m, err
}

// MessageType defines the type of message.
type MessageType string

const (
	TypePlayerData  MessageType = "playerData"
	TypeGameData    MessageType = "gameData"
	TypeEventData   MessageType = "eventData"
	TypeRequestData MessageType = "request"
)

// PlayerData contains data about a player.
type PlayerData struct {
	Version  string `json:"version"`
	Username string `json:"username"`
}

// ParsePlayerData returns player data from a byte slice.
func ParsePlayerData(b []byte) (d PlayerData, err error) {
	err = json.Unmarshal(b, &d)

	return d, err
}

// Serialise converts player data into a byte slice.
func (d PlayerData) Serialise() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return json.Marshal(Message{TypePlayerData, b})
}

// GameData contains a game's current state.
type GameData struct {
	Game backend.Game `json:"game"`
}

// ParseGameData returns game data from a byte slice.
func ParseGameData(b []byte) (d GameData, err error) {
	err = json.Unmarshal(b, &d)
	return d, err
}

// Serialise converts game data into a byte slice.
func (d GameData) Serialise() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Message{TypeGameData, b})
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
func ParseEventData(b []byte) (d EventData, err error) {
	err = json.Unmarshal(b, &d)
	return d, err
}

// Serialise converts event data into a byte slice.
func (d EventData) Serialise() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Message{TypeEventData, b})
}

// RequestData contains a request for data of a certain type.
type RequestData struct {
	Request MessageType `json:"request"`
}

// ParseRequestData returns request data from a byte slice.
func ParseRequestData(b []byte) (d RequestData, err error) {
	err = json.Unmarshal(b, &d)
	return d, err
}

// Serialise converts request data into a byte slice.
func (d RequestData) Serialise() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Message{TypeRequestData, b})
}
