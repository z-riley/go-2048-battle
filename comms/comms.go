package comms

// MessageType defines the type of message.
type MessageType string

const (
	MsgConnect    = "connect"
	MsgGameUpdate = "gridUpdate"
)

// Message contains data for multiplayer mode communication. The Type field
// should first be decoded, then the appropriate action can be taken to
// decode the Content field.
type Message struct {
	Type    MessageType `json:"type"`
	Content []byte      `json:"content"`
}

// PlayerData is for sharing data about the connected players.
type PlayerData struct {
	Version  string `json:"version"`
	Username string `json:"username"`
}

/*
Deserialising:
1. Decode the type field
2. Decide how to decode the other
*/
