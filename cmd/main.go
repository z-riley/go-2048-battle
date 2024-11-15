package main

import (
	"flag"
	"os"

	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/debug"
	"github.com/z-riley/go-2048-battle/log"
	"github.com/z-riley/go-2048-battle/screens"
	"github.com/z-riley/turdgl"
)

func main() {
	const path = "./assets/icon.png"
	icon, err := os.Open(path)
	if err != nil {
		log.Println("Failed to load window icon:", path)
	}

	// Create window
	win, err := turdgl.NewWindow(turdgl.WindowCfg{
		Title:  "2048 Battle",
		Width:  config.WinWidth,
		Height: config.WinHeight,
		Icon:   icon,
	})
	if err != nil {
		log.Fatalf("Failed to create new window: %v", err)
	}
	defer win.Destroy()

	if config.Debug {
		win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })
	}

	// Parse starting screen arg
	screenStr := flag.String("screen", string(screens.Title), "starting screen")
	flag.Parse()

	// Create screens
	screens.Init(win)
	screens.SetScreen(screens.ID(*screenStr), nil)

	debugWidget := debug.NewDebugWidget(win)

	// Main game loop
	for win.IsRunning() {
		screens.Update()

		if config.Debug {
			// Add debug overlay
			debugWidget.Update()
		}

		win.Update()
	}
}
