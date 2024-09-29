package main

import (
	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/debug"
	"github.com/z-riley/go-2048-battle/screen"
	"github.com/z-riley/turdgl"
)

func main() {
	// Create window
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "2048 Battle",
		Width:  config.WinWidth,
		Height: config.WinHeight,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Register window-level keybinds (for development only)
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })

	// Create screens
	screen.Init(win)
	screen.SetScreen(screen.Title, nil)
	screen.SetScreen(screen.Singleplayer, nil)

	debugWidget := debug.NewDebugWidget(win)

	// Main game loop
	for win.IsRunning() {
		// Update screen
		screen.CurrentScreen().Update()

		if config.Debug {
			// Add debug overlay
			debugWidget.Update()
		}

		win.Update()
	}
}
