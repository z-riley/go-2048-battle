package debug

import (
	"fmt"
	"time"

	"github.com/z-riley/go-2048-battle/common"
	"github.com/z-riley/gogl"
)

// Widget displays diagnostic information about the current window.
type Widget struct {
	win *gogl.Window

	location *gogl.Text
	fps      *gogl.Text
	frames   int              // for measuring FPS
	tick     <-chan time.Time // for measuring FPS
}

// NewDebugWidget constructs a new debug widget.
func NewDebugWidget(win *gogl.Window) *Widget {
	location := gogl.NewText("Loc: ", gogl.Vec{X: 1120, Y: 25}, common.FontPathMedium).
		SetAlignment(gogl.AlignBottomRight).
		SetSize(12)

	fps := gogl.NewText("FPS: -", gogl.Vec{X: 1180, Y: 25}, common.FontPathMedium).
		SetAlignment(gogl.AlignBottomRight).
		SetSize(12)

	return &Widget{
		win:      win,
		location: location,
		fps:      fps,
		frames:   0,
		tick:     time.Tick(time.Second),
	}
}

// Update updates and draws the debug widget.
func (t *Widget) Update() {
	t.frames++
	select {
	case <-t.tick:
		t.fps.SetText(fmt.Sprint("FPS: ", t.frames))
		t.frames = 0
	default:
	}

	t.location.SetText(fmt.Sprint("Loc: ", t.win.MouseLocation()))

	t.win.Draw(t.fps)
	t.win.Draw(t.location)
}
