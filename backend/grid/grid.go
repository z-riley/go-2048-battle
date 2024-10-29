package grid

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

const (
	GridLen    = 4
	gridWidth  = GridLen
	gridHeight = GridLen
)

// Grid contains the tiles for the game. Position {0,0} is the top left square.
type Grid struct {
	mu    sync.Mutex
	Tiles [gridWidth][gridHeight]Tile `json:"tiles"`

	LastMove Direction `json:"lastMove"`
}

// NewGrid constructs a new grid.
func NewGrid() *Grid {
	g := Grid{
		mu:    sync.Mutex{},
		Tiles: NewTiles(),
	}
	g.Reset()

	return &g
}

// Direction represents a direction that the player can move the tiles in.
type Direction string

const (
	DirUp    Direction = "up"
	DirDown  Direction = "down"
	DirLeft  Direction = "left"
	DirRight Direction = "right"
)

// Move attempts to move in the specified direction, spawning a new tile if appropriate.
func (g *Grid) Move(dir Direction) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	didMove, pointsGained := g.move(dir)
	if didMove {
		g.spawnTile()
	}
	g.LastMove = dir
	return pointsGained
}

// Reset resets the grid to a start-of-game state, spawning two '2' tiles in random locations.
func (g *Grid) Reset() {
	g.Tiles = NewTiles()
	// Place two '2' tiles in random positions
	type pos struct{ x, y int }
	tile1 := pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	tile2 := pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	for reflect.DeepEqual(tile1, tile2) {
		// Try again until they're unique
		tile2 = pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	}
	g.Tiles[tile1.x][tile1.y].Val = newTileVal()
	g.Tiles[tile2.x][tile2.y].Val = newTileVal()
}

// PseudoRandomReset resets the grid to a start-of-game state, spawning two '2' tiles in
// locations determined by the provided key. If this function is called twice with the same
// seed, the resulting grid would be the same.
func (g *Grid) PseudoRandomReset(seed string) {
	g.Tiles = NewTiles()

	// Hash the seed into a uint32
	hasher := fnv.New32a()
	if _, err := hasher.Write([]byte(seed)); err != nil {
		panic(err)
	}
	hash := hasher.Sum32()

	// Place two '2' tiles in random positions
	type pos struct{ x, y uint32 }
	tile1 := pos{(hash >> 0) & 0b11, (hash >> 2) & 0b11}
	tile2 := pos{(hash >> 4) & 0b11, (hash >> 6) & 0b11}

	// Ensure pos1 and pos2 are different
	if (tile1.x == tile2.x) && (tile1.y == tile2.y) {
		tile1.x = (tile1.x + 1) & 0b11 // wrap around with 0b11 to ensure it's in range 0-3
	}

	g.Tiles[tile1.x][tile1.y].Val = 2
	g.Tiles[tile2.x][tile2.y].Val = 2
}

// ClearCmbFlags clears the Cmb flag of every tile.
func (g *Grid) ClearCmbFlags() {
	for i := range g.Tiles {
		for j := range g.Tiles[i] {
			g.Tiles[i][j].Cmb = false
		}
	}
}

// NumTiles returns the number of non zero tiles on the grid.
func (g *Grid) NumTiles() int {
	n := 0
	for i := range g.Tiles {
		for j := range g.Tiles[i] {
			if g.Tiles[i][j].Val != 0 {
				n++
			}
		}
	}
	return n
}

// Outcome represents the outcome of a game.
type Outcome string

const (
	None Outcome = "none"
	Win  Outcome = "win"
	Lose Outcome = "lose"
)

// Outcome returns the current outcome of the grid.
func (g *Grid) Outcome() Outcome {
	switch {
	case g.isLoss():
		return Lose
	case g.HighestTile() >= 2048:
		return Win
	default:
		return None
	}
}

// spawnTile spawns a single new tile in a random location on the grid. The value of the
// tile is either 2 (90% chance) or 4 (10% chance).
func (g *Grid) spawnTile() {
	x, y := rand.Intn(gridWidth), rand.Intn(gridHeight)
	for g.Tiles[x][y].Val != emptyTile {
		// Try again until they're unique
		x, y = rand.Intn(gridWidth), rand.Intn(gridHeight)
	}

	g.Tiles[x][y].Val = newTileVal()
	g.Tiles[x][y].UUID = uuid.Must(uuid.NewV7())
}

// move attempts to move all tiles in the specified direction, combining them if appropriate.
// Returns true if any tiles were moved from the attempt, and the added score from any combinations.
func (g *Grid) move(dir Direction) (bool, int) {
	// Clear all of the "combined this turn" flags
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			g.Tiles[i][j].Cmb = false
		}
	}

	moved := false
	pointsGained := 0

	// Execute moves until grid can no longer move
	for {
		movedThisTurn := false
		for row := 0; row < gridHeight; row++ {
			var rowMoved bool
			var points int

			// The moveStep function only operates on a row, so to move vertically
			// we must transpose the grid before and after the move operation.
			if dir == DirUp || dir == DirDown {
				g.Tiles = Transpose(g.Tiles)
			}
			g.Tiles[row], rowMoved, points = moveStep(g.Tiles[row], dir)
			if points > 0 {
				pointsGained = points
			}
			if dir == DirUp || dir == DirDown {
				g.Tiles = Transpose(g.Tiles)
			}

			if rowMoved {
				movedThisTurn = true
				moved = true
			}
		}
		if !movedThisTurn {
			break
		}
	}

	return moved, pointsGained
}

