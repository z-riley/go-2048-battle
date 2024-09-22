package common

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/z-riley/go-2048-battle/backend/grid"
)

func TestGetTileMovements(t *testing.T) {
	// Generate example UUIDs to use
	id := make([]uuid.UUID, 4*4)
	for i := range id {
		id[i] = uuid.Must(uuid.NewV7())
	}

	before := [4][4]grid.Tile{
		{{Val: 0, UUID: id[0]}, {Val: 2, UUID: id[1]}, {Val: 2, UUID: id[2]}, {Val: 2, UUID: id[3]}},
		{{Val: 2, UUID: id[4]}, {Val: 0, UUID: id[5]}, {Val: 0, UUID: id[6]}, {Val: 2, UUID: id[7]}},
		{{Val: 0, UUID: id[8]}, {Val: 4, UUID: id[9]}, {Val: 2, UUID: id[10]}, {Val: 2, UUID: id[11]}},
		{{Val: 0, UUID: id[12]}, {Val: 0, UUID: id[13]}, {Val: 2, UUID: id[14]}, {Val: 2, UUID: id[15]}},
	}
	// after downwards move
	after := [4][4]grid.Tile{
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 2}}, // <-- new tile
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 0}, {Val: 2, UUID: id[1]}, {Val: 2, UUID: id[2]}, {Val: 4, UUID: uuid.Must(uuid.NewV7()), Cmb: true}},
		{{Val: 2, UUID: id[4]}, {Val: 4, UUID: id[9]}, {Val: 4, UUID: uuid.Must(uuid.NewV7()), Cmb: true}, {Val: 4, UUID: uuid.Must(uuid.NewV7()), Cmb: true}},
	}

	moves := generateAnimations(before, after, grid.DirDown)
	fmt.Println("Animations:")
	for _, move := range moves {
		fmt.Printf("%+v\n", move)
	}
}
