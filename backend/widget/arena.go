package widget

import (
	"sync"
)

type Arena struct {
	mu   sync.Mutex
	Grid *grid `json:"grid"`
}

// NewArena returns the game arena widget.
func NewArena() *Arena {
	a := Arena{
		mu:   sync.Mutex{},
		Grid: newGrid(),
	}

	return &a
}

// Direction represents a direction that the player can move the tiles in.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// Move attempts to move in the specified direction, spawning a new tile if appropriate.
func (a *Arena) Move(dir Direction) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	didMove, pointsGained := a.Grid.move(dir)
	if didMove {
		a.Grid.spawnTile()
	}
	return pointsGained
}

// ResetGrid resets the arena.
func (a *Arena) Reset() {
	a.Grid.resetGrid()
}

// Outcome represents the outcome of a game.
type Outcome int

const (
	None Outcome = iota
	Win
	Lose
)

// Outcome returns the current outcome of the arena.
func (a *Arena) Outcome() Outcome {
	switch {
	case a.Grid.isLoss():
		return Lose
	case a.Grid.highestTile() >= 2048:
		return Win
	default:
		return None
	}
}