// moveStep executes one part of the a move on a grid row. Call multiple times until false
// is returned to complete a full move. Returns the row after move, whether any tiles moved,
// and the number of points gained by the move.
func moveStep(g [gridWidth]Tile, dir Direction) ([gridWidth]Tile, bool, int) {
	// Iterate in the same direction as the move
	reverse := false
	if dir == DirRight || dir == DirDown {
		reverse = true
	}

	iter := NewIter(len(g), reverse)
	for iter.HasNext() {
		i := iter.Next()
		// Calculate the hypothetical next position for the tile
		newPos := i - 1
		if reverse {
			newPos = i + 1
		}

		// Skip if new position is not valid (on the grid)
		if newPos < 0 || newPos >= len(g) {
			continue
		}

		// Skip if source tile is empty
		if g[i].Val == emptyTile {
			continue
		}

		// Combine if similar tile exists at destination and end turn
		alreadyCombined := g[i].Cmb || g[newPos].Cmb
		if g[newPos].Val == g[i].Val && !alreadyCombined {
			g[newPos].Val += g[i].Val // update the new location
			g[newPos].Cmb = true
			g[newPos].UUID = uuid.Must(uuid.NewV7())
			valAfterCombine := g[newPos].Val
			g[i].Val = emptyTile // clear the old location
			g[i].UUID = uuid.Must(uuid.NewV7())
			return g, true, valAfterCombine

		} else if g[newPos].Val != emptyTile {
			// Move blocked by another tile
			continue
		}

		// Destination empty; move tile and end turn
		if g[newPos].Val == emptyTile {
			g[newPos] = g[i]
			g[i] = Tile{UUID: uuid.Must(uuid.NewV7())}
			return g, true, 0
		}
	}

	return g, false, 0
}

// isLoss returns true if the grid is in a losing state (gridlocked).
func (g *Grid) isLoss() bool {
	// False if any empty spaces exist
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			if g.Tiles[i][j].Val == emptyTile {
				return false
			}
		}
	}

	// False if any similar tiles exist next to each other
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth-1; j++ {
			if g.Tiles[i][j].Val == g.Tiles[i][j+1].Val {
				return false
			}
		}
	}
	t := Transpose(g.Tiles)
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth-1; j++ {
			if t[i][j].Val == t[i][j+1].Val {
				return false
			}
		}
	}

	return true
}

// HighestTile returns the value of the highest tile on the grid.
func (g *Grid) HighestTile() int {
	highest := 0
	for a := range gridHeight {
		for b := range gridWidth {
			if g.Tiles[a][b].Val > highest {
				highest = g.Tiles[a][b].Val
			}
		}
	}
	return highest
}

// Debug arranges the grid into a human readable Debug for debugging purposes.
func (g *Grid) Debug() string {
	var out string
	for row := 0; row < gridHeight; row++ {
		for col := range gridWidth {
			out += g.Tiles[row][col].paddedString() + "|"
		}
		out += "\n"
	}
	return out
}

// Transpose returns a transposed version of the grid.
func Transpose(matrix [gridWidth][gridHeight]Tile) [gridHeight][gridWidth]Tile {
	var transposed [gridHeight][gridWidth]Tile
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}
	return transposed
}

// clone returns a deep copy for debugging purposes.
func (g *Grid) clone() *Grid {
	newGrid := &Grid{}
	for a := range gridHeight {
		for b := range gridWidth {
			newGrid.Tiles[a][b] = g.Tiles[a][b]
		}
	}
	return newGrid
}

const emptyTile = 0

// Tile represents a single Tile on the grid.
type Tile struct {
	Val  int       `json:"val"`  // the value of the number on the tile
	Cmb  bool      `json:"cmb"`  // flag for whether tile was combined in the current turn
	UUID uuid.UUID `json:"uuid"` // unique ID for each tile
}

// NewTiles generates a fresh set of tiles.
func NewTiles() [gridHeight][gridWidth]Tile {
	t := [gridHeight][gridWidth]Tile{}
	for i := range t {
		for j := range t[i] {
			t[i][j].UUID = uuid.Must(uuid.NewV7())
		}
	}
	return t
}

// paddedString generates a padded version of the tile's value.
func (t *Tile) paddedString() string {
	s := fmt.Sprintf("%d", t.Val)
	switch len(s) {
	case 1:
		return "   " + s + "   "
	case 2:
		return "   " + s + "  "
	case 3:
		return "  " + s + "  "
	case 4:
		return "  " + s + " "
	case 5:
		return " " + s + " "
	case 6:
		return " " + s
	default:
		return s
	}
}

// Equal returns whether tile t1 is equal to t2.
func (t1 *Tile) Equal(t2 Tile) bool {
	return t1.Val == t2.Val &&
		t1.Cmb == t2.Cmb &&
		t1.UUID == t2.UUID
}

// EqualGrid returns whether grid g1 is equal to g2.
func EqualGrid(g1, g2 [gridWidth][gridHeight]Tile) bool {
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			if !g1[i][j].Equal(g2[i][j]) {
				return false
			}
		}
	}
	return true
}

// newTileVal generates the value of a new tile.
func newTileVal() int {
	if rand.Float64() >= 0.9 {
		return 4
	} else {
		return 2
	}
}
