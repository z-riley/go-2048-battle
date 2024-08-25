package common

import (
	"encoding/json"
	"fmt"

	game "github.com/z-riley/go-2048-battle"
	"github.com/z-riley/go-2048-battle/backend"
	"github.com/z-riley/turdgl"
)

const (
	tileSizePx float64 = 60
	arenaSize          = 4
)

// Arena displays the grid of a game.
type Arena struct {
	state backend.Game
	tiles [arenaSize][arenaSize]*turdgl.TextBox
}

// NewArena constructs a new arena widget.
func NewArena() *Arena {
	// Populate tiles
	tiles := [arenaSize][arenaSize]*turdgl.TextBox{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			tiles[j][i] = turdgl.NewTextBox(
				turdgl.NewRect(
					tileSizePx, tileSizePx,
					turdgl.Vec{
						X: 700 + float64(j)*tileSizePx*1.2,
						Y: 80 + float64(i)*tileSizePx*1.2,
					},
				),
				game.FontPath).SetTextAlignment(turdgl.AlignTopCentre)
		}
	}

	return &Arena{
		tiles: tiles,
	}
}

// Draw draws the arena to the frame buffer.
func (a *Arena) Draw(buf *turdgl.FrameBuffer) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			a.tiles[j][i].Draw(buf)
		}
	}
}

// Load updates the displayed arena to match the data provided.
func (a *Arena) Load(data []byte) {
	var game backend.Game
	if err := json.Unmarshal(data, &game); err != nil {
		panic(err)
	}

	// Update text on tiles
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			val := game.Grid.Tiles[i][j].Val
			if val == 0 {
				a.tiles[j][i].SetText("")

			} else {
				a.tiles[j][i].SetText(fmt.Sprint(val))
			}
		}
	}
}
