package widget

import (
	"fmt"
	"time"
)

type Timer struct {
	Time     time.Duration `json:"time"`
	isPaused bool
}

// NewTimer constructs a new timer. Resume must be called to start the timer.
func NewTimer() *Timer {
	t := Timer{
		Time:     time.Duration(0),
		isPaused: true,
	}

	go func() {
		for {
			time.Sleep(time.Second)
			if !t.isPaused {
				t.Time += time.Second
			}
		}
	}()

	return &t
}

// Resume resumes the timer.
func (t *Timer) Resume() *Timer {
	t.isPaused = false
	return t
}

// Pause pauses the timer.
func (t *Timer) Pause() *Timer {
	t.isPaused = true
	return t
}

// Reset starts the timer from zero.
func (t *Timer) Reset() *Timer {
	t.Time = 0
	return t
}

// Set sets the timer to the specified duration
func (t *Timer) Set(d time.Duration) *Timer {
	t.Time = d
	return t
}

// Duration returns the current time in time.Duration format.
func (t *Timer) Duration() time.Duration {
	return t.Time
}

// format formats a duration into the format "HH:MM:SS".
func format(t time.Duration) string {
	hours := int(t.Hours())
	minutes := int(t.Minutes()) % 60
	seconds := int(t.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
