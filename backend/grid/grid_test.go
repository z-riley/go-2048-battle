package grid

import (
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	input := Grid{
		Tiles: [4][4]Tile{
			{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 2}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
		},
	}
	dir := DirRight
	expected := Grid{
		Tiles: [4][4]Tile{
			{{Val: 0}, {Val: 0}, {Val: 2}, {Val: 4, Cmb: true}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
			{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
		},
	}

	got := input.clone()
	got.move(dir)
	if !gridsAreEqual(expected.Tiles, got.Tiles) {
		t.Errorf("Expected:\n<%v>\nGot:\n<%v>", expected.Debug(), got.Debug())
	}
}

func TestMoveStep(t *testing.T) {
	type tc struct {
		input    [4]Tile
		dir      Direction
		expected [4]Tile
		moved    bool
	}

	for n, tc := range []tc{
		// 2 2 2 2 --[left]--> 4 4 0 0
		{
			input:    [4]Tile{{Val: 2}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4, Cmb: true}, {Val: 0}, {Val: 2}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4, Cmb: true}, {Val: 0}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4, Cmb: true}, {Val: 2}, {Val: 0}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4, Cmb: true}, {Val: 2}, {Val: 0}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4, Cmb: true}, {Val: 2}, {Val: 2}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4, Cmb: true}, {Val: 2}, {Val: 2}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4, Cmb: true}, {Val: 4, Cmb: true}, {Val: 0}, {Val: 0}},
			moved:    true,
		},
		// 0 4 2 2 --[left]--> 4 4 0 0
		{
			input:    [4]Tile{{Val: 0}, {Val: 4}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4}, {Val: 0}, {Val: 2}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4}, {Val: 0}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4}, {Val: 2}, {Val: 0}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4}, {Val: 2}, {Val: 0}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4}, {Val: 2}, {Val: 2}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4}, {Val: 2}, {Val: 2}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4}, {Val: 4, Cmb: true}, {Val: 0}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 4}, {Val: 4, Cmb: true}, {Val: 0}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]Tile{{Val: 4}, {Val: 4, Cmb: true}, {Val: 0}, {Val: 0}},
			moved:    false,
		},
		// // 2 2 2 2 --[right]--> 4 4 0 0
		{
			input:    [4]Tile{{Val: 2}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 2}, {Val: 2}, {Val: 0}, {Val: 4, Cmb: true}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 2}, {Val: 2}, {Val: 0}, {Val: 4, Cmb: true}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 2}, {Val: 0}, {Val: 2}, {Val: 4, Cmb: true}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 2}, {Val: 0}, {Val: 2}, {Val: 4, Cmb: true}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 4, Cmb: true}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 4, Cmb: true}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 0}, {Val: 0}, {Val: 4, Cmb: true}, {Val: 4, Cmb: true}},
			moved:    true,
		},
		// // 0 2 2 2 --[right]--> 0 0 2 4
		{
			input:    [4]Tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 0}, {Val: 2}, {Val: 0}, {Val: 4, Cmb: true}},
			moved:    true,
		},
		{
			input:    [4]Tile{{Val: 0}, {Val: 2}, {Val: 0}, {Val: 4, Cmb: true}},
			dir:      DirRight,
			expected: [4]Tile{{Val: 0}, {Val: 0}, {Val: 2}, {Val: 4, Cmb: true}},
			moved:    true,
		},
	} {
		got, moved, _ := moveStep(tc.input, tc.dir)
		if !rowsAreEqual(tc.expected, got) {
			t.Errorf("[%d] \nExpected:\n<%v>\nGot:\n<%v>", n, tc.expected, got)
		}
		if tc.moved != moved {
			t.Errorf("[%d] \nExpected:\n<%v>\nGot:\n<%v>", n, tc.moved, moved)
		}
	}
}

func TestTranspose(t *testing.T) {
	input := [4][4]Tile{
		{{Val: 1}, {Val: 2}, {Val: 3}, {Val: 4}},
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 6}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 5}},
	}
	expected := [4][4]Tile{
		{{Val: 1}, {Val: 0}, {Val: 6}, {Val: 0}},
		{{Val: 2}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 3}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 4}, {Val: 0}, {Val: 0}, {Val: 5}},
	}
	got := transpose(input)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nExpected:\n<%v>\nGot:\n<%v>", expected, got)
	}
}

func TestIsLoss(t *testing.T) {
	type tc struct {
		input    Grid
		expected bool
	}

	tests := []tc{
		{
			input: Grid{
				Tiles: [4][4]Tile{
					{{Val: 2}, {Val: 0}, {Val: 8}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 4}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
				},
			},
			expected: false,
		},
		{
			input: Grid{
				Tiles: [4][4]Tile{
					{{Val: 4}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
				},
			},
			expected: false,
		},
		{
			input: Grid{
				Tiles: [4][4]Tile{
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
				},
			},
			expected: true,
		},
		{
			input: Grid{
				Tiles: [4][4]Tile{
					{{Val: 2}, {Val: 4}, {Val: 16}, {Val: 2}},
					{{Val: 8}, {Val: 32}, {Val: 64}, {Val: 16}},
					{{Val: 4}, {Val: 16}, {Val: 8}, {Val: 4}},
					{{Val: 2}, {Val: 8}, {Val: 4}, {Val: 2}},
				},
			},
			expected: true,
		},
		{
			input: Grid{
				Tiles: [4][4]Tile{
					{{Val: 4}, {Val: 16}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 32}, {Val: 4}, {Val: 2}},
					{{Val: 4}, {Val: 8}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 8}, {Val: 8}, {Val: 8}},
				},
			},
			expected: false,
		},
	}
	for i := range tests {
		actual := tests[i].input.isLoss()
		if tests[i].expected != actual {
			t.Errorf("Expected:\n<%v>\nGot:\n<%v>", tests[i].expected, actual)
		}
	}
}

// gridsAreEqual checks whether grids are equal, ignoring the UUID fields of tiles.
func gridsAreEqual(grid1, grid2 [4][4]Tile) bool {
	for i := range grid1 {
		if !rowsAreEqual(grid1[i], grid2[i]) {
			return false
		}
	}
	return true
}

// rowsAreEqual checks whether rows of tiles are equal, ignoring the UUID fields.
func rowsAreEqual(row1, row2 [4]Tile) bool {
	for i := range row1 {
		if row1[i].Val != row2[i].Val ||
			row1[i].Cmb != row2[i].Cmb {
			return false
		}
	}
	return true
}
