package common

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/z-riley/go-2048-battle/backend/grid"
)

func TestGenerateRowAnimations(t *testing.T) {
	type tc struct {
		name          string
		before, after [numTiles]grid.Tile
		dir           grid.Direction
		want          []rowAnimation
	}

	// Generate some UUIDs to use
	id := make([]uuid.UUID, 4*4)
	for i := range id {
		id[i] = uuid.Must(uuid.NewV7())
	}

	for _, tc := range []tc{
		{
			name: "Moving tiles",
			before: [4]grid.Tile{
				{Val: 0, Cmb: false, UUID: id[0]},
				{Val: 2, Cmb: false, UUID: id[1]},
				{Val: 0, Cmb: false, UUID: id[2]},
				{Val: 4, Cmb: false, UUID: id[3]},
			},
			after: [4]grid.Tile{
				{Val: 2, Cmb: false, UUID: id[1]},
				{Val: 4, Cmb: false, UUID: id[3]},
				{Val: 0, Cmb: false, UUID: id[6]},
				{Val: 0, Cmb: false, UUID: id[7]},
			},
			dir: grid.DirLeft,
			want: []rowAnimation{
				moveRowAnimation{
					origin: 1,
					dest:   0,
				},
				moveRowAnimation{
					origin: 3,
					dest:   1,
				},
			},
		},
		{
			name: "Combining 2 tiles with spawn",
			before: [4]grid.Tile{
				{Val: 0, Cmb: false, UUID: id[0]},
				{Val: 2, Cmb: false, UUID: id[1]},
				{Val: 0, Cmb: false, UUID: id[2]},
				{Val: 2, Cmb: false, UUID: id[3]},
			},
			after: [4]grid.Tile{
				{Val: 4, Cmb: true, UUID: id[4]},
				{Val: 0, Cmb: false, UUID: id[5]},
				{Val: 0, Cmb: false, UUID: id[6]},
				{Val: 2, Cmb: false, UUID: id[7]},
			},
			dir: grid.DirUp,
			want: []rowAnimation{
				newFromCombineRowAnimation{
					dest:   0,
					newVal: 4,
				},
				moveToCombineRowAnimation{
					origin: 1,
					dest:   0,
				},
				moveToCombineRowAnimation{
					origin: 3,
					dest:   0,
				},
				spawnRowAnimation{
					dest:   3,
					newVal: 2,
				},
			},
		},
		{
			name: "Combining 2 sets of 2 tiles",
			before: [4]grid.Tile{
				{Val: 2, Cmb: false, UUID: id[0]},
				{Val: 2, Cmb: false, UUID: id[1]},
				{Val: 4, Cmb: false, UUID: id[2]},
				{Val: 4, Cmb: false, UUID: id[3]},
			},
			after: [4]grid.Tile{
				{Val: 0, Cmb: false, UUID: id[4]},
				{Val: 0, Cmb: false, UUID: id[5]},
				{Val: 4, Cmb: true, UUID: id[6]},
				{Val: 8, Cmb: true, UUID: id[7]},
			},
			dir: grid.DirRight,
			want: []rowAnimation{
				newFromCombineRowAnimation{
					dest:   2,
					newVal: 4,
				},
				moveToCombineRowAnimation{
					origin: 0,
					dest:   2,
				},
				moveToCombineRowAnimation{
					origin: 1,
					dest:   2,
				},
				newFromCombineRowAnimation{
					dest:   3,
					newVal: 8,
				},
				moveToCombineRowAnimation{
					origin: 2,
					dest:   3,
				},
			},
		},
		{
			name: "Combining 4 similar tiles",
			before: [4]grid.Tile{
				{Val: 2, Cmb: false, UUID: id[0]},
				{Val: 2, Cmb: false, UUID: id[1]},
				{Val: 2, Cmb: false, UUID: id[2]},
				{Val: 2, Cmb: false, UUID: id[3]},
			},
			after: [4]grid.Tile{
				{Val: 0, Cmb: false, UUID: id[4]},
				{Val: 0, Cmb: false, UUID: id[5]},
				{Val: 4, Cmb: true, UUID: id[6]},
				{Val: 8, Cmb: true, UUID: id[7]},
			},
			dir: grid.DirDown,
			want: []rowAnimation{
				newFromCombineRowAnimation{
					dest:   2,
					newVal: 4,
				},
				moveToCombineRowAnimation{
					origin: 0,
					dest:   2,
				},
				moveToCombineRowAnimation{
					origin: 1,
					dest:   2,
				},
				newFromCombineRowAnimation{
					dest:   3,
					newVal: 8,
				},
				moveToCombineRowAnimation{
					origin: 2,
					dest:   3,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := generateRowAnimations(tc.before, tc.after, tc.dir)
			if !reflect.DeepEqual(sliceToMap(got), sliceToMap(tc.want)) {
				t.Fatalf("Got %v, want %v", got, tc.want)
			}
		})
	}

}

// sliceToMap converts a slice into a map where each key is an element of the slice.
func sliceToMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, element := range slice {
		m[element] = struct{}{}
	}
	return m
}
