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
	counters    []int
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

	size := int(window / granularity)

	t := &rollingWindow{
		window:      window,
		granularity: granularity,
		size:        size,
		position:    0,
		values:      make([]float64, size),
		counters:    make([]int, size),
		quitC:       make(chan struct{}),
		close:       false,
	}

	ticker := NewTicker(t.granularity)

	go t.cleanBuckets(ticker)

	return t, nil
}

func (r *rollingWindow) cleanBuckets(ticker TimeTicker) {
	for {
		select {
		case <-ticker.Chan():
			r.Lock()

			r.position++

			if r.position >= r.size {
				r.position = 0
			}

			r.values[r.position] = 0
			r.counters[r.position] = 0

			r.Unlock()
		case <-r.quitC:
			ticker.Stop()
			return
		}
	}
}

func (r *rollingWindow) Stop() {
	r.Lock()
	defer r.Unlock()

	r.close = true
	r.quitC <- struct{}{}
}

// Add adds a value to its given position in the rolling window.
// The given value is incremented with the position's value.
func (r *rollingWindow) Add(value float64) {
	r.Lock()
	defer r.Unlock()

	if r.close {
		return
	}

	r.counters[r.position]++
	r.values[r.position] += value
}

// Count returns the total number of transactions in the rolling window
func (r *rollingWindow) Count() int {
	r.RLock()
	defer r.RUnlock()

	total := 0

	for _, v := range r.counters {
		total += v
	}

	return total
}

// Average calculates the average inside rolling window.
func (r *rollingWindow) Average() float64 {
	r.RLock()
	defer r.RUnlock()

	total := 0.0

	for _, v := range r.values {
		total = total + v
	}

	if total == 0 {
		return total
	}

	return total / float64(r.size)
}

// Max returns the max value in the given rolling window.
func (r *rollingWindow) Max() float64 {
	r.RLock()
	defer r.RUnlock()

	max := r.values[0]

	for _, v := range r.values {
		if v > max {
			max = v
		}
	}

	return max
}

// Min returns the min value in the given rolling window.
func (r *rollingWindow) Min() float64 {
	r.RLock()
	defer r.RUnlock()

	min := r.values[0]

	for _, v := range r.values {
		if v <= min {
			min = v
		}
	}

	return min
}
