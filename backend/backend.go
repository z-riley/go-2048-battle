package backend

import (
	"encoding/json"
	"fmt"

	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/backend/store"
)

// Game contains all reuquired data for a 2048 game.
type Game struct {
	Grid  *grid.Grid `json:"grid"`
	Score *Score     `json:"currentScore"`
	Timer *Timer     `json:"time"`

	opts *Opts
}

// ops contains configuration for the game.
type Opts struct {
	SaveToDisk bool
	ResetKey   string
}

// NewGame returns the top-level struct for the game. If opts are nil, the
// default is used.
func NewGame(opts *Opts) *Game {
	if opts == nil {
		opts = &Opts{
			SaveToDisk: true,
			ResetKey:   "",
		}
	}

	g := &Game{
		Grid:  grid.NewGrid(),
		Score: NewScore(),
		Timer: NewTimer(),
		opts:  opts,
	}

	// Use pseudo random grid reset (for multiplayer synchronisation) if key provided
	if opts.ResetKey != "" {
		g.Grid.PseudoRandomReset(opts.ResetKey)
	}

	if g.opts.SaveToDisk {
		err := g.Load()
		if err != nil {
			fmt.Println("No save file found. Creating new one")
			if err := g.Save(); err != nil {
				panic(err)
			}
		}
	}

	return g
}

func (g *Game) PseudoRandomReset(seed string) *Game {
	g.Grid.PseudoRandomReset(seed)
	g.Score.Reset()
	g.Timer.Reset().Pause()
	return g
}

// Reset resets the game.
func (g *Game) Reset() *Game {
	g.Grid.Reset()
	g.Score.Reset()
	g.Timer.Reset().Pause()
	return g
}

// ExecuteMove carries out a move in the given direction.
func (g *Game) ExecuteMove(dir grid.Direction) {
	pointsGained := g.Grid.Move(dir)
	g.Score.AddToCurrent(pointsGained)

	if g.Grid.Outcome() == grid.Lose {
		g.Timer.Pause()
	} else {
		g.Timer.Resume()
	}

	// Note: the game should save on exit anyway but save after after move just in case
	if g.opts.SaveToDisk {
		go func() {
			if err := g.Save(); err != nil {
				panic(err)
			}
		}()
	}
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
func (g Game) Save() error {
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
	err = g.Deserialise(b)
	if err != nil {
		return err
	}
	// Cmb flags are required to be unset for the animations to work correctly
	g.Grid.ClearCmbFlags()
	return nil
}
