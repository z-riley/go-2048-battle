package backend

import (
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

// Reset sets the timer to zero.
func (t *Timer) Reset() *Timer {
	t.Time = 0
	return t
}

// Set sets the timer to the specified duration.
func (t *Timer) Set(d time.Duration) *Timer {
	t.Time = d
	return t
}

// Duration returns the current time.
func (t *Timer) Duration() time.Duration {
	return t.Time
}
