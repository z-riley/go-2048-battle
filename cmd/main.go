package main

import (
	"flag"
	"log"
	"os"

	"github.com/z-riley/go-2048-battle/config"
	"github.com/z-riley/go-2048-battle/debug"
	"github.com/z-riley/go-2048-battle/screen"
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

	// Register window-level keybinds (for development only)
	win.RegisterKeybind(turdgl.KeyLCtrl, turdgl.KeyPress, func() { win.Quit() })

	// Parse starting screen arg
	screenStr := flag.String("screen", string(screen.Title), "starting screen")
	flag.Parse()
	start := screen.ID(*screenStr)

	// Create screens
	screen.Init(win)
	screen.SetScreen(start, nil)

	debugWidget := debug.NewDebugWidget(win)

	// Main game loop
	for win.IsRunning() {
		screen.Update()

		if config.Debug {
			// Add debug overlay
			debugWidget.Update()
		}

		win.Update()
	}
}
