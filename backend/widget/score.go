package widget

// Score contains the current score and high score.
type Score struct {
	Current int `json:"current"`
	High    int `json:"high"`
}

// NewScore returns the current score widget.
func NewScore() *Score {
	return &Score{}
}

// CurrentScore returns the current score.
func (s *Score) CurrentScore() int {
	return s.Current
}

// AddToCurrent adds a value to the current score.
func (s *Score) AddToCurrent(score int) int {
	s.Current += score
	s.checkHighScore()
	return s.Current
}

// Reset resets the current score.
func (s *Score) Reset() {
	s.Current = 0
}

// checkHighScore refreshes the high score.
func (s *Score) checkHighScore() {
	if s.Current > s.High {
		s.High = s.Current
	}
}
