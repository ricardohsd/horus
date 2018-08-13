package horus

import "fmt"

type sma struct {
	window int
	values []float64
}

// NewSMA provides a slide window operation for Simple Moving Average.
// The given window must be an integer higher than 0.
func NewSMA(window int) (*sma, error) {
	if window == 0 {
		return nil, fmt.Errorf("window must be higher than 0")
	}

	return &sma{
		window: window,
		values: make([]float64, window),
	}, nil
}

func (s *sma) Add(value float64) {
	s.values = append(s.values[1:s.window], value)
}

func (s *sma) Average() float64 {
	avg := float64(0)
	for i := s.window; i > 0; i-- {
		avg = s.values[i-1] + avg
	}
	return avg / float64(s.window)
}
