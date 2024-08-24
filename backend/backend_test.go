package backend

import (
	"testing"
	"time"
)

func TestSerialiseDeserialise(t *testing.T) {
	// Create a game and let the timer change value
	game := NewGame()
	game.Timer.Reset().Resume()
	time.Sleep(100 * time.Millisecond)
	b, err := game.Serialise()
	if err != nil {
		t.Error(err)
	}

	// Create an empty game and set it to the previous game's state
	g := Game{}
	g.Deserialise(b)

	// Check that the two games are identical
	expected, err := game.Serialise()
	got, err := g.Serialise()
	if err != nil {
		t.Error(err)
	}
	if string(expected) != string(got) {
		t.Errorf("Expected:\n<%v>\nGot:\n<%v>", expected, got)
	}
}
