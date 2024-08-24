package widget

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"github.com/z-riley/go-2048-battle/backend/util"
)

const (
	gridWidth  = 4
	gridHeight = 4
)

// grid is the grid arena for the game. Position {0,0} is the top left square.
type grid struct {
	Tiles [gridWidth][gridHeight]tile `json:"tiles"`
}

// newGrid constructs a new grid.
func newGrid() *grid {
	g := grid{Tiles: [4][4]tile{}}
	g.resetGrid()

	return &g
}

// resetGrid resets the grid to a start-of-game state, spawning two '2' tiles in random locations.
func (g *grid) resetGrid() {
	g.Tiles = [gridWidth][gridHeight]tile{}
	// Place two '2' tiles in random positions
	type pos struct{ x, y int }
	tile1 := pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	tile2 := pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	for reflect.DeepEqual(tile1, tile2) {
		// Try again until they're unique
		tile2 = pos{rand.Intn(gridWidth), rand.Intn(gridHeight)}
	}
	g.Tiles[tile1.x][tile1.y].Val = 2
	g.Tiles[tile2.x][tile2.y].Val = 2
}

// spawnTile spawns a single new tile in a random location on the grid. The value of the
// tile is either 2 (90% chance) or 4 (10% chance).
func (g *grid) spawnTile() {
	val := 2
	if rand.Float64() >= 0.9 {
		val = 4
	}

	x, y := rand.Intn(gridWidth), rand.Intn(gridHeight)
	for g.Tiles[x][y].Val != emptyTile {
		// Try again until they're unique
		x, y = rand.Intn(gridWidth), rand.Intn(gridHeight)
	}

	g.Tiles[x][y].Val = val
}

// move attempts to move all tiles in the specified direction, combining them if appropriate.
// Returns whether true if any tiles were moved from the attempt, and the added score from any combinations.
func (g *grid) move(dir Direction) (bool, int) {
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
				g.Tiles = transpose(g.Tiles)
			}
			g.Tiles[row], rowMoved, points = moveStep(g.Tiles[row], dir)
			if points > 0 {
				pointsGained = points
			}
			if dir == DirUp || dir == DirDown {
				g.Tiles = transpose(g.Tiles)
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

	// Clear all of the "combined this turn" flags
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			g.Tiles[i][j].cmb = false
		}
	}

	return moved, pointsGained
}

// moveStep executes one part of the a move. Call multiple times until false is returned to
// complete a full move. Optional: variable to place the number of points gained by the step.
func moveStep(g [gridWidth]tile, dir Direction) ([gridWidth]tile, bool, int) {
	// Iterate in the same direction as the move
	reverse := false
	if dir == DirRight || dir == DirDown {
		reverse = true
	}
	iter := util.NewIter(len(g), reverse)

	for iter.HasNext() {
		// Calculate the hypothetical next position for the tile
		i := iter.Next()
		var newPos int
		if !reverse {
			newPos = i - 1
		} else {
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
		alreadyCombined := g[i].cmb || g[newPos].cmb
		if g[newPos].Val == g[i].Val && !alreadyCombined {
			g[newPos].Val += g[i].Val // update the new location
			g[newPos].cmb = true
			valAfterCombine := g[newPos].Val
			g[i].Val = emptyTile // clear the old location
			return g, true, valAfterCombine

		} else if g[newPos].Val != emptyTile {
			// Move blocked by another tile
			continue
		}

		// Destination empty; move tile and end turn
		if g[newPos].Val == emptyTile {
			g[newPos] = g[i]
			g[i] = tile{}
			return g, true, 0
		}
	}

	return g, false, 0
}

// isLoss returns true if the grid is in a losing state (gridlocked).
func (g *grid) isLoss() bool {
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
	t := transpose(g.Tiles)
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth-1; j++ {
			if t[i][j].Val == t[i][j+1].Val {
				return false
			}
		}
	}

	return true
}

// highestTile returns the value of the highest tile on the grid.
func (g *grid) highestTile() int {
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
func (g *grid) Debug() string {
	var out string
	for row := 0; row < gridHeight; row++ {
		for col := range gridWidth {
			out += g.Tiles[row][col].paddedString() + "|"
		}
		out += "\n"
	}
	return out
}

// transpose returns a transposed version of the grid.
func transpose(matrix [gridWidth][gridHeight]tile) [gridHeight][gridWidth]tile {
	var transposed [gridHeight][gridWidth]tile
	for i := 0; i < gridWidth; i++ {
		for j := 0; j < gridHeight; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}
	return transposed
}

// clone returns a deep copy for debugging purposes.
func (g *grid) clone() *grid {
	newGrid := &grid{}
	for a := range gridHeight {
		for b := range gridWidth {
			newGrid.Tiles[a][b] = g.Tiles[a][b]
		}
	}
	return newGrid
}

const emptyTile = 0

// tile represents a single tile on the grid.
type tile struct {
	Val int  `json:"val"` // the value of the number on the tile
	cmb bool // flag for whether tile was combined in the current turn
}

// paddedString generates a padded version of the tile's value.
func (t *tile) paddedString() string {
	if t.Val == 0 {
		return strings.Repeat(" ", 3)
	}

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
