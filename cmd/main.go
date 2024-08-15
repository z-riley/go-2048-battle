package main

import (
	"fmt"

	game "github.com/zac460/go-2048-battle"
	"github.com/zac460/go-2048-battle/debug"
	"github.com/zac460/go-2048-battle/screens"
	"github.com/zac460/turdgl"
)

type Screen int

const (
	TitleScreen Screen = iota
	MultiplayerMenuScreen
)

type UpdateFunc func()

var screen = TitleScreen

func main() {
	// Creat window
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "2048 Battle",
		Width:  game.Width,
		Height: game.Height,
	})
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	// Register window-level keybinds (for development only)
	win.RegisterKeybind(turdgl.KeyEscape, func() { win.Quit() })
	win.RegisterKeybind(turdgl.KeyLCtrl, func() { win.Quit() })

	// Creat screens
	titleScreen := screens.NewTitleScreen(win)
	debugWidget := debug.NewDebugWidget(win)

	// Main game loop
	for win.IsRunning() {
		switch screen {
		case TitleScreen:
			titleScreen.Update()
		case MultiplayerMenuScreen:
			// TODO
		default:
			panic(fmt.Sprint("unsupported screen:", screen))
		}
		if game.Debug {
			debugWidget.Update()
		}

		win.Update()
	}
}
