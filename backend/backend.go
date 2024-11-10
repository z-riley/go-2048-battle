package backend

import (
	"encoding/json"

	"github.com/z-riley/go-2048-battle/backend/grid"
	"github.com/z-riley/go-2048-battle/backend/store"
	"github.com/z-riley/go-2048-battle/log"
)

// Game contains all required data for a 2048 game.
type Game struct {
	Grid  *grid.Grid `json:"grid"`
	Score *Score     `json:"currentScore"`
	Timer *Timer     `json:"time"`

	store *store.Store
	opts  *Opts
}

// opts contains configuration for the game.
type Opts struct {
	SaveToDisk bool
}

// NewGame returns the top-level struct for the game. If opts are nil, the
// default is used.
func NewGame(opts *Opts) *Game {
	if opts == nil {
		opts = &Opts{
			SaveToDisk: true,
		}
	}

	g := &Game{
		Grid:  grid.NewGrid(),
		Score: NewScore(),
		Timer: NewTimer(),
		store: store.NewStore(".save.bruh"),
		opts:  opts,
	}

	if g.opts.SaveToDisk {
		err := g.Load()
		if err != nil {
			log.Println("No save file found. Creating new one")
			if err := g.Save(); err != nil {
				panic(err)
			}
		}
	}

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

	// Note: The game should save on exit anyway but save after move just in case
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
	return g.store.SaveBytes(j)
}

// Load loads the game state from the save file.
func (g *Game) Load() error {
	b, err := g.store.ReadBytes()
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
