package screen

import (
	"image/color"

	"github.com/zac460/turdgl"
)

type MultiplayerMenuScreen struct {
	win *turdgl.Window
}

// NewTitle Screen constructs a new multiplayer menu screen for the given window.
func NewMultiplayerMenuScreen(win *turdgl.Window) *MultiplayerMenuScreen {

	return &MultiplayerMenuScreen{
		win,
	}
}

// Update draws the multiplayer menu screen and updates its components. A pointer to
// the current screen state must be passed in so the screen can move to other screens.
func (t *MultiplayerMenuScreen) Update() {
	t.win.SetBackground(color.RGBA{46, 36, 27, 255})

}
