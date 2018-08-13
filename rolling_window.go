package horus

import (
	"fmt"
	"sync"
	"time"
)

type rollingWindow struct {
	sync.RWMutex
	window      time.Duration
	granularity time.Duration
	size        int
	position    int
	values      []float64
	quitC       chan struct{}
	close       bool
}

// NewRWindow provides a slide window operation for a moving average in a rolling window.
// The given window must be in a valid time.Duration higher than 0.
func NewRWindow(window time.Duration, granularity time.Duration) (*rollingWindow, error) {
	if window == 0 {
		return nil, fmt.Errorf("window must be higher than 0")
	}

	if granularity == 0 {
		return nil, fmt.Errorf("granularity must be higher than zero")
	}

	if granularity > window || window%granularity != 0 {
		return nil, fmt.Errorf("window must be a multiplier of granularity")
	}

	t := &rollingWindow{
		window:      window,
		granularity: granularity,
		size:        int(window / granularity),
		position:    0,
		values:      make([]float64, int(window/granularity)),
		quitC:       make(chan struct{}),
		close:       false,
	}

	ticker := NewTicker(t.granularity)

	go t.cleanBuckets(ticker)

	return t, nil
}

func (t *rollingWindow) cleanBuckets(ticker TimeTicker) {
	for {
		select {
		case <-ticker.Chan():
			t.Lock()

			t.position++

			if t.position >= t.size {
				t.position = 0
			}

			t.values[t.position] = 0

			t.Unlock()
		case <-t.quitC:
			ticker.Stop()
			return
		}
	}
}

func (t *rollingWindow) Stop() {
	t.Lock()
	defer t.Unlock()

	t.close = true
	t.quitC <- struct{}{}
}

// Add adds a value to its given position in the rolling window.
// The given value is incremented with the position's value.
func (t *rollingWindow) Add(value float64) {
	t.Lock()
	defer t.Unlock()

	if t.close {
		return
	}

	t.values[t.position] += value
}

// Average calculates the average inside rolling window.
func (t *rollingWindow) Average() float64 {
	t.RLock()
	defer t.RUnlock()

	total := 0.0

	for _, v := range t.values {
		total = total + v
	}

	if total == 0 {
		return total
	}

	return total / float64(t.size)
}

// Max returns the max value in the given rolling window.
func (t *rollingWindow) Max() float64 {
	t.RLock()
	defer t.RUnlock()

	max := t.values[0]

	for _, v := range t.values {
		if v > max {
			max = v
		}
	}

	return max
}

// Min returns the min value in the given rolling window.
func (t *rollingWindow) Min() float64 {
	t.RLock()
	defer t.RUnlock()

	min := t.values[0]

	for _, v := range t.values {
		if v <= min {
			min = v
		}
	}

	return min
}
