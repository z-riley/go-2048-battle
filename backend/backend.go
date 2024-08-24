package backend

import (
	"encoding/json"
	"fmt"

	"github.com/z-riley/go-2048-battle/backend/store"
	"github.com/z-riley/go-2048-battle/backend/widget"
)

type Game struct {
	Arena   *widget.Arena  `json:"arena"`
	Outcome widget.Outcome `json:"outcome"`
	Score   *widget.Score  `json:"currentScore"`
	Timer   *widget.Timer  `json:"time"`
}

// NewGame returns the top-level struct for the game.
func NewGame() *Game {
	g := &Game{
		Arena:   widget.NewArena(),
		Outcome: widget.None,
		Score:   widget.NewScore(),
		Timer:   widget.NewTimer(),
	}

	err := g.Load()
	if err != nil {
		fmt.Println("No save file found")
		if err := g.Save(); err != nil {
			panic(err)
		}
	}

	return g
}

// Serialise converts the current game state into JSON.
func (g *Game) Serialise() ([]byte, error) {
	return json.Marshal(g)
}

// Deserialise loads the JSON representation into memory.
func (g *Game) Deserialise(j []byte) error {
	return json.Unmarshal(j, &g)
}

// Save saves the game state to the save file.
func (g *Game) Save() error {
	j, err := g.Serialise()
	if err != nil {
		return err
	}
	return store.SaveBytes(j)
}

// Load loads the game state from the save file.
func (g *Game) Load() error {
	b, err := store.ReadBytes()
	if err != nil {
		return err
	}

	return g.Deserialise(b)
}

// ExecuteMove carries out a move (up, down, left, right).
func (g *Game) ExecuteMove(dir widget.Direction) {
	pointsGained := g.Arena.Move(dir)
	g.Score.AddToCurrent(pointsGained)

	if g.Outcome == widget.Lose {
		g.Timer.Pause()
	} else {
		g.Timer.Resume()
	}

	// Note: the game should save on exit anyway but save after after move just in case
	go func() {
		if err := g.Save(); err != nil {
			panic(err)
		}
	}()
}

// Reset resets the game.
func (g *Game) Reset() {
	// widget.SetCurrentScore(0)
	g.Arena.Reset()
	g.Score.Reset()
	g.Timer.Reset().Pause()
}
