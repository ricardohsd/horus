package horus

import "fmt"

// SWindow defines attributes for a fixed slide window
type SWindow struct {
	window int
	values []float64
}

// NewSWindow provides a fixed slide window operation that
// exposes average, max and min statistics.
// The given window must be an integer higher than 0.
func NewSWindow(window int) (*SWindow, error) {
	if window == 0 {
		return nil, fmt.Errorf("window must be higher than 0")
	}

	return &SWindow{
		window: window,
		values: make([]float64, window),
	}, nil
}

func (s *SWindow) Add(value float64) {
	s.values = append(s.values[1:s.window], value)
}

func (s *SWindow) Average() float64 {
	avg := 0.0

	for _, v := range s.values {
		avg = avg + v
	}

	return avg / float64(s.window)
}

func (s *SWindow) Max() float64 {
	max := s.values[0]

	for _, v := range s.values {
		if v > max {
			max = v
		}
	}

	return max
}

func (s *SWindow) Min() float64 {
	min := s.values[0]

	for _, v := range s.values {
		if v <= min {
			min = v
		}
	}

	return min
}
