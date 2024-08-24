package widget

import (
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	type tc struct {
		input    grid
		dir      Direction
		expected grid
	}

	for n, tc := range []tc{
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 2}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
				},
			},
			dir: DirRight,
			expected: grid{
				Tiles: [4][4]tile{
					{{Val: 0}, {Val: 0}, {Val: 2}, {Val: 4}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
				},
			},
		},
	} {
		got := tc.input.clone()
		got.move(DirRight)
		if !reflect.DeepEqual(tc.expected.Tiles, got.Tiles) {
			t.Errorf("[%d] \nExpected:\n<%v>\nGot:\n<%v>", n, tc.expected.Debug(), got.Debug())
		}
	}
}

func TestMoveStep(t *testing.T) {
	type tc struct {
		input    [4]tile
		dir      Direction
		expected [4]tile
		moved    bool
	}

	for n, tc := range []tc{
		// 2 2 2 2 --[left]--> 4 4 0 0
		{
			input:    [4]tile{{Val: 2}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4, cmb: true}, {Val: 0}, {Val: 2}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4, cmb: true}, {Val: 0}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4, cmb: true}, {Val: 2}, {Val: 0}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4, cmb: true}, {Val: 2}, {Val: 0}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4, cmb: true}, {Val: 2}, {Val: 2}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4, cmb: true}, {Val: 2}, {Val: 2}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4, cmb: true}, {Val: 4, cmb: true}, {Val: 0}, {Val: 0}},
			moved:    true,
		},
		// 0 4 2 2 --[left]--> 4 4 0 0
		{
			input:    [4]tile{{Val: 0}, {Val: 4}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4}, {Val: 0}, {Val: 2}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4}, {Val: 0}, {Val: 2}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4}, {Val: 2}, {Val: 0}, {Val: 2}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4}, {Val: 2}, {Val: 0}, {Val: 2}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4}, {Val: 2}, {Val: 2}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4}, {Val: 2}, {Val: 2}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4}, {Val: 4, cmb: true}, {Val: 0}, {Val: 0}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 4}, {Val: 4, cmb: true}, {Val: 0}, {Val: 0}},
			dir:      DirLeft,
			expected: [4]tile{{Val: 4}, {Val: 4, cmb: true}, {Val: 0}, {Val: 0}},
			moved:    false,
		},
		// // 2 2 2 2 --[right]--> 4 4 0 0
		{
			input:    [4]tile{{Val: 2}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirRight,
			expected: [4]tile{{Val: 2}, {Val: 2}, {Val: 0}, {Val: 4, cmb: true}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 2}, {Val: 2}, {Val: 0}, {Val: 4, cmb: true}},
			dir:      DirRight,
			expected: [4]tile{{Val: 2}, {Val: 0}, {Val: 2}, {Val: 4, cmb: true}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 2}, {Val: 0}, {Val: 2}, {Val: 4, cmb: true}},
			dir:      DirRight,
			expected: [4]tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 4, cmb: true}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 4, cmb: true}},
			dir:      DirRight,
			expected: [4]tile{{Val: 0}, {Val: 0}, {Val: 4, cmb: true}, {Val: 4, cmb: true}},
			moved:    true,
		},
		// // 0 2 2 2 --[right]--> 0 0 2 4
		{
			input:    [4]tile{{Val: 0}, {Val: 2}, {Val: 2}, {Val: 2}},
			dir:      DirRight,
			expected: [4]tile{{Val: 0}, {Val: 2}, {Val: 0}, {Val: 4, cmb: true}},
			moved:    true,
		},
		{
			input:    [4]tile{{Val: 0}, {Val: 2}, {Val: 0}, {Val: 4, cmb: true}},
			dir:      DirRight,
			expected: [4]tile{{Val: 0}, {Val: 0}, {Val: 2}, {Val: 4, cmb: true}},
			moved:    true,
		},
	} {
		got, moved, _ := moveStep(tc.input, tc.dir)
		if !reflect.DeepEqual(tc.expected, got) {
			t.Errorf("[%d] \nExpected:\n<%v>\nGot:\n<%v>", n, tc.expected, got)
		}
		if tc.moved != moved {
			t.Errorf("[%d] \nExpected:\n<%v>\nGot:\n<%v>", n, tc.moved, moved)
		}
	}
}

func TestTranspose(t *testing.T) {
	input := [4][4]tile{
		{{Val: 1}, {Val: 2}, {Val: 3}, {Val: 4}},
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 6}, {Val: 0}, {Val: 0}, {Val: 0}},
		{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 5}},
	}
	expected := [4][4]tile{
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
		input    grid
		expected bool
	}

	for _, tc := range []tc{
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 2}, {Val: 0}, {Val: 8}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 4}, {Val: 0}, {Val: 0}},
					{{Val: 0}, {Val: 0}, {Val: 0}, {Val: 0}},
				},
			},
			expected: false,
		},
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 4}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
				},
			},
			expected: false,
		},
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 4}, {Val: 2}, {Val: 4}},
					{{Val: 4}, {Val: 2}, {Val: 4}, {Val: 2}},
				},
			},
			expected: true,
		},
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 2}, {Val: 4}, {Val: 16}, {Val: 2}},
					{{Val: 8}, {Val: 32}, {Val: 64}, {Val: 16}},
					{{Val: 4}, {Val: 16}, {Val: 8}, {Val: 4}},
					{{Val: 2}, {Val: 8}, {Val: 4}, {Val: 2}},
				},
			},
			expected: true,
		},
		{
			input: grid{
				Tiles: [4][4]tile{
					{{Val: 4}, {Val: 16}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 32}, {Val: 4}, {Val: 2}},
					{{Val: 4}, {Val: 8}, {Val: 4}, {Val: 2}},
					{{Val: 2}, {Val: 8}, {Val: 8}, {Val: 8}},
				},
			},
			expected: false,
		},
	} {
		actual := tc.input.isLoss()
		if tc.expected != actual {
			t.Errorf("Expected:\n<%v>\nGot:\n<%v>", tc.expected, actual)
		}
	}
}
